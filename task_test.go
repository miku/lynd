package lynd

import "testing"

type noopTask struct{}

func (t noopTask) Requires() interface{} { return nil }
func (t noopTask) Run() error            { return nil }
func (t noopTask) Output() Target        { return nil }

type testTask struct {
	noopTask
	Val string `default:"Hello"`
}

type testTaskInt struct {
	noopTask
	Val int `default:"10"`
}

type testTaskFloat struct {
	noopTask
	Val float64 `default:"0.3"`
}

func TestSetDefaultsString(t *testing.T) {
	var cases = []struct {
		task testTask
		Val  string
	}{
		{testTask{}, "Hello"},
		{testTask{Val: "X"}, "X"},
	}

	for _, c := range cases {
		setDefaults(&c.task)
		if c.task.Val != c.Val {
			t.Errorf("got %s, want %s", c.task.Val, c.Val)
		}
	}
}

func TestSetDefaultsInt(t *testing.T) {
	var cases = []struct {
		task testTaskInt
		Val  int
	}{
		{testTaskInt{}, 10},
		{testTaskInt{Val: 20}, 20},
	}

	for _, c := range cases {
		setDefaults(&c.task)
		if c.task.Val != c.Val {
			t.Errorf("got %s, want %s", c.task.Val, c.Val)
		}
	}
}

func TestSetDefaultsFloat(t *testing.T) {
	var cases = []struct {
		task testTaskFloat
		Val  float64
	}{
		{testTaskFloat{}, 0.3},
		{testTaskFloat{Val: 0.7}, 0.7},
	}

	for _, c := range cases {
		setDefaults(&c.task)
		if c.task.Val != c.Val {
			t.Errorf("got %s, want %s", c.task.Val, c.Val)
		}
	}
}
