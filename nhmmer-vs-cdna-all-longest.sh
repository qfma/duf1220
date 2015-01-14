NOW=$(date +"%Y-%m-%d")
OUTFOLDER="$NOW-nhmmer-vs-cdna-longest"
for i in ./sequences/cdna/*.cdna.all.longest.fa;
	do 
		SPECIES=$(echo $(basename $i) | cut -f1 -d.);
		nhmmer --dfamtblout "duf1220-vs-$SPECIES.dfamtblout" \
			   --tblout "duf1220-vs-$SPECIES.tblout" \
			   -E 1e-10 \
			   --cpu=8 \
			   ./$NOW-duf1220-all-ensembl-nucl/$NOW-duf1220-all-ensembl.hmm \
			   $i;
	done;
mkdir $OUTFOLDER
mv *.tblout $OUTFOLDER
mv *.dfamtblout $OUTFOLDER
