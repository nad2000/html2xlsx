package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nad2000/html2xlsx"
	log "github.com/sirupsen/logrus"
)

// VERSION of the module...
const VERSION = "0.2.1"

func usage() {
	log.Infoln("USAGE:\n\n" + os.Args[0] + " <INPUT FILENAME> [<OUTPUT FILENAME>]")
}

func main() {
	if len(os.Args) < 2 {
		log.Errorln("Missing input filename.")
		usage()
		os.Exit(1)
	}
	filename := os.Args[1]
	if strings.Contains(filename, "-version") || strings.Contains(filename, "-V") {
		fmt.Println(VERSION)
		return
	}
	if strings.Contains(filename, "-h") || strings.Contains(filename, "-?") {
		usage()
		return
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.WithField("Filename", filename).Errorln("File does not exist.")
		os.Exit(2)
	}
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
	log.WithField("Output Filename", outputFilename).Infoln("The converted files stored.")
}
