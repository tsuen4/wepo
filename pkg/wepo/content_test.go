package wepo

import (
	"fmt"
	"testing"

	"github.com/tsuen4/wepo/pkg/wepo/config"
)

func TestSplitRow(t *testing.T) {
	type splitInput struct {
		str   string
		limit int
	}

	testCases := []struct {
		desc  string
		input splitInput
		want  []string
	}{
		{
			desc: "test 2 char split",
			input: splitInput{
				"123456",
				2,
			},
			want: []string{
				"12",
				"34",
				"56",
			},
		},
		{
			desc: "test 3 char split with line feed",
			input: splitInput{
				"123456",
				3,
			},
			want: []string{
				"123",
				"456",
			},
		},
		{
			desc: "test no split",
			input: splitInput{
				"123456",
				6,
			},
			want: []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			lines, err := splitRow(tc.input.str, tc.input.limit)
			if err != nil {
				if err != errNoNeedSplit {
					t.Fatalf("failed to %s: %s", tc.desc, err)
				}
			}
			if len(lines) != len(tc.want) {
				t.Fatalf("failed to %s: got: %v, want: %v", tc.desc, len(lines), len(tc.want))
			}
			for i := 0; i < len(lines); i++ {
				if lines[i] != tc.want[i] {
					t.Errorf("failed to %s: line: %d, got: %v, want: %v", tc.desc, i+1, lines, tc.want)
				}
			}
		})
	}
}

func TestAppendLine(t *testing.T) {
	type appendLineInput struct {
		lines []string
		limit int
	}

	testCases := []struct {
		desc  string
		input appendLineInput
		want  []string
	}{
		{
			desc: "test append line",
			input: appendLineInput{
				[]string{
					"123",
					"456",
					"78",
				},
				10,
			},
			want: []string{
				`123\n456`,
				"78",
			},
		},
		{
			desc: "test append line with charLimit value",
			input: appendLineInput{
				[]string{
					"123",
					"4567",
					"8901",
				},
				10,
			},
			want: []string{
				`123\n4567`,
				`8901`,
			},
		},
		{
			desc: "test append line with value exceeds the limit",
			input: appendLineInput{
				[]string{
					"00",
					"123456789012",
				},
				10,
			},
			want: []string{
				"00",
				"1234567890",
				"12",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			lines := []string{}
			for _, line := range tc.input.lines {
				var err error
				lines, err = appendLine(lines, line, tc.input.limit)
				if err != nil {
					t.Fatalf("failed to %s: %s", tc.desc, err)
				}
			}
			if len(lines) != len(tc.want) {
				fmt.Println("got:", lines)
				fmt.Println("want:", tc.want)
				t.Fatalf("failed to %s: got: %v, want: %v", tc.desc, len(lines), len(tc.want))
			}
			for i := 0; i < len(lines); i++ {
				if lines[i] != tc.want[i] {
					t.Errorf("failed to %s: line: %d, got: %v, want: %v", tc.desc, i+1, lines[i], tc.want[i])
				}
			}
		})
	}
}

func TestNewContents(t *testing.T) {
	type in struct {
		input string
		limit int
	}

	testCases := []struct {
		desc string
		in   in
		want []string
	}{
		{
			desc: "test new contents 1",
			in: in{
				`123 456 78`,
				10,
			},
			want: []string{
				"123 456 78",
			},
		},
		{
			desc: "test new contents with line feed",
			in: in{
				`123\n4567\n8901`,
				10,
			},
			want: []string{
				`123\n4567`,
				`\n8901`,
			},
		},
		{
			desc: "test new contents with empty line",
			in: in{
				`123\n4567\n\n8901`,
				10,
			},
			want: []string{
				`123\n4567`,
				`\n\n8901`,
			},
		},
		{
			desc: "test new contents with characters that need to be escaped",
			in: in{
				`{"content": "{input}"}`,
				1024,
			},
			want: []string{
				`{\"content\": \"{input}\"}`,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			client := wepo{
				cfg: &config.WepoConfig{
					CharLimit: tc.in.limit,
				},
			}

			lines, err := client.NewContents(tc.in.input)
			if err != nil {
				t.Fatalf("failed to %s: %s", tc.desc, err)
			}
			if len(lines) != len(tc.want) {
				t.Fatalf("failed to %s: got: %v, want: %v", tc.desc, len(lines), len(tc.want))
			}
			for i := 0; i < len(lines); i++ {
				if lines[i] != tc.want[i] {
					t.Errorf("failed to %s: line: %d,  got: %v, want: %v", tc.desc, i+1, lines[i], tc.want[i])
				}
			}
		})
	}
}
