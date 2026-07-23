package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Contact struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

type HOA struct {
	Fee float64 `json:"fee"`
}

type ListingHistory struct {
	Event        string     `json:"event"`
	Price        float64    `json:"price"`
	ListingType  string     `json:"listingType"`
	ListedDate   time.Time  `json:"listedDate"`
	RemovedDate  *time.Time `json:"removedDate"`
	DaysOnMarket int        `json:"daysOnMarket"`
}

type PropertyTaxes struct {
	Year  int     `json:"year"`
	Total float64 `json:"total"`
}

type SaleHistory struct {
	Event string    `json:"event"`
	Date  time.Time `json:"date"`
	Price float64   `json:"price"`
}

type TaxAssessments struct {
	Year         int     `json:"year"`
	Value        float64 `json:"value"`
	Land         float64 `json:"land"`
	Improvements float64 `json:"improvements"`
}

type ActiveListing struct {
	ID            string                    `json:"id"`
	AddressLine1  string                    `json:"addressLine1"`
	AddressLine2  *string                   `json:"addressLine2"`
	City          string                    `json:"city"`
	State         string                    `json:"state"`
	StateFips     string                    `json:"stateFips"`
	ZipCode       string                    `json:"zipCode"`
	County        string                    `json:"county"`
	CountyFips    string                    `json:"countyFips"`
	Latitude      float64                   `json:"latitude"`
	Longitude     float64                   `json:"longitude"`
	PropertyType  string                    `json:"propertyType"`
	Bedrooms      int                       `json:"bedrooms"`
	Bathrooms     float64                   `json:"bathrooms"`
	HouseSize     int                       `json:"squareFootage"`
	LotSize       int                       `json:"lotSize"`
	YearBuilt     int                       `json:"yearBuilt"`
	HOA           HOA                       `json:"hoa"`
	Status        string                    `json:"status"`
	Price         float64                   `json:"price"`
	ListedDate    time.Time                 `json:"listedDate"`
	RemovedDate   *time.Time                `json:"removedDate"`
	CreatedDate   time.Time                 `json:"createdDate"`
	LastSeenDate  time.Time                 `json:"lastSeenDate"`
	DaysOnMarket  int                       `json:"daysOnMarket"`
	ListingAgent  Contact                   `json:"listingAgent"`
	ListingOffice Contact                   `json:"listingOffice"`
	MLSName       string                    `json:"mlsName"`
	MLSNumber     string                    `json:"mlsNumber"`
	History       map[string]ListingHistory `json:"history"`
}

type PropertySale struct {
	ID               string                    `json:"id"`
	AssessorID       string                    `json:"assessorId"`
	LegalDescription string                    `json:"legalDescription"`
	ParcelNumber     string                    `json:"parcelNumber"`
	LastSalePrice    *float64                  `json:"lastSalePrice"`
	LastSaleDate     *time.Time                `json:"lastSaleDate"`
	TaxAssessments   map[string]TaxAssessments `json:"taxAssessments"`
	PropertyTaxes    map[string]PropertyTaxes  `json:"propertyTaxes"`
	History          map[string]SaleHistory    `json:"history"`
}

func fetchRentcastPages[T any](basePath string) ([]T, error) {
	apiKey := os.Getenv("RENTCAST_API_KEY")
	var allResults []T
	limit := 500
	offset := 0
	client := &http.Client{Timeout: 10 * time.Second}
	for {
		url := fmt.Sprintf("%s&limit=%d&offset=%d", basePath, limit, offset)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-Api-Key", apiKey)
		req.Header.Set("accept", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("API returned status: %s", resp.Status)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		var pageResults []T
		if err := json.Unmarshal(body, &pageResults); err != nil {
			return nil, err
		}
		allResults = append(allResults, pageResults...)
		if len(pageResults) < limit {
			break
		}
		offset += limit
	}
	return allResults, nil
}

func fetchActiveListings(zipCode string) ([]ActiveListing, error) {
	basePath := fmt.Sprintf("https://api.rentcast.io/v1/listings/sale?zipCode=%s&propertyType=Single-Family&status=Active", zipCode)
	return fetchRentcastPages[ActiveListing](basePath)
}

func fetchRecentSales(zipCode string, daysAgo int) ([]PropertySale, error) {
	basePath := fmt.Sprintf("https://api.rentcast.io/v1/properties?zipCode=%s&propertyType=Single-Family&saleDateRange=*:%d", zipCode, daysAgo)
	return fetchRentcastPages[PropertySale](basePath)
}
