package scrape

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
)

type HrefReader func(bodyReader io.ReadCloser) (hrefs []string)

func ReadHrefs(bodyReader io.ReadCloser) (hrefs []string) {
	defer bodyReader.Close()

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
				hrefs = append(hrefs, href)
			}
		}
	}

	return hrefs
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
