package app

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"

	"github.com/gofiber/storage/badger"
	"github.com/gofiber/storage/memory"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type AppOptions struct {
	Addr         string `yaml:"addr" json:"addr"`
	JWTSecret    string `yaml:"jwt_secret" json:"-"`
	SecureCookie bool   `yaml:"secure_cookie" json:"secure_cookie"`

	Database struct {
		Driver string `yaml:"driver" json:"driver"`
		Dsn    string `yaml:"dsn" json:"dsn"`
	} `yaml:"database" json:"database"`

	Storage struct {
		Driver string        `yaml:"driver" json:"driver"`
		Badger badger.Config `yaml:"badger" json:"badger"`
		Memory memory.Config `yaml:"memory" json:"memory"`
	} `yaml:"storage" json:"storage"`

	Logger zerolog.Logger `yaml:"-" json:"-"`
}

func (o AppOptions) userHome() string {
	s, err := os.UserHomeDir()
	if err != nil {
		o.Logger.Error().Err(err).Msg("failed to get home dir")
		return "./"
	}
	return s
}

func (o *AppOptions) Init() error {
	p := filepath.Join(o.userHome(), ".config", "skul", "skul.yaml")
	o.Logger.Debug().Str("config", p).Msg("loading config file")
	_ = o.LoadFromFile(p)
	if err := o.SetDefault(); err != nil {
		return err
	}
	if err := o.WriteToFile(p); err != nil {
		return err
	}
	return nil
}

func (o *AppOptions) LoadFromFile(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	if err := yaml.NewDecoder(file).Decode(o); err != nil {
		return err
	}
	return nil
}

func (o *AppOptions) WriteToFile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	if err := yaml.NewEncoder(file).Encode(o); err != nil {
		return err
	}
	return nil
}

func (o *AppOptions) SetDefault() error {
	if o.Addr == "" {
		o.Addr = ":8080"
	}
	if o.JWTSecret == "" {
		b := make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, b); err != nil {
			return err
		}
		o.JWTSecret = base64.RawURLEncoding.EncodeToString(b)
	}
	if o.Database.Driver == "" {
		o.Database.Driver = "sqlite3"
		o.Database.Dsn = ":memory:?cache=shared"
	}

	return nil
}
