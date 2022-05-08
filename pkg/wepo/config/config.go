package config

import (
	"flag"
	"fmt"
	"path/filepath"
	"strconv"

	"gopkg.in/ini.v1"
)

const CfgFileName = "config.ini"

// ini key
const (
	webhookURLKey = "webhook_url"
	payloadKey    = "payload"
	charLimitKey  = "char_limit"
)

// ini section name
var (
	section = ""
)

const (
	defaultCharLimit = 1024
	defaultPayload   = `{"content": "{input}"}`
)

func init() {
	flag.StringVar(&section, "s", "", fmt.Sprintf(`Section name of "%s" where "%s" is set`, CfgFileName, webhookURLKey))
}

// WepoConfig structure holds the parameters from the ini files.
type WepoConfig struct {
	URL       string
	CharLimit int
	Payload   string
}

const notSetMsg = `"%s" is not set in "%s"`

// New returns a *WepoConfig. Requires '{cfgDirPath}/config.ini'.
func New(cfgDirPath string) (*WepoConfig, error) {
	setting, err := ini.Load(filepath.Join(cfgDirPath, CfgFileName))
	if err != nil {
		return nil, err
	}

	url, err := url(setting, section)
	if err != nil {
		return nil, err
	}

	return &WepoConfig{
		URL:       url,
		CharLimit: charLimit(setting, section),
		Payload:   payload(setting, section),
	}, nil
}

func charLimit(setting *ini.File, sect string) int {
	charLimitStr := setting.Section(sect).Key(charLimitKey).String()
	// get from global
	if len(charLimitStr) == 0 {
		charLimitStr = setting.Section("").Key(charLimitKey).String()
	}
	// initial value
	charLimit, err := strconv.ParseInt(charLimitStr, 0, 64)
	if err != nil {
		// force set
		charLimit = defaultCharLimit
	}
	return int(charLimit)
}

func url(setting *ini.File, sect string) (string, error) {
	url := setting.Section(sect).Key(webhookURLKey).String()
	if len(url) == 0 {
		msg := fmt.Sprintf(notSetMsg, webhookURLKey, CfgFileName)
		if len(section) != 0 {
			msg = fmt.Sprintf("[%s] %s", section, msg)
		}
		return "", fmt.Errorf(msg)
	}
	return url, nil
}

func payload(setting *ini.File, sect string) string {
	payload := setting.Section(sect).Key(payloadKey).String()
	// get from global
	if len(payload) == 0 {
		payload = setting.Section("").Key(payloadKey).String()
	}
	// initial value
	if len(payload) == 0 {
		payload = defaultPayload
	}
	return payload
}
