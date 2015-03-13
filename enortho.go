package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// This converts the json output from the Ensembl REST API to a go struct
type Homologs []struct {
	TaxonomyLevel  string `json:"taxonomy_level"`
	ProteinId      string `json:"protein_id"`
	Species        string `json:"species"`
	Id             string `json:"id"`
	Type           string `json:"type"`
	MethodLinkType string `json:"method_link_type"`
}

// Takes an ID and returns the JSON output as Homolog structs
// The map[string][]map[string] part is necessary, because of the nested JSON
func GetOrthologous(id string) map[string][]map[string]Homologs {
	client := &http.Client{}

	var data map[string][]map[string]Homologs

	baseurl := "http://rest.ensembl.org"
	ext := "/homology/id/" + id + "?format=condensed"
	req, err := http.NewRequest("GET", baseurl+ext, nil)
	req.Header.Set("content-type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)

	return data
}

// Get us some simple print output for the homologs
func PrintSpeciesOrthologs(homologies Homologs, species []string) {
	for _, s := range species {
		for _, h := range homologies {
			if h.Species == s {
				fmt.Println(h.Id, s, h.Type)
			}
		}
	}
}

func main() {
	var (
		species = flag.String("species", "", "A comma seperated list of species")
	)
	flag.Parse()
	target_species := strings.Split(*species, ",")
	log.SetOutput(os.Stderr)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		id := scanner.Text()
		data := GetOrthologous(id)

		// Check if we found any homologs
		if len(data["data"]) != 0 {
			homologs := data["data"][0]["homologies"]
			PrintSpeciesOrthologs(homologs, target_species)
		}
		// Be polite, don't hammer the API
		time.Sleep(100 * time.Millisecond)

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
