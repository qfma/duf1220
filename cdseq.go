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

// Takes an ID and returns the fasta sequence from Ensembl
func GetSequence(id string) string {
	client := &http.Client{}

	baseurl := "http://rest.ensembl.org"
	ext := "/sequence/id/" + id + "?type=cds;multiple_sequences=1"
	req, err := http.NewRequest("GET", baseurl+ext, nil)
	// fmt.Println(baseurl + ext)
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

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		id := scanner.Text()
		seq := GetSequence(id)
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
