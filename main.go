package main

import (
	"flag"
	"fmt"
	"github.com/CNuge/go-fasta/fasta"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"strings"
)

func min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

type SubstringsMap map[string]int64

func (s SubstringsMap) Intersect(a SubstringsMap) int64 {
	newS := SubstringsMap{}
	var res int64
	for k, v := range s {
		if v2, ok := a[k]; ok {
			m := min(v, v2)
			res += m * int64(len(k))
			newS[k] = m
		}
	}
	return res
}

func circleString(seq string, ind, size int) string {
	if ind+size < len(seq) {
		return seq[ind:ind+size]
	} else {
		return fmt.Sprintf("%s%s", seq[ind:], seq[:(size-(len(seq)-ind))])
	}
}

func getKey(key string) string {
	return strings.Split(key, "|")[1]
}

func makeSubstringMap(seq string) SubstringsMap {
	res := SubstringsMap{}
	for i := 1; i <= len(seq); i++ {
		for j := 0; j < len(seq); j++ {
			res[circleString(seq, j, i)]++
		}
	}
	return res
}

func makeResFile(seqs []string, res map[string]map[string]float64) {
	f := xlsx.NewFile()
	sheet, err := f.AddSheet("Sheet1")
	if err != nil {
		log.Panicln("add xlsx sheet err ", err)
	}
	row := sheet.AddRow()
	row.AddCell()
	for _, v := range seqs {
		cell := row.AddCell()
		cell.Value = v
	}
	for i, v := range seqs {
		row := sheet.AddRow()
		first := row.AddCell()
		first.Value = v
		for k := 0; k < i; k++ {
			row.AddCell()
		}
		//c := row.AddCell()
		//c.Value = "1"
		for _, v2 := range seqs[i:] {
			c := row.AddCell()
			c.Value = fmt.Sprintf("%f", res[v][v2])
		}
	}
	resFile, err := os.OpenFile("res.xlsx", os.O_CREATE | os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Panicln(err)
	}
	if err := f.Write(resFile); err != nil {
		log.Panicln(err)
	}
	resFile.Close()
}

func main() {
	var (
		name string
	)
	flag.StringVar(&name, "file", "", "")
	flag.Parse()
	if name == "" {
		log.Panicln("empty name")
	}
	var keys []string
	seqs := fasta.Read(name)
	subMaps := make(map[string]SubstringsMap, len(seqs))
	for _, s := range seqs[:5] {
		keys = append(keys, getKey(s.Name))
		subMaps[getKey(s.Name)] = makeSubstringMap(s.Sequence)
	}
	
	res := make(map[string]map[string]float64, len(seqs))
	for i := range keys {
		res[keys[i]] = map[string]float64{}
		for j := i; j < len(keys); j++ {
			size := (len(seqs[i].Sequence) + len(seqs[j].Sequence))/2
			size = size * size * (size + 1)/2
			res[keys[i]][keys[j]] = float64(subMaps[keys[i]].Intersect(subMaps[keys[j]]))/float64(size)
		}
	}
	makeResFile(keys, res)
}
