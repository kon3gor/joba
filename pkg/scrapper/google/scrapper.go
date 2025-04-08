package google

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kon3gor/joba/pkg/scrapper"
	"golang.org/x/xerrors"
)

const (
	baseURL = "https://www.google.com/about/careers/applications"
)

type Scrapper struct {
	query     string
	pageLimit int
}

type Result struct {
	ID    string
	Title string
	Link  string
}

func (r Result) GetID() string {
	return r.ID
}

func (r Result) GetLink() string {
	return r.Link
}

func (r Result) String() string {
	return r.Title
}

func (r Result) GetAdditionalInfo() struct{} {
	return struct{}{}
}

func NewScrapper(query string, pageLimit int) scrapper.Scrapper {
	return &Scrapper{
		query:     query,
		pageLimit: pageLimit,
	}
}

func (s *Scrapper) Scrap() ([]scrapper.Result, error) {
	var err error
	var res []scrapper.Result

	total := make([]scrapper.Result, 0, 100)
	page := 1
	for err == nil && page <= s.pageLimit {
		res, err = s.scrap(page)
		if err == nil {
			if len(res) == 0 {
				break
			}

			total = append(total, res...)
		}
		page++
	}
	return total, nil
}

func (s *Scrapper) scrap(page int) ([]scrapper.Result, error) {
	req, err := http.NewRequest(http.MethodGet, s.query, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("page", strconv.Itoa(page))
	query.Set("sort_by", "date")

	req.URL.RawQuery = query.Encode()
	req.Header.Add("User-Agent", "kon3gor agent 0.0.1")
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, xerrors.Errorf("non 200 code: %d, err parsing: %w", res.StatusCode, err)
		}
		fmt.Println(string(body))
		return nil, xerrors.Errorf("non 200 code: %d", res.StatusCode)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	result := make([]scrapper.Result, 0)

	doc.Find(".lLd3Je").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".QJPWVe").Text()
		relLink, _ := s.Find(".VfPpkd-dgl2Hf-ppHlrf-sM5MNb").Find("a").Attr("href")
		result = append(result, Result{
			ID:    extractIDFromLink(relLink),
			Title: title,
			Link:  fmt.Sprintf("%s/%s", baseURL, relLink),
		})
	})

	return result, nil
}

func extractIDFromLink(link string) string {
	link, _ = strings.CutPrefix(link, "jobs/results/")
	id, _, _ := strings.Cut(link, "?")
	return id
}
