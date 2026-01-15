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

// Test script to find current season fixtures using different approaches

func main() {
	// Load .env file
	godotenv.Load(".env")

	apiKey := os.Getenv("API_FOOTBALL_KEY")
	if apiKey == "" {
		log.Fatal("API_FOOTBALL_KEY not found in .env file")
	}

	baseURL := "https://v3.football.api-sports.io"

	fmt.Println("=== Testing Current Season Fixtures ===\n")

	// Approach 1: Try season 2025
	fmt.Println("1. Testing season 2025...")
	testEndpoint(baseURL+"/fixtures?league=39&season=2025", apiKey)

	// Approach 2: Get upcoming fixtures (next 10)
	fmt.Println("\n2. Testing upcoming fixtures (next 10)...")
	testEndpoint(baseURL+"/fixtures?league=39&next=10", apiKey)

	// Approach 3: Get today's fixtures
	today := time.Now().Format("2006-01-02")
	fmt.Printf("\n3. Testing fixtures for today (%s)...\n", today)
	testEndpoint(baseURL+"/fixtures?league=39&date="+today, apiKey)

	// Approach 4: Get fixtures from today onwards (next 7 days)
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	nextWeek := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	fmt.Printf("\n4. Testing fixtures from %s to %s...\n", tomorrow, nextWeek)
	testEndpoint(baseURL+"/fixtures?league=39&from="+tomorrow+"&to="+nextWeek, apiKey)

	// Approach 5: Get last 5 fixtures (to see recent results)
	fmt.Println("\n5. Testing last 5 fixtures...")
	testEndpoint(baseURL+"/fixtures?league=39&last=5", apiKey)

	// Approach 6: Get fixtures for specific round (if season is active)
	fmt.Println("\n6. Testing current round fixtures...")
	testEndpoint(baseURL+"/fixtures?league=39&season=2025&round=Regular Season - 1", apiKey)
}

func testEndpoint(url, apiKey string) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("❌ Failed to create request: %v\n", err)
		return
	}

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

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("❌ Failed to parse JSON: %v\n", err)
		return
	}

	if resp.StatusCode == 200 {
		results := 0
		if r, ok := result["results"].(float64); ok {
			results = int(r)
		}

		fmt.Printf("✅ Success! Results: %d\n", results)

		if params, ok := result["parameters"].(map[string]interface{}); ok {
			fmt.Printf("   Parameters: %v\n", params)
		}

		// Print first fixture if available
		if response, ok := result["response"].([]interface{}); ok && len(response) > 0 {
			if fixture, ok := response[0].(map[string]interface{}); ok {
				if fixtureData, ok := fixture["fixture"].(map[string]interface{}); ok {
					fmt.Printf("   First fixture ID: %.0f\n", fixtureData["id"].(float64))
					fmt.Printf("   Date: %s\n", fixtureData["date"].(string))
				}
				if teams, ok := fixture["teams"].(map[string]interface{}); ok {
					if home, ok := teams["home"].(map[string]interface{}); ok {
						fmt.Printf("   Match: %s vs ", home["name"].(string))
					}
					if away, ok := teams["away"].(map[string]interface{}); ok {
						fmt.Printf("%s\n", away["name"].(string))
					}
				}
			}
		}
	} else {
		fmt.Printf("❌ Failed! Status: %d\n", resp.StatusCode)
		if errors := result["errors"]; errors != nil {
			fmt.Printf("   Errors: %v\n", errors)
		}
	}
}
