package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)


type Region struct {
	Chromosome string
	Start      int
	Stop       int
	Species    string
	SeqLength  int
	Skip       bool
}

// ByAge implements sort.Interface for []*Region based on
// the Start field.
type ByStart []*Region

func (a ByStart) Len() int           { return len(a) }
func (a ByStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStart) Less(i, j int) bool { return a[i].Start < a[j].Start }

// Adds a flanking region of a given size to each hit
func (r *Region) AddFlank(flank int) {
	if (r.Start - flank) < 0 {
		r.Start = 0
	} else {
		r.Start = r.Start - flank
	}
	if (r.Stop + flank) >= r.SeqLength {
		r.Stop = r.SeqLength
	} else {
		r.Stop = r.Stop + flank
	}
}

// Takes a region and returns the fasta sequence from ensembl
func GetRegionsFasta(regions []*Region, infile string) string {

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
			for _, r := range regions {
				if r.Chromosome == id {
					chr[id] = &buffer
				}
			}
		} else if _, ok := chr[id]; ok {
				seq := strings.ToUpper(scanner.Text())
				chr[id].WriteString(seq)
			}
	}

for _, r := range regions {

	if r.Stop >= utf8.RuneCountInString(chr[r.Chromosome].String()) {
		fmt.Println("Chunk too large! Using entire scaffold instead..")
		fasta = fasta + ">" + r.Chromosome + " " + strconv.Itoa(r.Start) + "-" + strconv.Itoa(r.Stop) + "\n"
		fasta = fasta + chr[r.Chromosome].String() + "\n"
	} else {
		fasta = fasta + ">" + r.Chromosome + " " + strconv.Itoa(r.Start) + "-" + strconv.Itoa(r.Stop) + "\n"
		// fmt.Println(r.Chromosome, utf8.RuneCountInString(chr[r.Chromosome].String()), start, stop, r.Start, r.Stop)
		fasta = fasta + chr[r.Chromosome].String()[r.Start:r.Stop] + "\n"
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
func ReadPSL(filepath string) []*Region {

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var regions []*Region
	var startsignal bool

	_, file := path.Split(filepath)
	species := strings.Split(strings.Split(file, ".")[0], "-")[2]

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// Treat all lines after ---- as fields
		if strings.HasPrefix(scanner.Text(), "-") {
			startsignal = true
			continue
		}
		// Once we have got the startsignal, start reading fields
		if startsignal {
			fields := strings.Fields(scanner.Text())
			chromosome := fields[13]
			start, _ := strconv.Atoi(fields[15])
			stop, _ := strconv.Atoi(fields[16])
			seqlen, _ := strconv.Atoi(fields[14])
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
				SeqLength:  seqlen,
				Skip:       false,
			}
			regions = append(regions, &hit)
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

// This function resolves overlaps between close regions on the same
// chromosome. Regions may overlap, because the flanking region extends into
// another region.
// Therefore situations like this can occur:
// Overlap r1 and r2 like this:
//            r2.Start----------------------r2.Stop
//r1.Start--------------------r1.Stop
//
// The functions orders regions on the same chromosome and then resolves the
// above overlaps.
func ResolveOverlappingRegions(regions []*Region, flank int) {

	ChromosomalRegions := make(map[string][]*Region)

	for _, r := range regions {
		ChromosomalRegions[r.Chromosome] = append(ChromosomalRegions[r.Chromosome], r)
	}

	for _, regions := range ChromosomalRegions {
		sort.Sort(ByStart(regions))

		if len(regions) >= 2 {
			for i := 0; i < len(regions)-1; i++ {
				r1 := regions[i]
				r2 := regions[i+1]
				// If the r2 region is within the r1 region, skip r2
				if r1.Stop > r2.Start && r1.Stop >= r2.Stop {
					r2.Skip = true
				} else if r1.Stop > r2.Start {
					// Overlap r1 and r2 like this:
					//            r2.Start----------------------r2.Stop
					//r1.Start--------------------r1.Stop
					fmt.Println("Overlap detected! Resolving..")
					r2.Start = r1.Stop + 1
					if r2.Start >= r2.SeqLength {
						r2.Start = r2.SeqLength
					}
				}
				fmt.Println(r1, r2)
			}
		}
	}

	return
}

func RemoveRegions(regions []*Region) []*Region {
	var newRegions []*Region

	for _, r := range regions {
		if r.Skip != true {
			newRegions = append(newRegions, r)
		}
	}
	return newRegions
}

func main() {
	var (
		psl = flag.String("psl", "", "A psl formatted file containing BLAT hits.")
		genome = flag.String("genome", "", "A tblout formatted file containing nhmmer hits.")
		fasta  = flag.String("fasta", "test", "The extended regions as fasta file")
		flank  = flag.Int("flank", 100000, "The length of the flanking region")
	)
	flag.Parse()
	log.SetOutput(os.Stdout)
	fmt.Println(*psl)
	regions := ReadPSL(*psl)
	fmt.Println(regions)
	for _, r := range regions {
		r.AddFlank(*flank)
	}

	ResolveOverlappingRegions(regions, *flank)
	validregions := RemoveRegions(regions)

	seqs := GetRegionsFasta(validregions, *genome)

	WriteRegions(*fasta, seqs)
}
