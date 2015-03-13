NOW=$(date +"%Y-%m-%d")
NOW="2015-02-20"
OUTFOLDER="$NOW-nhmmer-vs-dna"
for i in ./sequences/dna/dna-ensembl/*.fa;
	do 
		SPECIES=$(echo $(basename $i) | cut -f1 -d.);
		nhmmer --dfamtblout "duf1220-vs-$SPECIES.dfamtblout" \
			   --tblout "duf1220-vs-$SPECIES.tblout" \
			   -E 1e-10 \
			   --cpu=20 \
			   ./$NOW-duf1220-all-ensembl-nucl/$NOW-duf1220-all-ensembl.hmm \
			   $i;
	done;
mkdir $OUTFOLDER
mv *.tblout $OUTFOLDER
mv *.dfamtblout $OUTFOLDER
