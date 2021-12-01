package search

import "strings"

var tones = []string{}

func init() {
	numbers := []string{
		"¹", "²", "³", "⁴", "⁵", "⁶", "⁷", "⁸", "⁹",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	}
	aplphas := [][]string{
		{"a", "á", "à", "â", "ǎ", "ā", "a̍", "a̋"},
		{"e", "é", "è", "ê", "ě", "ē", "e̍", "e̋"},
		{"i", "í", "ì", "î", "ǐ", "ī", "i̍", "i̋"},
		{"i", "í", "ì", "î", "ǐ", "ī", "i̍", "i̋"},
		{"o", "ó", "ò", "ô", "ǒ", "ō", "o̍", "ő"},
		{"u", "ú", "ù", "û", "ǔ", "ū", "u̍", "ű"},
		{"m", "ḿ", "m̀", "m̂", "m̌", "m̄", "m̍", "m̋"},
		{"n", "ń", "ǹ", "n̂", "ň", "n̄", "n̍", "n̋"},
	}
	for _, n := range numbers {
		tones = append(tones, n, "")
	}
	for _, arr := range aplphas {
		for i, a := range arr {
			if i == 0 {
				continue
			}
			tones = append(tones, a, arr[0])
		}
	}
}

func clearTone(str string) string {
	return strings.NewReplacer(tones...).Replace(str)
}
