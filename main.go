package main

import (
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	datasetID := "rentcast_raw"
	zipCode := "85641"
	daysAgo := 180
	ctx := context.Background()
	fmt.Println("Starting Real Estate Pipeline for Zip Code:", zipCode)
	// Active Listings
	fmt.Println("\nFetching Active Listings...")
	listings, err := fetchActiveListings(zipCode)
	if err != nil {
		log.Fatalf("Failed to fetch active listings: %v", err)
	}
	err = StreamToBigQuery(ctx, projectID, datasetID, "active_listings", listings, ActiveListing{})
	if err != nil {
		log.Fatalf("Failed to upload active listings: %v", err)
	}
	fmt.Println("Successfully uploaded Active Listings!")
	// Recent Sales
	fmt.Printf("\nFetching Recent Sales (Last %d days)...\n", daysAgo)
	sales, err := fetchRecentSales(zipCode, daysAgo)
	if err != nil {
		log.Fatalf("Failed to fetch recent sales: %v", err)
	}
	err = StreamToBigQuery(ctx, projectID, datasetID, "recent_sales", sales, PropertySale{})
	if err != nil {
		log.Fatalf("Failed to upload recent sales: %v", err)
	}
	fmt.Println("Successfully uploaded Recent Sales!")
	fmt.Println("\nReal Estate Pipeline complete!")
}
