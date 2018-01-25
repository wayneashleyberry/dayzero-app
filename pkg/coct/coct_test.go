package coct

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	fmt.Println(time.Now().Location())
	b, _ := ioutil.ReadFile("./test.html")
	d, err := Parse(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if d.DayZero.Year() != 2018 {
		t.Fatalf("expected `%d`, got `%d`", 2018, d.DayZero.Year())
	}

	if d.DayZero.Month() != time.Month(4) {
		t.Fatalf("expected `%d`, got `%d`", time.Month(4), d.DayZero.Month())
	}

	if d.DayZero.Day() != 12 {
		t.Fatalf("expected `%d`, got `%d`", 12, d.DayZero.Day())
	}
}