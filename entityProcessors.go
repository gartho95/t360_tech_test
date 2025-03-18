package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

func processCompanies(companyToVehicle map[string][]string) {
	var wg sync.WaitGroup

	for company, vehicles := range companyToVehicle {
		wg.Add(1)
		go func(company string, vehicles []string) {
			defer wg.Done()
			processVehicles(company, vehicles)
		}(company, vehicles)
	}

	wg.Wait()
}

func processVehicles(company string, vehicles []string) {
	var wg sync.WaitGroup

	for _, vehicle := range vehicles {
		wg.Add(1)
		go func(vehicle string) {
			defer wg.Done()
			searchResponse, timeout := sendSearchRequest(company, vehicle)
			if timeout {
				return
			}

			buildAndPublishVehicle(company, searchResponse)
		}(vehicle)
	}

	wg.Wait()
}
func buildAndPublishVehicle(company string, searchResponse SearchResponse) {

	newPacket := searchResponse

	newPacket.Reference = uuid.NewString()
	newPacket.Vrm = searchResponse.Vrm
	newPacket.ContraventionDate = searchResponse.ContraventionDate
	newPacket.IsHirerVehicle = searchResponse.IsHirerVehicle
	newPacket.LeaseCompany.Companyname = company
	newPacket.LeaseCompany.AddressLine1 = searchResponse.LeaseCompany.AddressLine1
	newPacket.LeaseCompany.AddressLine2 = searchResponse.LeaseCompany.AddressLine2
	newPacket.LeaseCompany.AddressLine3 = searchResponse.LeaseCompany.AddressLine3
	newPacket.LeaseCompany.AddressLine4 = searchResponse.LeaseCompany.AddressLine4
	newPacket.LeaseCompany.Postcode = searchResponse.LeaseCompany.Postcode

	publishMessage("positive_searches", newPacket)

}

func sendSearchRequest(company string, vehicle string) (SearchResponse, bool) {

	requestUrl := fmt.Sprintf("%s/%s", apiBaseURL, company)

	currentTimeFormatted := time.Now().UTC().Format(time.RFC3339)
	requestBody := SearchRequest{
		Vrm:               vehicle,
		ContraventionDate: currentTimeFormatted,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", requestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("Request timeout exceeded for %s, %s", company, vehicle)
		return SearchResponse{}, true
	}
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading body: %v", err)
	}
	var searchResponse SearchResponse

	err = json.Unmarshal(responseBody, &searchResponse)
	if err != nil {
		log.Fatalf("Error reading body: %v", err)
	}
	return searchResponse, false
}

type SearchRequest struct {
	Vrm               string `json:"vrm"`
	ContraventionDate string `json:"contravention_date"`
}
type SearchResponse struct {
	Reference         string    `json:"reference"`
	Vrm               string    `json:"vrm"`
	ContraventionDate time.Time `json:"contravention_date"`
	IsHirerVehicle    bool      `json:"is_hirer_vehicle"`
	LeaseCompany      struct {
		Companyname  string  `json:"companyname"`
		AddressLine1 *string `json:"address_line1"`
		AddressLine2 *string `json:"address_line2"`
		AddressLine3 *string `json:"address_line3"`
		AddressLine4 *string `json:"address_line4"`
		Postcode     *string `json:"postcode"`
	} `json:"lease_company"`
}
