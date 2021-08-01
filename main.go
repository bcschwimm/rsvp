package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type contactList struct {
	Name   string
	Number string
}

type csvTemplate []string

var templateHeaders = csvTemplate{"Name", "Number"}

// template produces a uploadable csv template
// for a user to fill in and upload to trigger a mass text
func (c csvTemplate) produceTemplate() {
	file, err := os.Create("template.csv")
	if err != nil {
		fmt.Println("Error Generating Template:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write takes a []string to write each line
	err = writer.Write(c)
	if err != nil {
		fmt.Println("Error Writing Headers to Template", err)
	}
	fmt.Println("Template Generated")
}

// sendText takes a message string and texts each Name/Number
// from a populated contactList struct
func (c contactList) sendText(message string) {
	fmt.Println(message, "WIP")
}

// populateTemplate returns a []contactList
// populated with Name,Number data to text message
func populateTemplate(fileName string) (data []contactList) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file template.csv", err)
	}

	r := csv.NewReader(file)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		rowData := contactList{
			Name:   record[0],
			Number: record[1],
		}
		data = append(data, rowData)
	}
	return
}

func main() {
	d := populateTemplate("template.csv")
	fmt.Println(d)
}
