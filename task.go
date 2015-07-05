package lynd

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/miku/structs"
)

// funcMap maps string keys to functions from interface{} to (string, error).
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
		d := time.Duration(-weekday) * 24 * time.Hour
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

// SetDefaults evaluates the default struct tag if the field has the zero
// value. A pointer to a task must be passed in, since they this methods
// potentially alters field values. Subsequent calls to SetDefaults should not
// change the task, since any zero value has been filled on the first call or
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
		var err error
		for key, f := range defaultsFuncs {
			if v == key {
				v, err = f(v)
				if err != nil {
					return err
				}
				break
			}
		}
		return setFieldValue(field, v)
	}
	return nil
}

// Adjust performs further (and final) adjustments to task parameters. While
// SetDefaults will set default values, only if no value was given for a task,
// Adjust will alter existing values. For example, dates can be mapped to
// certain intervals.
func Adjust(task Task) error {
	s := structs.New(task)
	for _, field := range s.Fields() {
		v := field.Tag("adjust")
		if v == "" {
			continue
		}
		var err error
		for key, f := range adjustFuncs {
			if v == key {
				v, err = f(field.Value())
				if err != nil {
					return err
				}
				break
			}
		}
		return setFieldValue(field, v)
	}
	return nil
}
