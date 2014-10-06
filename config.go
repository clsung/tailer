package tailer

// Config decide what we listen to, publish, and target
type Config struct {
	Publisher string `json:"publisher"`
	URL       string `json:"url,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
}
