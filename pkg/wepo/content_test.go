package wepo

import (
	"testing"
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
				"123\n456\n78",
			},
		},
		{
			desc: "test append line with dividing value",
			input: appendLineInput{
				[]string{
					"123",
					"4567",
					"8901",
				},
				10,
			},
			want: []string{
				"123\n4567",
				"8901",
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

func TestContents(t *testing.T) {
	type contentInput struct {
		input string
		cfg   wepoConfig
	}

	cfg := wepoConfig{
		CharLimit: 10,
	}

	testCases := []struct {
		desc  string
		input contentInput
		want  []string
	}{
		{
			desc: "test new contents 1",
			input: contentInput{
				"123 456 78",
				cfg,
			},
			want: []string{
				"123 456 78",
			},
		},
		{
			desc: "test new contents with line feed",
			input: contentInput{
				"123\n4567\n8901",
				cfg,
			},
			want: []string{
				"123\n4567",
				"8901",
			},
		},
		{
			desc: "test new contents with empty line",
			input: contentInput{
				"123\n4567\n\n8901",
				cfg,
			},
			want: []string{
				"123\n4567\n",
				"8901",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			lines, err := tc.input.cfg.Contents(tc.input.input)
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
