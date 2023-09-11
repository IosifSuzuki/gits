package model

import (
	"fmt"
	"gits/internal/model/constant"
	"gits/internal/model/errs"
	"gopkg.in/yaml.v3"
	"net"
	"os"
	"time"
)

const confPath = "configs/config.yml"

type Config struct {
	Server  Server `yaml:"server"`
	Storage Storage
	Cache   Cache  `yaml:"cache"`
	AWS     AWS    `yaml:"aws"`
	IpInfo  IpInfo `yaml:"ip_info"`
}

type Server struct {
	Host      string
	Port      string
	SecretKey string `yaml:"secret_key"`
}

type Storage struct {
	Username string
	Password string
	Name     string
	Host     string
	Port     string
	SSLMode  string
}

type Cache struct {
	Host       string
	Port       string
	Password   string
	SessionTTL time.Duration `yaml:"session_ttl"`
}

type AWS struct {
	S3Bucket string `yaml:"s3_bucket"`
	S3Region string `yaml:"s3_region"`
}

type IpInfo struct {
	Token string `yaml:"token"`
}

func NewConfig() (*Config, error) {
	var config Config
	file, err := os.OpenFile(confPath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, config.SetValueFromENV()
}

func (c *Config) SetValueFromENV() error {
	if host, ok := os.LookupEnv(constant.HostServerKey); ok {
		c.Server.Host = host
	} else {
		return errs.ConfiguredBadEnvError
	}
	if port, ok := os.LookupEnv(constant.HostPortKey); ok {
		c.Server.Port = port
	} else {
		return errs.ConfiguredBadEnvError
	}

	if username, ok := os.LookupEnv(constant.DataBaseUserKey); ok {
		c.Storage.Username = username
	} else {
		return errs.ConfiguredBadEnvError
	}
	if password, ok := os.LookupEnv(constant.DataBasePassword); ok {
		c.Storage.Password = password
	} else {
		return errs.ConfiguredBadEnvError
	}
	if name, ok := os.LookupEnv(constant.DataBaseName); ok {
		c.Storage.Name = name
	} else {
		return errs.ConfiguredBadEnvError
	}
	if host, ok := os.LookupEnv(constant.DataBaseHostKey); ok {
		c.Storage.Host = host
	} else {
		return errs.ConfiguredBadEnvError
	}
	if port, ok := os.LookupEnv(constant.DataBasePortKey); ok {
		c.Storage.Port = port
	} else {
		return errs.ConfiguredBadEnvError
	}
	if sslMode, ok := os.LookupEnv(constant.DataBaseSSLMode); ok {
		c.Storage.SSLMode = sslMode
	} else {
		return errs.ConfiguredBadEnvError
	}

	if host, ok := os.LookupEnv(constant.CacheHostKey); ok {
		c.Cache.Host = host
	} else {
		return errs.ConfiguredBadEnvError
	}
	if port, ok := os.LookupEnv(constant.CachePortKey); ok {
		c.Cache.Port = port
	} else {
		return errs.ConfiguredBadEnvError
	}
	if password, ok := os.LookupEnv(constant.CachePasswordKey); ok {
		c.Cache.Password = password
	} else {
		return errs.ConfiguredBadEnvError
	}
	return nil
}

func (s *Server) Address() string {
	return net.JoinHostPort(s.Host, s.Port)
}

func (s *Storage) Address() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Europe/London",
		s.Host,
		s.Username,
		s.Password,
		s.Name,
		s.Port,
		s.SSLMode,
	)
}

func (c *Cache) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
