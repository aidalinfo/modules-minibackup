package utils

type BackupArgs struct {
	Mongo       *Mongo `json:"Mongo,omitempty"`
	Path        string `json:"path,omitempty"`
	GlacierMode bool   `json:"Glaciermode,omitempty"`
}

type Mongo struct {
	Host      string   `json:"host,omitempty"`
	Port      string   `json:"port,omitempty"`
	User      string   `json:"user,omitempty"`
	Password  string   `json:"password,omitempty"`
	SSL       bool   `json:"ssl,omitempty"`
}
