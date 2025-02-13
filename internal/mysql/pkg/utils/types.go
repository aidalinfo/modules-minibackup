package utils

type BackupArgs struct {
	Mysql       *Mysql `json:"mysql,omitempty"`
	Path        string `json:"path,omitempty"`
	GlacierMode bool   `json:"Glaciermode,omitempty"`
}

type Mysql struct {
	All       bool     `json:"all,omitempty"`
	Databases []string `json:"databases,omitempty"`
	Host      string   `json:"host,omitempty"`
	Port      string   `json:"port,omitempty"`
	User      string   `json:"user,omitempty"`
	Password  string   `json:"password,omitempty"`
	SSL       string   `json:"ssl,omitempty"`
}
