package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"
)

// Rule represents a single rule with its matching criteria and date extraction regex.
type Rule struct {
	VendorMatchRegex      string `json:"vendorMatchRegex"`
	AdditionalMatchRegex  string `json:"additionalMatchRegex"`
	DateExtractionRegex   string `json:"dateExtractionRegex"`
}

// Config holds all the rules loaded from a JSON file.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadRules loads the rules from a JSON file.
func LoadRules(filename string) (*Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// ExtractDate extracts and formats the date from a string using the given regex.
func ExtractDate(text, regex string) (string, error) {
	re := regexp.MustCompile(regex)
	matches := re.FindStringSubmatch(text)
	if len(matches) < 2 {
		return "", fmt.Errorf("date not found")
	}
	parsedDate, err := time.Parse("January 2, 2006", matches[1]) // Adjust this based on the expected date format in the text
	if err != nil {
		return "", err
	}
	return parsedDate.Format("20060102"), nil
}

// ApplyRules applies the rules to the OCR output text and suggests a filename.
func ApplyRules(ocrText string, rules []Rule) (string, error) {
	for _, rule := range rules {
		if vendorMatched, _ := regexp.MatchString(rule.VendorMatchRegex, ocrText); vendorMatched {
			if additionalMatched, _ := regexp.MatchString(rule.AdditionalMatchRegex, ocrText); additionalMatched {
				date, err := ExtractDate(ocrText, rule.DateExtractionRegex)
				if err != nil {
					continue // If date extraction fails, try the next rule
				}
				return fmt.Sprintf("%s - %s.pdf", date, "Google Cloud"), nil // Adjust the vendor name accordingly
			}
		}
	}
	return "", fmt.Errorf("no matching rule found")
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <rules.json> <ocr-output.txt>", os.Args[0])
	}

	rulesFile := os.Args[1]
	ocrFile := os.Args[2]

	// Load rules from JSON file
	config, err := LoadRules(rulesFile)
	if err != nil {
		log.Fatalf("Error loading rules: %v", err)
	}

	// Read OCR output text
	ocrText, err := ioutil.ReadFile(ocrFile)
	if err != nil {
		log.Fatalf("Error reading OCR output: %v", err)
	}

	// Apply rules and suggest a filename
	suggestedFilename, err := ApplyRules(string(ocrText), config.Rules)
	if err != nil {
		log.Fatalf("Error applying rules: %v", err)
	}

	fmt.Println("Suggested Filename:", suggestedFilename)
}

