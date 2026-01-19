package main

import (
	"flag"
	"fmt"
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

	fmt.Println("Database file:", dbFile)
	if utils.FileExists(dbFile) {
		fmt.Println("Database file exists.")
	} else {
		fmt.Println("Database file does not exist.")
		fmt.Println("Creating necessary directories...")
		err := os.MkdirAll(filepath.Dir(dbFile), os.ModePerm)
		if err != nil {
			log.Fatal("Error creating directories:", err)
		}
		fmt.Println("Directories created.")
		fmt.Println("Creating database...")
		_, err = utils.CreateDatabase(dbFile)
		if err != nil {
			log.Fatal("Error creating database:", err)
		}
		fmt.Println("Database created successfully.")
	}

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
}

// printVersion prints the version information for the application.
func printVersion() {
	println("Light Backup version:", VERSION)
	println("Author:", AUTHOR)
	println("License:", LICENSE)
}
