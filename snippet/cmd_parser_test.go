package snippet

import (
	"fmt"
	"testing"
)

func TestZshParse(t *testing.T) {
	tests := []struct {
		line   string
		expect string
	}{
		{`: 1538942302:0;echo "hello"`, `echo "hello"`},
		{`: 1538942295:0;ls`, `ls`},
		{`: 1538942296:0;rm zsh_history`, `rm zsh_history`},
		// Edges
		{`;`, ``},
		{`; `, ` `},
	}

	for i, tt := range tests {
		parser := ZshCmdParser{}
		got := parser.Parse(tt.line)
		if got != tt.expect {
			t.Errorf("Case %d: Expected `%s` but got `%s`", i, tt.expect, got)
		}
	}
}

func TestBashParse(t *testing.T) {
	tests := []struct {
		line   string
		expect string
	}{
		{"cat /root/.bash_history", "cat /root/.bash_history"},
		{"echo 'hello'", "echo 'hello'"},
		{"ls", "ls"},
		{"file /usr/local/share/", "file /usr/local/share/"},
		// Edges, this implementation acts as a passthrough
		{"", ""},
		{" ", " "},
	}

	for i, tt := range tests {
		parser := BashCmdParser{}
		got := parser.Parse(tt.line)
		if got != tt.expect {
			t.Errorf("Case %d: Expected `%s` but got `%s`", i, tt.expect, got)
		}
	}
}

func TestFishParse(t *testing.T) {
	tests := []struct {
		line   string
		expect string
	}{
		{"- cmd: cat /some/file", "cat /some/file"},
		{"- cmd: who", "who"},
		{"- cmd: ls -al (PWD)", "ls -al (PWD)"},
		// It should ignore lines without the cmd prefix
		{"when: 1538944231", ""},
		{"    - /bin/sh", ""},
		// Edges
		{"", ""},
		{" ", ""},
	}

	for i, tt := range tests {
		parser := FishCmdParser{}
		got := parser.Parse(tt.line)
		if got != tt.expect {
			t.Errorf("Case %d: Expected `%s` but got `%s`", i, tt.expect, got)
		}
	}
}

func TestGetCmdParser(t *testing.T) {
	tests := []struct {
		shellType string
		expect    CommandParser
		err       error
	}{
		{"zsh", ZshCmdParser{}, nil},
		{"bash", BashCmdParser{}, nil},
		{"fish", FishCmdParser{}, nil},
		// Edges
		{"", nil, fmt.Errorf("")},
		{" ", nil, fmt.Errorf("")},
		{"foo", nil, fmt.Errorf("")},
	}
	for i, tt := range tests {
		got, err := GetCmdParser(tt.shellType)
		if got != tt.expect {
			t.Errorf("Case %d: Expected `%s` but got `%s`", i, tt.expect, got)
		}

		if tt.err == nil && err != nil {
			t.Errorf("Case %d: Wasn't supposed to return an error but did: %+v", i, err)
		}

		if tt.err != nil && err == nil {
			t.Errorf("Case %d: Was supposed to return an error but didn't: %+v", i, err)
		}
	}
}
