NOW=$(date +"%Y-%m-%d")
INFOLDER="2015-01-21-BLAT-CM-vs-dna"
OUTFOLDER="$NOW-all-BLAT-CM-regions-plus-flank"
for psl in $INFOLDER/*.psl;
    do
        SPECIES=$(echo $(basename $psl) | cut -f1 -d.);
        echo "Extracting CM regions for $SPECIES..."
        sleep 5
        ./extract-CM-regions-from-psl-ensembl -psl=$psl -fasta=$SPECIES.fa -flank=100000
    done;

mkdir $OUTFOLDER
mv *.fa $OUTFOLDER

# echo "Merging output files..."
# cat $OUTFOLDER/*-domains.cdna.fa > $OUTFOLDER/$NOW-duf1220-all-ensembl.cdna.fa
# echo "Aligning merged output..."
# mafft --auto $OUTFOLDER/$NOW-duf1220-all-ensembl.cdna.fa > $OUTFOLDER/$NOW-duf1220-all-ensembl.cdna.aln
# echo "Making nucleotide HMM..."
# hmmbuild -n duf1220_nucl --dna --cpu 8 $OUTFOLDER/$NOW-duf1220-all-ensembl.hmm $OUTFOLDER/$NOW-duf1220-all-ensembl.cdna.aln
