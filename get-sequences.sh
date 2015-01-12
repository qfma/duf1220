NOW=$(date +"%Y-%m-%d")

# Ensembl species
SPECIES=( "homo_sapiens"
		  "gorilla_gorilla"
		  "pongo_abelii"
		  "pan_troglodytes"
		  "nomascus_leucogenys"
		  "macaca_mulatta"
		  "tarsius_syrichta"
		  "microcebus_murinus"
		  "otolemur_garnettii"
		  "callithrix_jacchus"
		  "papio_anubis"
		  "papio_hamadryas"
		  "chlorocebus_sabaeus"
		  )

VERSION="78"
OUTFOLDER="seq$VERSION"

for s in "${SPECIES[@]}";
	do
		# Capitalize every first letter
		cs=$(echo $s | awk '{print toupper(substr($0,1,1)) tolower(substr($0,2)) }')
		echo "Downloading sequences for $cs "
		ftp_pep="ftp://ftp.ensembl.org/pub/release-$VERSION/fasta/$s/pep/"
		ftp_cdna="ftp://ftp.ensembl.org/pub/release-$VERSION/fasta/$s/cdna/"
		ftp_dna="ftp://ftp.ensembl.org/pub/release-$VERSION/fasta/$s/dna/"
		wget -r --accept "*.pep.all.fa.gz" --level 2 $ftp_pep
		wget -r --accept "*.cdna.all.fa.gz" --level 2 $ftp_cdna
		wget -r --accept "*.dna.toplevel.fa.gz" --level 2 $ftp_dna

	done;

#mkdir $OUTFOLDER $OUTFOLDER/pep $OUTFOLDER/cdna-all $OUTFOLDER/dna
