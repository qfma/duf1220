## DUF1220 analysis
This document describes the methods and every command used to analyze the distribution of the DUF1220([PF06758](http://pfam.xfam.org/family/duf1220)) protein domains across various primates.
The analysis and pipeline is based on some of the methods described by [O'Bleness et al.](http://g3journal.org/content/2/9/977). The main difference to th O'Bleness paper is,
that BLAT search has been replaced by Hidden Markov Models (HMM).

If you find some of these scripts useful, please cite: 

[Phylogenetic Analysis Supports a Link between DUF1220 Domain Number and Primate Brain Expansion](http://gbe.oxfordjournals.org/content/7/8/2083)  
Fabian Zimmer and Stephen H. Montgomery*


## Software requirements

The analysis was performed using Linux/OSX and has not been tested on Windows.
In order to run the scripts you need:

- A usable Python 2.7 installation
- A usable [Google Go](http://golang.org) installation
- An internet connection for querying the Ensembl REST API (doi:10.1093/bioinformatics/btu613)

## Method summary

In order to characterize the distribution of DUF1220 domains in NBPF genes, we obtained 12 primate genomes from Ensembl [REF] version 78 (*Homo sapiens*, *Pan troglodytes*, *Gorilla gorilla*, *Pongo abelii*, *Nomascus leucogenys*, *Macaca mulatta*, *Tarsius syrichta*, *Microcebus murinus*, *Otolemur garnettii*, *Callithrix jacchus*, *Papio anubis*, *Chlorocebus sabaeus*). We used HMMER [REF] to build a Hidden Markov Model (HMM) from the DUF1220([PF06758](http://pfam.xfam.org/family/duf1220)) seed alignment stored in the [PFAM database](http://pfam.xfam.org/family/duf1220#tabview=tab3). We extracted the longest isoforms for all proteomes and searched them using the DUF1220 HMM (hmmsearch, E-value < 1e-10), providing a first estimate of the DUF1220 in these species (Table REF, Peptide counts)

We used the positional information from the HMMER output to extract the corresponding cDNA regions and built a DUF1120 nucleotide profile HMM, allowing us to detect DUF1120 in nucleotide sequences across a broad phylogenetic range. We confirmed that the nucleotide HMM was effective in detecting DUF1120 domains by searching the longest isoforms for all cDNA sequences of all 12 Ensembl species (nhmmer, E-value < 1e-10), recovering a similar number of DUF1120 domains in comparison to the peptide analysis (Table REF, cDNA counts). The largest difference in domain numbers is found in the human dataset, where  a large set of untranslated pseudogenes contains DUF1220 domains.

In order to make sure that we include all DUF1120 domains present in a given gene, we searched the complete genomic DNA for all 12 species using the same nHMMER model. In comparison to the number of DUF1120 domains in peptide and cDNA sequences, the genomic search returns a much higher count (Table REF). The main reason for this increase is that many of the DUF1120 domains are located in regions with no feature annotation by Ensembl, and are therefore not included in peptide or cDNA sequences (Table REF). This is in contrast to those domains that are in exonic regions, which are similar to the peptide and cDNA counts.

Because we are primarily interested in those DUF1220 domains that are located in NBPF genes, we extracted the CM (Conserved Mammal) promoter region as defined by [O'Bleness et al.](http://g3journal.org/content/2/9/977) upstream of human genes NBPF4, NBPF6 and NBPF7 and aligned them using MAFFT (mafft --globalpair). This alignment was again used to build a nucleotide HMM for the CM promoter, which we used to search the genomic region of all DUF1220 containing genes, including 1000bp up- and downstream for significant CM promoter hits (nhmmer, E-value < 1e-10), shown in Table REF, CDS with CM promoter. We then subdivided exonic domains in those that are within genes with a CM promoter and those that are not.

## Method issues, shortcomings etc.

The good:

- HMM very sensitive
- HMM built from alignments, better suited for detection of distant DUF1120 domains
- Clear separation between DUF1120 domains in CM genes and those domains in genes without CM promoter
- Pipeline will be open source, easy to repeat

The bad:
- Counts vary depending on what is included or excluded (eg. pseudogenes)
- HMM detection and E-value of 10e-10 is strict and may underestimate the number of DUF1120 domains
- Function of the DUF1120 domain unknown

The ugly:

- Genome quality for primates extremely variable. Entire genes may disappear depending on assembly version.
- Some genomes only available as Scaffolds, positional information is unavailable.

## Command history and data used

### Primate sequences used

**[Ensembl](http://www.ensembl.org/index.html)** version 78:

1. Homo sapiens (Humans)
2. Pan troglodytes (Common Chimpanzees)
3. Gorilla gorilla (Gorilla)
4. Pongo abelii (Orangutan)
5. Nomascus leucogenys (Gibbon)
6. Macaca mulatta (Rhesus Macaque)
7. Tarsius syrichta (Tarsier)
8. Microcebus murinus (Mouse lemur)
9. Otolemur garnettii (Bushbaby)
10. Callithrix jacchus (Marmoset)
11. Papio anubis (Olive baboon)
12. Chlorocebus sabaeus (Green Vervet monkey)

### DUF1220 proteome annotation

The DUF1220 ([PF06758](http://pfam.xfam.org/family/duf1220)) seed alignment was downloaded from the PFAM database [here](http://pfam.xfam.org/family/duf1220#tabview=tab3).
The the profile Hidden Markov Model was built using HMMER 3.1b1.

``` bash
hmmbuild -n duf1220_seed --cpu 8  PF06758_seed.hmm PF06758_seed.txt
```

The PF06758_seed.hmm was used to search the proteomes of all 12 Ensembl primate proteomes using ```hmmsearch``` with the E-value threshold of 1e-10.
Only the longest isoform was used for all proteomes.

``` bash
# Get longest isoforms
./get-longest-isoform.py sequences/pep sequences/pep
# Run hmmsearch against all proteomes
./hmmsearch-vs-pep-all-longest.sh
```

``` bash
# Count the number of annotated domains
./get-peptide-domain-count.sh
Callithrix_jacchus: 12
Chlorocebus_sabaeus: 22
Gorilla_gorilla: 38
Homo_sapiens: 246
Macaca_mulatta: 21
Microcebus_murinus: 0
Nomascus_leucogenys: 5
Otolemur_garnettii: 0
Pan_troglodytes: 37
Papio_anubis: 27
Pongo_abelii: 28
Tarsius_syrichta: 2
```

### Building nucleotide alignments of DUF1220 from protein domain annotation

We then extracted the positional information of the protein domains from the domtblout files generated by ```HMMER```. The positions were used to obtain the corresponding nucleotide sequence from the cDNA of the protein translations and
subsequently built a nucleotide HMM.

``` bash
./get-all-positions.sh
```

The cDNA sequences are aligned using ```mafft``` and this alignment is then used to create a HMM with ```hmmbuild```.

### Verification of domains using cDNA

We now verified the DUF1220 counts using the cDNA of the 12 Ensembl species. The E-value cutoff is again 1e-10.

``` bash
# Make longest isoforms
./get-longest-isoform.py sequences/cdna sequences/cdna
# Run nhmmer
./nhmmer-vs-cdna-all-longest.sh
# Count the results
./get-cdna-domain-count.sh
Callithrix_jacchus: 15
Chlorocebus_sabaeus: 24
Gorilla_gorilla: 41
Homo_sapiens: 292
Macaca_mulatta: 22
Microcebus_murinus: 2
Nomascus_leucogenys: 6
Otolemur_garnettii: 2
Pan_troglodytes: 41
Papio_anubis: 29
Pongo_abelii: 33
Tarsius_syrichta: 4
```

### Pseudogenes in humans

The increase in number of domains from the peptide counts in human can be explained by the presence on pseudogenes, that are not translated:

```text
ENST00000619932,ENSG00000272150,6,transcribed_unprocessed_pseudogene
ENST00000617931,ENSG00000268043,11,protein_coding
ENST00000577412,ENSG00000266338,6,protein_coding
ENST00000609741,ENSG00000273136,10,protein_coding
ENST00000453025,ENSG00000227001,3,unprocessed_pseudogene
ENST00000318220,ENSG00000142794,5,protein_coding
ENST00000613157,ENSG00000196427,4,protein_coding
ENST00000615421,ENSG00000269713,9,protein_coding
ENST00000612480,ENSG00000275131,1,unprocessed_pseudogene
ENST00000583271,ENSG00000270231,8,unprocessed_pseudogene
ENST00000369373,ENSG00000162825,67,protein_coding
ENST00000449715,ENSG00000203825,1,unprocessed_pseudogene
ENST00000621744,ENSG00000271383,45,protein_coding
ENST00000612520,ENSG00000213240,5,protein_coding
ENST00000444082,ENSG00000227242,5,unprocessed_pseudogene
ENST00000425093,ENSG00000231382,2,unprocessed_pseudogene
ENST00000294652,ENSG00000186086,4,protein_coding
ENST00000357046,ENSG00000243967,2,transcribed_unprocessed_pseudogene
ENST00000615281,ENSG00000263956,7,protein_coding
ENST00000401087,ENSG00000179571,6,unprocessed_pseudogene
ENST00000583866,ENSG00000271425,42,protein_coding
ENST00000445758,ENSG00000215864,2,unprocessed_pseudogene
ENST00000430580,ENSG00000219481,7,protein_coding
ENST00000619423,ENSG00000270629,32,protein_coding
ENST00000590707,ENSG00000205449,2,transcribed_unprocessed_pseudogene
Protein_coding_domains,254
Pseudogene_domains,38
```

Statistics were generated using:

```bash
./get-cdna-statistics.py 2015-02-20-nhmmer-vs-cdna-longest
```


### DUF1120 in genomic DNA regions

The first two analysis, in peptide and cDNA sequences, always used the longest isoform for a given gene and therefore may not capture all DUF1120 domains present, because some shorter, non-overlapping isoforms could contain DUF1120 domains. These domains would then not be counted using only the longest isoform. We therefore searched the entire genomic DNA for all 12 species and then determined if the domains are within an exonic region, or without any feature annotation.

```bash
# Uses the same nucleotide HMM as for the cDNA analysis
# Runs against complete genomes
./nhmmer-vs-dna.sh

```
Next we query the Ensembl REST API to return us information about the features present for a given DNA region. We use the information provided in the *.tblout files produced by nHMMER.

```bash
./get-all-ensembl-features-for-domains.sh > DNA-domain-counts.txt
```

As a result of the above analysis we obtain a DUF1120 count on a per gene basis:

```text
# Marmoset example
# Per gene estimates
ENSCJAG00000018042 2
ENSCJAG00000018105 2
ENSCJAG00000005113 2
ENSCJAG00000019954 1
ENSCJAG00000015447 2
ENSCJAG00000020055 2
ENSCJAG00000020356 4
#Species Domains Infeature WithoutFeature
Callithrix_jacchus 75 15 60
```

### Counting only DUF1220 domains that are close to the CM promoter

All of the DUF1220 domain counts provided by the three analyses above (peptide, cDNA, DNA) provide data for all genes containing DUF1220 domains. However, we are primarily interested in those DUF1220 domains that occur in NBPF genes, which are characterized by [O'Bleness et al.](http://g3journal.org/content/2/9/977) as those genes that contain a CM (conserved mammal promoter). The CM promoter is around 900bp long and is positioned either upstream of the gene or in the fourth intron. We extracted the 900 bp upstream region of the human genes NBPF4, NBPF6 and NBPF7, which where then used to built a nucleotide HMM

```bash
# Get human CM promoter regions for NBPF4, NBPF6 and NBPF7
./get-CM-promoter-regions.sh
```

The promoter regions were aligned using ```mafft --globalpair``` and another nucleotide HMM was built using ```hmmbuild```.

In the next step we tried to identify only those genes that contained a CM promoter either upstream or in an intron. We therefore obtained the genomic regions for all CDS genes with DUF1120 domains including 1000bp up- and downstream. These regions where then searched using the CM promoter HMM.

```bash
# Find CM genes
./find-cm-genes.sh
```

The resulting gene list contains shows which genes contain a CM promoter.
```text
ENSG00000186086 2 # Two promoters
ENSG00000196427 1 # One promoter
ENSG00000203825 0 # None
ENSG00000205449 1 # One promoter
ENSG00000213240 1 # One promoter
```

We used this information to compile a fourth DUF1120 domain count, including only exonic DUF1120 domains in genes that have a valid CM promoter.

## Orthology analyis

In the last analysis we obtained a list of homologs for each DUF1120 containing gene in all 12 species from Ensembl COMPARA in order to determine if the increase in DUF1120 domains was primarily driven by an increase in domain number of some genes.

```bash
./get-dna-orthologs.sh
```

Here an example output for human gene ENSG00000162825
```text
ENSG00000270629 homo_sapiens within_species_paralog
ENSG00000142794 homo_sapiens within_species_paralog
ENSG00000273136 homo_sapiens within_species_paralog
ENSG00000271425 homo_sapiens within_species_paralog
ENSG00000268043 homo_sapiens within_species_paralog
ENSG00000266338 homo_sapiens within_species_paralog
ENSG00000219481 homo_sapiens within_species_paralog
ENSG00000271383 homo_sapiens within_species_paralog
ENSG00000269713 homo_sapiens within_species_paralog
ENSG00000196427 homo_sapiens within_species_paralog
ENSG00000186086 homo_sapiens within_species_paralog
ENSG00000271254 homo_sapiens within_species_paralog
ENSG00000263956 homo_sapiens within_species_paralog
ENSGGOG00000005749 gorilla_gorilla ortholog_one2one
ENSCJAG00000020356 callithrix_jacchus ortholog_one2one
ENSPPYG00000010655 pongo_abelii ortholog_one2one
ENSCSAG00000018389 chlorocebus_sabaeus ortholog_one2one
ENSMMUG00000031757 macaca_mulatta ortholog_one2one
ENSOGAG00000028082 otolemur_garnettii ortholog_one2many
ENSPTRG00000023855 pan_troglodytes ortholog_one2many
ENSPTRG00000039633 pan_troglodytes ortholog_one2many
ENSPANG00000018417 papio_anubis ortholog_one2one
ENSTSYG00000007889 tarsius_syrichta ortholog_one2many
```

### License

The MIT License (MIT)

Copyright (c) 2015 Fabian Zimmer

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

