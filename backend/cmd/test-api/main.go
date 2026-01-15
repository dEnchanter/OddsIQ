package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Test script to check what API-Football endpoints are available

func main() {
	// Load .env file (when running from backend/ directory)
	godotenv.Load(".env")

	apiKey := os.Getenv("API_FOOTBALL_KEY")
	if apiKey == "" {
		log.Fatal("API_FOOTBALL_KEY not found in .env file")
	}

	baseURL := "https://v3.football.api-sports.io"

	fmt.Println("=== Testing API-Football Endpoints ===\n")

	// Test 1: Check status/quota
	fmt.Println("1. Checking API Status & Quota...")
	testEndpoint(baseURL+"/status", apiKey)

	// Test 2: Get available leagues
	fmt.Println("\n2. Getting Premier League info...")
	testEndpoint(baseURL+"/leagues?id=39", apiKey)

	// Test 3: Get current season fixtures
	fmt.Println("\n3. Getting upcoming Premier League fixtures...")
	testEndpoint(baseURL+"/fixtures?league=39&season=2024&next=5", apiKey)

	// Test 4: Check if odds endpoint is available
	fmt.Println("\n4. Checking Odds endpoint...")
	testEndpoint(baseURL+"/odds?league=39&season=2024", apiKey)

	// Test 5: Check bookmakers
	fmt.Println("\n5. Checking Bookmakers endpoint...")
	testEndpoint(baseURL+"/odds/bookmakers", apiKey)

	// Test 6: Check available bets
	fmt.Println("\n6. Checking Bets endpoint...")
	testEndpoint(baseURL+"/odds/bets", apiKey)
}

func testEndpoint(url, apiKey string) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("❌ Failed to create request: %v\n", err)
		return
	}

	// API-Football requires this header
	req.Header.Add("x-apisports-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Failed to read response: %v\n", err)
		return
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("❌ Failed to parse JSON: %v\n", err)
		return
	}

	// Check response
	if resp.StatusCode == 200 {
		results := int(result["results"].(float64))
		fmt.Printf("✅ Success! Status: %d, Results: %d\n", resp.StatusCode, results)

		// Print quota info if available
		if params, ok := result["parameters"].(map[string]interface{}); ok {
			fmt.Printf("   Parameters: %v\n", params)
		}

		// Print first result if available
		if response, ok := result["response"].([]interface{}); ok && len(response) > 0 {
			prettyJSON, _ := json.MarshalIndent(response[0], "   ", "  ")
			fmt.Printf("   First result:\n   %s\n", string(prettyJSON))
		}
	} else {
		fmt.Printf("❌ Failed! Status: %d\n", resp.StatusCode)
		errors := result["errors"]
		fmt.Printf("   Errors: %v\n", errors)
	}

	// Print remaining requests
	if paging, ok := result["paging"].(map[string]interface{}); ok {
		fmt.Printf("   Paging: %v\n", paging)
	}
}
