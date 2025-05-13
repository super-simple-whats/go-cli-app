package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	APIKey    string
	hooksHost string
	hooksPath string
)

func loadEnv() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Set variables from environment
	APIKey = os.Getenv("API_KEY")
	hooksHost = os.Getenv("HOOKS_HOST")
	hooksPath = os.Getenv("HOOKS_PATH")

	// Validate required environment variables
	if APIKey == "" || hooksHost == "" ||
		hooksPath == "" {
		log.Fatal("Missing required environment variables")
	}
}
