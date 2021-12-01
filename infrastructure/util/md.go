package util

import (
	"app/infrastructure/log"
	"bytes"
	"github.com/yuin/goldmark"
)

func MdToHtml(md []byte) string {
	var buf bytes.Buffer
	if err := goldmark.Convert(md, &buf); err != nil {
		log.Error(err.Error())
	}
	return buf.String()
}
