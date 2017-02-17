#!/usr/local/bin/python2.7
# -*- coding: utf-8 -*-
#
#       get-cdna-statistics.py
#==============================================================================
import argparse
import sys
import os
from collections import defaultdict
#==============================================================================
#Command line options==========================================================
#==============================================================================
parser = argparse.ArgumentParser()
parser.add_argument("dfamtblout", type=str, nargs="+",
                    help="A file orfolder containing dfamtblout file")
if len(sys.argv) == 1:
    parser.print_help()
    sys.exit(1)
args = parser.parse_args()
#==============================================================================


def read_dfamtblout(source):
    hits = {}
    try:
        with open(source, "r") as infile:
            for line in infile:
                # Ignore comments
                if not line.startswith("#"):
                    fields = line.split()
                    target = fields[0]
                    if target in hits:
                        hits[target]["domains"] += 1
                    else:
                        gene = {"tid": target,
                                "biotype": fields[17].split(":")[1],
                                "domains": 1,
                                "gid": fields[16].split(":")[1]}
                        hits[target] = gene
            return hits
    except IOError:
        print "File does not exit!"


def list_files(current_dir):
    file_list = []
    for path, subdirs, files in os.walk(current_dir):  # Walk directory tree
        for name in files:
            f = os.path.join(path, name)
            file_list.append(f)
    return file_list

def sum_domains(dfamtblout):
    protein_coding_domains = 0
    pseudogene_domains = 0
    for i in dfamtblout:
        if dfamtblout[i]["biotype"] == "protein_coding":
            protein_coding_domains += dfamtblout[i]["domains"]
        else:
            pseudogene_domains += dfamtblout[i]["domains"]
    return protein_coding_domains, pseudogene_domains

#==============================================================================

def main():

    if len(args.dfamtblout) == 1:
        if os.path.isdir(args.dfamtblout[0]):
            infiles = list_files(args.dfamtblout[0])
        else:
            infiles = args.dfamtblout
    else:
        infiles = args.dfamtblout

    for dfamtblout in [f for f in infiles if f.endswith(".dfamtblout")]:
        print dfamtblout
        fname = ".".join(dfamtblout.split("/")[-1].split(".")[:-1]+["stats", "dfamtblout"])
        with open(fname, "w") as outfile:
            dfamtblout = read_dfamtblout(dfamtblout)
            for i in dfamtblout:
                outfile.write(",".join([i, dfamtblout[i]["gid"], str(dfamtblout[i]["domains"]), dfamtblout[i]["biotype"]])+"\n")
            protein_coding_domains, pseudogene_domains = sum_domains(dfamtblout)
            outfile.write("Protein_coding_domains," + str(protein_coding_domains)+"\n")
            outfile.write("Pseudogene_domains,"+str(pseudogene_domains)+"\n")


if __name__ == "__main__":
    main()
