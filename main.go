package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	metric "github.com/rootsongjc/k8s-app-monitor-test/service"
	chart "github.com/wcharczuk/go-chart"
)

var m metric.Metric

func main() {
	listenPort := fmt.Sprintf(":%s", listenPort())
	fmt.Printf("Listening on %s\n", listenPort)
	http.HandleFunc("/", drawChart)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Get : %v", err)
	}
}
func drawChart(res http.ResponseWriter, req *http.Request) {
	port := os.Getenv("APP_PORT")
	service := os.Getenv("SERVICE_NAME")
	if len(port) == 0 {
		port = "3000"
	}
	if len(service) == 0 {
		service = "localhost"
	}
	resp, err := http.Get("http://" + service + ":" + port + "/metrics")
	checkError(err)
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&m)
	checkError(err)
	sbc := chart.BarChart{
		Title: "AppName:" + m.AppName + "\nDomain:" + m.Domain + "\nHost:" + m.Host,
		TitleStyle: chart.Style{
			Show:                true,
			TextHorizontalAlign: 1,
		},
		Height:   512,
		BarWidth: 60,
		XAxis: chart.Style{
			Show: true,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Bars: []chart.Value{
			{Value: m.FailRatio, Label: "FailRatio"},
			{Value: float64(m.FailAmount), Label: "FailAmount"},
			{Value: float64(m.AccessAmount), Label: "AccessAmount"},
			{Value: float64(m.MaxConcurrent), Label: "MaxConcurrent"},
			{Value: float64(m.MinLatency), Label: "MinLatency"},
			{Value: float64(m.AvgLatency), Label: "AvgLatency"},
		},
	}

	res.Header().Set("Content-Type", "image/png")
	err = sbc.Render(chart.PNG, res)
	if err != nil {
		fmt.Printf("Error rendering chart: %v\n", err)
	}
}

func listenPort() string {
	if len(os.Getenv("PORT")) > 0 {
		return os.Getenv("PORT")
	}
	return "8888"
}
