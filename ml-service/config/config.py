import os
from typing import List
from dotenv import load_dotenv

load_dotenv()


class Config:
    """ML Service configuration"""

    # Database
    DATABASE_URL: str = os.getenv("DATABASE_URL", "postgresql://localhost:5432/oddsiq")

    # Model
    MODEL_PATH: str = os.getenv("MODEL_PATH", "./models")
    MODEL_VERSION: str = os.getenv("MODEL_VERSION", "v1.0")
    FEATURE_STORE_PATH: str = os.getenv("FEATURE_STORE_PATH", "./feature_store")

    # Service
    PORT: int = int(os.getenv("PORT", "8001"))
    ENV: str = os.getenv("ENV", "development")
    LOG_LEVEL: str = os.getenv("LOG_LEVEL", "INFO")

    # Training
    TRAINING_SEASONS: List[int] = [
        int(s) for s in os.getenv("TRAINING_SEASONS", "2021,2022,2023").split(",")
    ]
    MIN_TRAINING_SAMPLES: int = int(os.getenv("MIN_TRAINING_SAMPLES", "1000"))
    TEST_SIZE: float = float(os.getenv("TEST_SIZE", "0.2"))
    RANDOM_STATE: int = int(os.getenv("RANDOM_STATE", "42"))

    # XGBoost Hyperparameters
    XGBOOST_N_ESTIMATORS: int = int(os.getenv("XGBOOST_N_ESTIMATORS", "200"))
    XGBOOST_MAX_DEPTH: int = int(os.getenv("XGBOOST_MAX_DEPTH", "6"))
    XGBOOST_LEARNING_RATE: float = float(os.getenv("XGBOOST_LEARNING_RATE", "0.1"))


config = Config()
