package config

import (
	"fmt"
	"testing"

	"gopkg.in/ini.v1"
)

const (
	TEST_INI_FILE              = "../../../test/data/config.ini"
	TEST_INI_FILE_EMPTY_GLOBAL = "../../../test/data/config_empty_global.ini"
)

const UNEXPECTED_ERROR_MSG = "failed to test '%s' is unexpected value: got: %v, want: %v"

const (
	GLOBAL_PAYLOAD = `{"content": "{input}"}`
	GLOBAL_LIMIT   = 1024
)

func loadINI(t *testing.T, fileName string) *ini.File {
	t.Helper()
	setting, err := ini.Load(fileName)
	if err != nil {
		t.Fatalf("failed to load file from '%s': %s", fileName, err)
	}
	return setting
}

func getTestConfig(t *testing.T) *ini.File {
	t.Helper()
	return loadINI(t, TEST_INI_FILE)
}
func getTestConfigEmptyGlobal(t *testing.T) *ini.File {
	t.Helper()
	return loadINI(t, TEST_INI_FILE_EMPTY_GLOBAL)
}

func TestCharLimit(t *testing.T) {
	limit := 0

	cfg := getTestConfig(t)

	testCases := []struct {
		desc    string
		section string
		want    int
	}{
		{
			desc:    fmt.Sprintf("test read '%s': global value", charLimitKey),
			section: "",
			want:    GLOBAL_LIMIT,
		},
		{
			desc:    fmt.Sprintf("test read '%s': empty value", charLimitKey),
			section: "sec1",
			want:    GLOBAL_LIMIT,
		},
		{
			desc:    fmt.Sprintf("test read '%s': setting value", charLimitKey),
			section: "sec2",
			want:    10,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			limit = charLimit(cfg, tc.section)
			if limit != tc.want {
				t.Fatalf(UNEXPECTED_ERROR_MSG, charLimitKey, limit, tc.want)
			}
		})
	}
}

func TestPayload(t *testing.T) {
	got := ""

	cfg := getTestConfig(t)

	testCases := []struct {
		desc    string
		section string
		want    string
	}{
		{
			desc:    fmt.Sprintf("test read '%s': global value", payloadKey),
			section: "",
			want:    GLOBAL_PAYLOAD,
		},
		{
			desc:    fmt.Sprintf("test read '%s': setting value", payloadKey),
			section: "sec1",
			want:    `{"content": "prefix {input} suffix"}`,
		},
		{
			desc:    fmt.Sprintf("test read '%s': empty value", payloadKey),
			section: "sec2",
			want:    GLOBAL_PAYLOAD,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got = payload(cfg, tc.section)
			if got != tc.want {
				t.Fatalf(UNEXPECTED_ERROR_MSG, payloadKey, got, tc.want)
			}
		})
	}
}

func TestURL(t *testing.T) {
	cfg := getTestConfig(t)

	testCases := []struct {
		desc    string
		section string
		want    string
		isError bool
	}{
		{
			desc:    fmt.Sprintf("test read '%s': global value", WebhookURLKey),
			section: "",
			want:    "https://example.com",
			isError: false,
		},
		{
			desc:    fmt.Sprintf("test read '%s': setting value", WebhookURLKey),
			section: "sec1",
			want:    "https://example.com/sec1",
			isError: false,
		},
		{
			desc:    fmt.Sprintf("test read '%s': empty value", WebhookURLKey),
			section: "sec3",
			want:    "",
			isError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := url(cfg, tc.section)
			if tc.isError {
				if err == nil {
					t.Fatalf("failed to test: no error has occurred")
				}
			} else {
				if got != tc.want {
					t.Fatalf(UNEXPECTED_ERROR_MSG, WebhookURLKey, got, tc.want)
				}
			}
		})
	}
}

func TestConfigEmptyGlobal(t *testing.T) {

	cfg := getTestConfigEmptyGlobal(t)

	limit := charLimit(cfg, "")
	if limit != GLOBAL_LIMIT {
		t.Fatalf(UNEXPECTED_ERROR_MSG, charLimitKey, limit, GLOBAL_LIMIT)
	}

	payload := payload(cfg, "")
	if payload != GLOBAL_PAYLOAD {
		t.Fatalf(UNEXPECTED_ERROR_MSG, payloadKey, payload, GLOBAL_PAYLOAD)
	}

	_, err := url(cfg, "")
	if err == nil {
		t.Fatalf("failed to test: no error has occurred")
	}
}
