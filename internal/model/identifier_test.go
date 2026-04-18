package model

import (
	"testing"
)

func TestParseIdentifierType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    IdentifierType
		wantErr bool
	}{
		{"SPEC lowercase", "spec", TypeSpec, false},
		{"SPEC uppercase", "SPEC", TypeSpec, false},
		{"SPEC mixed case", "Spec", TypeSpec, false},
		{"CONTRACT", "CONTRACT", TypeContract, false},
		{"TEST", "TEST", TypeTest, false},
		{"DESIGN", "DESIGN", TypeDesign, false},
		{"invalid type", "INVALID", "", true},
		{"empty string", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIdentifierType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIdentifierType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseIdentifierType(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewIdentifier(t *testing.T) {
	id := NewIdentifier("SPEC-001", TypeSpec, "Test Spec", "docs/test.md", 10)

	if id.ID != "SPEC-001" {
		t.Errorf("ID = %q, want %q", id.ID, "SPEC-001")
	}
	if id.Type != TypeSpec {
		t.Errorf("Type = %v, want %v", id.Type, TypeSpec)
	}
	if id.Title != "Test Spec" {
		t.Errorf("Title = %q, want %q", id.Title, "Test Spec")
	}
	if id.Source != "docs/test.md" {
		t.Errorf("Source = %q, want %q", id.Source, "docs/test.md")
	}
	if id.Line != 10 {
		t.Errorf("Line = %d, want %d", id.Line, 10)
	}
	if id.RawRef != "SPEC-001" {
		t.Errorf("RawRef = %q, want %q", id.RawRef, "SPEC-001")
	}
	if id.Links == nil {
		t.Error("Links should not be nil")
	}
}

func TestIdentifierSet_Add_Get_Has(t *testing.T) {
	set := NewIdentifierSet()
	spec := NewIdentifier("SPEC-001", TypeSpec, "Test", "test.md", 1)

	set.Add(spec)

	if !set.Has("SPEC-001") {
		t.Error("Has(SPEC-001) = false, want true")
	}
	if set.Has("NONEXISTENT") {
		t.Error("Has(NONEXISTENT) = true, want false")
	}

	got, ok := set.Get("SPEC-001")
	if !ok {
		t.Error("Get(SPEC-001) ok = false, want true")
	}
	if got.ID != "SPEC-001" {
		t.Errorf("Get(SPEC-001).ID = %q, want %q", got.ID, "SPEC-001")
	}
}

func TestIdentifierSet_Count(t *testing.T) {
	set := NewIdentifierSet()
	if set.Count() != 0 {
		t.Errorf("Count() = %d, want 0", set.Count())
	}

	set.Add(NewIdentifier("SPEC-001", TypeSpec, "", "", 0))
	set.Add(NewIdentifier("TEST-001", TypeTest, "", "", 0))

	if set.Count() != 2 {
		t.Errorf("Count() = %d, want 2", set.Count())
	}
}

func TestIdentifierSet_All(t *testing.T) {
	set := NewIdentifierSet()
	set.Add(NewIdentifier("SPEC-001", TypeSpec, "", "", 0))
	set.Add(NewIdentifier("TEST-001", TypeTest, "", "", 0))

	all := set.All()
	if len(all) != 2 {
		t.Errorf("len(All()) = %d, want 2", len(all))
	}
}

func TestIdentifierSet_Merge(t *testing.T) {
	set1 := NewIdentifierSet()
	set1.Add(NewIdentifier("SPEC-001", TypeSpec, "", "", 0))

	set2 := NewIdentifierSet()
	set2.Add(NewIdentifier("TEST-001", TypeTest, "", "", 0))

	set1.Merge(set2)

	if set1.Count() != 2 {
		t.Errorf("After merge Count() = %d, want 2", set1.Count())
	}
}

func TestIdentifier_AddLink(t *testing.T) {
	id := NewIdentifier("SPEC-001", TypeSpec, "", "", 0)
	id.AddLink("TEST-001")
	id.AddLink("TEST-002")

	if len(id.Links) != 2 {
		t.Errorf("len(Links) = %d, want 2", len(id.Links))
	}
}

func TestNewAnnotation(t *testing.T) {
	ann := NewAnnotation(TypeSpec, "SPEC-001", "test.go", "@spec SPEC-001", "context", 10)

	if ann.Type != TypeSpec {
		t.Errorf("Type = %v, want %v", ann.Type, TypeSpec)
	}
	if ann.Ref != "SPEC-001" {
		t.Errorf("Ref = %q, want %q", ann.Ref, "SPEC-001")
	}
	if ann.Source != "test.go" {
		t.Errorf("Source = %q, want %q", ann.Source, "test.go")
	}
	if ann.Line != 10 {
		t.Errorf("Line = %d, want %d", ann.Line, 10)
	}
}

func TestAnnotation_ToIdentifier(t *testing.T) {
	ann := NewAnnotation(TypeSpec, "SPEC-001", "test.go", "@spec SPEC-001", "context", 10)
	id := ann.ToIdentifier()

	if id.ID != "SPEC-001" {
		t.Errorf("ID = %q, want %q", id.ID, "SPEC-001")
	}
	if id.Type != TypeSpec {
		t.Errorf("Type = %v, want %v", id.Type, TypeSpec)
	}
	if id.Source != "test.go" {
		t.Errorf("Source = %q, want %q", id.Source, "test.go")
	}
	if id.RawRef != "@spec SPEC-001" {
		t.Errorf("RawRef = %q, want %q", id.RawRef, "@spec SPEC-001")
	}
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Rule:    "test-rule",
		Message: "test message",
	}

	expected := "[test-rule] test message"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

func TestNewValidationResult(t *testing.T) {
	result := NewValidationResult()

	if result.Valid != false {
		t.Error("Valid should be false initially")
	}
	if result.Errors == nil {
		t.Error("Errors should not be nil")
	}
	if result.Warnings == nil {
		t.Error("Warnings should not be nil")
	}
}

func TestValidationResult_AddError(t *testing.T) {
	result := NewValidationResult()
	result.AddError("rule1", "message", "source", "link", "code")

	if result.Valid {
		t.Error("Valid should be false after AddError")
	}
	if len(result.Errors) != 1 {
		t.Errorf("len(Errors) = %d, want 1", len(result.Errors))
	}
}

func TestValidationResult_AddWarning(t *testing.T) {
	result := NewValidationResult()
	result.AddWarning("rule", "msg", "", "", "")

	if len(result.Warnings) != 1 {
		t.Errorf("len(Warnings) = %d, want 1", len(result.Warnings))
	}
	if result.Valid {
		t.Error("Valid should be false (default)")
	}
}

func TestValidationResult_Sort(t *testing.T) {
	result := NewValidationResult()
	result.AddError("bbb", "msg2", "", "", "")
	result.AddError("aaa", "msg1", "", "", "")
	result.AddError("aaa", "msg2", "", "", "")

	result.Sort()

	if result.Errors[0].Rule != "aaa" {
		t.Errorf("Errors[0].Rule = %q, want 'aaa'", result.Errors[0].Rule)
	}
}
