package utils

type BackupArgs struct {
	Fs       *Fs `json:"fs,omitempty"`
	Path        string `json:"path,omitempty"`
	GlacierMode bool   `json:"Glaciermode,omitempty"`
}

type Fs struct {
	Paths       []string `json:"paths,omitempty"`
}
