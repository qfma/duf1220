NOW=$(date +"%Y-%m-%d")
INFOLDER="$NOW-BLAT-CM-vs-dna"
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

