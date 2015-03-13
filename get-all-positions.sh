NOW=$(date +"%Y-%m-%d")
INFOLDER="$NOW-hmmsearch-vs-pep-all-longest"
OUTFOLDER="$NOW-duf1220-all-ensembl-nucl"
for CDNA in sequences/cdna/*.cdna.all.fa; do 
		SPECIES=$(echo $(basename $CDNA) | cut -f1 -d.);
		ASSEMBLY=$(echo $(basename $CDNA) | cut -f2 -d.);
		RELEASE=$(echo $(basename $CDNA) | cut -f3 -d.);
		echo "Detecting nucleotide domain positions in $SPECIES..."
		./get-nucleotide-domain-positions -domtbl $INFOLDER/duf1220-vs-$SPECIES.domtblout -cdna $CDNA -nucl duf1220-$SPECIES-domains.cdna.fa
done

mkdir $OUTFOLDER
mv *-domains.cdna.fa $OUTFOLDER

echo "Merging output files..."
cat $OUTFOLDER/*-domains.cdna.fa > $OUTFOLDER/$NOW-duf1220-all-ensembl.cdna.fa
echo "Aligning merged output..."
mafft --auto $OUTFOLDER/$NOW-duf1220-all-ensembl.cdna.fa > $OUTFOLDER/$NOW-duf1220-all-ensembl.cdna.aln
echo "Making nucleotide HMM..."
hmmbuild -n duf1220_nucl --dna --cpu 8 $OUTFOLDER/$NOW-duf1220-all-ensembl.hmm $OUTFOLDER/$NOW-duf1220-all-ensembl.cdna.aln
