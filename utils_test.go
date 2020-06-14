package dsvg

import "testing"

func TestParseUnits(t *testing.T) {
	res, err := ParseUnits("3.55in")
	if err != nil {
		t.Fatal("failed to parse units: ", err)
	}

	t.Log(res)

	if res == 0 {
		t.Fatal("Got bad number from parse units")
	}
}
