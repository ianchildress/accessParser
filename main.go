package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type LoggedRequest struct {
	IpAddress string `json:"ip"`
	FileURL   string `json:"url"`
}

func main() {
	// Setup some useful vars
	totalLines := 0
	badLines := 0 // Used to identify problem lines

	// Fetch command line argument for file path
	inFile := flag.String("in", "", "Path to the file for parsing.")
	outFile := flag.String("out", "out.json", "Path to place the parsed json file.")

	flag.Parse()
	if *inFile == "" {
		log.Fatalln("No file path given. You must specify a file for parsing.")
	}

	// Open the file
	file, err := os.Open(*inFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a slice to store our parsed data
	logBucket := []LoggedRequest{}

	// Iterate through the file line by line
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		totalLines++ // We start with line 1
		logged, err := Parse(scanner.Text())
		if err != nil {
			badLine := fmt.Sprintf("Failed to parse line %v.\n%s\n", totalLines, scanner.Text())
			badLines++
			log.Println(badLine)
		}
		logBucket = append(logBucket, logged)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if err := Write(logBucket, *outFile); err != nil {
		log.Fatal(err)
	}

	goodLines := totalLines - badLines
	log.Printf("Successfully parsed %v lines and skipped %v bad lines to file %s.", goodLines, badLines, *outFile)

}

// Parse receives a string and pulls the ip address
// and file path from the access log.
func Parse(s string) (LoggedRequest, error) {
	if len(s) == 0 {
		return LoggedRequest{}, fmt.Errorf("Failed to parse line.")
	}
	// Get the ip address
	a := strings.Split(s, " - - ")
	if len(a) < 2 {
		return LoggedRequest{}, fmt.Errorf("Failed to parse line.")
	}

	// Get the image url
	b := strings.Split(s, `GET `)
	if len(b) < 2 {
		return LoggedRequest{}, fmt.Errorf("Failed to parse line.")
	}
	c := strings.Split(b[1], ` HTTP`)
	if len(c) < 2 {
		return LoggedRequest{}, fmt.Errorf("Failed to parse line.")
	}

	// Trim the whitespace and assign to logged
	var logged LoggedRequest
	logged.IpAddress = strings.TrimSpace(a[0])
	logged.FileURL = strings.TrimSpace(c[0])

	if logged.IpAddress == "" || logged.FileURL == "" {
		return LoggedRequest{}, fmt.Errorf("Failed to parse line.")
	}

	return logged, nil
}

// Write receives the slice of LoggedRequest and the output path
// and writes the contents to file in json format.
func Write(logBucket []LoggedRequest, path string) error {
	fp, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to create %v. Err: %v.", path, err)
		return err
	}
	defer fp.Close()

	encoder := json.NewEncoder(fp)
	if err = encoder.Encode(logBucket); err != nil {
		log.Fatalf("Unable to encode Json file. Err: %v.", err)
		return err
	}

	return nil
}
