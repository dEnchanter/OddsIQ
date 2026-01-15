package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dEnchanter/OddsIQ/backend/internal/models"
)

// MLClient handles communication with the Python ML service
type MLClient struct {
	baseURL    string
	httpClient *http.Client
}

// PredictionRequest represents a request to the ML service
type PredictionRequest struct {
	HomeTeamID int    `json:"home_team_id"`
	AwayTeamID int    `json:"away_team_id"`
	MatchDate  string `json:"match_date"`
	FixtureID  *int   `json:"fixture_id,omitempty"`
}

// BatchPredictionRequest represents a batch prediction request
type BatchPredictionRequest struct {
	Fixtures []PredictionRequest `json:"fixtures"`
}

// PredictionResponse represents the ML service response
type PredictionResponse struct {
	FixtureID        *int               `json:"fixture_id"`
	HomeTeamID       int                `json:"home_team_id"`
	AwayTeamID       int                `json:"away_team_id"`
	ModelVersion     string             `json:"model_version"`
	Predictions      PredictionProbs    `json:"predictions"`
	PredictedOutcome string             `json:"predicted_outcome"`
	Confidence       float64            `json:"confidence"`
	FeaturesUsed     int                `json:"features_used"`
	PredictedAt      string             `json:"predicted_at"`
}

// PredictionProbs holds probability predictions
type PredictionProbs struct {
	HomeWinProb float64 `json:"home_win_prob"`
	DrawProb    float64 `json:"draw_prob"`
	AwayWinProb float64 `json:"away_win_prob"`
}

// BatchPredictionResponse represents the batch response
type BatchPredictionResponse struct {
	Predictions  []PredictionResponse `json:"predictions"`
	ModelVersion string               `json:"model_version"`
	Count        int                  `json:"count"`
	PredictedAt  string               `json:"predicted_at"`
}

// ModelMetricsResponse represents model performance metrics
type ModelMetricsResponse struct {
	ModelVersion     string   `json:"model_version"`
	TrainingDate     string   `json:"training_date"`
	Accuracy         float64  `json:"accuracy"`
	BaselineAccuracy float64  `json:"baseline_accuracy"`
	Improvement      float64  `json:"improvement"`
	ConfigName       *string  `json:"config_name"`
	FeatureCount     int      `json:"feature_count"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status       string `json:"status"`
	Service      string `json:"service"`
	Version      string `json:"version"`
	ModelVersion string `json:"model_version"`
}

// NewMLClient creates a new ML service client
func NewMLClient(baseURL string) *MLClient {
	return &MLClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// HealthCheck checks if the ML service is healthy
func (c *MLClient) HealthCheck(ctx context.Context) (*HealthResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call ML service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service unhealthy: status %d", resp.StatusCode)
	}

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &health, nil
}

// Predict gets a prediction for a single fixture
func (c *MLClient) Predict(ctx context.Context, fixture *models.Fixture) (*models.Prediction, error) {
	reqBody := PredictionRequest{
		HomeTeamID: fixture.HomeTeamID,
		AwayTeamID: fixture.AwayTeamID,
		MatchDate:  fixture.MatchDate.Format("2006-01-02"),
		FixtureID:  &fixture.ID,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/predict", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call ML service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service error: status %d", resp.StatusCode)
	}

	var predResp PredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&predResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to internal Prediction model
	prediction := &models.Prediction{
		FixtureID:        fixture.ID,
		ModelVersion:     predResp.ModelVersion,
		HomeWinProb:      predResp.Predictions.HomeWinProb,
		DrawProb:         predResp.Predictions.DrawProb,
		AwayWinProb:      predResp.Predictions.AwayWinProb,
		PredictedOutcome: predResp.PredictedOutcome,
		ConfidenceScore:  predResp.Confidence,
		Features: map[string]interface{}{
			"features_used": predResp.FeaturesUsed,
		},
		PredictedAt: time.Now(),
	}

	return prediction, nil
}

// PredictBatch gets predictions for multiple fixtures
func (c *MLClient) PredictBatch(ctx context.Context, fixtures []*models.Fixture) ([]*models.Prediction, error) {
	requests := make([]PredictionRequest, len(fixtures))
	for i, f := range fixtures {
		requests[i] = PredictionRequest{
			HomeTeamID: f.HomeTeamID,
			AwayTeamID: f.AwayTeamID,
			MatchDate:  f.MatchDate.Format("2006-01-02"),
			FixtureID:  &f.ID,
		}
	}

	reqBody := BatchPredictionRequest{Fixtures: requests}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/predict/batch", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call ML service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service error: status %d", resp.StatusCode)
	}

	var batchResp BatchPredictionResponse
	if err := json.NewDecoder(resp.Body).Decode(&batchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to internal Prediction models
	predictions := make([]*models.Prediction, len(batchResp.Predictions))
	for i, predResp := range batchResp.Predictions {
		fixtureID := 0
		if predResp.FixtureID != nil {
			fixtureID = *predResp.FixtureID
		}

		predictions[i] = &models.Prediction{
			FixtureID:        fixtureID,
			ModelVersion:     predResp.ModelVersion,
			HomeWinProb:      predResp.Predictions.HomeWinProb,
			DrawProb:         predResp.Predictions.DrawProb,
			AwayWinProb:      predResp.Predictions.AwayWinProb,
			PredictedOutcome: predResp.PredictedOutcome,
			ConfidenceScore:  predResp.Confidence,
			Features: map[string]interface{}{
				"features_used": predResp.FeaturesUsed,
			},
			PredictedAt: time.Now(),
		}
	}

	return predictions, nil
}

// GetModelMetrics retrieves model performance metrics
func (c *MLClient) GetModelMetrics(ctx context.Context) (*ModelMetricsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/model/metrics", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call ML service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service error: status %d", resp.StatusCode)
	}

	var metrics ModelMetricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &metrics, nil
}

// ReloadModel triggers model reload on the ML service
func (c *MLClient) ReloadModel(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/model/reload", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call ML service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ML service error: status %d", resp.StatusCode)
	}

	return nil
}
