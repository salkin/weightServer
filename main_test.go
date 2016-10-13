package main

import (
	"testing"
)

func TestPostInflux(t *testing.T) {
	c, err := influxConnect()
	w := Weight{
		Date:   "05/01/2016",
		Value:  78.5,
		Person: "Brian",
	}

	err = influxPost(c, w)
	if err != nil {
		t.FailNow()
	}
}

func TestGetInflux(t *testing.T) {

	c, _ := influxConnect()
	w, _ := influxGet(c, "Brian")

	for _, we := range w {
		t.Log("Date %s. Weight: %f", we.Date, we.Value)

	}
	if len(w) != 1 {
		t.Fatal("Brian has %u weights", len(w))
	}
}
