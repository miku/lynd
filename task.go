package lynd

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
)

// funcMap maps string keys to functions from interface{} to (string, error).
// TODO(miku): make types for string and field parameter funcs?
type funcMap map[string]func(interface{}) (string, error)

// defaultFuncs contain custom functions, that may be invoked during defaults
// evaluation. Most useful example might be a "today".
var defaultsFuncs = funcMap{
	"today": func(_ interface{}) (string, error) {
		return time.Now().Format("2006-01-02"), nil
	},
	"yesterday": func(_ interface{}) (string, error) {
		return time.Now().Add(-24 * time.Hour).Format("2006-01-02"), nil
	},
}

// adjustFunc may alter already existing field values.
var adjustFuncs = funcMap{
	"weekly": func(v interface{}) (string, error) {
		s, ok := v.(string)
		if !ok {
			return "", fmt.Errorf("must be a string")
		}
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return "", err
		}
		weekday := shiftWeekday(int(t.Weekday()))
		d := time.Duration(-weekday+1) * 24 * time.Hour
		return t.Add(d).Format("2006-01-02"), nil
	},
}

// shiftWeekday makes a week begin on Monday.
func shiftWeekday(weekday int) int {
	if weekday == 0 {
		return 7
	}
	return weekday - 1
}

// setFieldValue will try to set the value of a given field to v, which is
// given as string and converted to the field type as required.
func setFieldValue(field *structs.Field, v string) error {
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
		return fmt.Errorf("cannot set field value for field kind %s", field.Kind())
	}
	return nil
}

// SetDefaults evaluates the default struct tag only if the field has the zero
// value. A pointer to a task must be passed in, since this method potentially
// alters field values. Subsequent calls to SetDefaults will not change the
// task, since either any zero value has been filled on the first call or
// SetDefaults returned an error.
func SetDefaults(task Task) error {
	s := structs.New(task)
	for _, field := range s.Fields() {
		if !field.IsZero() {
			continue
		}
		v := field.Tag("default")
		if v == "" {
			continue
		}
		if f, ok := defaultsFuncs[v]; ok {
			var err error
			v, err = f(v)
			if err != nil {
				return err
			}
		}
		return setFieldValue(field, v)
	}
	return nil
}

// Adjust performs further (and final) adjustments to task parameters. While
// SetDefaults will set default values, only if no value was given for a
// parameter, Adjust will alter only existing values. For example, dates can
// be mapped to certain intervals.
func Adjust(task Task) error {
	s := structs.New(task)
	for _, field := range s.Fields() {
		v := field.Tag("adjust")
		if v == "" {
			continue
		}
		var err error
		if f, ok := adjustFuncs[v]; ok {
			v, err = f(field.Value())
			if err != nil {
				return err
			}
			return setFieldValue(field, v)
		}
	}
	return nil
}

func Init(task Task) error {
	if err := SetDefaults(task); err != nil {
		return err
	}
	return Adjust(task)
}

// ParameterMap returns a map of the significant parameters for a task.
func ParameterMap(task Task) map[string]string {
	s := structs.New(task)
	m := make(map[string]string)
	for _, f := range s.Fields() {
		if !f.IsExported() || f.Tag("significant") == "false" {
			continue
		}
		switch f.Kind() {
		case reflect.String, reflect.Int, reflect.Int64, reflect.Float64:
			m[f.Name()] = fmt.Sprintf("%s", f.Value())
		default:
			continue
		}
	}
	return m
}

// TODO(miku): replace with something real.
func mapToSlug(m map[string]string) string {
	var parts []string
	for k, v := range m {
		parts = append(parts, k)
		parts = append(parts, v)
	}
	return strings.Join(parts, "-")
}

// pkgName returns the lowest level package name or panic, if this name cannot
// be determined. TODO(miku): why not use the whole pkg hierarchy as
// locations?
func pkgName(t reflect.Type) string {
	parts := strings.Split(t.PkgPath(), "/")
	if len(parts) == 0 {
		panic("invalid pkg path")
	}
	return parts[len(parts)-1]

}

// TaskID returns a string, that uniquely identifies a task. The ID will
// consist of the task name (its type) and a slugified version of its
// significant parameters.
func TaskID(task Task) string {
	t := reflect.TypeOf(task)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	pmap := ParameterMap(task)
	return fmt.Sprintf("%s/%s/%s", pkgName(t), t.Name(), strings.ToLower(mapToSlug(pmap)))
}
