NOW=$(date +"%Y-%m-%d")
OUTFOLDER="$NOW-blat-CM-vs-dna"
for i in ./sequences/dna/*.toplevel.fa;
    do
        SPECIES=$(echo $(basename $i) | cut -f1 -d.);
        blat $i 2015-01-12-human-NBPF4-CM-promoter-900.fa CM-vs-$SPECIES.psl -minScore=100 -minIdentity=50
    done;
mkdir $OUTFOLDER
mv *.psl $OUTFOLDER
