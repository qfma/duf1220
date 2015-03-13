package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// A struct to store a domain from tblout HMMER output
type Domain struct {
	Chromosome string
	Alifrom    int // The start of the HMM alignment within the query
	Alito      int // The stop of the HMM alignment within the query
	// Seqstart   int // The start of the query sequence on the chromosome
	// Seqstop    int // The stop of the query sequence on the chromosome
	Species   string
	SeqLength int
	Features  []Feature
}

// A single coordinate for a given protein ID as returned by
// http://rest.ensembl.org/documentation/info/assembly_translation
type Feature struct {
	Source        string `json:"source"`
	LogicName     string `json:"logic_name"`
	Version       string `json:"version"`
	FeatureType   string `json:"feature_type"`
	ExternalName  string `json:"external_name"`
	Description   string `json:"description"`
	Assemblyname  string `json:"assembly_name"`
	End           int    `json:"end"`
	Start         int    `json:"end"`
	Strand        int    `json:"strand"`
	SeqRegionName string `json:"seq_region_name"`
	Id            string `json:"id"`
	Biotype       string `json:"biotype"`
}

// GetDomainFeature takes a domain struct and returns the feature annotated in
// Ensembl for the corresponding sequence region on a Chromosome
func GetDomainFeature(d *Domain, feature string) {

	var features []Feature
	client := &http.Client{}

	baseurl := "http://rest.ensembl.org"
	// Was needed when using non-Ensembl species
	// domainstart := strconv.Itoa(d.Seqstart + d.Alifrom)
	// domainstop := strconv.Itoa(d.Seqstart + d.Alito)
	ext := "/overlap/region/" + d.Species + "/" + d.Chromosome + ":" + strconv.Itoa(d.Alifrom) + "-" + strconv.Itoa(d.Alito) + "?feature=" + feature

	req, err := http.NewRequest("GET", baseurl+ext, nil)
	req.Header.Set("content-type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&features)

	d.Features = features
}

// Parse a tblout file from HMMER
func ReadTblout(filepath string) []*Domain {

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var domains []*Domain

	_, file := path.Split(filepath)
	species := strings.Split(strings.Split(file, ".")[0], "-")[2]

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// # denotes the start of a comment and can be ignored
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}

		fields := strings.Fields(scanner.Text())
		chromosome := fields[0]
		start, _ := strconv.Atoi(fields[6])
		stop, _ := strconv.Atoi(fields[7])
		seqlen, _ := strconv.Atoi(fields[10])
		// Depending on the strand, start is larger or smaller than stop.
		// switch variables if this is the case
		if start > stop {
			start, stop = stop, start
		}
		hit := Domain{
			Chromosome: chromosome,
			Alifrom:    start,
			Alito:      stop,
			Species:    species,
			SeqLength:  seqlen,
		}
		domains = append(domains, &hit)
	}

	return domains
}

func CountDomainsPerGene(domains []*Domain) map[string]int {
	frequencyForGene := map[string]int{}

	for _, d := range domains {
		time.Sleep(100 * time.Millisecond)
		GetDomainFeature(d, "exon")
		if len(d.Features) != 0 {
			GetDomainFeature(d, "gene")
			// fmt.Printf("%+v\n", d.Features[0].Id)
			frequencyForGene[d.Features[0].Id] += 1
		}
	}
	return frequencyForGene
}

func main() {
	var (
		tblout = flag.String("tblout", "", "A tblout formatted file containing nhmmer hits.")
	)
	flag.Parse()
	log.SetOutput(os.Stdout)
	domains := ReadTblout(*tblout)
	_, file := path.Split(*tblout)
	species := strings.Split(strings.Split(file, ".")[0], "-")[2]

	genes := CountDomainsPerGene(domains)
	// fmt.Println("Gene", "DUFCount")
	for k, v := range genes {
		fmt.Println(k, v)
	}

	var counter int

	for _, d := range domains {
		time.Sleep(100 * time.Millisecond)
		GetDomainFeature(d, "exon")
		if len(d.Features) == 2 {
			continue
		} else if len(d.Features) == 0 {
			counter += 1
		} else if len(d.Features) == 1 {
			continue
		}
		// fmt.Printf("%+v\n", d)
	}
	fmt.Println(species, len(domains), len(domains)-counter, counter)

}
