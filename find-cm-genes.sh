NOW=$(date +"%Y-%m-%d")
OUTDIR="$NOW-genes-1000bp-updown"


# Extract the gene names from the tblout files and return the sequence using
# ENSEMBL
for tblout in $NOW-nhmmer-vs-cdna-longest/*.tblout; do 
		SPECIES=$(echo $(basename $tblout) | cut -f1 -d.);
		cat $tblout | \
		grep -v "#" | \
		awk '{print $18}' | \
		cut -f2 -d":" | \
		sort | uniq | ./gseq-flank
done

# Move the results into a folder
mkdir $OUTDIR
mv *.1000bp.updown.fa $OUTDIR

# Now run nhmmer for each gene against the CM promoter model and determine
# if a CM hit is present
for i in $OUTDIR/*.1000bp.updown.fa; do
	GENE=$(echo $(basename $i) | cut -f1 -d.);
	nhmmer --dfamtblout "$GENE.dfamtblout" \
		   --tblout "$GENE.tblout" \
		   -E 1e-10 \
		   --cpu=22 \
		   ./$NOW-CM-promoter/$NOW-human-NBPF-CM-promoter-900.hmm \
		   $i;
done

# Move the results into a folder
mv *.dfamtblout $OUTDIR
mv *.tblout $OUTDIR

# Print the CM counts for each gene
for i in $OUTDIR/*.tblout; do
	GENE=$(echo $(basename $i) | cut -f1 -d.);
	GENECOUNT=$(cat $i| grep -v "#" | wc -l)
	echo $GENE $GENECOUNT
done