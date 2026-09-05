package main

import "testing"

func TestImpactedCommandSurface(t *testing.T) {
	command := commandNamed(docsCmd, "impacted")
	if command == nil || command.RunE == nil {
		t.Fatal("docs impacted command is absent from the executable CLI contract")
	}
	if command.Flags().Lookup("base") == nil {
		t.Error("docs impacted command is missing --base")
	}
}
