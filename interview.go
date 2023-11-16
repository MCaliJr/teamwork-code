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
			return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, dc := range counts {
			if err := writer.Write([]string{dc.Domain, strconv.Itoa(dc.Count)}); err != nil {
					return err
			}
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
