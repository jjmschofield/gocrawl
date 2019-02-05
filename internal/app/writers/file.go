package writers

import (
	"encoding/json"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
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
	err := os.Mkdir(w.FilePath, os.ModePerm)
	pageFile, err := os.Create(path.Join(w.FilePath, "pages.jsonl"))

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
