"""
Prediction API Endpoints

Provides endpoints for making match outcome predictions across multiple markets:
- 1X2 (Match Result): Home Win / Draw / Away Win
- Over/Under 2.5 Goals
- BTTS (Both Teams to Score)
"""
import os
import pickle
from fastapi import APIRouter, HTTPException, Query
from pydantic import BaseModel
from typing import List, Dict, Any, Optional
from datetime import datetime
import numpy as np

from config.config import config
from app.database.connection import get_fixtures_for_training
from app.features.feature_builder import extract_features_for_fixture

router = APIRouter()

# Model cache for each market
_model_cache: Dict[str, Any] = {}

# Market configuration
MARKETS = {
    '1x2': {
        'model_file': 'xgboost_v1.pkl',
        'description': 'Match Result (Home/Draw/Away)',
        'outcomes': ['home_win', 'draw', 'away_win'],
    },
    'over_under': {
        'model_file': 'xgboost_over_under_model.pkl',
        'description': 'Over/Under 2.5 Goals',
        'outcomes': ['under_2_5', 'over_2_5'],
    },
    'btts': {
        'model_file': 'xgboost_btts_model.pkl',
        'description': 'Both Teams to Score',
        'outcomes': ['no', 'yes'],
    }
}


def get_model(market: str = '1x2'):
    """Load and cache the trained model for a specific market"""
    global _model_cache

    if market not in MARKETS:
        raise HTTPException(
            status_code=400,
            detail=f"Unknown market: {market}. Available: {list(MARKETS.keys())}"
        )

    if market in _model_cache:
        return _model_cache[market]

    model_file = MARKETS[market]['model_file']
    model_path = os.path.join(config.MODEL_PATH, model_file)

    if not os.path.exists(model_path):
        raise HTTPException(
            status_code=503,
            detail=f"Model not found for {market} at {model_path}. Please train the model first."
        )

    with open(model_path, 'rb') as f:
        _model_cache[market] = pickle.load(f)

    print(f"[OK] Model loaded for {market}: {_model_cache[market].get('version', 'v1.0')}")
    return _model_cache[market]


def get_all_models():
    """Load all available models"""
    available = {}
    for market in MARKETS:
        try:
            get_model(market)
            available[market] = True
        except HTTPException:
            available[market] = False
    return available


def clear_model_cache(market: Optional[str] = None):
    """Clear model cache"""
    global _model_cache
    if market:
        _model_cache.pop(market, None)
    else:
        _model_cache = {}


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


class MarketPrediction(BaseModel):
    """Prediction for a single market"""
    market: str
    description: str
    probabilities: Dict[str, float]
    predicted_outcome: str
    confidence: float


class MultiMarketPredictionResponse(BaseModel):
    """Multi-market prediction response"""
    fixture_id: Optional[int]
    home_team_id: int
    away_team_id: int
    match_date: str
    predictions: Dict[str, MarketPrediction]
    features_used: int
    predicted_at: str


class PredictionResponse(BaseModel):
    """Legacy single prediction response (1X2 only)"""
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


class AllModelsMetricsResponse(BaseModel):
    """Metrics for all market models"""
    markets: Dict[str, ModelMetricsResponse]
    available_markets: List[str]


def extract_fixture_features(home_team_id: int, away_team_id: int, match_date: datetime, fixture_id: Optional[int] = None):
    """Extract features for a fixture"""
    seasons = [2022, 2023, 2024]
    all_fixtures = get_fixtures_for_training(seasons)

    fixture = {
        'id': fixture_id or 0,
        'home_team_id': home_team_id,
        'away_team_id': away_team_id,
        'match_date': match_date,
        'season': match_date.year if match_date.month >= 8 else match_date.year - 1,
        'home_score': None,
        'away_score': None,
        'status': 'NS',
    }

    return extract_features_for_fixture(fixture, all_fixtures)


def make_1x2_prediction(features: Dict, model_data: Dict) -> MarketPrediction:
    """Make 1X2 (match result) prediction"""
    model = model_data['model']
    feature_names = model_data['feature_names']

    feature_vector = [float(features.get(col, 0.0) or 0.0) for col in feature_names]
    X = np.array([feature_vector])

    probabilities = model.predict_proba(X)[0]

    outcome_map = {0: 'home_win', 1: 'draw', 2: 'away_win'}
    predicted_class = int(np.argmax(probabilities))

    return MarketPrediction(
        market='1x2',
        description='Match Result (Home/Draw/Away)',
        probabilities={
            'home_win': float(probabilities[0]),
            'draw': float(probabilities[1]),
            'away_win': float(probabilities[2]),
        },
        predicted_outcome=outcome_map[predicted_class],
        confidence=float(probabilities[predicted_class])
    )


