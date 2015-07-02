package lynd

import (
	"fmt"
	"os"

	"github.com/dchest/safefile"
)

type Target interface {
	Exists() bool
}

type LocalTarget struct {
	Path string
}

func (target LocalTarget) Exists() bool {
	_, err := os.Stat(target.Path)
	return !os.IsNotExist(err)
}

type Task interface {
	Requires() interface{}
	Run() error
	Output() Target
}

func In(task Task) Arbiter {
	return Arbiter{Value: task.Requires()}
}

func Out(task Task) Arbiter {
	return Arbiter{Value: task.Output()}
}

type Arbiter struct {
	Value interface{}
}

func (ar Arbiter) Path() string {
	switch v := ar.Value.(type) {
	case LocalTarget:
		return v.Path
	case []LocalTarget:
		if len(v) == 0 {
			panic("Path not allowed on empty requirements")
		}
		return v[0].Path
	case []Target:
		if len(v) == 0 {
			panic("Path not allowed on empty requirements")
		}
		for _, t := range v {
			target, ok := t.(LocalTarget)
			if !ok {
				continue
			}
			return target.Path
		}
		panic("0-length requirements")
	default:
		panic(fmt.Sprintf("Path not supported: %s", v))
	}
}

func (ar Arbiter) PathList() []string {
	switch value := ar.Value.(type) {
	case LocalTarget:
		return []string{value.Path}
	case []LocalTarget:
		var list []string
		for _, target := range value {
			list = append(list, target.Path)
		}
		return list
	case []Target:
		var list []string
		for _, t := range value {
			target, ok := t.(LocalTarget)
			if !ok {
				continue
			}
			list = append(list, target.Path)
		}
		return list
	default:
		panic(fmt.Sprintf("PathList not supported: %s", value))
	}
}

func (ar Arbiter) CreateFile() (*safefile.File, error) {
	switch value := ar.Value.(type) {
	default:
		panic(fmt.Sprintf("CreateFile not supported: %s", value))
	case LocalTarget:
		return safefile.Create(value.Path, 0644)
	}
}

func (ar Arbiter) File() (*os.File, error) {
	switch value := ar.Value.(type) {
	default:
		panic(fmt.Sprintf("File not supported: %s", value))
	case LocalTarget:
		return os.Open(value.Path)
	}
}
