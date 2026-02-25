package config

import (
	"fmt"
	"net/url"

	"gopkg.in/yaml.v3"
)

// URL is a YAML-unmarshalable wrapper around net/url.URL.
type URL struct {
	value *url.URL
}

// ParseURL parses raw string into config URL wrapper.
func ParseURL(raw string) (URL, error) {
	if raw == "" {
		return URL{}, nil
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return URL{}, err
	}

	return URL{value: parsed}, nil
}

// MustParseURL parses URL and panics on failure.
func MustParseURL(raw string) URL {
	parsed, err := ParseURL(raw)
	if err != nil {
		panic(err)
	}
	return parsed
}

// URL returns underlying *url.URL.
func (u *URL) URL() *url.URL {
	if u == nil {
		return nil
	}
	return u.value
}

// String returns string representation of underlying URL.
func (u *URL) String() string {
	if u == nil || u.value == nil {
		return ""
	}
	return u.value.String()
}

// UnmarshalYAML parses YAML string into URL.
func (u *URL) UnmarshalYAML(value *yaml.Node) error {
	var raw string
	if err := value.Decode(&raw); err != nil {
		return err
	}

	if raw == "" {
		u.value = nil
		return nil
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	u.value = parsed
	return nil
}
