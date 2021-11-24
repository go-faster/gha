package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/go-faster/errors"
)

const layout = "2006-01-02-15"

type DateVar struct {
	Date *time.Time
}

func (d DateVar) String() string {
	if d.Date.IsZero() {
		return ""
	}
	return d.Date.Format(layout)
}

func (d DateVar) Set(s string) error {
	v, err := time.ParseInLocation(layout, s, time.UTC)
	if err != nil {
		return errors.Wrap(err, "parse time")
	}

	*d.Date = v
	return nil
}

func dateFlag(date *time.Time, name, usage string) {
	flag.Var(DateVar{Date: date}, name, usage)
}

func GetURL(date time.Time) string {
	return fmt.Sprintf("https://data.gharchive.org/%s.json.gz",
		date.Format(layout),
	)
}

func main() {
	var arg struct {
		Date time.Time
	}
	dateFlag(&arg.Date, "date", "date to download")
	flag.Parse()

	link := GetURL(arg.Date)
	start := time.Now()
	res, err := http.Head(link)
	if err != nil {
		panic(err)
	}
	fmt.Println("HEAD", link, res.Status, time.Since(start).Round(time.Millisecond))
}
