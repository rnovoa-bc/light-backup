package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/rnovoa-bc/light-backup/utils"
)

var VERSION = "1.0.0"
var AUTHOR = "Ramon Novoa <rnovoa@bnc.cat>"
var LICENSE = "GPLv3"

// GLOBAL VARIABLES
var dbFile string

func main() {
	// Initialize some application settings if needed

	// Get the directory containing the executable
	baseFolder, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dbFile = filepath.Join(baseFolder, "db", "light-backup.db")

	// Process command-line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Define subcommands
	backupCommand := flag.NewFlagSet("backup", flag.ExitOnError)
	restoreCommand := flag.NewFlagSet("restore", flag.ExitOnError)
	helpCommand := flag.NewFlagSet("help", flag.ExitOnError)
	versionCommand := flag.NewFlagSet("version", flag.ExitOnError)
	setupCommand := flag.NewFlagSet("setup", flag.ExitOnError)
	newDestination := flag.NewFlagSet("new-destination", flag.ExitOnError)
	newJob := flag.NewFlagSet("new-job", flag.ExitOnError)
	testSource := flag.NewFlagSet("test-source", flag.ExitOnError)

	switch os.Args[1] {
	case "backup":
		backupCommand.Parse(os.Args[2:])
		// Implement backup logic here
		println("Backup command executed")
	case "restore":
		restoreCommand.Parse(os.Args[2:])
		// Implement restore logic here
		println("Restore command executed")
	case "help":
		helpCommand.Parse(os.Args[2:])
		printUsage()
	case "version":
		versionCommand.Parse(os.Args[2:])
		printVersion()
	case "setup":
		// Create and initialize the database
		shouldInit := setupCommand.Bool("initialize", false, "Initialize the database schema")
		setupCommand.Parse(os.Args[2:])
		log.Println("Should initialize:", *shouldInit)
		err := utils.CreateDatabase(dbFile, shouldInit)
		if err != nil {
			log.Fatal("Error initializing database:", err)
		}
	case "new-destination":
		region := newDestination.String("region", "", "Region of the destination")
		bucket := newDestination.String("bucket", "", "Bucket name of the destination")
		accessKey := newDestination.String("access_key", "", "Access key for the destination")
		secretKey := newDestination.String("secret_key", "", "Secret key for the destination")
		newDestination.Parse(os.Args[2:])

		if *region == "" || *bucket == "" || *accessKey == "" || *secretKey == "" {
			log.Fatal("All parameters (region, bucket, access_key, secret_key) are required")
		}

		db, err := utils.OpenDatabase(dbFile)
		if err != nil {
			log.Fatal("Error opening database:", err)
		}
		defer db.Close()

		id, err := utils.AddDestination(db, *region, *bucket, *accessKey, *secretKey)
		if err != nil {
			log.Fatal("Error adding new destination:", err)
		}
		log.Printf("New destination [%v] added successfully", id)
	case "new-job":
		name := newJob.String("name", "", "Name of the job")
		sourcePath := newJob.String("source_path", "", "Path of the source")
		sourceType := newJob.String("source_type", "local", "Type of the source (local, nfs, etc.)")
		sourceOptions := newJob.String("source_options", "", "Options for the source (JSON encoded)")
		destinationID := newJob.Int("destination", 0, "ID of the destination")
		hashAlgorithm := newJob.String("hash_algorithm", "xxhash", "Hash algorithm to use (blake3, xxhash)")
		maxDuration := newJob.Int("max_duration", 604800, "Maximum duration of the job in seconds")
		newJob.Parse(os.Args[2:])
		if *name == "" || *sourcePath == "" || *destinationID == 0 {
			log.Fatal("Parameters name, source_path, and destination_id are required")
		}
		db, err := utils.OpenDatabase(dbFile)
		if err != nil {
			log.Fatal("Error opening database:", err)
		}
		defer db.Close()
		id, err := utils.AddJob(db, *name, *sourcePath, *sourceType, *sourceOptions, *destinationID, *hashAlgorithm, *maxDuration)
		if err != nil {
			log.Fatal("Error adding new job:", err)
		}
		log.Printf("New job [%v] added successfully", id)
	case "test-source":
		sourcePath := testSource.String("path", "", "Path of the source to test")
		sourceType := testSource.String("type", "local", "Type of the source (local, nfs, etc.)")
		sourceOptions := testSource.String("options", "", "Options for the source (JSON encoded)")
		testSource.Parse(os.Args[2:])
		if *sourcePath == "" {
			log.Fatal("Source path is required")
		}
		log.Printf("Testing source with path: %s, type: %s, options: %s", *sourcePath, *sourceType, *sourceOptions)
		// Implement source testing logic here
		err := utils.WalkDirectory(*sourcePath)
		if err != nil {
			log.Fatal("Error testing source:", err)
		}
	default:
		printUsage()
		os.Exit(1)
	}

}

// printUsage prints the usage information for the application.
func printUsage() {
	println("Usage: light-backup [command]")
	println("Commands:")
	println("  backup    Perform a backup operation")
	println("  restore   Restore from a backup")
	println("  help      Show this help message")
	println("  setup     Setup the database")
	println("    --initialize   Initialize the database schema")
}

// printVersion prints the version information for the application.
func printVersion() {
	println("Light Backup version:", VERSION)
	println("Author:", AUTHOR)
	println("License:", LICENSE)
}
