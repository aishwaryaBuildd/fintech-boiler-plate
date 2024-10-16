package utils

// VideoCipherClient struct definition
type VideoCipherClient struct {
	url    string
	secret string
}

// NewVideoCipherClient initializes and returns a new VideoCipherClient
func NewVideoCipherClient(url, secret string) *VideoCipherClient {
	return &VideoCipherClient{url: url, secret: secret}
}
