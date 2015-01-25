package steam

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const (
	urlGroupAnnounceDelete = "https://steamcommunity.com/groups/%s/announcements/delete/%d"
	urlGroupAnnounce       = "https://steamcommunity.com/groups/%s/announcements"
	urlGroupAnnounceCreate = "https://steamcommunity.com/groups/%s/announcements/create"
)

type GroupAnnounce struct {
	g     *group
	s     *goquery.Selection
	Title string
	URL   string
	Id    uint64
}

type GroupAnnounces struct {
	List []*GroupAnnounce
}

func (g *group) Announces(page int) (*GroupAnnounces, error) {
	var v url.Values
	if page > 0 {
		v = url.Values{"p": {string(page)}}
	}
	doc, err := g.s.htmlGet(fmt.Sprintf(urlGroupAnnounce, g.custom), v)
	if err != nil {
		return nil, err
	}
	r := new(GroupAnnounces)
	annos := doc.Find("#announcementsContainer > .anouncement")
	r.List = make([]*GroupAnnounce, annos.Length())

	annos.Each(func(i int, s *goquery.Selection) {
		r.List[i], err = g.loadAnnounce(s)
	})
	return r, nil
}

func (g *group) AnnounceCreate(title, body string) error {
	cdoc, err := g.s.htmlGet(g.toURL(urlGroupAnnounceCreate), nil)
	if err != nil {
		return err
	}
	token, ok := cdoc.Find("form#post_announcement_form input[name=sessionID]").Attr("value")
	if !ok {
		return errors.New("failed")
	}
	v := url.Values{
		"sessionID": {token},
		"action":    {"post"},
		"headline":  {title},
		"body":      {body},
	}
	res, err := g.s.service.PostForm(g.toURL(urlGroupAnnounce), v)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (g *group) loadAnnounce(s *goquery.Selection) (*GroupAnnounce, error) {
	class, _ := s.Attr("class")
	if class != "group_content" {
		return nil, errors.New("this is not group tag")
	}
	a := new(GroupAnnounce)
	a.g = g
	a.s = s
	title := s.Find(".large_title").First()
	a.Title = title.Text()
	a.URL, _ = title.Attr("href")
	idid, _ := s.Find(".bodytext").Attr("id")
	a.Id, _ = strconv.ParseUint(idid[6:], 10, 64)
	return a, nil
}

func (a *GroupAnnounce) Delete() error {
	session, _ := a.g.s.sessionId()
	v := url.Values{"sessionID": {session}}
	resp, err := a.g.s.httpGet(fmt.Sprintf(urlGroupAnnounceDelete, a.g.Group, a.Id), v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
