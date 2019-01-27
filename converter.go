package html2xlsx

import (
	"archive/zip"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Convert opens zip file and goes throug all the .xls (with HTML content)
// convrerts them into .xlsx and stores into a new output zipfile
func Convert(fileName, outputFilename string) {
	// Open a zip archive for reading.
	r, err := zip.OpenReader(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "__MACOSX") || filepath.Ext(f.Name) != ".xls" {
			continue
		}
		log.WithFields(log.Fields{"Filename": f.Name}).Infoln("Converting...")
		// rc, err := f.Open()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// _, err = io.CopyN(os.Stdout, rc, 68)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// rc.Close()
		// fmt.Println()
	}
}