def make_binary_prediction(features: Dict, model_data: Dict, market: str) -> MarketPrediction:
    """Make binary prediction (O/U or BTTS)"""
    model = model_data['model']
    feature_names = model_data['feature_names']

    feature_vector = [float(features.get(col, 0.0) or 0.0) for col in feature_names]
    X = np.array([feature_vector])

    probabilities = model.predict_proba(X)[0]

    if market == 'over_under':
        outcome_map = {0: 'under_2_5', 1: 'over_2_5'}
        prob_names = ['under_2_5', 'over_2_5']
        description = 'Over/Under 2.5 Goals'
    else:  # btts
        outcome_map = {0: 'no', 1: 'yes'}
        prob_names = ['no', 'yes']
        description = 'Both Teams to Score'

    predicted_class = int(np.argmax(probabilities))

    return MarketPrediction(
        market=market,
        description=description,
        probabilities={
            prob_names[0]: float(probabilities[0]),
            prob_names[1]: float(probabilities[1]),
        },
        predicted_outcome=outcome_map[predicted_class],
        confidence=float(probabilities[predicted_class])
    )


def make_multi_market_prediction(
    home_team_id: int,
    away_team_id: int,
    match_date: datetime,
    fixture_id: Optional[int] = None,
    markets: Optional[List[str]] = None
) -> Dict[str, Any]:
    """
    Make predictions for multiple markets

    Args:
        home_team_id: Home team database ID
        away_team_id: Away team database ID
        match_date: Match date
        fixture_id: Optional fixture ID
        markets: List of markets to predict (default: all available)

    Returns:
        Multi-market prediction dictionary
    """
    if markets is None:
        markets = ['1x2', 'over_under', 'btts']

    # Extract features once
    features = extract_fixture_features(home_team_id, away_team_id, match_date, fixture_id)

    predictions = {}
    features_used = 0

    for market in markets:
        try:
            model_data = get_model(market)
            features_used = max(features_used, len(model_data.get('feature_names', [])))

            if market == '1x2':
                predictions[market] = make_1x2_prediction(features, model_data)
            else:
                predictions[market] = make_binary_prediction(features, model_data, market)
        except HTTPException:
            # Model not available, skip
            continue

    return {
        'fixture_id': fixture_id,
        'home_team_id': home_team_id,
        'away_team_id': away_team_id,
        'match_date': match_date.isoformat(),
        'predictions': predictions,
        'features_used': features_used,
        'predicted_at': datetime.now().isoformat(),
    }


def make_prediction(
    home_team_id: int,
    away_team_id: int,
    match_date: datetime,
    fixture_id: Optional[int] = None
) -> Dict[str, Any]:
    """
    Make a prediction for a single fixture (legacy 1X2 only)
    """
    model_data = get_model('1x2')
    model = model_data['model']
    feature_names = model_data['feature_names']
    model_version = model_data.get('model_version', model_data.get('version', 'v1.0'))

    features = extract_fixture_features(home_team_id, away_team_id, match_date, fixture_id)

    feature_vector = [float(features.get(col, 0.0) or 0.0) for col in feature_names]
    X = np.array([feature_vector])

    probabilities = model.predict_proba(X)[0]

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


# ============= ENDPOINTS =============

