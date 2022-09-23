package config

import (
	"fmt"
	"strconv"

	"gopkg.in/ini.v1"
)

const CfgFileName = "config.ini"

var cfgFilePath string

// ini key
const (
	WebhookURLKey = "webhook_url"
	payloadKey    = "payload"
	charLimitKey  = "char_limit"
)

const (
	defaultCharLimit = 1024
	defaultPayload   = `{"content": "{input}"}`
)

// WepoConfig structure holds the parameters from the ini files.
type WepoConfig struct {
	URL       string
	CharLimit int
	Payload   string
}

const notSetMsg = `"%s" is not set in "%s"`

// New returns a *WepoConfig. Requires 'config.ini'.
func New(iniPath, section string) (*WepoConfig, error) {
	cfgFilePath = iniPath
	setting, err := ini.Load(cfgFilePath)
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
	url := setting.Section(sect).Key(WebhookURLKey).String()
	if len(url) == 0 {
		msg := fmt.Sprintf(notSetMsg, WebhookURLKey, cfgFilePath)
		if len(sect) != 0 {
			msg = fmt.Sprintf("[%s] %s", sect, msg)
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
