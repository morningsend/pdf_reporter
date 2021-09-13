package controller

import (
	"log"
	"net/http"
	"pdf_reporter/pdfco"
	"pdf_reporter/renderer"
	"pdf_reporter/report"
)

type Controller struct {
	rend *renderer.Renderer
	pdf  *pdfco.Client
}

func New(r *renderer.Renderer, pdf *pdfco.Client) *Controller {
	return &Controller{
		rend: r,
		pdf:  pdf,
	}
}

var _ http.Handler = (*Controller)(nil)

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	report := report.Report{
		Title:         "Cute Cats",
		CoverImageURL: "https://images.unsplash.com/photo-1561948955-570b270e7c36?crop=entropy&cs=tinysrgb&fit=crop&fm=jpg&h=1500&ixid=MnwxfDB8MXxyYW5kb218MHx8fHx8fHx8MTYzMTI3NjE0OA&ixlib=rb-1.2.1&q=80&utm_campaign=api-credit&utm_medium=referral&utm_source=unsplash_source&w=1200",
	}
	switch url {
	case "/html":
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Add("content-type", "text/html")
		if err := c.rend.Render(w, report); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	case "/pdf":
		log.Print("/pdf")
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		htmlString, err := c.rend.RenderToString(report)
		log.Print("rendered html: %s", htmlString)
		if err != nil {
			w.Header().Add("content-type", "text/html")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		job, err := c.pdf.NewHTMLToPDFJob(htmlString)
		if err != nil {
			w.Header().Set("content-type", "text/html")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Printf("error: %s", err)
			return
		}
		log.Print("creating pdf job")
		log.Printf("success, downloading file")
		w.Header().Set("content-type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename="+job.FileName)

		if err := job.DownloadPDF(r.Context(), w); err != nil {
			w.Header().Set("content-type", "text/html")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
