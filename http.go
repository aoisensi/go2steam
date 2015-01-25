package steam

import (
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func (s *steam) htmlPost(url string, data url.Values) (*goquery.Document, error) {
	resp, err := s.service.PostForm(url, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return goquery.NewDocumentFromReader(resp.Body)
}

func (s *steam) htmlGet(url string, v url.Values) (*goquery.Document, error) {

	resp, err := s.httpGet(url, v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return goquery.NewDocumentFromReader(resp.Body)
}

func (s *steam) httpGet(url string, v url.Values) (*http.Response, error) {
	if v != nil {
		url = url + "?" + v.Encode()
	}
	return s.service.Get(url)
}
