package conf

import (
	"testing"

	"github.com/spf13/viper"
)

func TestEmptyEnvTurnsOff(t *testing.T) {
	viper.Reset()
	t.Setenv("PUBLIC_IP_URL", "")
	t.Setenv("TIMEOUT", "30")
	cfg := read(t.TempDir() + "/config.yaml")
	if cfg.PublicIPURL != "" {
		t.Errorf("empty env should override the default, got %q", cfg.PublicIPURL)
	}
	if cfg.Timeout != 30 || cfg.Theme != "cosmo" {
		t.Errorf("got timeout %d theme %q", cfg.Timeout, cfg.Theme)
	}
}
