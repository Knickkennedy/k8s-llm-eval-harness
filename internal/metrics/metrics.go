package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EvalPassRate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "llm_eval_pass_rate",
		Help: "Pass rate for each benchmark category (0.0 to 1.0)",
	}, []string{"category", "model"})

	EvalScore = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "llm_eval_score",
		Help: "Score for each individual benchmark",
	}, []string{"benchmark_id", "category", "model"})

	EvalDurationMs = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "llm_eval_duration_ms",
		Help: "Duration in milliseconds for each benchmark run",
	}, []string{"benchmark_id", "category", "model"})

	EvalRunTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "llm_eval_runs_total",
		Help: "Total number of evaluation runs",
	}, []string{"category", "model", "result"})

	EvalLastRunTimestamp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "llm_eval_last_run_timestamp",
		Help: "Unix timestamp of the last evaluation run",
	}, []string{"category", "model"})
)
