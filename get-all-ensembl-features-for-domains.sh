NOW=$(date +"%Y-%m-%d")
NOW="2015-02-20"
echo "Species Domains Infeature WithoutFeature"
for i in $1/*.tblout; do
		./extract-ensembl-feature-for-domain -tblout $i
done


