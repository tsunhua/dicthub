package cc

import (
	"log"

	"github.com/liuzl/gocc"
)

func init() {
	dir := "./static/cc"
	gocc.Dir = &dir
}

func Convert2All(in string) []string {
	all := make(map[string]bool)
	all[in] = true
	
	{
		t2s, err := gocc.New("t2s")
		if err != nil {
			log.Fatal(err)
		}
		out, err := t2s.Convert(in)
		if err != nil {
			log.Fatal(err)
		}
		all[out] = true
	}

	{
		s2t, err := gocc.New("s2t")
		if err != nil {
			log.Fatal(err)
		}
		out, err := s2t.Convert(in)
		if err != nil {
			log.Fatal(err)
		}
		all[out] = true
	}

	rz := make([]string, len(all))
	for k := range all {
		rz = append(rz, k)
	}
	return rz
}
