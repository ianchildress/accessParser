package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
)

type LoggedRequest struct {
	IpAddress string `json:"ip"`
	FileURL   string `json:"url"`
}

func main() {
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
		logged, err := Parse(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		logBucket = append(logBucket, logged)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if err := Write(logBucket, *outFile); err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully parsed file to %s", *outFile)

}

// Parse receives a string and pulls the ip address
// and file path from the access log.
func Parse(s string) (LoggedRequest, error) {
	// Get the ip address
	a := strings.Split(s, " - - ")

	// Get the image url
	b := strings.Split(s, `GET `)
	c := strings.Split(b[1], ` HTTP`)

	// Trim the whitespace and assign to logged
	var logged LoggedRequest
	logged.IpAddress = strings.TrimSpace(a[0])
	logged.FileURL = strings.TrimSpace(c[0])
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
