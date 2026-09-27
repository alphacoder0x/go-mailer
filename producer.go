package main

import (
	"encoding/csv"
	"fmt"
	"os"
)
func loadRecipient(filePath string) error{
    
	f,err := os.Open(filePath)
	if err != nil {
		return err
	}

	r  := csv.NewReader(f)
	records,err := r.ReadAll()

	if err != nil{
		return  err
	}

	for _,record := range records[1:] {
		fmt.Println(record[1])
		//send -> consumer -> channels
	}

	return  nil

}