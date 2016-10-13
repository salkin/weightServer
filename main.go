package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Weight struct {
	Date   string  `json: "date"`
	Value  float64 `json: "value"`
	Person string
}

func main() {

	log.SetOutput(os.Stdout)
	influxCreateDb()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", index)
	router.HandleFunc("/api/v1/weight/{name}", weight)
	log.Fatal(http.ListenAndServe(":3000", router))

}

func influxCreateDb() {
	c, _ := influxConnect()
	cmd := "CREATE DATABASE " + InfluxDB
	q := client.Query{
		Command:  cmd,
		Database: InfluxDB,
	}
	_, err := c.Query(q)
	if err != nil {
		log.Fatal("Could not create DB")
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	str := `<h1>simple REST server for daily Weight followup</h1> 
	<p />
	Register your weight through RESP API:
	<p />
	Example: <br/>
	POST /api/v1/weight/Brian  <br/>
	{ </br>
	  "date": "dd/mm/yyyy" </br> 
	  "value": 67.0 </br>
	} </br>
	<p />
	GET /api/v1/weight/Brian`
	fmt.Fprint(w, str)

}

func weight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	v := mux.Vars(r)
	name := v["name"]

	switch r.Method {
	case "GET":
		fmt.Println("Serving: ")
		c, _ := influxConnect()
		weights, _ := influxGet(c, name)
		s := `{
	"Person": "` + name + `",
	"Weights": [
`
		for i, we := range weights {
			if i != 0 {
				s += ","
			}
			s += `
			{ "Date": "` + we.Date + `",
			        "Value": ` + strconv.FormatFloat(we.Value, 'f', 1, 64) + ` 
			}`
		}
		s += `]
		}`
		fmt.Fprint(w, s)

	case "POST":
		var weight Weight
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&weight)
		if err != nil {
			http.Error(w, `{ "error": "invalid request" }`, http.StatusInternalServerError)
			log.Fatal(v)
			return
		}
		weight.Person = name
		c, err := influxConnect()
		if err != nil {
			log.Print("Could not connect to influx")
		}
		err = influxPost(c, weight)
		if err != nil {
			log.Print("Error", err)
		}
	}

}

const (
	InfluxDB = "weight"
	username = "admin"
	password = "admin"
)

func influxConnect() (client.Client, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://influx:8086",
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Fatalln("Error: ", err)
		return nil, err
	}
	return c, nil
}

func influxGet(c client.Client, person string) ([]Weight, error) {
	cmd := "SELECT date,value FROM weight WHERE person='" + person + "'"
	q := client.Query{
		Command:  cmd,
		Database: InfluxDB,
	}
	res, err := c.Query(q)
	if err != nil {
		log.Fatal("Failed query: ", err)
		return nil, err
	}
	var w []Weight

	results := res.Results[0]
	for _, row := range results.Series[0].Values {
		val, _ := row[2].(json.Number).Float64()
		weig := Weight{
			Person: person,
			Date:   row[1].(string),
			Value:  float64(val),
		}
		w = append(w, weig)
	}

	return w, nil
}

func influxPost(c client.Client, w Weight) error {
	tags := map[string]string{
		"person": w.Person,
	}
	fields := map[string]interface{}{
		"value": w.Value,
		"date":  w.Date,
	}
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  InfluxDB,
		Precision: "us",
	})
	t, err := time.Parse("02/01/2006", w.Date)
	if err != nil {
		log.Print("Invalid time given:", err)
	}
	p, err := client.NewPoint("weight",
		tags,
		fields,
		t,
	)

	bp.AddPoint(p)
	err = c.Write(bp)
	if err != nil {
		log.Print(err)
	}
	return err
}
