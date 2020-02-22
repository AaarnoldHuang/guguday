package Module

type tomlconfig struct {
	DB database `toml:"database"`
}

type database struct {
	Server   string
	Ports    []int
	DBname   string
	DBuser   string
	DBpasswd string
	ConnMax  int `toml:"connection_max"`
	Enabled  bool
}
