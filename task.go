package lynd

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/miku/structs"
)

// funcMap maps string keys to functions from string to (string, error).
type funcMap map[string]func(string) (string, error)

// defaultFuncs contain custom functions, that may be invoked during defaults
// evaluation. Most useful example might be a "today". Maybe all special
// methods should carry a prefix, like `s:today`.
var defaultsFuncs = funcMap{
	"today": func(value string) (string, error) {
		return time.Now().Format("2006-01-02"), nil
	},
	"yesterday": func(value string) (string, error) {
		return time.Now().Add(-24 * time.Hour).Format("2006-01-02"), nil
	},
	// more here? ...
	//
	// "weekly"
	// "weekly:<weekday>"
	// "random"
	// "func:someCustomFuncOnStruct"
	//
	// ...
	//
	// custom func could be circumvented with a generic content based filepath
	// e.g.
	//      allow AutoTargets to be generated from the output of a task
	//      task.Run will be always called, but the task may report completion, if nothing has changed
	//      e.g. a ftp mirror and the output of the task is just a list of mirrored files.
	//
}

// setDefaults evaluates the default struct tag if the field has the zero
// value. A pointer to a task must be passed in, since they this methods
// potentially alters field values.
func setDefaults(task Task) error {
	s := structs.New(task)
	for _, field := range s.Fields() {
		if !field.IsZero() {
			continue
		}
		v := field.Tag("default")
		if v == "" {
			continue
		}
		var err error
		for key, f := range defaultsFuncs {
			if v == key {
				v, err = f(v)
				if err != nil {
					return err
				}
			}
		}
		switch field.Kind() {
		case reflect.String:
			err := field.Set(v)
			if err != nil {
				return err
			}
		case reflect.Int, reflect.Int64:
			i, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			err = field.Set(i)
			if err != nil {
				return err
			}
		case reflect.Float64:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}
			err = field.Set(f)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("cannot set default for %s", field.Kind())
		}
	}
	return nil
}
