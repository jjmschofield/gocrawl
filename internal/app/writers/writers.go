package writers

import (
	"encoding/json"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"log"
	"os"
	"path"
	"sync"
)

type Writer func(in chan pages.Page, wg *sync.WaitGroup)

type FileWriter struct {
	FilePath string
}

func (w *FileWriter) Write(in chan pages.Page, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)

	err := os.Mkdir(w.FilePath, os.ModePerm)

	pageFile, err := os.Create(path.Join(w.FilePath, "pages.jsonl"))
	defer pageFile.Close()
	pageEncoder := json.NewEncoder(pageFile)

	if err != nil {
		log.Panicf("Can't open file to write results to!, %v", err)
	}

	for page := range in {
		writePage(page, pageEncoder)
	}
}

func writePage(page pages.Page, pageEncoder *json.Encoder) {
	err := pageEncoder.Encode(page)
	if err != nil {
		log.Panicf("Can't write entry! %v", err)
	}
}
