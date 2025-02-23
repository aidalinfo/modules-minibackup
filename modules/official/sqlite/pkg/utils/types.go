package utils

type BackupArgs struct {
	Sqlite       *Sqlite `json:"sqlite,omitempty"`
	Path        string `json:"path,omitempty"`
	GlacierMode bool   `json:"Glaciermode,omitempty"`
}

type Sqlite struct {
	Paths       []string `json:"paths,omitempty"`
}
