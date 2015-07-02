package lynd

import (
	"io/ioutil"
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

type simpleTask struct{}

func (task simpleTask) Requires() interface{} {
	return LocalTarget{Path: "/hello/world"}
}

func (task simpleTask) Run() error     { return nil }
func (task simpleTask) Output() Target { return nil }

type taskWithTwoReqs struct{}

func (task taskWithTwoReqs) Requires() interface{} {
	return []Target{
		LocalTarget{Path: "/hello/world"},
		LocalTarget{Path: "/hello/moon"},
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
		p := In(c.task).Path()
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
		p := In(c.task).PathList()
		if !reflect.DeepEqual(p, c.pathlist) {
			t.Errorf("got %v, want %v", p, c.pathlist)
		}
	}
}

type taskWithOutput struct{}

func (task taskWithOutput) Requires() interface{} { return nil }
func (task taskWithOutput) Run() error            { return nil }
func (task taskWithOutput) Output() Target {
	return LocalTarget{Path: "/tmp/Hello"}
}

func TestTaskWithOutput(t *testing.T) {
	var cases = []struct {
		task    Task
		content string
	}{
		{taskWithOutput{}, "Hello World"},
	}

	for _, c := range cases {
		f, err := Out(c.task).CreateFile()
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		_, err = f.WriteString(c.content)
		if err != nil {
			t.Error(err)
		}

		err = f.Commit()
		if err != nil {
			t.Error(err)
		}

		file, err := Out(c.task).File()
		if err != nil {
			t.Error(err)
		}
		content, err := ioutil.ReadAll(file)
		if err != nil {
			t.Error(err)
		}
		if string(content) != c.content {
			t.Errorf("got %s, want %s", string(content), c.content)
		}
	}
}
