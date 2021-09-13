package pdfco

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type JobStatus string

const (
	JobWorking JobStatus = "working"
	JobSuccess JobStatus = "success"
	BaseURL              = "https://api.pdf.co/v1"
)

type Job struct {
	DownloadURL string
	Status      string
	FileName    string
}

func (t *Job) GetStatus() (JobStatus, error) {
	return JobWorking, nil
}

func (j *Job) DownloadPDF(ctx context.Context, file io.Writer) error {
	log.Print("downloading pdf")
	request, err := http.NewRequest("GET", j.DownloadURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	request = request.WithContext(ctx)
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("response status: %s", resp.Status)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		log.Printf("warning: expected content-type to be application/json, got %s", resp.Header.Get("Content-Type"))
	}

	defer resp.Body.Close()

	log.Print("copying pdf to request body")
	n, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d byres received", n)
	return nil
}

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

func (c *Client) NewHTMLToPDFJob(html string) (*Job, error) {

	payload, err := json.Marshal(ConvertFromHTMLRequest{
		Async:       false,
		HTML:        html,
		PaperSize:   "A4",
		FileName:    "report.pdf",
		Orientation: "Portrait",
	})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("GET",
		fmt.Sprintf("%s/pdf/convert/from/html?x-api-key=%s", BaseURL, c.apiKey),
		bytes.NewBuffer(payload),
	)

	if err != nil {
		return nil, err
	}

	request.Header.Add("x-api-key", c.apiKey)

	log.Printf("sending request: %s", request.URL.String())
	resp, err := http.DefaultClient.Do(request)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: %s", resp.Status)
	}
	defer resp.Body.Close()
	jobResponse := ConvertFromHTMLResponse{}
	err = json.NewDecoder(resp.Body).Decode(&jobResponse)
	if err != nil {
		return nil, err
	}
	log.Printf("job success: url %s", jobResponse.URL)

	return &Job{
		Status:      string(JobSuccess),
		DownloadURL: jobResponse.URL,
		FileName:    jobResponse.FileName,
	}, nil
}

type pdfcoApiClient struct {
	apiKey  string
	baseURL string
}

type ConvertFromHTMLRequest struct {
	HTML        string `json:"html"`
	Async       bool   `json:"async"`
	PaperSize   string `json:"paperSize"`
	FileName    string `json:"name"`
	Header      string `json:"header"`
	Footer      string `json:"footer"`
	Orientation string `json:"orientation"`
}

type ConvertFromHTMLResponse struct {
	URL              string `json:"url"`
	PageCount        int    `json:"pageCount"`
	Error            bool   `json:"error"`
	Status           int    `json:"status"`
	FileName         string `json:"name"`
	RemainingCredits int    `json:"remainingCredits"`
}
