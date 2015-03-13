#!/usr/local/bin/python2.7
# -*- coding: utf-8 -*-
#
#       get-longest-isoform.py
#==============================================================================
import argparse
import sys
import os
#==============================================================================
#Command line options==========================================================
#==============================================================================
parser = argparse.ArgumentParser()
parser.add_argument("genomes", type=str, nargs="+",
                    help="A list of input genomes in fasta format\
                          or an input folder")
parser.add_argument("outdir", type=str, default="./",
                    help="The output directory for filtered genomes\
                          [default: ./]")
if len(sys.argv) == 1:
    parser.print_help()
    sys.exit(1)
args = parser.parse_args()
#==============================================================================


def read_fasta(source):
    SG = {}
    c = 0
    try:
        with open(source, "r") as file:
            for line in file.readlines():
                if line[0] == ">":
                    name = line[1:].rstrip().split()[0]
                    # Correct for missing names in Macaca_fascicularis
                    if not name.startswith("ENS"):
                        c += 1
                        name = "peptide_"+str(c)
                    header = line[1:].rstrip()
                    SG[name] = ["", header]
                else:
                    SG[name][0] += line.rstrip()
            return SG
    except IOError:
        print "File does not exit!"


def is_number(s):
    try:
        int(s)
        return True
    except ValueError:
        return False


def valid_chromosome(seq_header):

    chrom = seq_header.split()[2].split(":")[2]

    names = ["GL",        # Callithrix
             "ACF",       # Callithrix
             "AQ",        # Chlorocebus
             "KE",        # Chlorocebus
             "unplaced",  # Gorilla
             "cutchr",    # Gorilla
             "scaffold",      # Microcebus, Tarsius
             "GeneScaffold",  # Microcebus, Tarsius
             "ADF",           # Nomascus
             "AAQ",           # Otolemur
             "AAC",           # Pan
             "AHZ",           # Papio
             "JH",            # Papio
             "_random",       # Pongo
             "X",
             "Y",
             "A",
             "a",
             "b",
             "B",
             "MT"
             ]
    if is_number(chrom):
        return True

    for n in names:
        if chrom.startswith(n):
            return True
        # For Pongo
        elif chrom.endswith(n):
            return True
    return False


def find_isoforms(genome):
    isoforms = {}
    for seq in genome:
        if valid_chromosome(genome[seq][1]):
            gname = genome[seq][1].split()[3].split(":")[-1]
            if gname not in isoforms:
                isoforms[gname] = [(seq, len(genome[seq][0]))]
            else:
                isoforms[gname].append((seq, len(genome[seq][0])))
    return isoforms


def write_longest_isoform(isoforms, genome, fname):
    with open(fname, "w") as outfile:
        for seq in isoforms:
            if len(isoforms[seq]) > 1:
                m = max(zip(*isoforms[seq])[1])
                longest = [p[0] for p in isoforms[seq] if p[1] is m][0]
                header = ">" + genome[longest][1]+"\n"
                outfile.write(header)
                outfile.write(genome[longest][0]+"\n")
            else:
                longest = isoforms[seq][0][0]
                header = ">" + genome[longest][1]+"\n"
                outfile.write(header)
                outfile.write(genome[longest][0]+"\n")


def list_files(current_dir):
    file_list = []
    for path, subdirs, files in os.walk(current_dir):  # Walk directory tree
        for name in files:
            f = os.path.join(path, name)
            file_list.append(f)
    return file_list


#==============================================================================

def main():
    if not os.path.isdir(args.outdir):
        os.mkdir(args.outdir)

    if len(args.genomes) == 1:
        if os.path.isdir(args.genomes[0]):
            infiles = list_files(args.genomes[0])
        else:
            infiles = args.genomes
    else:
        infiles = args.genomes

    for genome in [f for f in infiles if f.endswith(".fa")]:
        print genome

        fname = ".".join(genome.split("/")[-1].split(".")[:-1]+["longest", "fa"])
        genome = read_fasta(genome)
        isoforms = find_isoforms(genome)
        write_longest_isoform(isoforms, genome, args.outdir+"/"+fname)

if __name__ == "__main__":
    main()
