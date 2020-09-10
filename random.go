package main

import (
	"fmt"
	"github.com/CNuge/go-fasta/fasta"
	"math"
	"math/rand"
)

var alphabet = "ACDEFGHIKLMNPQRSTVWY"

func RandomSeqs() fasta.Fasta {
	res := fasta.Fasta{}
	for i := 0; i < 10; i++ {
		var s string
		for j := 0; j < 600; j++ {
			s = s + string(alphabet[int(math.Floor(rand.Float64()*float64(len(alphabet))))])
		}
		res = append(res, fasta.Seq{
			Name:     fmt.Sprintf("sp|RAND%d|RAND%d", i, i),
			Sequence: s,
		})
	}
	return res
}
