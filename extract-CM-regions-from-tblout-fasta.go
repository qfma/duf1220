package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	// "io/ioutil"
	"log"
	// "net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Region struct {
	Chromosome string
	Start      int
	Stop       int
	Species    string
}

// Takes a region and a flank returns the fasta sequence from ensembl
func GetRegionsFasta(regions []Region, flank int, infile string) string {

	f, err := os.Open(infile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var id string
	var fasta string
	var buffer bytes.Buffer

	chr := make(map[string]*bytes.Buffer)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), ">") {
			id = strings.Fields(scanner.Text()[1:])[0]
			// fmt.Println(id)
			for _, r := range regions {
				if r.Chromosome == id {
					fmt.Println(r.Chromosome, id)
					chr[id] = &buffer
				}
			}
		} else if _, ok := chr[id]; ok {
			seq := strings.ToUpper(scanner.Text())
			chr[id].WriteString(seq)
		}
	}

	for _, r := range regions {

		var start, stop int

		if (r.Start - flank) < 0 {
			start = 0
		} else {
			start = r.Start - flank
		}

		chunk := (r.Stop + flank) - start

		if chunk >= utf8.RuneCountInString(chr[r.Chromosome].String()) {
			fmt.Println("Chunk too large! Using entire scaffold instead..")
			fasta = fasta + ">" + r.Chromosome + " " + strconv.Itoa(r.Start) + "-" + strconv.Itoa(r.Stop) + "\n"
			fasta = fasta + chr[r.Chromosome].String() + "\n"
		} else {
			stop = r.Stop + flank
			fasta = fasta + ">" + r.Chromosome + " " + strconv.Itoa(start) + "-" + strconv.Itoa(stop) + "\n"
			fmt.Println(r.Chromosome, utf8.RuneCountInString(chr[r.Chromosome].String()), start, stop, r.Start, r.Stop)
			fasta = fasta + chr[r.Chromosome].String()[start:stop] + "\n"
		}
	}
	return fasta
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
	return seqmap
}

// Reads a nhmmer tblout output file and for each hit returns a region struct
// that is stored in a slice
func ReadTblout(filepath string) []Region {

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var regions []Region

	_, file := path.Split(filepath)
	species := strings.Split(strings.Split(file, ".")[0], "-")[2]

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// # denotes the start of a comment and can be ignored
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		} else {
			fields := strings.Fields(scanner.Text())
			chromosome := fields[0]
			start, _ := strconv.Atoi(fields[6])
			stop, _ := strconv.Atoi(fields[7])
			// Depending on the strand, start is larger or smaller than stop.
			// switch variables if this is the case
			if start > stop {
				start, stop = stop, start
			}
			hit := Region{
				Chromosome: chromosome,
				Start:      start,
				Stop:       stop,
				Species:    species,
			}
			regions = append(regions, hit)
		}

	}

	return regions
}

// Write the extracted regions as a single fasta file.
func WriteRegions(fname string, seqs string) {

	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(seqs)

	return
}

// As regions are extended by a defined flanking size, we have to check if
// regions now overlap and correct this. Otherwise domain annotation may produce
// duplicates.
func CheckOverlap(regions []Region, flank int) {

	for _, r1 := range regions {
		for _, r2 := range regions {
			if r1 != r2 && r1.Chromosome == r2.Chromosome {
				if r2.Start+flank > r1.Start+flank && r2.Start+flank < r1.Stop+flank {
					fmt.Println("Overlap detected!")
					fmt.Println(r1, r2)
				} else if r2.Stop+flank < r1.Stop+flank && r2.Stop+flank > r1.Start+flank {
					fmt.Println("Overlap detected!")
					fmt.Println(r1, r2)
				}
			}
		}
	}
	return
}

func main() {
	var (
		tblout = flag.String("tblout", "", "A tblout formatted file containing nhmmer hits.")
		genome = flag.String("genome", "", "A tblout formatted file containing nhmmer hits.")
		fasta  = flag.String("fasta", "test", "The extended regions as fasta file")
		flank  = flag.Int("flank", 100000, "The length of the flanking region")
	)
	flag.Parse()
	log.SetOutput(os.Stdout)
	regions := ReadTblout(*tblout)
	CheckOverlap(regions, *flank)

	seqs := GetRegionsFasta(regions, *flank, *genome)

	WriteRegions(*fasta, seqs)
}
