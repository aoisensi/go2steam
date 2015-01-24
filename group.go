package steam

import (
	"errors"
	"fmt"
	"net/url"
)

const (
	urlGroup               = "https://steamcommunity.com/groups/%s"
	urlGroupAnnounce       = "https://steamcommunity.com/groups/%s/announcements"
	urlGroupAnnounceCreate = "https://steamcommunity.com/groups/%s/announcements/create"
)

type group struct {
	s      *steam
	custom string
	Group
}

type Group interface {
	AnnounceCreate(title, body string) error
}

func (s *steam) LoadGroupFromCustom(custom string) (Group, error) {
	g := new(group)
	g.s = s
	g.custom = custom
	return g, nil
}

func (g *group) AnnounceCreate(title, body string) error {
	cdoc, err := g.s.htmlGet(g.toURL(urlGroupAnnounceCreate))
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

func (g *group) toURL(url string) string {
	return fmt.Sprintf(url, g.custom)
}
