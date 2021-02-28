package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	hibp "github.com/mattevans/pwned-passwords"
)

const (
	defaultThreadsNumber = 5
	defaultStoreExpiry   = 1 * time.Hour
)

var (
	version = "unknown"
)

func perror(msg string) {
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, msg)
}

func arg_fail(msg string) {
	perror(msg)
	perror("Usage:")
	flag.PrintDefaults()
	os.Exit(2)
}

type CLIArgs struct {
	expire  time.Duration
	threads uint
}

type ScanResult struct {
	Username    string
	Compromised bool
}

type ScanInput struct {
	Username string
	Password string
}

func parseArgs() *CLIArgs {
	var args CLIArgs
	flag.UintVar(&args.threads, "threads", defaultThreadsNumber, "number of threads for network requests")
	flag.DurationVar(&args.expire, "expire", defaultStoreExpiry, "cache TTL")
	flag.Parse()
	if args.threads < 1 {
		args.threads = 1
	}
	if args.expire < 0 {
		args.expire = 0
	}
	return &args
}

func worker(client *hibp.Client, input <-chan *ScanInput, output chan<- *ScanResult) {
	for item := range input {
		pwned, err := client.Pwned.Compromised(item.Password)
		if err != nil {
			log.Printf("hibp client error: %v. Login %s wasn't processed.", err, item.Username)
			continue
		}

		output <- &ScanResult{item.Username, pwned}
	}
}

func collector(results <-chan *ScanResult) {
	for res := range results {
		if res.Compromised {
			fmt.Println(res.Username)
			log.Printf("Found breach for login \"%s\"!", res.Username)
		}
	}
}

func main() {
	args := parseArgs()

	log.Printf("passcheck %s starting...", version)

	client := hibp.NewClient(args.expire)
	inlet := make(chan *ScanInput, args.threads*2)
	outlet := make(chan *ScanResult, args.threads*2)

	var workersWG, collectorWG sync.WaitGroup

	collectorWG.Add(1)
	go func() {
		defer collectorWG.Done()
		collector(outlet)
	}()

	for i := uint(0); i < args.threads; i++ {
		workersWG.Add(1)
		go func() {
			defer workersWG.Done()
			worker(client, inlet, outlet)
		}()
	}

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
		inlet <- &ScanInput{record[0], record[1]}
	}

	close(inlet)
	workersWG.Wait()
	close(outlet)
	collectorWG.Wait()

	log.Printf("passcheck %s finished.", version)
}
