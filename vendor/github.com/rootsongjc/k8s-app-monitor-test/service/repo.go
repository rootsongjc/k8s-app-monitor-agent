package service

import (
	"errors"
	"math/rand"
	"os"
	"strings"
)

type inMemoryMetricRepository struct {
	AppMetric Metric
}

// NewRepository creates a new in-memory Metric repository
func newInMemoryRepository() *inMemoryMetricRepository {
	repo := &inMemoryMetricRepository{}
	repo.newFakeMetric()
	return repo
}

func (repo *inMemoryMetricRepository) newMetric() (err error) {
	repo.newFakeMetric()
	return err
}

func (repo *inMemoryMetricRepository) getMetric() (metric Metric, err error) {
	return repo.AppMetric, err
}

func (repo *inMemoryMetricRepository) getAppMetric(appName string) (metric Metric, err error) {
	if strings.Compare(repo.AppMetric.AppName, appName) == 0 {
		metric = repo.AppMetric
	} else {
		err = errors.New("Could not find metric in repository")
	}
	return metric, err
}

func (repo *inMemoryMetricRepository) newFakeMetric() {
	repo.AppMetric = Metric{PerformanceIndex{rand.Float64(), rand.Int63n(100), rand.Int63n(100), rand.Int63n(100), rand.Int63n(100),
		rand.Int63n(100)}, os.Getenv("HOSTNAME"), "test-app", "test-domain"}
}
