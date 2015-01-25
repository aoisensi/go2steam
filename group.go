package steam

import "fmt"

const (
	urlGroup = "https://steamcommunity.com/groups/%s"
)

type group struct {
	s      *steam
	custom string
	Group
}

type Group interface {
	AnnounceCreate(title, body string) error
	Announces(page int) (*GroupAnnounces, error)
}

func (s *steam) LoadGroupFromCustom(custom string) (Group, error) {
	g := new(group)
	g.s = s
	g.custom = custom
	return g, nil
}

func (g *group) toURL(url string) string {
	return fmt.Sprintf(url, g.custom)
}
