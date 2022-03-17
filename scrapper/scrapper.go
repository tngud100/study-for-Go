package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedjob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

func Scrape(term string) {
	var baseURL string = "http://kr.indeed.com/jobs?q=" + term + "&limit=50"
	var jobs []extractedjob
	c := make(chan []extractedjob)
	totalPages := getPages(baseURL)

	for i := 0; i < totalPages; i++ {
		go getPage(i, baseURL, c)
	}

	for i := 0; i < totalPages; i++ {
		extractedjob := <-c
		jobs = append(jobs, extractedjob...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func writeJobs(jobs []extractedjob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"http://kr.indeed.com/jobs?q=python&limit=50" + job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func getPage(page int, url string, mainC chan<- []extractedjob) {
	var jobs []extractedjob
	c := make(chan extractedjob)
	pageURL := url + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".mosaic-provider-jobcards")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractjob(card, c)
	})
	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	mainC <- jobs
}

func extractjob(card *goquery.Selection, c chan<- extractedjob) {
	id, _ := card.Attr(".data-jk > a")
	title := CleanString(card.Find(".jobtitle-newJob").Text())
	location := CleanString(card.Find(".companyLocation").Text())
	salary := CleanString(card.Find(".companyName").Text())
	summary := CleanString(card.Find(".job-snippet").Text())
	c <- extractedjob{id: id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary}
}

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages(url string) int {
	pages := 0
	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})
	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with status:", res.StatusCode)
	}
}
