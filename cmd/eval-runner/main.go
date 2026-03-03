package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/knickkennedy/k8s-llm-eval-harness/internal/benchmarks"
	"github.com/knickkennedy/k8s-llm-eval-harness/internal/metrics"
	"github.com/knickkennedy/k8s-llm-eval-harness/internal/ollama"
	"github.com/knickkennedy/k8s-llm-eval-harness/internal/scorer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ollamaURL    := getEnv("OLLAMA_HOST", "http://ollama.ollama:11434")
	model        := getEnv("OLLAMA_MODEL", "mistral")
	metricsPort  := getEnv("METRICS_PORT", "9091")
	benchmarkDir := getEnv("BENCHMARK_DIR", "./evaluations/benchmarks")

	log.Printf("Starting eval runner against %s model=%s", ollamaURL, model)
	log.Printf("Loading benchmarks from %s", benchmarkDir)

	// Load suites from YAML files at runtime
	suites, err := benchmarks.LoadSuitesFromDir(benchmarkDir)
	if err != nil {
		log.Fatalf("Failed to load benchmark suites: %v", err)
	}
	log.Printf("Loaded %d benchmark suites", len(suites))

	client := ollama.NewClient(ollamaURL)

	// Health check Ollama before running evals
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := client.HealthCheck(ctx); err != nil {
		log.Fatalf("Ollama health check failed: %v", err)
	}
	log.Println("Ollama health check passed")

	// Start metrics server in background
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Printf("Metrics server listening on :%s", metricsPort)
		if err := http.ListenAndServe(":"+metricsPort, nil); err != nil {
			log.Fatalf("Metrics server failed: %v", err)
		}
	}()

	// Run all suites
	for _, suite := range suites {
		result := runSuite(client, suite, model)
		recordMetrics(result, model)
		printSuiteResult(result)
	}

	log.Println("All evaluations complete")

	// Keep metrics server alive so Prometheus can scrape before pod exits
	time.Sleep(60 * time.Second)
}

func runSuite(client *ollama.Client, suite benchmarks.Suite, model string) benchmarks.SuiteResult {
	log.Printf("Running suite: %s (%d benchmarks)", suite.Name, len(suite.Benchmarks))

	result := benchmarks.SuiteResult{
		Suite:      suite,
		TotalCount: len(suite.Benchmarks),
	}

	var totalDuration int64

	for _, benchmark := range suite.Benchmarks {
		log.Printf("  Running benchmark: %s", benchmark.ID)

		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
		start := time.Now()

		resp, err := client.Generate(ctx, model, benchmark.Prompt)
		elapsed := time.Since(start)
		cancel()

		var benchResult benchmarks.Result
		if err != nil {
			log.Printf("  ERROR: %v", err)
			benchResult = benchmarks.Result{
				Benchmark:    benchmark,
				Passed:       false,
				Score:        0.0,
				DurationMs:   elapsed.Milliseconds(),
				ErrorMessage: err.Error(),
			}
		} else {
			benchResult = scorer.Score(benchmark, resp, elapsed)
		}

		result.Results = append(result.Results, benchResult)
		if benchResult.Passed {
			result.PassCount++
		}
		totalDuration += benchResult.DurationMs

		status := "PASS"
		if !benchResult.Passed {
			status = "FAIL"
		}
		log.Printf("  %s [%s] score=%.2f duration=%dms", status, benchmark.ID, benchResult.Score, benchResult.DurationMs)
	}

	if result.TotalCount > 0 {
		result.PassRate = float64(result.PassCount) / float64(result.TotalCount)
		result.AvgDurationMs = totalDuration / int64(result.TotalCount)
	}

	return result
}

func recordMetrics(result benchmarks.SuiteResult, model string) {
	category := string(result.Suite.Category)
	now := float64(time.Now().Unix())

	metrics.EvalPassRate.WithLabelValues(category, model).Set(result.PassRate)
	metrics.EvalLastRunTimestamp.WithLabelValues(category, model).Set(now)

	for _, r := range result.Results {
		metrics.EvalScore.WithLabelValues(r.Benchmark.ID, category, model).Set(r.Score)
		metrics.EvalDurationMs.WithLabelValues(r.Benchmark.ID, category, model).Set(float64(r.DurationMs))

		resultLabel := "pass"
		if !r.Passed {
			resultLabel = "fail"
		}
		metrics.EvalRunTotal.WithLabelValues(category, model, resultLabel).Inc()
	}
}

func printSuiteResult(result benchmarks.SuiteResult) {
	fmt.Printf("\n=== %s Suite Results ===\n", result.Suite.Name)
	fmt.Printf("Pass Rate: %.1f%% (%d/%d)\n", result.PassRate*100, result.PassCount, result.TotalCount)
	fmt.Printf("Avg Duration: %dms\n", result.AvgDurationMs)
	for _, r := range result.Results {
		status := "✓"
		if !r.Passed {
			status = "✗"
		}
		fmt.Printf("  %s [%s] score=%.2f\n", status, r.Benchmark.ID, r.Score)
	}
	fmt.Println()
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
