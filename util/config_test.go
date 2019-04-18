package util

import "testing"

func TestGetConfig(t *testing.T) {
	config := GetConfig()

	if config.Salt != "imW$OcsNQwy7XtVld@p&Nr#0mkN&qN33$M5*4ZNzmYe%95e&qUdE0f7!Lr0mPMoI" {
		t.Errorf("salt config not loaded")
	}

	if config.DB.Name != "scoreplus_test" || config.DB.Host != "localhost" || config.DB.Port != 5432 ||
		config.DB.Username != "scoreplus_user" || config.DB.Password != "8BcD2T3W6xpPjlXA9D#loRfQTz%p^zhT" {
		t.Errorf("DB config not loaded")
	}
}
