package main

import (
	"flag"
	"github.com/jjmschofield/GoCrawl/internal/app/fetch"
	"log"
	"net/url"
)

func main(){
	crawlUrlRaw := flag.String("word", "foo", "a string")

	flag.Parse()

	crawlUrl, err := url.Parse(*crawlUrlRaw)

	if err != nil{
		log.Panic(err)
	}


	links, err := fetch.GetLinks(*crawlUrl)

	if err != nil{
		log.Panic(err)
	}

	log.Print(links)
}