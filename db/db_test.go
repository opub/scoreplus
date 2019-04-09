package db

import "testing"

func TestConnect(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Errorf("error in Connect: %s", err)
	}

	err = db.Ping()
	if err != nil {
		t.Errorf("ping failed: %s", err)
	}
}