@router.post("/predict", response_model=PredictionResponse)
async def predict(request: PredictionRequest):
    """
    Generate 1X2 prediction for a single fixture (legacy endpoint)

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


@router.post("/predict/multi", response_model=MultiMarketPredictionResponse)
async def predict_multi_market(
    request: PredictionRequest,
    markets: Optional[str] = Query(None, description="Comma-separated list of markets (1x2,over_under,btts)")
):
    """
    Generate predictions for multiple markets

    - **home_team_id**: Database ID of home team
    - **away_team_id**: Database ID of away team
    - **match_date**: Match date in ISO format (YYYY-MM-DD)
    - **markets**: Optional comma-separated list of markets (default: all)

    Available markets:
    - **1x2**: Match result (Home Win / Draw / Away Win)
    - **over_under**: Over/Under 2.5 Goals
    - **btts**: Both Teams to Score
    """
    try:
        match_date = datetime.fromisoformat(request.match_date.replace('Z', '+00:00'))
    except ValueError:
        raise HTTPException(status_code=400, detail="Invalid date format. Use ISO format (YYYY-MM-DD)")

    market_list = None
    if markets:
        market_list = [m.strip() for m in markets.split(',')]
        invalid = [m for m in market_list if m not in MARKETS]
        if invalid:
            raise HTTPException(
                status_code=400,
                detail=f"Invalid markets: {invalid}. Available: {list(MARKETS.keys())}"
            )

    try:
        result = make_multi_market_prediction(
            home_team_id=request.home_team_id,
            away_team_id=request.away_team_id,
            match_date=match_date,
            fixture_id=request.fixture_id,
            markets=market_list
        )
        return MultiMarketPredictionResponse(**result)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/predict/batch", response_model=BatchPredictionResponse)
async def predict_batch(request: BatchPredictionRequest):
    """
    Generate 1X2 predictions for multiple fixtures (legacy endpoint)
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

    model_data = get_model('1x2')

    return BatchPredictionResponse(
        predictions=predictions,
        model_version=model_data.get('model_version', model_data.get('version', 'v1.0')),
        count=len(predictions),
        predicted_at=datetime.now().isoformat()
    )


@router.get("/model/metrics", response_model=ModelMetricsResponse)
async def get_model_metrics(market: str = Query('1x2', description="Market type (1x2, over_under, btts)")):
    """
    Get model performance metrics for a specific market
    """
    try:
        model_data = get_model(market)
        metrics = model_data.get('metrics', {})

        return ModelMetricsResponse(
            model_version=model_data.get('model_version', model_data.get('version', 'v1.0')),
            training_date=model_data.get('training_date', datetime.now().isoformat()),
            accuracy=metrics.get('accuracy', 0.0),
            baseline_accuracy=metrics.get('baseline_accuracy', metrics.get('baseline', 0.39)),
            improvement=metrics.get('improvement', metrics.get('accuracy', 0.0) - metrics.get('baseline_accuracy', metrics.get('baseline', 0.39))),
            config_name=metrics.get('config_name', model_data.get('config')),
            feature_count=len(model_data.get('feature_names', []))
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/model/metrics/all", response_model=AllModelsMetricsResponse)
async def get_all_model_metrics():
    """
    Get performance metrics for all available market models
    """
    metrics = {}
    available = []

    for market in MARKETS:
        try:
            model_data = get_model(market)
            model_metrics = model_data.get('metrics', {})

            metrics[market] = ModelMetricsResponse(
                model_version=model_data.get('model_version', model_data.get('version', 'v1.0')),
                training_date=model_data.get('training_date', datetime.now().isoformat()),
                accuracy=model_metrics.get('accuracy', 0.0),
                baseline_accuracy=model_metrics.get('baseline_accuracy', model_metrics.get('baseline', 0.5)),
                improvement=model_metrics.get('improvement', 0.0),
                config_name=model_metrics.get('config_name', model_data.get('config')),
                feature_count=len(model_data.get('feature_names', []))
            )
            available.append(market)
        except HTTPException:
            continue

    return AllModelsMetricsResponse(
        markets=metrics,
        available_markets=available
    )


@router.get("/markets")
async def list_markets():
    """
    List all available betting markets and their status
    """
    status = {}
    for market, info in MARKETS.items():
        model_file = info['model_file']
        model_path = os.path.join(config.MODEL_PATH, model_file)
        status[market] = {
            'description': info['description'],
            'outcomes': info['outcomes'],
            'model_available': os.path.exists(model_path),
        }
    return {'markets': status}


@router.post("/model/reload")
async def reload_model(market: Optional[str] = Query(None, description="Market to reload (or all if not specified)")):
    """
    Reload model(s) from disk

    Use this after retraining to load the new model.
    """
    clear_model_cache(market)

    if market:
        model_data = get_model(market)
        return {
            "status": "reloaded",
            "market": market,
            "model_version": model_data.get('model_version', model_data.get('version', 'v1.0')),
            "feature_count": len(model_data.get('feature_names', []))
        }
    else:
        reloaded = []
        for m in MARKETS:
            try:
                get_model(m)
                reloaded.append(m)
            except HTTPException:
                pass
        return {
            "status": "reloaded",
            "markets": reloaded
        }
