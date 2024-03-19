package okconf

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

var _ Config = (*ExampleConfig)(nil)

func TestLoadDefault(t *testing.T) {
	cfg := Load[ExampleConfig]()

	if cfg.Name != "Foo" {
		t.Fail()
	}

	if cfg.StartedAt.IsZero() {
		t.Fail()
	}

	if cfg.Nested.Count != 100 {
		t.Fail()
	}
}

func TestFromJSON(t *testing.T) {
	json := "{\"started_at\":\"2024-01-01T13:37:00.000000Z\",\"nested\":{\"Count\":200}}"
	f, err := os.CreateTemp("", "test-conf-*.json")
	if err != nil {
		t.Error(err)
	}

	defer f.Close()

	f.WriteString(json)

	p, err := filepath.Abs(f.Name())
	if err != nil {
		t.Error(err)
	}

	cfg, err := FromJSON[ExampleConfig](p)
	if err != nil {
		t.Error(err)
	}

	if cfg.Name != "Foo" {
		t.Fail()
	}

	if cfg.StartedAt.String() != "2024-01-01 13:37:00 +0000 UTC" {
		t.Fail()
	}

	if cfg.Nested.Count != 200 {
		t.Fail()
	}
}

func TestSaveJSON(t *testing.T) {
	cfg := Load[ExampleConfig]()

	tmp, err := os.CreateTemp("", "save-conf-*.json")
	if err != nil {
		t.Error(err)
	}

	if err := SaveJSON(cfg, tmp.Name()); err != nil {
		t.Error(err)
	}
}

func TestFromYAML(t *testing.T) {
	yml := "name: Bar\nstartedAt: 2024-01-01T13:37:00.000000Z\nnested:\n    count: 100\n"
	f, err := os.CreateTemp("", "test-conf-*.yml")
	if err != nil {
		t.Error(err)
	}

	defer f.Close()

	f.WriteString(yml)

	p, err := filepath.Abs(f.Name())
	if err != nil {
		t.Error(err)
	}

	cfg, err := FromYAML[ExampleConfig](p)
	if err != nil {
		t.Error(err)
	}

	if cfg.Name != "Bar" {
		t.Fail()
	}

	if cfg.StartedAt.String() != "2024-01-01 13:37:00 +0000 UTC" {
		t.Fail()
	}

	if cfg.Nested.Count != 100 {
		t.Fail()
	}
}

func TestSaveYaml(t *testing.T) {
	cfg := Load[ExampleConfig]()

	tmp, err := os.CreateTemp("", "save-conf-*.yaml")
	if err != nil {
		t.Error(err)
	}

	if err := SaveYAML(cfg, tmp.Name()); err != nil {
		t.Error(err)
	}
}

type ExampleConfig struct {
	Name      string       `yaml:"name" json:"name"`
	StartedAt time.Time    `yaml:"startedAt" json:"started_at"`
	Nested    NestedConfig `yaml:"nested" json:"nested"`
}

type NestedConfig struct {
	Count int `yaml:"count" json:"count"`
}

func (e ExampleConfig) Default() Config {
	return ExampleConfig{
		Name:      "Foo",
		StartedAt: time.Now().UTC(),
		Nested: NestedConfig{
			Count: 100,
		},
	}
}
