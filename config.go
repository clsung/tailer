package tailer

// Config describe what we want and publish
type Config struct {
	FileGlob      string `json:"fileglob,omitempty"`
	Pattern       string `json:"pattern,omitempty"`
	IngorePattern string `json:"ignore,omitempty"`
}
