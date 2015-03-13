NOW=$(date +"%Y-%m-%d")
OUTFOLDER="$NOW-hmmsearch-vs-pep-all-longest"
for i in ./sequences/pep/*.pep.all.longest.fa;
	do
		SPECIES=$(echo $(basename $i) | cut -f1 -d.);
		hmmsearch --domtblout "duf1220-vs-$SPECIES.domtblout" --tblout "duf1220-vs-$SPECIES.tblout" --domE 1e-10 -E 1e-10 --cpu=8 ./pfam_hmm/PF06758_seed.hmm $i;
	done;
mkdir $OUTFOLDER
mv *.tblout $OUTFOLDER
mv *.domtblout $OUTFOLDER
