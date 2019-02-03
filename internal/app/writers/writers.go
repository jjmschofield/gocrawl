package writers

import (
	"encoding/json"
	"github.com/jjmschofield/GoCrawl/internal/app/pages"
	"log"
	"os"
	"sync"
)

type Writer func(in chan pages.Page, wg *sync.WaitGroup)

func StdoutWriter(in chan pages.Page, wg *sync.WaitGroup){
	defer wg.Done()
	wg.Add(1)

	pageEncoder := json.NewEncoder(os.Stdout)


	for page := range in {
		err := pageEncoder.Encode(page)

		if err != nil {
			log.Panicln("Can't write entry!")
		}
	}
}

type FileWriter struct {
	FilePath string
}

func (w *FileWriter) Write (in chan pages.Page, wg *sync.WaitGroup){
	defer wg.Done()
	wg.Add(1)

	jsonlFile, err := os.Create(w.FilePath)
	defer jsonlFile.Close()

	pageEncoder := json.NewEncoder(jsonlFile)

	if err != nil {
		log.Panicln("Can't open file to write results to!")
	}

	for page := range in {
		err := pageEncoder.Encode(page)

		if err != nil {
			log.Panicln("Can't write entry!")
		}
	}
}
