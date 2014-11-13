package tailer

// Config describe what we want and publish
type Config struct {
	FileGlob      string `json:"fileglob,omitempty"`
	Match         string `json:"match,omitempty"`
	IngorePattern string `json:"ignore,omitempty"`
}
