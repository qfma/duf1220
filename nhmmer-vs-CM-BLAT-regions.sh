NOW=$(date +"%Y-%m-%d")
OUTFOLDER="$NOW-nhmmer-vs-BLAT-CM-regions-plus-flank"
INFOLDER="$NOW-all-BLAT-CM-regions-plus-flank"
for i in $INFOLDER/*.fa;
	do
		SPECIES=$(echo $(basename $i) | cut -f1 -d. | cut -f3 -d"-");
		nhmmer --tblout "duf1220-vs-CM-region-$SPECIES.tblout" \
			   -E 1e-10 \
			   --cpu=8 \
			   ./2015-01-14-duf1220-all-ensembl-nucl/2015-01-14-duf1220-all-ensembl.hmm \
			   $i;
	done;
mkdir $OUTFOLDER
mv *.tblout $OUTFOLDER
mv *.dfamtblout $OUTFOLDER
