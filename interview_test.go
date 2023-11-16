package customerimporter

import (
	"encoding/csv"
	"reflect"
	"testing"
	"strings"
	"os"
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

func TestSortAndSave(t *testing.T) {
	domainCounts := map[string]int{
			"github.io": 5,
			"cyberchimps.com": 2,
			"faceSmile.net": 7,
	}
	filename := "test_sorted_domains.csv"

	if err := sortAndSave(domainCounts, filename); err != nil {
			t.Fatalf("sortAndSave returned an error: %v", err)
	}

	// Check saved file content
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
