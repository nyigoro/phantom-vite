package main

import "testing"

func TestWriteTempScript(t *testing.T) {
	file, err := writeTempScript("https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file)
}
