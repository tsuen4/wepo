package wepo

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

// PostDiscord sends content to Wepoconfig.URL
func (cfg wepoConfig) PostDiscord(content string) error {

	body := strings.ReplaceAll(cfg.Payload, "{input}", content)
	resp, err := http.Post(cfg.URL, "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status code err: got: %d, want: %d", resp.StatusCode, http.StatusNoContent)
	}
	return nil
}
