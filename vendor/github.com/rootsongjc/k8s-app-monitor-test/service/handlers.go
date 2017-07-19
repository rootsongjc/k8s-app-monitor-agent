package service

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func getAppMetricHandler(formatter *render.Render, repo metricRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		repo.newMetric()
		vars := mux.Vars(req)
		appName := vars["appname"]
		m, err := repo.getAppMetric(appName)
		if err != nil {
			formatter.JSON(w, http.StatusNotFound, err.Error())
		} else {
			var mdr Metric
			mdr.copyMetric(m)
			formatter.JSON(w, http.StatusOK, &mdr)
		}
	}
}

func getMetricsHandler(formatter *render.Render, repo metricRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		repo.newMetric()
		repoMetric, err := repo.getMetric()
		if err == nil {
			metric := repoMetric
			formatter.JSON(w, http.StatusOK, metric)
		} else {
			formatter.JSON(w, http.StatusNotFound, err.Error())
		}
	}
}
