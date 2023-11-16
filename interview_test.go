package customerimporter

import (
    "reflect"
    "testing"
		"strings"
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
