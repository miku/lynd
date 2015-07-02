package common

import "testing"

func TestExecutable(t *testing.T) {
	var cases = []struct {
		name   string
		exists bool
	}{
		{"ls", true},
		{"lslslslsls", false},
	}

	for _, c := range cases {
		task := Executable{Name: c.name}
		exists := task.Output().Exists()
		if exists != c.exists {
			t.Errorf("got %s, want %s", exists, c.exists)
		}
	}
}
