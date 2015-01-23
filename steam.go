package steam

type steam struct {
	Steam
}

type Steam interface {
}

func NewSteam() Steam {
	return new(steam)
}
