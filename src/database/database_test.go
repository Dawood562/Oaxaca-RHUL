//go:build integration

package database

import (
	"testing"
)

func TestDBAuth(t *testing.T) {
	user, name, pass := fetchDBAuth()
	if user == "-1" || name == "-1" || pass == "-1" {
		t.Fail()
	}
}
