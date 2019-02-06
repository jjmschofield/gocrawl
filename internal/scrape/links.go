package scrape

import (
	"fmt"
	"github.com/jjmschofield/gocrawl/internal/links"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/url"
)

func fetchLinks(target url.URL) (linkMap map[string]links.Link, err error) {
	bodyReader, err := fetchBody(target)

	if err != nil {
		return nil, err
	}

	return readLinks(target, bodyReader), nil
}

func readLinks(srcUrl url.URL, bodyReader io.ReadCloser) (outLinks map[string]links.Link) {
	defer bodyReader.Close()

	outLinks = make(map[string]links.Link)

	tokens := html.NewTokenizer(bodyReader)

	for {
		tokenType := tokens.Next()
		token := tokens.Token()

		isAnchor, eof := tokenIsTagType("a", tokenType, token)

		if eof {
			break
		} else if isAnchor {
			href, err := getAttrValFromToken("href", token)
			if err == nil {
				link, err := links.FromHref(srcUrl, href)
				if err == nil {
					outLinks[link.Id] = link
				}
			}
		}
	}

	return outLinks
}

func tokenIsTagType(tagType string, tokenType html.TokenType, token html.Token) (isAnchor bool, eof bool) {
	switch {
	case tokenType == html.ErrorToken:
		log.Printf("invalid token %s", token.Data)
		return false, true
	case tokenType == html.EndTagToken && token.Data == "html":
		return false, true
	case tokenType == html.StartTagToken && token.Data == tagType:
		return true, false
	default:
		return false, false
	}
}

func getAttrValFromToken(attrKey string, token html.Token) (val string, err error) {
	for _, attr := range token.Attr {
		if attr.Key == attrKey {
			return attr.Val, nil
		}
	}
	return "", fmt.Errorf("attribute key %s not found on token", attrKey)
}
