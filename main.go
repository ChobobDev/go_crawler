package main

import (
	"fmt"
	"log"
	"net/http"
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

var source_URL string = "https://kr.indeed.com/취업?q=python&limit=50"

func main() {
	var jobs []parsedJob
	totalPgs := getPages()

	for i := 0; i < totalPgs; i++ {
		parsedJob := getPage(i)
		jobs = append(jobs, parsedJob...)

	}
	fmt.Println(jobs)

}
func getPage(page int) []parsedJob {
	var jobs []parsedJob
	pageURL := source_URL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	pullCards := doc.Find(".jobsearch-SerpJobCard")

	pullCards.Each(func(i int, card *goquery.Selection) {
		job := extractCard(card)
		jobs = append(jobs, job)

	})

	return jobs

}

func extractCard(card *goquery.Selection) parsedJob {
	id, _ := card.Attr("data-jk")
	title := cleanStr(card.Find(".title>a").Text())
	location := cleanStr(card.Find(".sjcl").Text())
	salary := cleanStr(card.Find(".salaryText").Text())
	summary := cleanStr(card.Find(".summary").Text())
	return parsedJob{id: id, title: title, location: location, salary: salary, summary: summary}

}

func cleanStr(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages() int {
	pages := 0
	res, err := http.Get(source_URL)
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
		log.Fatalln("Request failed w Status:", res.StatusCode)

	}

}
