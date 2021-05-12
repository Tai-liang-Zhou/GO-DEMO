package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Video struct {
	Title string
	Link  string
}

func main() {

	fName := "videos.json"
	file, err := os.Create(fName)

	if err != nil {
		log.Fatal("Cannot create file %q: %s\n", fName, err)
		return
	}

	defer file.Close()

	// Instantiate default collector
	c := colly.NewCollector()

	videos := make([]Video, 0, 200)

	// detailCollector := c.Clone()
	var formData = map[string]string{
		"agree":  "yes",
		"submit": "是，我已年滿18歲。Yes, I am.",
	}

	var count int = 0

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// fmt.Printf("link : http://wwwddd.eyny.com/%s , %s\n", link, e.Text)
		if strings.Index(e.Text, "MG") > -1 || strings.Index(e.Text, "Mg") > -1 || strings.Index(e.Text, "mg") > -1 {
			title := e.Text
			fmt.Printf("link : http://wwwddd.eyny.com/%s , %s\n", link, title)

			video := Video{
				Title: title,
				Link:  "http://wwwddd.eyny.com/" + link,
			}

			videos = append(videos, video)
		}

	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")
		if strings.Index(e.Text, "下一頁") > -1 {

			if count == 3 {
				return
			} else {
				count += 1
				c.Post("http://wwwddd.eyny.com/"+link, formData)
			}
		}

	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.190 Safari/537.36")
		log.Println("visiting", r.URL.String())
	})

	c.Post("http://www.eyny.com/forum-2-1.html", formData) // 進到該url 執行POST

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	// Dump json to the standard output
	enc.Encode(videos)

}
