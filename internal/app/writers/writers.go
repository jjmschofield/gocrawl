package writers

import (
	"encoding/csv"
	"encoding/json"
	"github.com/jjmschofield/GoCrawl/internal/app/links"
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

	linkFile, err := os.Create(path.Join(w.FilePath, "links.jsonl"))
	defer linkFile.Close()
	linkEncoder := json.NewEncoder(linkFile)

	edgeFile, err := os.Create(path.Join(w.FilePath, "edges.csv"))
	defer edgeFile.Close()
	edgeEncoder := csv.NewWriter(edgeFile)

	if err != nil {
		log.Panicf("Can't open file to write results to!, %v", err)
	}

	for page := range in {
		writePage(page, pageEncoder)
		writeAllLinks(page, linkEncoder)
		writeEdges(page, edgeEncoder)
	}
}

func writePage(page pages.Page, pageEncoder *json.Encoder) {
	err := pageEncoder.Encode(page)
	if err != nil {
		log.Panicf("Can't write entry! %v", err)
	}
}

func writeAllLinks(page pages.Page, linkEncoder *json.Encoder) {
	toWrite := append(page.OutLinks.Internal, page.OutLinks.InternalFile...)
	toWrite = append(toWrite, page.OutLinks.External...)
	toWrite = append(toWrite, page.OutLinks.ExternalFile...)
	toWrite = append(toWrite, page.OutLinks.Tel...)
	toWrite = append(toWrite, page.OutLinks.Mailto...)
	toWrite = append(toWrite, page.OutLinks.Unknown...)
	writeLinks(toWrite, linkEncoder)
}

func writeLinks(l []links.Link, encoder *json.Encoder) {
	for _, link := range l {
		err := encoder.Encode(link)

		if err != nil {
			log.Panicf("Can't write link! %v", err)
		}
	}
}

func writeEdges(page pages.Page, edgeEncoder *csv.Writer) {
	for _, outPage := range page.OutPages.Internal {
		err := edgeEncoder.Write([]string{page.Id, outPage.Id})

		if err != nil {
			log.Panicf("Can't write edge! %v", err)
		}
	}
}
