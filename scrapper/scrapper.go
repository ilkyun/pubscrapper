package scrapper

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedPaper struct {
	title    string
	journal  string
	impact   string
	pubdate  string
	authors  string
	abstract string
	issn     string
}

//Scrape papers
func Scrape(term string, retmax string, impThr float64, reldate string) {

	var impPapers []extractedPaper
	//	var jobs []extractedPaper
	//c := make(chan []extractedPaper)

	ids := getIds(term, retmax, reldate)
	extractedPapers := getInfo(ids)
	for i := 0; i < len(extractedPapers); i++ {
		ifs := checkIF(extractedPapers[i].journal, extractedPapers[i].issn)
		imp, _ := strconv.ParseFloat(ifs, 64)
		if imp >= impThr {
			extractedPapers[i].impact = ifs
			impPapers = append(impPapers, extractedPapers[i])
		}

	}
	writePapers(impPapers)
	fmt.Println("We scrapped " + strconv.Itoa(len(extractedPapers)) + " papers with terms and got " + strconv.Itoa(len(impPapers)) + " papers filtered by IF!")
}

func getIds(term string, retmax string, reldate string) string {
	baseURL := "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&term=" + term + "&retmax=" + retmax + "&reldate=" + reldate
	var ids string = ""
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)
	doc, err2 := goquery.NewDocumentFromReader(res.Body)
	checkErr(err2)
	doc.Find("id").Each(func(i int, s *goquery.Selection) {
		ids = ids + "," + s.Text()
	})
	ids = trimLeftChar(ids)

	return ids
}

func getInfo(ids string) []extractedPaper {
	var papers []extractedPaper
	var paper extractedPaper
	var searchCards []*goquery.Selection
	title := ""
	impact := ""
	journal := ""
	pubdate := ""
	abstract := ""
	issn := ""
	//c := make(chan extractedPaper)
	baseURL := "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pubmed&id=" + ids + "&retmode=xml&rettype=abstract"
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	doc, err2 := goquery.NewDocumentFromReader(res.Body)
	checkErr(err2)

	doc.Find("PubmedArticle").Each(func(i int, s *goquery.Selection) {
		searchCards = append(searchCards, s)
	})
	for i := 0; i < len(searchCards); i++ {
		title = searchCards[i].Find("ArticleTitle").Text()
		searchCards[i].Find("Journal").Each(func(i int, s *goquery.Selection) {
			journal = s.Find("Title").Text()
			pubdate = CleanString(s.Find("Pubdate").Text())
			issn = s.Find("ISSN").Text()
		})
		authorN := ""
		searchCards[i].Find("AuthorList").Find("Author").Each(func(i int, s *goquery.Selection) {
			authorN = authorN + ", " + s.Find("LastName").Text() + " " + s.Find("ForeName").Text()
		})
		authorN = trimLeftChar(authorN)
		abstract = searchCards[i].Find("AbstractText").Text()
		paper = extractedPaper{title: title, impact: impact, journal: journal, pubdate: pubdate, authors: authorN, abstract: abstract, issn: issn}
		papers = append(papers, paper)
	}
	return papers
}

func checkIF(journal string, issn string) string {
	csvFile, err := os.Open("JCR_ISSN_2018.csv")
	impact := "0"
	journal = strings.ToUpper(CleanString(journal))
	checkErr(err)
	rdr := csv.NewReader(bufio.NewReader(csvFile))
	rows, _ := rdr.ReadAll()
	for _, row := range rows {
		if row[2] == issn {
			impact = row[4]
			break
		} else if strings.ToUpper(row[1]) == journal {
			impact = row[4]
			break
		}
	}
	return impact
}
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status: res.StatusCode")
	}
}

//CleanString makes string clean
func CleanString(str string) string {
	str = strings.ReplaceAll(str, ".", "")
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func writePapers(papers []extractedPaper) {
	file, err := os.Create("results.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Title", "Journal", "Impact factor", "PubDate", "Authors", "Abstract"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, paper := range papers {
		jobSlice := []string{paper.title, paper.journal, paper.impact, paper.pubdate, paper.authors, paper.abstract}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}
