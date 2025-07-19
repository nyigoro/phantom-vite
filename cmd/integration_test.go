package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestAgentCommand(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "agent", "summarize this document")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Agent command failed: %v\nOutput: %s", err, output)
	}

	expectedOutput := "This page is a sample domain used for illustrative purposes."
	if !strings.Contains(string(output), expectedOutput) {
		t.Errorf("Expected output \"%s\" not found in:\n%s", expectedOutput, output)
	}
}

func TestAgentCommandGeneralPrompt(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "agent", "hello world")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Agent command failed: %v\nOutput: %s", err, output)
	}

	expectedOutput := "[Agent Reply] I received: hello world"
	if !strings.Contains(string(output), expectedOutput) {
		t.Errorf("Expected output \"%s\" not found in:\n%s", expectedOutput, output)
	}
}