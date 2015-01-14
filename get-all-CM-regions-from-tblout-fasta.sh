NOW=$(date +"%Y-%m-%d")
INFOLDER="2015-01-12-nhmmer-CM-vs-dna"
OUTFOLDER="$NOW-all-CM-regions-plus-flank"
GENOMES="sequences/dna"

NONENSEMBL=("Deubetonia_madagascariensis", "Macaca_fascicularis",
            "Pan_paniscus", "Papio_hamadryas", "Saimiri_boliviensis")

./extract-CM-regions-from-tblout-fasta \
    -tblout=$INFOLDER/CM-vs-Saimiri_boliviensis.tblout\
    -fasta=CM-vs-Saimiri_boliviensis.fa \
    -genome=sequences/dna/Saimiri_boliviensis.SaiBol1.0.dna_rm.toplevel.fa \
    -flank=100000 \

./extract-CM-regions-from-tblout-fasta \
    -tblout=$INFOLDER/CM-vs-Deubetonia_madagascariensis.tblout\
    -fasta=CM-vs-Deubetonia_madagascariensis.fa \
    -genome=sequences/dna/Deubetonia_madagascariensis.dna.toplevel.fa \
    -flank=100000 \

./extract-CM-regions-from-tblout-fasta \
    -tblout=$INFOLDER/CM-vs-Macaca_fascicularis.tblout\
    -fasta=CM-vs-Macaca_fascicularis.fa \
    -genome=sequences/dna/Macaca_fascicularis.MacFas_5.0.76.dna.toplevel.fa \
    -flank=100000 \

./extract-CM-regions-from-tblout-fasta \
    -tblout=$INFOLDER/CM-vs-Pan_paniscus.tblout\
    -fasta=CM-vs-Pan_paniscus.fa \
    -genome=sequences/dna/Pan_paniscus.dna.toplevel.fa \
    -flank=100000 \

./extract-CM-regions-from-tblout-fasta \
    -tblout=$INFOLDER/CM-vs-Papio_hamadryas.tblout\
    -fasta=CM-vs-Papio_hamadryas.fa \
    -genome=sequences/dna/Papio_hamadryas.Pham_1.0.dna_rm.toplevel.fa \
    -flank=100000 \


