package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

// loadRecipient reads the CSV at filePath and sends each row to ch as a
// Recipient. The first row is treated as a header (Name,Email) and skipped.
// ch is always closed on return so the workers know there is no more work.
func loadRecipient(filePath string, ch chan Recipient) error {

	defer close(ch)

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer f.Close()
	r := csv.NewReader(f)
	records, err := r.ReadAll()

	if err != nil {
		return err
	}

	f.Close()

	// records[1:] skips the header row. Each send blocks until a worker is free.
	for _, record := range records[1:] {
		fmt.Println("producing data ============>")

		ch <- Recipient{
			Name:  record[0],
			Email: record[1],
		}
		//send -> consumer -> channels
	}

	return nil

}
