package acceptance

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Scenario struct {
	ID      string
	Feature string
	Name    string
	Tags    []string
	Steps   []string
}

func LoadScenarios(dir string) ([]Scenario, error) {
	root, err := repoRoot()
	if err != nil {
		return nil, err
	}
	files, err := filepath.Glob(filepath.Join(root, dir, "*.feature"))
	if err != nil {
		return nil, err
	}
	var scenarios []Scenario
	for _, file := range files {
		parsed, err := parseFeatureFile(file)
		if err != nil {
			return nil, err
		}
		scenarios = append(scenarios, parsed...)
	}
	return scenarios, nil
}

func RequireScenario(dir string, id string) (Scenario, error) {
	scenarios, err := LoadScenarios(dir)
	if err != nil {
		return Scenario{}, err
	}
	for _, scenario := range scenarios {
		if scenario.ID == id {
			return scenario, nil
		}
	}
	return Scenario{}, fmt.Errorf("scenario %s not found in %s", id, dir)
}

func (s Scenario) HasTag(tag string) bool {
	tag = strings.TrimPrefix(tag, "@")
	for _, candidate := range s.Tags {
		if strings.TrimPrefix(candidate, "@") == tag {
			return true
		}
	}
	return false
}

func parseFeatureFile(path string) ([]Scenario, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var feature string
	var pendingTags []string
	var current *Scenario
	var scenarios []Scenario
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case line == "" || strings.HasPrefix(line, "#"):
			continue
		case strings.HasPrefix(line, "Feature:"):
			feature = strings.TrimSpace(strings.TrimPrefix(line, "Feature:"))
		case strings.HasPrefix(line, "@"):
			pendingTags = append([]string{}, strings.Fields(line)...)
		case strings.HasPrefix(line, "Scenario:"):
			if current != nil {
				scenarios = append(scenarios, *current)
			}
			current = &Scenario{
				Feature: feature,
				Name:    strings.TrimSpace(strings.TrimPrefix(line, "Scenario:")),
				Tags:    pendingTags,
				ID:      scenarioID(pendingTags),
			}
			pendingTags = nil
		case current != nil && isStep(line):
			current.Steps = append(current.Steps, line)
		}
	}
	if current != nil {
		scenarios = append(scenarios, *current)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return scenarios, nil
}

func isStep(line string) bool {
	for _, prefix := range []string{"Given ", "When ", "Then ", "And ", "But "} {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

func scenarioID(tags []string) string {
	for _, tag := range tags {
		value := strings.TrimPrefix(tag, "@")
		if strings.HasPrefix(value, "SCN-") {
			return value
		}
	}
	return ""
}

func repoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		next := filepath.Dir(dir)
		if next == dir {
			return "", fmt.Errorf("repo root not found from %s", dir)
		}
		dir = next
	}
}
