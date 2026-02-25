package config

// HydraConfig represents Hydra OAuth2/OIDC server configuration.
type HydraConfig struct {
	AdminURL *URL `yaml:"admin_url" validate:"required"` // Hydra Admin API URL
}
