NOW=$(date +"%Y-%m-%d")
OUTFOLDER="$NOW-CM-promoter"
CMFILENAME="$NOW-human-NBPF-CM-promoter-900"
mkdir $OUTFOLDER

./get-human-NBPF-promoter-regions -fasta=$CMFILENAME.fa -cm -upstream=900
mafft --globalpair --thread 8 $CMFILENAME.fa > $CMFILENAME.aln
hmmbuild -n NBPF_CM --dna --cpu 8 $CMFILENAME.hmm $CMFILENAME.aln

mv *.fa $OUTFOLDER
mv *.aln $OUTFOLDER
mv *.hmm $OUTFOLDER

