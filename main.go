package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var googleCloudProjectID string
var baseURL string
var apiBaseURL string

func main() {
	initialiseEnv()

	companyToVehicle, err := vehicleWebCrawler(baseURL)
	if err != nil {
		log.Fatalf("Error crawling company vehicles: %v", err)
		return
	}
	processCompanies(companyToVehicle)
}

func initialiseEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("No .env file found")
	}
	baseURL = os.Getenv("BASE_URL")
	apiBaseURL = os.Getenv("API_BASE_URL")

	googleCloudProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT_ID")
}
