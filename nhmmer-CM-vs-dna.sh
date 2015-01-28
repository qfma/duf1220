NOW=$(date +"%Y-%m-%d")
NOW="2015-01-23"
OUTFOLDER="$NOW-nhmmer-CM-vs-dna"
for i in ./sequences/dna/*assembly.fa;
    do
        SPECIES=$(echo $(basename $i) | cut -f1 -d.);
        nhmmer --dfamtblout "CM-vs-$SPECIES.dfamtblout" \
               --tblout "CM-vs-$SPECIES.tblout" \
               -E 1e-10 \
               --cpu=22 \
               ./$NOW-CM-promoter/$NOW-human-NBPF-CM-promoter-900.hmm \
               $i;
    done;
mkdir $OUTFOLDER
mv *.tblout $OUTFOLDER
mv *.dfamtblout $OUTFOLDER
