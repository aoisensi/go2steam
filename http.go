package steam

import (
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

func (s *steam) htmlGet(url string) (*goquery.Document, error) {
	resp, err := s.service.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return goquery.NewDocumentFromReader(resp.Body)
}
