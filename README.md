# k8s-llm-eval-harness

A production-grade, GitOps-managed LLM evaluation harness running on Kubernetes. Automatically benchmarks LLM behavior across safety, reasoning, and factuality dimensions — exposing results as Prometheus metrics and visualizing model quality trends over time in Grafana.

## Why This Exists

Deploying a model is only half the problem. Knowing whether it's behaving correctly, safely, and consistently is the other half. This harness treats model evaluation as an infrastructure problem — scheduled, automated, observable, and declarative.

## Stack
- **Inference**: [k8s-llm-inference-stack](https://github.com/knickkennedy/k8s-llm-inference-stack) (Ollama + Mistral 7B)
- **Eval Runner**: Custom Go binary running benchmark suites on a schedule
- **Metrics**: Custom Prometheus exporter surfacing eval scores
- **Orchestration**: Kubernetes CronJob — runs evals, exits, repeats
- **GitOps**: ArgoCD with Kustomize environment overlays
- **Observability**: Grafana dashboard showing model quality over time

## Benchmark Categories
- **Safety**: Tests refusal behavior on harmful prompts, consistency under adversarial rephrasing
- **Reasoning**: Logic, multi-step problem solving, and chain-of-thought accuracy
- **Factuality**: Factual accuracy and hallucination detection

## Getting Started
See [Bootstrap Guide](bootstrap/README.md) for setup instructions.

## Architecture
Coming soon.
