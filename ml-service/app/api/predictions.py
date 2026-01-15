"""
Prediction API Endpoints

Provides endpoints for making match outcome predictions
"""
import os
import pickle
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import List, Dict, Any, Optional
from datetime import datetime

from config.config import config
from app.database.connection import get_fixtures_for_training, get_all_fixtures
from app.features.feature_builder import extract_features_for_fixture

router = APIRouter()

# Global model cache
_model_cache = None


def get_model():
    """Load and cache the trained model"""
    global _model_cache

    if _model_cache is not None:
        return _model_cache

    model_path = os.path.join(config.MODEL_PATH, 'xgboost_v1.pkl')

    if not os.path.exists(model_path):
        raise HTTPException(
            status_code=503,
            detail=f"Model not found at {model_path}. Please train the model first."
        )

    with open(model_path, 'rb') as f:
        _model_cache = pickle.load(f)

    print(f"[OK] Model loaded: {_model_cache.get('model_version', 'unknown')}")
    return _model_cache


def clear_model_cache():
    """Clear model cache (useful after retraining)"""
    global _model_cache
    _model_cache = None


# Request/Response Models
class PredictionRequest(BaseModel):
    """Single prediction request"""
    home_team_id: int
    away_team_id: int
    match_date: str  # ISO format date string
    fixture_id: Optional[int] = None


class BatchPredictionRequest(BaseModel):
    """Batch prediction request"""
    fixtures: List[PredictionRequest]


class PredictionResponse(BaseModel):
    """Prediction response"""
    fixture_id: Optional[int]
    home_team_id: int
    away_team_id: int
    model_version: str
    predictions: Dict[str, float]
    predicted_outcome: str
    confidence: float
    features_used: int
    predicted_at: str


class BatchPredictionResponse(BaseModel):
    """Batch prediction response"""
    predictions: List[PredictionResponse]
    model_version: str
    count: int
    predicted_at: str


class ModelMetricsResponse(BaseModel):
    """Model performance metrics"""
    model_version: str
    training_date: str
    accuracy: float
    baseline_accuracy: float
    improvement: float
    config_name: Optional[str]
    feature_count: int


def make_prediction(
    home_team_id: int,
    away_team_id: int,
    match_date: datetime,
    fixture_id: Optional[int] = None
) -> Dict[str, Any]:
    """
    Make a prediction for a single fixture

    Args:
        home_team_id: Home team database ID
        away_team_id: Away team database ID
        match_date: Match date
        fixture_id: Optional fixture ID

    Returns:
        Prediction dictionary
    """
    # Load model
    model_data = get_model()
    model = model_data['model']
    feature_names = model_data['feature_names']
    model_version = model_data.get('model_version', 'v1.0')

    # Get historical fixtures for context
    seasons = [2022, 2023, 2024]
    all_fixtures = get_fixtures_for_training(seasons)

    # Create a fixture dict for feature extraction
    fixture = {
        'id': fixture_id or 0,
        'home_team_id': home_team_id,
        'away_team_id': away_team_id,
        'match_date': match_date,
        'season': match_date.year if match_date.month >= 8 else match_date.year - 1,
        'home_score': None,
        'away_score': None,
        'status': 'NS',  # Not started
    }

    # Extract features
    features = extract_features_for_fixture(fixture, all_fixtures)

    # Prepare feature vector
    feature_vector = []
    for col in feature_names:
        value = features.get(col, 0.0)
        if value is None:
            value = 0.0
        feature_vector.append(float(value))

    # Make prediction
    import numpy as np
    X = np.array([feature_vector])
    probabilities = model.predict_proba(X)[0]

    # Map to outcomes
    outcome_map = {0: 'home_win', 1: 'draw', 2: 'away_win'}
    predicted_class = int(np.argmax(probabilities))
    predicted_outcome = outcome_map[predicted_class]
    confidence = float(probabilities[predicted_class])

    return {
        'fixture_id': fixture_id,
        'home_team_id': home_team_id,
        'away_team_id': away_team_id,
        'model_version': model_version,
        'predictions': {
            'home_win_prob': float(probabilities[0]),
            'draw_prob': float(probabilities[1]),
            'away_win_prob': float(probabilities[2]),
        },
        'predicted_outcome': predicted_outcome,
        'confidence': confidence,
        'features_used': len(feature_names),
        'predicted_at': datetime.now().isoformat(),
    }


@router.post("/predict", response_model=PredictionResponse)
async def predict(request: PredictionRequest):
    """
    Generate prediction for a single fixture

    - **home_team_id**: Database ID of home team
    - **away_team_id**: Database ID of away team
    - **match_date**: Match date in ISO format (YYYY-MM-DD)
    """
    try:
        match_date = datetime.fromisoformat(request.match_date.replace('Z', '+00:00'))
    except ValueError:
        raise HTTPException(status_code=400, detail="Invalid date format. Use ISO format (YYYY-MM-DD)")

    try:
        result = make_prediction(
            home_team_id=request.home_team_id,
            away_team_id=request.away_team_id,
            match_date=match_date,
            fixture_id=request.fixture_id
        )
        return PredictionResponse(**result)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/predict/batch", response_model=BatchPredictionResponse)
async def predict_batch(request: BatchPredictionRequest):
    """
    Generate predictions for multiple fixtures

    Accepts a list of fixtures and returns predictions for each.
    """
    predictions = []

    for fixture in request.fixtures:
        try:
            match_date = datetime.fromisoformat(fixture.match_date.replace('Z', '+00:00'))
            result = make_prediction(
                home_team_id=fixture.home_team_id,
                away_team_id=fixture.away_team_id,
                match_date=match_date,
                fixture_id=fixture.fixture_id
            )
            predictions.append(result)
        except Exception as e:
            # Include error in response but continue with other fixtures
            predictions.append({
                'fixture_id': fixture.fixture_id,
                'home_team_id': fixture.home_team_id,
                'away_team_id': fixture.away_team_id,
                'model_version': 'error',
                'predictions': {'error': str(e)},
                'predicted_outcome': 'error',
                'confidence': 0.0,
                'features_used': 0,
                'predicted_at': datetime.now().isoformat(),
            })

    model_data = get_model()

    return BatchPredictionResponse(
        predictions=predictions,
        model_version=model_data.get('model_version', 'v1.0'),
        count=len(predictions),
        predicted_at=datetime.now().isoformat()
    )


@router.get("/model/metrics", response_model=ModelMetricsResponse)
async def get_model_metrics():
    """
    Get current model performance metrics

    Returns accuracy, baseline comparison, and configuration details.
    """
    try:
        model_data = get_model()
        metrics = model_data.get('metrics', {})

        return ModelMetricsResponse(
            model_version=model_data.get('model_version', 'v1.0'),
            training_date=model_data.get('training_date', datetime.now().isoformat()),
            accuracy=metrics.get('accuracy', 0.0),
            baseline_accuracy=metrics.get('baseline_accuracy', 0.39),
            improvement=metrics.get('accuracy', 0.0) - metrics.get('baseline_accuracy', 0.39),
            config_name=metrics.get('config_name'),
            feature_count=len(model_data.get('feature_names', []))
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/model/reload")
async def reload_model():
    """
    Reload the model from disk

    Use this after retraining to load the new model.
    """
    clear_model_cache()
    model_data = get_model()

    return {
        "status": "reloaded",
        "model_version": model_data.get('model_version', 'v1.0'),
        "feature_count": len(model_data.get('feature_names', []))
    }
