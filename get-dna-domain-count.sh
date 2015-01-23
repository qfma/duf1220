NOW=$(date +"%Y-%m-%d")
for dfam in $1/*.tblout;
	do
		SPECIES=$(echo $(basename $dfam) | cut -f1 -d. | cut -f5 -d-)
		domains=$(cat $dfam | grep -v "#" | wc -l);
		echo $SPECIES: $domains
	done;
