package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type contactList struct {
	Name   string
	Number string
}

type csvTemplate []string

var templateHeaders = csvTemplate{"Name", "Number"}

// twilio variables
var (
	accountSid   = os.Getenv("TWILIO_SID")
	authToken    = os.Getenv("TWILIO_AUTH")
	twilioNumber = os.Getenv("TWILIO_FROM")
	urlStr       = "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
)

const helpMessage = "Enter -u /filepath/to/template.csv to upload a template and message all users\nEnter -d to download a blank template"

// template produces a uploadable csv template
// for a user to fill in and upload to trigger a mass text
func (c csvTemplate) produceTemplate() {
	file, err := os.Create("template.csv")
	if err != nil {
		fmt.Println("Error: Generating Template:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write takes a []string to write each line
	err = writer.Write(c)
	if err != nil {
		fmt.Println("Error: Writing Headers to Template", err)
	}
	fmt.Println("Template Generated")
}

// sendText takes a message string and texts the Name/Number
// from a populated contactList struct
func (c contactList) sendText(message string) {

	// set url values struct to encode
	// for post request to twilio api
	v := url.Values{}
	v.Set("To", c.Number)
	v.Set("From", twilioNumber)
	v.Set("Body", message)

	rb := *strings.NewReader(v.Encode())
	client := &http.Client{}

	req, err := http.NewRequest("POST", urlStr, &rb)
	if err != nil {
		fmt.Println("Error: Twilio Post Request", err)
	}

	defer req.Body.Close()

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: Sending HTTP Request", err)
	}

	defer resp.Body.Close()

	fmt.Println("Message Sent Status:", resp.Status)
}

// populateTemplate returns a []contactList
// populated with Name,Number data to text message
func populateTemplate(fileName string) (data []contactList) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error: opening file template.csv", err)
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

	if len(os.Args) == 1 {
		fmt.Println(helpMessage)
	} else {
		switch userChoice := os.Args[1]; userChoice {
		case "-d":
			templateHeaders.produceTemplate()
		case "-u":
			d := populateTemplate(os.Args[2])
			fmt.Println(d)
		default:
			fmt.Println(helpMessage)
		}
	}
}
