package tailer

type Config struct {
	Publisher string `json:"publisher"`
	URL       string `json:"url,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
}
