package main

import "testing"

func TestPrefixedFind(t *testing.T) {
	expected := "Stue"

	fields := []string{
		"id=4c31cfdd3e8798598870393eeff0f79d",
		"rm=",
		"ve=05",
		"md=Chromecast",
		"ic=/setup/icon.png",
		"fn=Stue",
		"ca=4101",
		"st=1",
		"bs=FA8FCA93F0C4",
		"rs=]",
		"Addr:192.168.0.102",
	}

	value := findField(fields, "fn=")
	if value != expected {
		t.Errorf("Expected %v but got %v", expected, value)
	}
}