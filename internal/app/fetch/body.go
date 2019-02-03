package fetch

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type PageBodyFetcher func(url.URL) (bodyReader io.ReadCloser, err error)

func Body(targetUrl url.URL) (bodyReader io.ReadCloser, err error) {
	response, err := http.Get(targetUrl.String())

	if err != nil {
		log.Printf("get for %s failed %s", targetUrl.String(), err)
		return nil, err
	}

	if !isAllowedContentType(response) {
		err := fmt.Errorf("non-htmldocs Content-Type from %s", targetUrl.String())
		log.Print(err)
		return nil, err
	}

	// TODO - Other validations eg UTF-8 check?

	return response.Body, nil
}

func isAllowedContentType(response *http.Response) bool {
	contentType := response.Header.Get("Content-Type")

	if !strings.Contains(contentType, "text/html") {
		return false
	}

	return true
}
