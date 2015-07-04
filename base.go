package lynd

import (
	"errors"
	"io"
)

var (
	errNoDependencies = errors.New("no dependencies")
	errNotImplemented = errors.New("not implemented")
)

type Target interface {
	// Name serves as identifier or locator, e.g. the filename.
	Name() string
	// Exists determines, whether this target need to be built or not.
	Exists() bool
	// Input/output for target. If a task is able to work with the usual
	// interfaces only, we do not need to change the business logic, when we
	// switch a target implementation from local files to e.g. hdfs or other
	// storage options. Embed `lynd.WithoutIO`, if a target does not support IO.
	io.ReadWriteCloser
}

type Task interface {
	// Requires should return a Task, []Task or map[string]Task.
	Requires() interface{}
	// Run runs the business logic.
	Run() error
	// Output is a single target.
	Output() Target
}

// Dependencies return a list of dependent tasks. Non-recursive.
func Dependencies(task Task) []Task {
	req := task.Requires()
	switch t := req.(type) {
	case Task:
		return []Task{t}
	case []Task:
		return t
	case map[string]Task:
		var ts []Task
		for _, v := range t {
			ts = append(ts, v)
		}
		return ts
	default:
		return []Task{}
	}
}

// MustInput returns the first output of the required task and panics, if
// there is no target to return.
func MustInput(task Task) Target {
	t, err := Input(task)
	if err != nil {
		panic(err)
	}
	return t
}

// Input returns the first target and an error, if no target is found.
func Input(task Task) (Target, error) {
	deps := Dependencies(task)
	if len(deps) == 0 {
		return &Failed{}, errNoDependencies
	}
	return deps[0].Output(), nil
}

// Inputs returns a list of targets from the tasks' requirements.
func Inputs(task Task) []Target {
	var targets []Target
	for _, dep := range Dependencies(task) {
		targets = append(targets, dep.Output())
	}
	return targets
}

// Build build a task and all its dependencies, if required.
func Build(task Task) error {
	if task.Output().Exists() {
		return nil
	}

	for _, t := range Dependencies(task) {
		err := Build(t)
		if err != nil {
			return err
		}
	}

	err := task.Run()
	if err != nil {
		return err
	}

	return nil
}
