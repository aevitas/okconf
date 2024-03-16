package okconf

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	errPathNotAbs     = errors.New("specified path is not an absolute path")
	errBuffNotWritten = errors.New("could not write full buffer to file; config may be corrupt on disk")
)

// loads the default configuration as returned by config.Default
func Load[T Config]() *T {
	return defaultCfg[T]()
}

// loads config from the specified JSON file, using the default configuration as a starting point
func FromJSON[T Config](file string) (*T, error) {
	fp, err := filepath.Abs(file)
	if err != nil {
		return new(T), err
	}

	f, err := os.Open(fp)
	if err != nil {
		return new(T), err
	}

	defer f.Close()

	return fromJSONStream[T](defaultCfg[T](), f)
}

// saves the
func SaveJSON[T Config](cfg T, file string) error {
	if !filepath.IsAbs(file) {
		return errPathNotAbs
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	nw, err := f.Write(buf)
	if err != nil {
		return err
	}

	if nw != len(buf) {
		return errBuffNotWritten
	}

	return nil
}

func fromJSONStream[T Config](conf *T, reader io.Reader) (*T, error) {
	dec := json.NewDecoder(reader)
	err := dec.Decode(&conf)
	if err != nil {
		return new(T), err
	}

	return conf, nil
}

func FromYAML[T Config](file string) (*T, error) {
	fp, err := filepath.Abs(file)
	if err != nil {
		return new(T), err
	}

	f, err := os.Open(fp)
	if err != nil {
		return new(T), err
	}

	defer f.Close()

	return fromYAMLStream(defaultCfg[T](), f)
}

func fromYAMLStream[T Config](cfg *T, r io.Reader) (*T, error) {
	dec := yaml.NewDecoder(r)

	if err := dec.Decode(cfg); err != nil {
		return new(T), err
	}

	return cfg, nil
}

func defaultCfg[T Config]() *T {
	var cfg T
	cfg = cfg.Default().(T)

	return &cfg
}

type Config interface {
	Default() Config
}
