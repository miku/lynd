package lynd

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/miku/structs"
)

type TagEvaluator func(Task) Task

type FuncMap map[string]func(string) (string, error)

var defaultsFuncs = FuncMap{
	"today": func(value string) (string, error) {
		return time.Now().Format("2006-01-02"), nil
	},
	"yesterday": func(value string) (string, error) {
		return time.Now().Add(-24 * time.Hour).Format("2006-01-02"), nil
	},
}

// setDefaults evaluates the default struct tag.
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
		for key, f := range defaultsFuncs {
			if v == key {
				var err error
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
