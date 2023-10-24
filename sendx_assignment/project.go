package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

var (
	crawledPages = make(map[string]crawledPageData)
	accessLog    = make(map[string][]time.Time)
	lock         = &sync.RWMutex{}
)

type crawledPageData struct {
	Content   string
	TimeStamp time.Time
}

func crawlHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	isPaying := r.URL.Query().Get("isPaying") == "true"

	lock.RLock()
	cachedPage, exists := crawledPages[url]
	lock.RUnlock()

	// Log access time for the URL
	logAccess(url)

	if exists && time.Since(cachedPage.TimeStamp).Minutes() <= 60 {
		// Update the timestamp for the same URL
		lock.Lock()
		cachedPage.TimeStamp = time.Now()
		lock.Unlock()

		w.Write([]byte(cachedPage.Content)) // Return the cached page
	} else {
		if isPaying {
			go crawlAndCache(url, isPaying)
		} else {
			go crawlAndCache(url, isPaying)
		}
		w.Write([]byte("Crawling in progress..."))
	}
}

func crawlAndCache(url string, isPaying bool) {
	c := colly.NewCollector()
	var pageContent string

	c.OnHTML("html", func(e *colly.HTMLElement) {
		pageContent = e.Text
	})

	err := c.Visit(url)
	if err != nil {
		fmt.Printf("Error crawling URL %s: %v\n", url, err)
		return
	}

	lock.Lock()
	crawledPages[url] = crawledPageData{Content: pageContent, TimeStamp: time.Now()}
	lock.Unlock()
}

func logAccess(url string) {
	lock.Lock()
	defer lock.Unlock()

	if _, exists := accessLog[url]; !exists {
		accessLog[url] = []time.Time{}
	}
	accessLog[url] = append(accessLog[url], time.Now())
}

func accessLogHandler(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	defer lock.RUnlock()

	html := "<h1>Access Log</h1>"

	for url, accessTimes := range accessLog {
		html += fmt.Sprintf("<h2>URL: %s</h2>", url)
		html += "<ul>"
		for _, accessTime := range accessTimes {
			html += fmt.Sprintf("<li>%s</li>", accessTime.Format("2006-01-02 15:04:05"))
		}
		html += "</ul>"
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	defer lock.RUnlock()

	html := "<h1>Crawled Data</h1>"

	for url, data := range crawledPages {
		entry := fmt.Sprintf("<div class='entry'><strong>URL:</strong> %s<br><strong>Timestamp:</strong> %s</div>", url, data.TimeStamp.Format("2006-01-02 15:04:05"))
		html += entry
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func main() {
	fmt.Println("Starting the server on :8080...")
	http.HandleFunc("/crawl", crawlHandler)
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/accesslog", accessLogHandler) // New route for access log
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(":8080", nil)
}
