package okconf

import (
	"log"
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

func TestLoadJSON(t *testing.T) {
	json := "{\"StartedAt\":\"2024-01-01T13:37:00.000000Z\",\"Nested\":{\"Count\":200}}"
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

	cfg, err := LoadJSON[ExampleConfig](p)
	if err != nil {
		t.Error(err)
	}

	if cfg.Name != "Foo" {
		t.Fail()
	}

	s := cfg.StartedAt.String()
	log.Print(s)
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

	SaveJSON(cfg, tmp.Name())
}

type ExampleConfig struct {
	Name      string
	StartedAt time.Time
	Nested    NestedConfig
}

type NestedConfig struct {
	Count int
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
