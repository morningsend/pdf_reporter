package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pdf_reporter/controller"
	"pdf_reporter/pdfco"
	"pdf_reporter/renderer"
	"syscall"
)

var (
	templateFile = flag.String("template", "template/index.html", "path to template file")
	apiKey       = flag.String("api-key", "", "PDF.co API key")
)

func main() {
	fmt.Println("hello world")
	rend, err := renderer.NewRenderer(*templateFile)
	if err != nil {
		log.Fatal(err)
	}

	pdf := pdfco.NewClient(*apiKey)

	handler := controller.New(rend, pdf)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := http.ListenAndServe("localhost:8080", handler); err != nil {
			log.Fatal(err)
		}
	}()

	<-sigs
}
