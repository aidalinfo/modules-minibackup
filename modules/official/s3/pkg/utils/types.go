package utils

type BackupArgs struct {
	S3       *S3 `json:"S3,omitempty"`
	Path        string `json:"path,omitempty"`
	GlacierMode bool   `json:"Glaciermode,omitempty"`
}

type S3 struct {
	All        bool     `json:"all,omitempty"`
	Bucket     []string `json:"bucket,omitempty"`
	Endpoint   string   `json:"endpoint,omitempty"`
	PathStyle  bool     `json:"pathStyle,omitempty"`
	Region     string   `json:"region,omitempty"`
	ACCESS_KEY string   `json:"ACCESS_KEY,omitempty"`
	SECRET_KEY string   `json:"SECRET_KEY,omitempty"`
}