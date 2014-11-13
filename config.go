package tailer

// Config describe what we want and publish
type Config struct {
	FileGlob      string `json:"fileglob,omitempty"`
	Match         string `json:"match,omitempty"`
	IgnorePattern string `json:"ignore,omitempty"`
}
