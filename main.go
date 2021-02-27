package main

import (
	"fmt"
	"os"
	"time"
	"log"
	"encoding/csv"
	"io"

	hibp "github.com/mattevans/pwned-passwords"
)

const (
	storeExpiry = 1 * time.Hour
)

var (
	version = "unknown"
)

func main() {
	client := hibp.NewClient(storeExpiry)

	log.Printf("passcheck %s starting...", version)
	log.Println("Reading CSV file with login,password pairs from STDIN...")

	r := csv.NewReader(os.Stdin)
	r.FieldsPerRecord = 2
	r.LazyQuotes = true

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		login, password := record[0], record[1]

		pwned, err := client.Pwned.Compromised(password)
		if err != nil {
			log.Printf("hibp client error: %v. Login %s wasn't processed.", err, login)
			continue
		}

		if pwned {
			fmt.Println(login)
			log.Printf("Found breach for login \"%s\"!", login)
		}
	}
	log.Printf("passcheck %s finished.", version)
}
