package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"

	"cloud.google.com/go/bigquery"
)

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	datasetID := "rentcast_raw"
	zipCode := os.Getenv("ZIP_CODE")
	if zipCode == "" {
		log.Fatal("Invalid ZIP_CODE environment variable")
	}
	daysAgo, err := strconv.Atoi(os.Getenv("DAYS_AGO"))
	if err != nil {
		log.Fatalf("Invalid DAYS_AGO environment variable: %v", err)
	}
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create BigQuery client: %v", err)
	}
	defer client.Close()
	log.Printf("Starting real estate pipeline for zip code: %s", zipCode)
	var wg sync.WaitGroup
	wg.Add(2)
	go processActiveListings(ctx, &wg, client, datasetID, zipCode)
	go processRecentSales(ctx, &wg, client, datasetID, zipCode, daysAgo)
	wg.Wait()
	log.Println("Real estate pipeline complete!")
}

func processActiveListings(ctx context.Context, wg *sync.WaitGroup, client *bigquery.Client, datasetID string, zipCode string) {
	defer wg.Done()
	log.Println("Fetching active_listings...")
	listings, err := fetchActiveListings(zipCode)
	if err != nil {
		log.Fatalf("Failed to fetch active listings: %v", err)
	}
	err = StreamToBigQuery(ctx, client, datasetID, "active_listings", listings, ActiveListing{})
	if err != nil {
		log.Fatalf("Failed to upload active_listings: %v", err)
	}
	log.Println("Successfully uploaded active_listings!")
}

func processRecentSales(ctx context.Context, wg *sync.WaitGroup, client *bigquery.Client, datasetID string, zipCode string, daysAgo int) {
	defer wg.Done()
	log.Printf("Fetching recent_sales (Last %d days)...", daysAgo)
	sales, err := fetchRecentSales(zipCode, daysAgo)
	if err != nil {
		log.Fatalf("Failed to fetch recent sales: %v", err)
	}
	err = StreamToBigQuery(ctx, client, datasetID, "recent_sales", sales, PropertySale{})
	if err != nil {
		log.Fatalf("Failed to upload recent_sales: %v", err)
	}
	log.Println("Successfully uploaded recent_sales!")
}
