package benchmarks

// Category represents a benchmark category
type Category string

const (
	CategorySafety     Category = "safety"
	CategoryReasoning  Category = "reasoning"
	CategoryFactuality Category = "factuality"
)

// Benchmark represents a single evaluation prompt with expected behavior
type Benchmark struct {
	ID                string   `yaml:"id"`
	Category          Category `yaml:"category"`
	Prompt            string   `yaml:"prompt"`
	ShouldRefuse      bool     `yaml:"should_refuse"`
	ExpectedKeywords  []string `yaml:"expected_keywords,omitempty"`
	ForbiddenKeywords []string `yaml:"forbidden_keywords,omitempty"`
	Description       string   `yaml:"description"`
}

// Result represents the outcome of running a single benchmark
type Result struct {
	Benchmark     Benchmark
	Response      string
	Passed        bool
	Score         float64
	DurationMs    int64
	ErrorMessage  string
}

// Suite represents a collection of benchmarks
type Suite struct {
	Name       string      `yaml:"name"`
	Category   Category    `yaml:"category"`
	Benchmarks []Benchmark `yaml:"benchmarks"`
}

// SuiteResult represents the outcome of running a full suite
type SuiteResult struct {
	Suite         Suite
	Results       []Result
	PassCount     int
	TotalCount    int
	PassRate      float64
	AvgDurationMs int64
}
