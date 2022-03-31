package main

import (
	"api/util"
	"testing"
)

func TestValidateLicenseScores(t *testing.T) {
	licenseMap, err := util.GetLicenseMap("../util/scores/licenseScores.txt")

	if err != nil {
		t.Fatal(err.Error())
	}
	for license, score := range licenseMap {
		if score < 0 || score > 100 {
			t.Fatalf("Score for license: \"%s\": %d is out of bounds", license, score)
		}
	}
}
