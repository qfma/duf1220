echo "====PEP-all-longest===="
echo "======================="

NOW=$(date +"%Y-%m-%d")
NOW="2015-02-20"
for dfam in $NOW-hmmsearch-vs-pep-all-longest-ensembl-range/*.domtblout;
	do 
		domains=$(cat $dfam | grep -v "#" | wc -l);
		echo $dfam: $domains
	done;
