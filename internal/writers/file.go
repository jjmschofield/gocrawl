package writers

import (
	"encoding/json"
	"github.com/jjmschofield/gocrawl/internal/pages"
	"log"
	"os"
	"path"
)

type FileWriter struct {
	FilePath string
	pageFile *os.File
	encoder  *json.Encoder
}

func (w *FileWriter) Start(in chan pages.Page) {
	workingDir, err := os.Getwd()

	if err != nil {
		log.Panicf("Can't get working directory!, %v", err)
	}

	filePath := path.Join(workingDir, w.FilePath)

	err = os.MkdirAll (filePath, os.ModePerm)

	if err != nil {
		log.Panicf("Can't create directory!, %v", err)
	}


	pageFile, err := os.Create(path.Join(filePath, "pages.jsonl"))

	if err != nil {
		log.Panicf("Can't open file to write results to!, %v", err)
	}

	w.encoder = json.NewEncoder(pageFile)

	for page := range in {
		w.write(page)
	}
}

func (w *FileWriter) write(page pages.Page) {
	err := w.encoder.Encode(page)

	if err != nil {
		log.Panicf("Can't write entry! %v", err)
	}
}
