NOW=$(date +"%Y-%m-%d")

# Ensembl species
SPECIES=( "homo_sapiens"
		  "pan_troglodytes"
		  "macaca_mulatta"
		  )

OUTFOLDER="ensembl-48to78"
mkdir $OUTFOLDER

for VERSION in {48..78};
	do
		for s in "${SPECIES[@]}";
			do
				# Capitalize every first letter
				cs=$(echo $s | awk '{print toupper(substr($0,1,1)) tolower(substr($0,2)) }')
				echo "Downloading peptides for $cs version $VERSION"
				ftp_pep="ftp://ftp.ensembl.org/pub/release-$VERSION/fasta/$s/pep/"
				wget -r --accept "*.pep.all.fa.gz" --level 2 $ftp_pep
				mv ftp.ensembl.org/pub/release-$VERSION/fasta/$s/pep/*.pep.all.fa.gz $OUTFOLDER/$s-v$VERSION.pep.all.fa.gz
			done;
	done;

