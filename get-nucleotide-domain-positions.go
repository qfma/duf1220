package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Read a fasta alignment from an infile
// The header is used as the key for the dictionary
// No sanity checks are performed as of now
func ReadFastaAlignment(fname string) map[string]string {

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var id string
	aln := make(map[string]string)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), ">") {
			id = strings.Fields(scanner.Text()[1:])[0]
			aln[id] = ""
		} else {
			// Remove potential space characters with nothing
			seq := strings.Replace(strings.ToUpper(scanner.Text()), " ", "", -1)
			aln[id] = aln[id] + seq
		}

	}

	return aln
}

type Protein struct {
	PID, GID, TID string
	Length        int
	Domains       []Domain
}

type Domain struct {
	ID, Name    string
	Start, Stop int
}

func ReadDomtblout(fname string) map[string]*Protein {

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	proteins := make(map[string]*Protein)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// # denotes the start of a comment and can be ignored
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		} else {
			fields := strings.Fields(scanner.Text())
			pid := fields[0]
			start, _ := strconv.Atoi(fields[17])
			stop, _ := strconv.Atoi(fields[18])
			domain := Domain{
				ID:    fields[4],
				Name:  fields[3],
				Start: start,
				Stop:  stop,
			}
			if protein, ok := proteins[pid]; ok {
				protein.Domains = append(protein.Domains, domain)
			} else {
				proteins[pid] = new(Protein)
				proteins[pid].PID = pid
				proteins[pid].GID = strings.Split(fields[24], ":")[1]
				proteins[pid].TID = strings.Split(fields[25], ":")[1]
				proteins[pid].Domains = append(proteins[pid].Domains, domain)
			}
		}

	}

	return proteins
}

func GetNucleotidePosition(proteins map[string]*Protein, seqs map[string]string) map[string]string {

	nucleotideDomains := make(map[string]string)
	var nucID string
	usedGenes := make(map[string]bool)

	for _, p := range proteins {
		// Ignore genes that have been used already. Happens because of different splice variants.
		if usedGenes[p.GID] {
			continue
		}
		for _, domain := range p.Domains {
			// fmt.Println(nucID)
			nucID = p.TID + " " + strconv.Itoa(domain.Start) + ":" + strconv.Itoa(domain.Stop)
			nucleotideDomains[nucID] = seqs[p.TID][domain.Start*3 : domain.Stop*3]
		}
		usedGenes[p.GID] = true
	}

	return nucleotideDomains
}

// Takes a map with header/sequences and writes a sequential Phylip formated
// output file. Currently, illegal characters are replaced from the header,
// but the header is not truncated at 10 characters
func WriteFasta(fname string, seqs map[string]string) {

	f, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for id, s := range seqs {
		f.WriteString(">" + id + "\n")
		f.WriteString(s + "\n")
	}
	return
}

func main() {
	var (
		domtblout = flag.String("domtbl", "", "A domtblout formated file produced by HMMER3 hmmsearch. Must contain *peptide* information.")
		cdna      = flag.String("cdna", "", "ENSEMBL cdna fasta file. Must correspond the species that was searched with hmmsearch.")
		nucl      = flag.String("nucl", "", "Fasta formated output file storing the nucleotide sequences corresponding to domains detected and provided via the domtblout file.")
	)
	flag.Parse()
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lmicroseconds)

	fmt.Println("Reading domain input...")
	proteins := ReadDomtblout(*domtblout)
	fmt.Println("Reading cdna input...")
	seqs := ReadFastaAlignment(*cdna)
	nucldomainseqs := GetNucleotidePosition(proteins, seqs)
	WriteFasta(*nucl, nucldomainseqs)
}
