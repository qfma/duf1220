package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Takes a GID Ensembl identifier and returns the gene plus a defined
// upstream region of this gene in fasta format.
func GetUpstreamRegion(gid string, upstream int) string {
	client := &http.Client{}
	baseurl := "http://rest.ensembl.org"
	// ext := "/sequence/id/" + gid + "?expand_5prime=1000;expand_3prime=1000"
	ext := "/sequence/id/" + gid + "?expand_5prime=" + strconv.Itoa(upstream)

	// fmt.Println(baseurl + ext)
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

// Takes a filename, sequence identifier and seq and writes it as a fasta format
func WriteFastaPromotor(fname string, length int, seqs map[string]string) {

	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for id, s := range seqs {
		fmt.Println(id, len(s))
		f.WriteString(">" + id + " " + strconv.Itoa(length) + "bp upstream promoter" + "\n")
		f.WriteString(s[:length] + "\n")
	}

	return
}

// Takes a DNA string (all capital letters) and converts it into the reverse
// complement of the sequence.
func ReverseComplementDNA(DNA string) string {
	reverse := map[rune]rune{
		'A': 'T',
		'C': 'G',
		'G': 'C',
		'T': 'A',
	}
	rcDNA := make([]rune, len(DNA))
	start := len(DNA)

	for _, c := range DNA {
		// quietly skip invalid UTF-8
		if c != utf8.RuneError {
			start--
			// Reverse the DNA bases
			rcDNA[start] = reverse[c]
		}
	}
	return string(rcDNA[start:])
}

// Takes a DNA string (all capital letters) and converts it into the reverse
// of the sequence.
func ReverseDNA(DNA string) string {

	rcDNA := make([]rune, len(DNA))
	start := len(DNA)

	for _, c := range DNA {
		// quietly skip invalid UTF-8
		if c != utf8.RuneError {
			start--
			// Reverse the DNA bases
			rcDNA[start] = c
		}
	}
	return string(rcDNA[start:])
}

// Parses a fasta file from a string and returns a map where each key is a GID
// and each entry is the sequence of the entry
func ReadFastaFromString(fasta string) map[string]string {
	// Make a new scanner from the fasta file string
	scanner := bufio.NewScanner(strings.NewReader(fasta))
	var id, strand string
	// Store the sequences and id
	seqmap := make(map[string]string)
	// keep track of the strandedness
	strandmap := make(map[string]string)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), ">") {
			fields := strings.Fields(scanner.Text())
			strand = strings.Split(fields[1], ":")[5]
			id = fields[0][1:] + ":" + strand
			strandmap[id] = strand
		} else {
			// Remove potential space characters with nothing
			seq := strings.Replace(strings.ToUpper(scanner.Text()), " ", "", -1)
			seqmap[id] = seqmap[id] + seq
		}
	}

	// Check if the strand is the reverse and translate accordingly into reverse
	// complement
	// for id, seq := range seqmap {
	// 	if strandmap[id] == "-1" {
	// 		// seqmap[id] = ReverseComplementDNA(seq)
	// 		seqmap[id] = seq
	// 	} else {
	// 		// seqmap[id] = ReverseDNA(seq)
	// 		seqmap[id] = seq
	// 	}
	// }

	return seqmap
}

// Wrapper to get the promoter regions for a given list of gids and a defined
// upstream size.
func GetPromoterRegions(gids []string, upstream int) map[string]string {

	PromoterRegions := make(map[string]string)

	for _, gid := range gids {
		// Returns a fasta file
		promoter := GetUpstreamRegion(gid, upstream)
		promotermap := ReadFastaFromString(promoter)
		for k, v := range promotermap {
			PromoterRegions[k] = v
		}
	}

	return PromoterRegions
}

func main() {
	var (
		fasta    = flag.String("fasta", "", "Fasta formated output file storing the nucleotide sequences of the upstream promoter region")
		upstream = flag.Int("upstream", 900, "The length of the upstream region")
		cm       = flag.Bool("cm", false, "Use NBPF genes with CM promoter")
		evi5     = flag.Bool("evi5", false, "Use NBPF genes with EVI5 promoter")
	)
	flag.Parse()
	log.SetOutput(os.Stdout)

	// All of these are protein coding genes and not pseudogenes
	// These are the three genes that have a CM promoter
	humanCM := []string{
		"ENSG00000196427", // NBPF4
		"ENSG00000186086", // NBPF6
		"ENSG00000215864", // NBPF7
	}

	// All of these are protein coding genes and not pseudogenes
	// These genes have a EVI5 promoter and a CM as the fourth intron
	// See O'bleness et al. paper
	humanNHPF := []string{
		"ENSG00000219481", // NBPF1
		"ENSG00000142794", // NBPF3
		"ENSG00000162825", // NBPF8
		"ENSG00000269713", // NBPF9
		"ENSG00000271425", // NBPF10
		"ENSG00000263956", // NBPF11
		"ENSG00000268043", // NBPF12
		"ENSG00000270629", // NBPF14
		"ENSG00000266338", // NBPF15 (16 is now merged)
	}

	if *cm == true {
		PromoterRegions := GetPromoterRegions(humanCM, *upstream)
		WriteFastaPromotor(*fasta, *upstream, PromoterRegions)
	}

	if *evi5 == true {
		PromoterRegions := GetPromoterRegions(humanNHPF, *upstream)
		WriteFastaPromotor(*fasta, *upstream, PromoterRegions)
	}
}
