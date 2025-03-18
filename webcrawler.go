package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func vehicleWebCrawler(url string) (map[string][]string, error) {

	companyToVehicle := make(map[string][]string)

	// Far from optimal - long term would have these stored in some kind of config/mapping DB
	companyNameToRequestName := make(map[string]string)
	companyNameToRequestName["ACME Company Ltd"] = "acmelease"
	companyNameToRequestName["Lease Company Ltd"] = "leasecompany"
	companyNameToRequestName["Fleet Company Ltd"] = "fleetcompany"
	companyNameToRequestName["Hire Company Ltd"] = "leasecompany"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching URL:", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error parsing HTML:", err)
	}

	doc.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		vrm := row.Find("td:nth-child(1) span").Text()
		company := row.Find("td:nth-child(2)").Text()
		if companyRequestName, ok := companyNameToRequestName[company]; ok {
			companyToVehicle[companyRequestName] = append(companyToVehicle[companyRequestName], vrm)
		} else {
			log.Printf("Company %s not found in our mappings", company)
		}
	})
	return companyToVehicle, nil
}
