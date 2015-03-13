package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Takes an ID and returns the fasta sequence of the CDS from Ensembl
func GetSequence(id string) string {
	client := &http.Client{}

	baseurl := "http://rest.ensembl.org"
	ext := "/sequence/id/" + id + "?type=cds;multiple_sequences=1"
	req, err := http.NewRequest("GET", baseurl+ext, nil)
	req.Header.Set("content-type", "text/x-fasta")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	seq, err := ioutil.ReadAll(resp.Body)
	return string(seq)
}

func main() {

	log.SetOutput(os.Stderr)

	// Make a scanner to receive a stream of IDs from Stdin
	// pass the IDs to the Ensembl REST API and write the cds sequence
	// as fasta
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		id := scanner.Text()
		seq := GetSequence(id)
		// Be polite, don't hammer the API
		time.Sleep(100 * time.Millisecond)
		d := []byte(seq)
		err := ioutil.WriteFile(id+".cds.all.fa", d, 0644)
		if err != nil {
			panic(err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
