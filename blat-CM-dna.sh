NOW=$(date +"%Y-%m-%d")
NOW=2015-01-23
OUTFOLDER="$NOW-BLAT-CM-vs-dna"
for i in ./sequences/dna/*assembly.fa;
    do
        SPECIES=$(echo $(basename $i) | cut -f1 -d.);
        blat $i $NOW-CM-promoter/$NOW-human-NBPF4-CM-promoter-900.fa CM-vs-$SPECIES.psl -minScore=100 -minIdentity=50
    done;
mkdir $OUTFOLDER
mv *.psl $OUTFOLDER
