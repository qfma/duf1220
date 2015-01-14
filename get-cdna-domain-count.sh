NOW=$(date +"%Y-%m-%d")
echo "====CDNA-all-longest===="
echo "============================="
for dfam in $NOW-nhmmer-vs-cdna-longest/*.dfamtblout;
	do 
		SPECIES=$(echo $(basename $dfam) | cut -f1 -d. | cut -f3 -d-)
		domains=$(cat $dfam | grep -v "#" | wc -l);
		echo $SPECIES: $domains
	done;
