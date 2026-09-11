package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	taskDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "processor_task_duration_seconds",
			Help: "Task execution time in seconds",
		},
		[]string{"translator"},
	)

	tasksTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "processor_tasks_total",
			Help: "Total number of processed tasks",
		},
		[]string{"translator"},
	)
)

func initMetrics() {
	prometheus.MustRegister(taskDuration)
	prometheus.MustRegister(tasksTotal)

	http.Handle("/metrics", promhttp.Handler())

	go http.ListenAndServe(":2112", nil)
}
