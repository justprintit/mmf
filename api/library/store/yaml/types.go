package yaml

type User struct {
	Username string
	Name     string  `yaml:"-"`
	Avatar   string  `yaml:",omitempty"`
	Groups   []Group `yaml:",omitempty"`
}

type Group struct {
	Id        int
	Name      string  `yaml:"-"`
	Subgroups []Group `yaml:",omitempty"`
}