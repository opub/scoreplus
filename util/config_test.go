package util

import "testing"

func TestGetConfig(t *testing.T) {
	config, err := GetConfig()
	if err != nil {
		t.Errorf("error in GetConfig: %s", err)
	}

	if config.Salt == "" {
		t.Errorf("no salt config loaded")
	}

	if config.DB.Name == "" || config.DB.Host == "" || config.DB.Port == 0 ||
		config.DB.Username == "" || config.DB.Password == "" {
		t.Errorf("no DB config loaded")
	}
}
