package common

import (
	"fmt"
	"os/exec"

	"github.com/miku/lynd"
)

type ZeroTarget struct{}

func (target ZeroTarget) Exists() bool { return false }

type Executable struct {
	Name    string
	Message string
}

func (task Executable) Requires() interface{} { return nil }

func (task Executable) Run() error {
	return fmt.Errorf("%s not found - %s", task.Name, task.Message)
}

func (task Executable) Output() lynd.Target {
	p, err := exec.LookPath(task.Name)
	if err != nil {
		return ZeroTarget{}
	}
	return lynd.LocalTarget{Path: p}
}
