package main

import (
	"fmt"
	"testing"

	"github.com/oomph-ac/api/database"
)

func TestAuthDB(t *testing.T) {
	res, err := database.ObtainAuth("DEV")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(res)

	_, err = database.ObtainAuth("thisIsAnInvalidKeyThatShouldNeverBeUsed")
	if err == nil {
		// Since it's an invalid key, this should never happen
		t.Fail()
	}
}
