package model

import (
	"testing"
)

func TestLinkType_Constants(t *testing.T) {
	if LinkImplements != "implements" {
		t.Errorf("LinkImplements = %q, want %q", LinkImplements, "implements")
	}
	if LinkTests != "tests" {
		t.Errorf("LinkTests = %q, want %q", LinkTests, "tests")
	}
	if LinkReferences != "references" {
		t.Errorf("LinkReferences = %q, want %q", LinkReferences, "references")
	}
	if LinkAnnotates != "annotates" {
		t.Errorf("LinkAnnotates = %q, want %q", LinkAnnotates, "annotates")
	}
}

func TestReverseLinkType(t *testing.T) {
	tests := []struct {
		name string
		link LinkType
		want LinkType
	}{
		{"tests -> implements", LinkTests, LinkImplements},
		{"implements -> tests", LinkImplements, LinkTests},
		{"references -> references", LinkReferences, LinkReferences},
		{"annotates -> annotates", LinkAnnotates, LinkAnnotates},
		{"unknown -> same", LinkType("unknown"), LinkType("unknown")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReverseLinkType(tt.link)
			if got != tt.want {
				t.Errorf("ReverseLinkType(%v) = %v, want %v", tt.link, got, tt.want)
			}
		})
	}
}

func TestNewLink(t *testing.T) {
	link := NewLink("SPEC-001", "TEST-001", LinkTests, "test.go", 10)

	if link.From != "SPEC-001" {
		t.Errorf("From = %q, want %q", link.From, "SPEC-001")
	}
	if link.To != "TEST-001" {
		t.Errorf("To = %q, want %q", link.To, "TEST-001")
	}
	if link.Type != LinkTests {
		t.Errorf("Type = %v, want %v", link.Type, LinkTests)
	}
	if link.Source != "test.go" {
		t.Errorf("Source = %q, want %q", link.Source, "test.go")
	}
	if link.Line != 10 {
		t.Errorf("Line = %d, want %d", link.Line, 10)
	}
}
