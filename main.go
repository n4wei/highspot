package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/n4wei/highspot/collection"
	"github.com/n4wei/highspot/models"
)

const (
	logFormat             = log.Ldate | log.Ltime | log.Lshortfile | log.LUTC
	defaultFilePermission = 0666
)

func main() {
	// Parse command line flags
	var mixtapeFile, changesFile, outputFile string
	flag.StringVar(&mixtapeFile, "m", "", "filepath to the JSON mixtape file")
	flag.StringVar(&changesFile, "c", "", "filepath to the JSON changes file")
	flag.StringVar(&outputFile, "o", "./output.json", "filepath to write changed mixtape JSON file")
	flag.Parse()

	if mixtapeFile == "" || changesFile == "" {
		handleFlagError(errors.New("missing required flags -m and -c"))
	}

	// Read mixtape file
	mixtape := &models.Mixtape{}
	err := readFromFile(mixtapeFile, mixtape)
	handleError(err)

	// Read changes file
	changes := &models.Changes{}
	err = readFromFile(changesFile, changes)
	handleError(err)

	// Create the object used to apply changes to mixtape
	logger := log.New(os.Stdout, "", logFormat)
	collection := collection.New(mixtape, logger)

	// Apply changes to mixtape
	err = collection.ApplyChanges(changes)
	handleError(err)

	// Write mixtape to file
	err = writeToFile(mixtape, outputFile)
	handleError(err)
}

func readFromFile(filepath string, object interface{}) error {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, object)
	if err != nil {
		return fmt.Errorf("error unmarshaling %s to JSON: %v", filepath, err)
	}

	return nil
}

func writeToFile(mixtape *models.Mixtape, filepath string) error {
	bytes, err := json.Marshal(mixtape)
	if err != nil {
		return fmt.Errorf("error marshaling Mixtape object to JSON: %v", err)
	}

	return ioutil.WriteFile(filepath, bytes, defaultFilePermission)
}

func handleFlagError(err error) {
	fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
	flag.PrintDefaults()
	os.Exit(1)
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}
