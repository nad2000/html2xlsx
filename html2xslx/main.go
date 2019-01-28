package main

import (
	"os"
	"path/filepath"
	"strings"

	html2xlsx "github.com/nad2000/html2xslx"
	log "github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) < 2 {
		log.Errorln("Missing input filename.")
		log.Infoln("USAGE: " + os.Args[0] + " <INPUT FILENAME> [<OUTPUT FILENAME>]")
		os.Exit(-1)
	}
	filename := os.Args[1]
	var outputFilename string
	if len(os.Args) > 2 {
		outputFilename = os.Args[2]
	} else {
		dir, fn := filepath.Split(filename)
		ext := filepath.Ext(fn)
		outputFilename = filepath.Join(dir,
			strings.TrimSuffix(fn, ext)+"_OUTPUT"+ext)
	}
	html2xlsx.Convert(filename, outputFilename)
}
