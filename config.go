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

// loads the default state for configuration T, as returned by the configuration's Default method
func Load[T Config]() *T {
	return defaultCfg[T]()
}

// loads configuration of type T from the specified JSON file, applying it as a layer on top of T's Default configuration
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

// flattens and saves the configuration to the specified JSON file
func SaveJSON[T Config](cfg T, file string) error {
	buf, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := saveFile(file, buf); err != nil {
		return err
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

// loads configuration of type T from the specified YAML file, applying it as a layer on top of T's Default configuration
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

func saveFile(file string, buf []byte) error {
	if !filepath.IsAbs(file) {
		return errPathNotAbs
	}

	f, err := os.Create(file)
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

// flattens and saves the specified configuration to a YAML file
func SaveYAML[T Config](cfg T, file string) error {
	buf, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := saveFile(file, buf); err != nil {
		return err
	}

	return nil
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

// base interface for okconf configuration
type Config interface {
	Default() Config
}
