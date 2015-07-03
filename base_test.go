package lynd

import (
	"reflect"
	"testing"
)

func TestLocalTargetExists(t *testing.T) {
	var cases = []struct {
		target LocalTarget
		exists bool
	}{
		{LocalTarget{Path: "/etc"}, true},
		{LocalTarget{Path: "/"}, true},
		{LocalTarget{Path: "/a/b/c/d/e/f/kind/of/unlikely"}, false},
	}
	for _, c := range cases {
		exists := c.target.Exists()
		if exists != c.exists {
			t.Errorf("got %s, want %s", exists, c.exists)
		}
	}
}

type externalTask struct {
	Path string
}

func (task externalTask) Requires() []Task { return nil }
func (task externalTask) Run() error       { return nil }
func (task externalTask) Output() Target {
	return LocalTarget{Path: task.Path}
}

type simpleTask struct{}

func (task simpleTask) Requires() []Task {
	return []Task{externalTask{Path: "/hello/world"}}
}

func (task simpleTask) Run() error     { return nil }
func (task simpleTask) Output() Target { return nil }

type taskWithTwoReqs struct{}

func (task taskWithTwoReqs) Requires() []Task {
	return []Task{
		externalTask{Path: "/hello/world"},
		externalTask{Path: "/hello/moon"},
	}
}

func (task taskWithTwoReqs) Run() error     { return nil }
func (task taskWithTwoReqs) Output() Target { return nil }

func TestInputAsPath(t *testing.T) {
	var cases = []struct {
		task Task
		path string
	}{
		{simpleTask{}, "/hello/world"},
		{taskWithTwoReqs{}, "/hello/world"},
	}

	for _, c := range cases {
		p := LocalInput(c.task).Path
		if p != c.path {
			t.Errorf("got %s, want %s", p, c.path)
		}
	}
}

func TestInputAsPathList(t *testing.T) {
	var cases = []struct {
		task     Task
		pathlist []string
	}{
		{simpleTask{}, []string{"/hello/world"}},
		{taskWithTwoReqs{}, []string{"/hello/world", "/hello/moon"}},
	}

	for _, c := range cases {
		p := LocalPaths(c.task)
		if !reflect.DeepEqual(p, c.pathlist) {
			t.Errorf("got %v, want %v", p, c.pathlist)
		}
	}
}
