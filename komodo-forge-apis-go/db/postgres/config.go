package postgres

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"time"
)

type Config struct {
	ConnString          string
	Host                string
	Port                int
	User                string
	Password            string
	Database            string
	SSLMode             string
	ApplicationName     string
	ConnectTimeout      time.Duration
	MaxConns            int32
	MinConns            int32
	MaxConnLifetime     time.Duration
	MaxConnIdleTime     time.Duration
	HealthCheckPeriod   time.Duration
	StatementCacheCap   int
	TLSConfig           *tls.Config
	PreferSimpleProtocol bool
}

func (cfg Config) connectionString() (string, error) {
 	if cfg.ConnString != "" {
 		return cfg.ConnString, nil
 	}

 	if cfg.Host == "" {
 		return "", errors.New("postgres host is required")
 	}
 	if cfg.Port == 0 {
 		cfg.Port = 5432
 	}
 	if cfg.User == "" {
 		return "", errors.New("postgres user is required")
 	}
 	if cfg.Database == "" {
 		return "", errors.New("postgres database is required")
 	}

 	url := &url.URL{
 		Scheme: "postgres",
 		User:   url.UserPassword(cfg.User, cfg.Password),
 		Host:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
 		Path:   cfg.Database,
 	}

 	query := url.Query()
 	if cfg.SSLMode != "" {
 		query.Set("sslmode", cfg.SSLMode)
 	}
 	if cfg.ApplicationName != "" {
 		query.Set("application_name", cfg.ApplicationName)
 	}
 	url.RawQuery = query.Encode()

 	return url.String(), nil
}
