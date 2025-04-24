package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"triple-s/api"
	"triple-s/config"
	"triple-s/storage"
	"triple-s/storage/buckets"
	"triple-s/utils"
)

func main() {
	// Define command-line flags
	port := flag.String("port", "8080", "Port number (e.g., 8080)") // Removed leading colon
	storageDir := flag.String("dir", "./data", "Path to the storage directory")
	help := flag.Bool("help", false, "Show help screen")

	// Parse the flags
	flag.Parse()

	// Show help and exit if --help is provided
	if *help {
		utils.PrintUsage()
		os.Exit(0)
	}

	// Create a new config instance
	cfg := config.NewConfig(*port, *storageDir)

	// Initialize storage and root directory
	storage.InitStorage()
	if err := buckets.CreateRootDirectory(); err != nil {
		log.Fatalf("Error creating root directory: %v\n", err)
		showHelpAndExit()
	}

	// Attempt to start the server
	log.Printf("Starting server on port %s, storing files in %s\n", *port, *storageDir)
	if err := api.StartServer(cfg); err != nil {
		log.Printf("Error starting server: %v\n", err)
		showHelpAndExit()
	}
}

// showHelpAndExit prints the help screen and exits the program with a non-zero status.
func showHelpAndExit() {
	fmt.Println("An error occurred. Please review the options below:")
	utils.PrintUsage()
	os.Exit(1)
}
