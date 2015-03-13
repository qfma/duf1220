#NOW=$(date +"%Y-%m-%d")
NOW="2015-02-20"
OUTDIR="$NOW-dna-orthologs"
INDIR="$NOW-nhmmer-vs-dna"

# Pass all the genes to enortho and find the Ensembl Compara orthologs for it
cat $INDIR/$NOW-dna-domain-counts.txt | grep "EN" | cut -f1 -d" " | 
while read GENE; 
	do 
		echo $GENE | ./enortho -species homo_sapiens,gorilla_gorilla,callithrix_jacchus,pongo_abelii,chlorocebus_sabaeus,macaca_mulatta,microcebus_murinus,nomascus_leucogenys,otolemur_garnettii,pan_troglodytes,papio_anubis,tarsius_syrichta > $GENE.ortho 
	done


mkdir $OUTDIR 
mv *.ortho $OUTDIR
