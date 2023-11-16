// package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package customerimporter

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// readCSV reads data from a CSV file using concurrency.
func readCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return processCSVLines(file)
}

// processCSVLines processes lines from an open file concurrently.
func processCSVLines(file *os.File) ([][]string, error) {
	var records [][]string
	var mutex sync.Mutex
	var wg sync.WaitGroup

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
			data := scanner.Text()
			wg.Add(1)
			go processLine(data, &records, &mutex, &wg)
	}

	if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("scanner error: %w", err)
	}

	wg.Wait()
	return records, nil
}

// processLine processes a single line of CSV data using goroutine 
func processLine(data string, records *[][]string, mutex *sync.Mutex, wg *sync.WaitGroup) error {
	defer wg.Done()
	reader := csv.NewReader(strings.NewReader(data))
	record, err := reader.Read()
	if err != nil {
			return fmt.Errorf("reader error: %w", err)
	}

	mutex.Lock()
	*records = append(*records, record)
	mutex.Unlock()

	return nil
}

// countEmailDomains uses concurrency for faster execution
func countEmailDomains(records [][]string) (map[string]int, error) {
	if len(records) == 0 {
			return nil, errors.New("no records provided")
	}

	emailColumn, err := findEmailColumn(records[0])
	if err != nil {
			return nil, fmt.Errorf("error finding email column: %w", err)
	}

	var wg sync.WaitGroup
	domainCounts := make(map[string]int)
	mutex := &sync.Mutex{}

	for _, record := range records {
			wg.Add(1)
			go func(rec []string) {
					defer wg.Done()
					if len(rec) <= emailColumn {
							return
					}
					email := rec[emailColumn]
					parts := strings.Split(email, "@")
					if len(parts) != 2 {
							return
					}
					domain := parts[1]

					mutex.Lock()
					domainCounts[domain]++
					mutex.Unlock()
			}(record)
	}

	wg.Wait()
	return domainCounts, nil
}

// findEmailColumn identifies index of the "email" column
func findEmailColumn(record []string) (int, error) {
	for i, field := range record {
			if strings.Contains(field, "@") || strings.ToLower(field) == "email" {
					return i, nil
			}
	}
	return -1, errors.New("email column not found")
}

type DomainCount struct {
	Domain string
	Count  int
}

// sortDomains sorts the domain counts alphabetically by domain name and returns them in a sorted slice of DomainCount.
func sortDomains(domainCounts map[string]int) []DomainCount {
	var counts []DomainCount
	for domain, count := range domainCounts {
			counts = append(counts, DomainCount{Domain: domain, Count: count})
	}

	sort.Slice(counts, func(i, j int) bool {
			return counts[i].Domain < counts[j].Domain // alphabetical order
	})

	return counts
}

// saveToFile saves the sorted domain counts to a specified CSV file.
func saveToFile(sortedDomains []DomainCount, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("error creating output file: %w", err)
    }
    defer file.Close()

    buffer := bufio.NewWriter(file)
    writer := csv.NewWriter(buffer)
    defer writer.Flush()

    for _, dc := range sortedDomains {
        if err := writer.Write([]string{dc.Domain, strconv.Itoa(dc.Count)}); err != nil {
            return err
        }
    }

    return buffer.Flush()
}

// ProcessCustomers reads the CSV, counts and sorts email domains.
// If an outputFile is provided, it saves the result to the file.
func ProcessCustomers(inputFile string, outputFile ...string) ([]DomainCount, error) {
	records, err := readCSV(inputFile)
	if err != nil {
			return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	domainCounts, err := countEmailDomains(records)
	if err != nil {
			return nil, fmt.Errorf("failed to count email domains: %w", err)
	}

	sortedDomains := sortDomains(domainCounts)

	// Save to file if an outputFile name is provided
	if len(outputFile) > 0 {
			if err := saveToFile(sortedDomains, outputFile[0]); err != nil {
					return sortedDomains, fmt.Errorf("failed to save domains: %w", err)
			}
	}

	return sortedDomains, nil
}
