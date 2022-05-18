package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] == "status" {
		fmt.Fprintf(w, "<h1>The connection is %v</h1>", Arduino.isConnected)
	} else {
		fmt.Fprintf(w, `<h1>Error 404</h1><p><a href="/status">/status</a></p>`)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[1:]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "%s", p.Body)
}

func loadPage(title string) (*Page, error) {
	color.Yellow("loadpage, got: %s", title)
	if title == "" {
		title = "index.html"
	}
	filename := title
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
