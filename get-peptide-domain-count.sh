echo "====PEP-all-longest===="
echo "======================="

NOW=$(date +"%Y-%m-%d")
for domtblout in $NOW-hmmsearch-vs-pep-all-longest/*.domtblout;
	do 
		SPECIES=$(echo $(basename $domtblout) | cut -f1 -d. | cut -f3 -d-)
		# Sum the tblout reg column, non overlapping regions.
		domains=$(cat $domtblout | grep -v "#" | wc -l);
		echo $SPECIES: $domains
	done;
