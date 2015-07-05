README
======

* IO arbiter: `Input(task)`, `Inputs(task)`, `Target` embeds `io.ReadWriteCloser`
* shellout with [clam](http://github.com/miku/clam)
* CLI generator
* TSV writer and more open, read, write utils
* parallel task builder

Observations and premises
-------------------------

A few observations on design choices.

Use plain types for parameters
------------------------------

Task parameters must be expressable as string, because on the command line
they start their life as strings. So use the base types for parameters, too,
like `string`, `int`, `int64` or `float64`.

A date will start life as string. It is more convenient to write:

    SomeTask{Date: "2010-01-01"}

than

    t, err := time.Parse("2006-01-02", "2010-01-01")
    if err != nil {
        log.Fatal(err)
    }
    SomeTask{Date: t}

For any date algebra, write helpers that can handle date strings (TODO):

    dr := dateRange.Weekly{From: "2010-01-01", To: "2011-01-01"}
    for _, date := range dr {
        log.Println(date)
    }

One downside of plain types is the need to write extra tags for validation.
For example: a string parameter must be a valid date. Or that the number of
workers must be a positive integer, ...

Probably pack this validation into struct tags as well and hope the magic ends here.

    type MyTask struct {
        Date     string `valid:"Date"`
        DateHour string `valid:"DateHour"`
        Worker   int    `valid:"NonZeroInteger"`
    }

Then again, these validations would not be expressable without a lot of
different types, whereas one can always add validation rules.
