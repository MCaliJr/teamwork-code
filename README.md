# Customer Importer

## Overview

The `customerimporter` package is designed to process a CSV file containing customer data and return a sorted list of email domains along with the count of customers associated with each domain. This package is optimized for performance and is capable of handling large datasets efficiently.

## Features

- **Data Processing**: Reads customer data from a CSV file and extracts email domains.
- **Sorting**: Sorts the email domains alphabetically, along with the count of customers for each domain.
- **Performance**: Optimized for large datasets, suitable for files with thousands to millions of lines.
- **Error Handling**: Robust error handling throughout the data processing pipeline.
- **Optional File Saving**: Ability to save the sorted output to a CSV file.

## Usage

### Importing the Package

Import the `customerimporter` package into your Go project:

```go
import "path/to/customerimporter"
```

### Functionality

The primary function in this package is `ProcessCustomers`, which can be used as follows:

```go
import "path/to/customerimporter"

func main() {
    inputFile := "path/to/your/input.csv" // Input CSV file with customer data

    // To process data and get the result without saving to a file
    result, err := customerimporter.ProcessCustomers(inputFile)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Sorted Domain Counts:", result)

    // To process data and save the result to a file
    outputFile := "path/to/your/output.csv" // Desired output CSV file
    _, err = customerimporter.ProcessCustomers(inputFile, outputFile)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Data processed and saved to:", outputFile)
}
```

## Testing

### Overview

Comprehensive tests have been written to ensure the reliability and correctness of the package's functionality.

### Test Cases

- **TestReadCSV**: Validates the correct reading and parsing of the CSV file.
- **TestCountEmailDomains**: Ensures accurate counting of email domains from the customer data.
- **TestSortDomains**: Checks the sorting functionality of email domains.
- **TestSaveToFile**: Verifies the ability to save sorted data to an output file.
- **TestProcessCustomers**: A holistic test that covers the entire processing pipeline, from reading the CSV file to returning/saving sorted domain counts.

### Running Tests

To run the tests, use the following command in the project directory:

```bash
go test
```

## Additional Information

This package is optimized for performance, making it suitable for large datasets. Error handling and logging are integral parts of the package, ensuring robustness in various scenarios.
