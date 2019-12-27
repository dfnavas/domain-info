package middleware

import (
	"log"

	"github.com/anaskhan96/soup"
)

// TagXtractor :
type TagXtractor struct {
	Get func(string) (string, string)
}

func getTitleAndLogo(url string) (string, string) {
	resp, err := soup.Get("https://" + url)
	if err != nil {
		log.Println("Error:", err)
		return "unavailable", "unavailable"
	}
	doc := soup.HTMLParse(resp)
	titleNode := doc.Find("title")
	logoNode := doc.Find("link", "rel", "shortcut")
	var title string
	var logo string
	if titleNode.Error != nil {
		title = "Not found"
		log.Println("Error:", titleNode.Error)
	} else {
		title = titleNode.Text()
	}
	if logoNode.Error != nil {
		logo = "Not found"
		log.Println("Error:", logoNode.Error)

	} else {
		logo = logoNode.Attrs()["href"]
	}
	return title, logo
}

// CreateTagXtractor :
func CreateTagXtractor() *TagXtractor {
	return &TagXtractor{
		Get: getTitleAndLogo,
	}
}