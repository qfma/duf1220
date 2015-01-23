echo "====PEP-all-longest===="
echo "======================="

NOW=$(date +"%Y-%m-%d")
for dfam in $NOW-hmmsearch-vs-pep-all-longest/*.domtblout;
	do 
		SPECIES=$(echo $(basename $dfam) | cut -f1 -d. | cut -f3 -d-)
		domains=$(cat $dfam | grep -v "#" | wc -l);
		echo $SPECIES: $domains
	done;
