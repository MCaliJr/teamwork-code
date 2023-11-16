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

func sortAndSave(domainCounts map[string]int, filename string) error {
	var counts []DomainCount
	for domain, count := range domainCounts {
			counts = append(counts, DomainCount{Domain: domain, Count: count})
	}

	sort.Slice(counts, func(i, j int) bool {
			return counts[i].Count > counts[j].Count // desc order
	})

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer file.Close()

	buffer := bufio.NewWriter(file)
	writer := csv.NewWriter(buffer)
	defer writer.Flush()

	for _, dc := range counts {
			if err := writer.Write([]string{dc.Domain, strconv.Itoa(dc.Count)}); err != nil {
					return err
			}
	}
	
	if err := buffer.Flush(); err != nil {
		return fmt.Errorf("error flushing buffer to file: %w", err)
	}

	return nil
}

func ProcessCustomers(inputFile, outputFile string) error {
	records, err := readCSV(inputFile)
	if err != nil {
			return fmt.Errorf("failed to read CSV: %w", err)
	}

	domainCounts, err := countEmailDomains(records)
	if err != nil {
			return fmt.Errorf("failed to count email domains: %w", err)
	}

	if err := sortAndSave(domainCounts, outputFile); err != nil {
			return fmt.Errorf("failed to sort and save domains: %w", err)
	}

	return nil
}
