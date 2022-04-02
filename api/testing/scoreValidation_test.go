package main

import (
	"api/util"
	"testing"
)

func TestValidatelicenseScoring(t *testing.T) {
	licenseMap, err := util.GetLicenseMap("../util/scores/licenseScoring.csv")

	if err != nil {
		t.Fatal(err.Error())
	}
	for license, score := range licenseMap {
		if score < 0 || score > 100 {
			t.Fatalf("Score for license: \"%s\": %f is out of bounds", license, score)
		}
	}
}

func TestValidateCategoryWeights(t *testing.T) {
	err := util.GetActivityScoringData("../util/scores/activityScoring.csv")

	if err != nil {
		t.Fatal(err.Error())
	}
}
