package coct

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Dashboard struct {
	DayZero     time.Time   `json:"dayzero"`
	City        City        `json:"city"`
	Dams        Dams        `json:"dams"`
	CapeTonians CapeTonians `json:"capetonians"`
	Other       []Project   `json:"other"`
	Timestamp   time.Time   `json:"timestamp"`
}

type Project struct {
	Description string `json:"description"`
	Percentage  int    `json:"percentage"`
	Status      int    `json:"status"`
}

type Trend struct {
	Amount    int  `json:"amount"`
	Direction bool `json:"direction"`
}

type Dams struct {
	Level float64 `json:"level"`
	Trend Trend   `json:"trend"`
}

type CapeTonians struct {
	Amount float64 `json:"amount"`
	Trend  Trend   `json:"trend"`
}

type City struct {
	Progress int       `json:"progress"`
	Projects []Project `json:"projects"`
}

func Get() (io.Reader, error) {
	var client = &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Get("http://coct.co/water-dashboard/")
	if err != nil {
		return bytes.NewReader([]byte("")), err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return bytes.NewReader(body), nil
}

func Parse(r io.Reader) (Dashboard, error) {
	var d Dashboard
	d.Timestamp = time.Now()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return d, err
	}

	// Day Zero

	dayZero, err := getDayZero(doc)
	if err == nil {
		d.DayZero = dayZero
	}

	// Dam Level

	level, err := getDamLevel(doc)
	if err == nil {
		d.Dams.Level = level
	}

	// CapeTonian Amount

	amount, err := getCapeTonianAmount(doc)
	if err == nil {
		d.CapeTonians.Amount = amount
	}

	return d, nil
}

func getCapeTonianAmount(doc *goquery.Document) (float64, error) {
	amountS := doc.Find(".percentage_label").Eq(2).Text()
	amountS = strings.Replace(amountS, "%", "", -1)

	return strconv.ParseFloat(amountS, 64)
}

func getDamLevel(doc *goquery.Document) (float64, error) {
	levelS := doc.Find(".percentage_label").Eq(1).Text()
	levelS = levelS[0:4]

	return strconv.ParseFloat(levelS, 64)
}

func getDayZero(doc *goquery.Document) (time.Time, error) {
	h3 := doc.Find("h3").First().Text()
	h3 = strings.TrimSpace(h3)
	h3 = strings.Replace(h3, " ", "", -1)

	dayS := h3[0:2]
	monthS := h3[3:5]
	yearS := h3[6:10]

	day, err := strconv.Atoi(dayS)
	if err != nil {
		return time.Now(), err
	}

	month, err := strconv.Atoi(monthS)
	if err != nil {
		return time.Now(), err
	}

	year, err := strconv.Atoi(yearS)
	if err != nil {
		return time.Now(), err
	}

	loc, err := time.LoadLocation("Africa/Johannesburg")
	if err != nil {
		return time.Now(), err
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc), nil
}
