package apachetomcat

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

// Config represents the receiver config settings within the collectors config.yaml
type Config struct {
	Endpoint string `mapstructure:"endpoint"`
	Interval string `mapstructure:"interval"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// Validate checks if the receiver configuration is valid
func (c *Config) Validate() error {
	if c.Endpoint == "" {
		return errors.New("no endpoint was provided")
	}

	var err error
	res, err := url.Parse(c.Endpoint)
	if err != nil {
		return fmt.Errorf("unable to parse endpoint %s: %w", c.Endpoint, err)
	}

	interval, err := time.ParseDuration(c.Interval)
	if err != nil {
		return fmt.Errorf("unable to parse interval %s: %w", c.Interval, err)
	}

	if interval.Seconds() < 1 {
		err = errors.Join(err, errors.New("when defined, the interval has to be set to at least 1 second (1s)"))
	}

	if c.Username == "" {
		err = errors.Join(err, errors.New("username not provided and is required"))
	}

	if c.Password == "" {
		err = errors.Join(err, errors.New("password not provided and is required"))
	}

	if res.Scheme != "http" && res.Scheme != "https" {
		err = errors.Join(err, errors.New("endpoint url scheme must be http or https"))
	}

	return err
}
