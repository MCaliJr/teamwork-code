package customerimporter

import (
	"encoding/csv"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestReadCSV(t *testing.T) {
    filename := "test_customers.csv"
    expected := [][]string{
        {"first_name", "last_name", "email", "gender", "ip_address"},
        {"John", "Doe", "thatGuyDoe@faceSmile.net", "Male", "53.191.87.821"},
        {"Mildred", "Hernandez", "mhernandez0@github.io", "Female", "38.194.51.128"},
        {"Bonnie", "Ortiz", "bortiz1@cyberchimps.com", "Female", "197.54.209.129"},
        {"Dennis", "Henry", "dhenry2@hubpages.com", "Male", "155.75.186.217"},
    }

    records, err := readCSV(filename)
    if err != nil {
        t.Errorf("readCSV returned an error: %v", err)
    }

		// Convert records to maps for comparison with ignored order
		expectedMap := make(map[string]struct{})
		for _, record := range expected {
				expectedMap[strings.Join(record, ",")] = struct{}{}
		}

		recordsMap := make(map[string]struct{})
		for _, record := range records {
				recordsMap[strings.Join(record, ",")] = struct{}{}
		}

		if !reflect.DeepEqual(expectedMap, recordsMap) {
				t.Errorf("readCSV did not return the expected records")
		}
}

func TestCountEmailDomains(t *testing.T) {
	input := [][]string{
			{"first_name", "last_name", "email", "gender", "ip_address"},
			{"John", "Doe", "thatGuyDoe@faceSmile.net", "Male", "53.191.87.821"},
			{"Mildred", "Hernandez", "mhernandez0@github.io", "Female", "38.194.51.128"},
			{"Another", "GitUser", "someuser@github.io", "Male", "45.22.321.128"},
			{"Bonnie", "Ortiz", "bortiz1@cyberchimps.com", "Female", "197.54.209.129"},
	}
	expected := map[string]int{
			"github.io": 2,
			"cyberchimps.com": 1,
			"faceSmile.net": 1,
	}

	result, err := countEmailDomains(input)
	if err != nil {
			t.Fatalf("countEmailDomains returned an error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
			t.Errorf("countEmailDomains returned %v, expected %v", result, expected)
	}
}

func TestSortDomains(t *testing.T) {
	domainCounts := map[string]int{
			"github.io": 5,
			"cyberchimps.com": 2,
			"faceSmile.net": 7,
	}
	expected := []DomainCount{
			{Domain: "cyberchimps.com", Count: 2},
			{Domain: "faceSmile.net", Count: 7},
			{Domain: "github.io", Count: 5},
	}

	result := sortDomains(domainCounts)
	if !reflect.DeepEqual(result, expected) {
			t.Errorf("sortDomains returned %v, expected %v", result, expected)
	}
}

func TestSaveToFile(t *testing.T) {
	sortedDomains := []DomainCount{
			{Domain: "faceSmile.net", Count: 7},
			{Domain: "github.io", Count: 5},
			{Domain: "cyberchimps.com", Count: 2},
	}
	filename := "test_save_to_file.csv"

	if err := saveToFile(sortedDomains, filename); err != nil {
			t.Fatalf("saveToFile returned an error: %v", err)
	}

	// Verify the file content
	file, err := os.Open(filename)
	if err != nil {
			t.Fatalf("Failed to open the file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
			t.Fatalf("Failed to read from the file: %v", err)
	}

	expectedLines := []string{
			"faceSmile.net,7",
			"github.io,5",
			"cyberchimps.com,2",
	}

	for i, line := range lines {
			joinedLine := strings.Join(line, ",")
			if joinedLine != expectedLines[i] {
					t.Errorf("Line %d of file is incorrect, got: %s, want: %s", i, joinedLine, expectedLines[i])
			}
	}
}

func TestProcessCustomers(t *testing.T) {
	inputCSV := "test_customers.csv"
	outputCSV := "test_output.csv"

	// Expected result based on the provided input
	expectedDomainCounts := []DomainCount{
			{Domain: "cyberchimps.com", Count: 1},
			{Domain: "faceSmile.net", Count: 1},
			{Domain: "github.io", Count: 1},
			{Domain: "hubpages.com", Count: 1},
	}

	// Test without saving to file
	result, err := ProcessCustomers(inputCSV)
	if err != nil {
			t.Fatalf("ProcessCustomers without file saving returned an error: %v", err)
	}
	if !reflect.DeepEqual(result, expectedDomainCounts) {
			t.Errorf("ProcessCustomers without file saving returned %v, expected %v", result, expectedDomainCounts)
	}

	// Test with saving to file
	_, err = ProcessCustomers(inputCSV, outputCSV)
	if err != nil {
			t.Fatalf("ProcessCustomers with file saving returned an error: %v", err)
	}

	// Verify the file content
	file, err := os.Open(outputCSV)
	if err != nil {
			t.Fatalf("Failed to open the output file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
			t.Fatalf("Failed to read from the output file: %v", err)
	}

	// Validate each line in the output file
	for i, line := range lines {
			expectedLine := expectedDomainCounts[i].Domain + "," + strconv.Itoa(expectedDomainCounts[i].Count)
			joinedLine := strings.Join(line, ",")
			if joinedLine != expectedLine {
					t.Errorf("Line %d of output file is incorrect, got: %s, want: %s", i, joinedLine, expectedLine)
			}
	}

	if len(lines) > len(expectedDomainCounts) {
			t.Errorf("Output file has more lines (%d) than expected (%d)", len(lines), len(expectedDomainCounts))
	}
}
