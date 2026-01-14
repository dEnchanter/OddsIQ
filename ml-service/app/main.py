from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager

from app.api import predictions
from config.config import config


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan events"""
    # Startup
    print(f"🚀 ML Service starting on port {config.PORT}")
    print(f"📊 Model version: {config.MODEL_VERSION}")
    print(f"💾 Model path: {config.MODEL_PATH}")

    yield

    # Shutdown
    print("Shutting down ML Service...")


app = FastAPI(
    title="OddsIQ ML Service",
    description="Machine Learning service for sports betting predictions",
    version="0.1.0",
    lifespan=lifespan
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Configure appropriately for production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Health check endpoint
@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "service": "oddsiq-ml-service",
        "version": "0.1.0",
        "model_version": config.MODEL_VERSION
    }

# Include routers
app.include_router(predictions.router, prefix="/api", tags=["predictions"])


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=config.PORT,
        reload=config.ENV == "development"
    )
