package lynd

import (
	"log"
	"os"
)

type Target interface {
	Exists() bool
}

type Task interface {
	// Requires zero or more other tasks. Use slice for now.
	Requires() []Task
	// Run runs business logic.
	Run() error
	// Output is something that either exists or not.
	Output() Target
}

type LocalTarget struct {
	Path string
}

func (target LocalTarget) Exists() bool {
	_, err := os.Stat(target.Path)
	return !os.IsNotExist(err)
}

func Input(task Task) Target {
	inputs := Inputs(task)
	if len(inputs) == 0 {
		panic("no inputs to task")
	}
	return inputs[0]
}

func Inputs(task Task) []Target {
	var targets []Target
	for _, dep := range task.Requires() {
		targets = append(targets, dep.Output())
	}
	return targets
}

func LocalInput(task Task) LocalTarget {
	inputs := LocalInputs(task)
	if len(inputs) == 0 {
		panic("no inputs to task")
	}
	return inputs[0]
}

func LocalInputs(task Task) []LocalTarget {
	var targets []LocalTarget
	for _, dep := range task.Requires() {
		t, ok := dep.Output().(LocalTarget)
		if !ok {
			continue
		}
		targets = append(targets, t)
	}
	return targets
}

func LocalInputPath(task Task) string {
	return LocalInput(task).Path
}

func LocalInputPaths(task Task) []string {
	var paths []string
	for _, li := range LocalInputs(task) {
		paths = append(paths, li.Path)
	}
	return paths
}

func LocalOutput(task Task) LocalTarget {
	return task.Output().(LocalTarget)
}

func LocalOutputPath(task Task) string {
	return LocalOutput(task).Path
}

func Build(task Task) error {
	log.Printf("building %+v (pending)", task)
	if task.Output().Exists() {
		log.Printf("task %+v (done)", task)
		return nil
	}
	for _, t := range task.Requires() {
		err := Build(t)
		if err != nil {
			return err
		}
	}
	err := task.Run()
	if err != nil {
		return err
	}
	log.Println("finished task")
	return nil
}
