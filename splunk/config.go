package splunk

import (
	"log"
)

type Config struct {
	URL                string
	Username           string
	Password           string
	InsecureSkipVerify bool
}

// Client() returns a new client for accessing Splunk.
func (c *Config) Client() (*Client, error) {
	client := New(c.URL, c.Username, c.Password, c.InsecureSkipVerify)
	log.Printf("[INFO] Splunk Client configured for: %s@%s", c.Username, c.URL)
	return client, nil
}
