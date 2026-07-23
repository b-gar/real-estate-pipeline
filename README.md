# Real Estate Data Pipeline

A highly optimized, serverless ETL data pipeline written in Go. This pipeline dynamically fetches real estate market data (Active Listings and Recent Sales) from [Rentcast](https://www.rentcast.io/) and streams it directly into Google Cloud BigQuery for analytics.

## Architecture

This pipeline is designed specifically for **Google Cloud Run Jobs** and **Google Cloud Scheduler**, allowing it to execute on a recurring schedule. 

- **Concurrency:** Uses Go's `sync.WaitGroup` to fetch and upload multiple datasets simultaneously, drastically reducing execution time.
- **Security:** The resulting Docker container is built on an ultra-secure `distroless` image containing zero shells or package managers. 
- **Secret Management:** API keys are injected at runtime via Google Secret Manager.

## Environment Variables

The pipeline requires the following environment variables to run. In a production environment, these should be passed dynamically by Cloud Scheduler or the Cloud Run UI:

- `GOOGLE_CLOUD_PROJECT`: Your GCP Project ID
- `ZIP_CODE`: The target Zip Code to query
- `DAYS_AGO`: How far back to query recent sales
- `RENTCAST_API_KEY`: The API key to authenticate with Rentcast

## Deployment

This repository includes a 3-stage `cloudbuild.yaml` file designed to run in Google Cloud Build. 

Pushing to the `main` branch will automatically:
1. Compile the Go binary into a hardened Distroless Docker image
2. Push the image to Google Artifact Registry
3. Deploy (or update) the Cloud Run Job and bind the Secret Manager API keys
