// package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package customerimporter

import (
	"bufio"
	"encoding/csv"
	"os"
	"sync"
	"strings"
	"errors"
)

// read data from a CSV file using concurrency for faster read time
func readCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
			return nil, err
	}
	defer file.Close()

	var records [][]string
	var mutex sync.Mutex
	var wg sync.WaitGroup

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
			wg.Add(1)
			go func(data string) {
					defer wg.Done()
					reader := csv.NewReader(strings.NewReader(data))
					record, err := reader.Read()
					if err != nil {
							// TODO handle / log error
							return
					}

					mutex.Lock()
					records = append(records, record)
					mutex.Unlock()
			}(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
			return nil, err
	}

	wg.Wait()
	return records, nil
}

func countEmailDomains(records [][]string) (map[string]int, error) {
	if len(records) == 0 {
			return nil, errors.New("no records provided")
	}

	emailColumn, err := findEmailColumn(records[0])
	if err != nil {
			return nil, err
	}

	domainCounts := make(map[string]int)
	for _, record := range records {
			if len(record) <= emailColumn {
					continue // skip row without email data
			}
			email := record[emailColumn]
			parts := strings.Split(email, "@")
			if len(parts) != 2 {
					continue // skip malformed email addresses
			}
			domain := parts[1]
			domainCounts[domain]++
	}
	return domainCounts, nil
}

// identify the index of email column
func findEmailColumn(record []string) (int, error) {
	for i, field := range record {
			if strings.Contains(field, "@") || strings.ToLower(field) == "email" {
					return i, nil
			}
	}
	return -1, errors.New("email column not found")
}