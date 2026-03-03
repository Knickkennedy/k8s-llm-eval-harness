package benchmarks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadSuite loads a benchmark suite from a YAML file
func LoadSuite(path string) (Suite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Suite{}, fmt.Errorf("reading benchmark file %s: %w", path, err)
	}

	var suite Suite
	if err := yaml.Unmarshal(data, &suite); err != nil {
		return Suite{}, fmt.Errorf("parsing benchmark file %s: %w", path, err)
	}

	if suite.Name == "" {
		return Suite{}, fmt.Errorf("benchmark file %s missing required field: name", path)
	}

	if suite.Category == "" {
		return Suite{}, fmt.Errorf("benchmark file %s missing required field: category", path)
	}

	return suite, nil
}

// LoadSuitesFromDir loads all YAML benchmark suites from a directory
func LoadSuitesFromDir(dir string) ([]Suite, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading benchmark directory %s: %w", dir, err)
	}

	var suites []Suite
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".yaml") && !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		suite, err := LoadSuite(path)
		if err != nil {
			return nil, fmt.Errorf("loading suite from %s: %w", path, err)
		}

		suites = append(suites, suite)
	}

	if len(suites) == 0 {
		return nil, fmt.Errorf("no benchmark YAML files found in %s", dir)
	}

	return suites, nil
}
