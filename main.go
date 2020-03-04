package main

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

type parsedJob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

var sourceURL string = "https://kr.indeed.com/취업?q=python&limit=50"

func main() {
	c := make(chan []parsedJob)
	var jobs []parsedJob
	totalPgs := getPages()

	for i := 0; i < totalPgs; i++ {
		go getPage(i, c)

	}

	for i := 0; i < totalPgs; i++ {
		parsedJob := <-c
		jobs = append(jobs, parsedJob...)
	}

	createcsv(jobs)
	fmt.Println("Extraction Success", len(jobs))

}

func getPage(page int, mainC chan<- []parsedJob) {
	var jobs []parsedJob
	c := make(chan parsedJob)
	pageURL := sourceURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	pullCards := doc.Find(".jobsearch-SerpJobCard")

	pullCards.Each(func(i int, card *goquery.Selection) {
		go extractCard(card, c)

	})

	for i := 0; i < pullCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs

}

func extractCard(card *goquery.Selection, c chan<- parsedJob) {
	id, _ := card.Attr("data-jk")
	title := cleanStr(card.Find(".title>a").Text())
	location := cleanStr(card.Find(".sjcl").Text())
	salary := cleanStr(card.Find(".salaryText").Text())
	summary := cleanStr(card.Find(".summary").Text())
	c <- parsedJob{id: id, title: title, location: location, salary: salary, summary: summary}

}

func cleanStr(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages() int {
	pages := 0
	res, err := http.Get(sourceURL)
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

func createcsv(jobs []parsedJob) {
	file, err := os.Create("list_of_jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSline := []string{job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSline)
		checkErr(jwErr)

	}

}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed w Status:", res.StatusCode)

	}

}
