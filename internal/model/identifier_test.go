// Package model provides testing utilities for the model module.

// Spec: docs/internal/model/spec.md
// Test: docs/internal/model/testing.md
package model

import (
	"testing"
)

// @test TEST-INT_MOD-019
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

// @test TEST-INT_MOD-020
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

// @test TEST-INT_MOD-001
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

// @test TEST-INT_MOD-002
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

// @test TEST-INT_MOD-003
func TestIdentifierSet_All(t *testing.T) {
	set := NewIdentifierSet()
	set.Add(NewIdentifier("SPEC-001", TypeSpec, "", "", 0))
	set.Add(NewIdentifier("TEST-001", TypeTest, "", "", 0))

	all := set.All()
	if len(all) != 2 {
		t.Errorf("len(All()) = %d, want 2", len(all))
	}
}

// @test TEST-INT_MOD-004
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

// @test TEST-INT_MOD-005
func TestIdentifier_AddLink(t *testing.T) {
	id := NewIdentifier("SPEC-001", TypeSpec, "", "", 0)
	id.AddLink("TEST-001")
	id.AddLink("TEST-002")

	if len(id.Links) != 2 {
		t.Errorf("len(Links) = %d, want 2", len(id.Links))
	}
}

// @test TEST-INT_MOD-022
func TestNewAnnotation(t *testing.T) {
	ann := NewAnnotation(TypeSpec, "SPEC-001", "test.go", "@implement SPEC-001", "context", 10)

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

// @test TEST-INT_MOD-006
func TestAnnotation_ToIdentifier(t *testing.T) {
	ann := NewAnnotation(TypeSpec, "SPEC-001", "test.go", "@implement SPEC-001", "context", 10)
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
	if id.RawRef != "@implement SPEC-001" {
		t.Errorf("RawRef = %q, want %q", id.RawRef, "@implement SPEC-001")
	}
}

// @test TEST-INT_MOD-007
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

// @test TEST-INT_MOD-021
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

// @test TEST-INT_MOD-008
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

// @test TEST-INT_MOD-009
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

// @test TEST-INT_MOD-010
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

// @test TEST-INT_MOD-023
func TestNewIdentifierWithDescribe(t *testing.T) {
	id := NewIdentifierWithDescribe("SPEC-001", TypeSpec, "Test Spec", "Validates JWT tokens", "docs/test.md", 10)

	if id.ID != "SPEC-001" {
		t.Errorf("ID = %q, want %q", id.ID, "SPEC-001")
	}
	if id.Type != TypeSpec {
		t.Errorf("Type = %v, want %v", id.Type, TypeSpec)
	}
	if id.Title != "Test Spec" {
		t.Errorf("Title = %q, want %q", id.Title, "Test Spec")
	}
	if id.Describe != "Validates JWT tokens" {
		t.Errorf("Describe = %q, want %q", id.Describe, "Validates JWT tokens")
	}
	if id.Source != "docs/test.md" {
		t.Errorf("Source = %q, want %q", id.Source, "docs/test.md")
	}
	if id.Line != 10 {
		t.Errorf("Line = %d, want %d", id.Line, 10)
	}
	if id.Origin != OriginDoc {
		t.Errorf("Origin = %v, want %v", id.Origin, OriginDoc)
	}
}

// @test TEST-INT_MOD-024
func TestNewAnnotationWithComment(t *testing.T) {
	ann := NewAnnotationWithComment(TypeSpec, "SPEC-001", "test.go", "@implement SPEC-001", "context", "Validates authentication", 10)

	if ann.Type != TypeSpec {
		t.Errorf("Type = %v, want %v", ann.Type, TypeSpec)
	}
	if ann.FunctionComment != "Validates authentication" {
		t.Errorf("FunctionComment = %q, want %q", ann.FunctionComment, "Validates authentication")
	}
}

// @test TEST-INT_MOD-011
func TestAnnotation_ToIdentifier_WithFunctionComment(t *testing.T) {
	ann := NewAnnotationWithComment(TypeSpec, "SPEC-001", "test.go", "@implement SPEC-001", "context", "Validates JWT tokens", 10)
	id := ann.ToIdentifier()

	if id.Describe != "Validates JWT tokens" {
		t.Errorf("Describe = %q, want %q", id.Describe, "Validates JWT tokens")
	}
}

// @test TEST-INT_MOD-012
func TestAnnotation_ToIdentifier_WithoutFunctionComment(t *testing.T) {
	ann := NewAnnotation(TypeSpec, "SPEC-001", "test.go", "@implement SPEC-001", "context", 10)
	id := ann.ToIdentifier()

	if id.Describe != "" {
		t.Errorf("Describe = %q, want empty string", id.Describe)
	}
}

// @test TEST-INT_MOD-013
func TestIdentifierSet_GetAll(t *testing.T) {
	set := NewIdentifierSet()
	spec1 := NewIdentifier("SPEC-001", TypeSpec, "", "file1.md", 1)
	spec2 := NewIdentifier("SPEC-001", TypeSpec, "", "file2.md", 2)

	set.Add(spec1)
	set.Add(spec2)

	all := set.GetAll("SPEC-001")
	if len(all) != 2 {
		t.Errorf("GetAll(SPEC-001) returned %d items, want 2", len(all))
	}

	none := set.GetAll("NONEXISTENT")
	if len(none) != 0 {
		t.Errorf("GetAll(NONEXISTENT) returned %d items, want 0", len(none))
	}
}

// @test TEST-INT_MOD-014
func TestIdentifierSet_AllIdentifiers(t *testing.T) {
	set := NewIdentifierSet()
	set.Add(NewIdentifier("SPEC-001", TypeSpec, "", "", 0))
	set.Add(NewIdentifier("SPEC-001", TypeSpec, "", "", 0))
	set.Add(NewIdentifier("TEST-001", TypeTest, "", "", 0))

	all := set.AllIdentifiers()
	if len(all) != 3 {
		t.Errorf("AllIdentifiers() returned %d items, want 3", len(all))
	}
}

// @test TEST-INT_MOD-015
func TestIdentifierSet_ByOrigin(t *testing.T) {
	set := NewIdentifierSet()
	docSpec := NewIdentifier("SPEC-001", TypeSpec, "", "doc.md", 0)
	docSpec.Origin = OriginDoc
	codeSpec := NewIdentifier("SPEC-002", TypeSpec, "", "code.go", 0)
	codeSpec.Origin = OriginCode

	set.Add(docSpec)
	set.Add(codeSpec)

	docOnly := set.ByOrigin(OriginDoc)
	if len(docOnly) != 1 {
		t.Errorf("ByOrigin(OriginDoc) returned %d items, want 1", len(docOnly))
	}

	codeOnly := set.ByOrigin(OriginCode)
	if len(codeOnly) != 1 {
		t.Errorf("ByOrigin(OriginCode) returned %d items, want 1", len(codeOnly))
	}
}

// @test TEST-INT_MOD-016
func TestIdentifierSet_HasOrigin(t *testing.T) {
	set := NewIdentifierSet()
	docSpec := NewIdentifier("SPEC-001", TypeSpec, "", "doc.md", 0)
	docSpec.Origin = OriginDoc
	set.Add(docSpec)

	if !set.HasOrigin("SPEC-001", OriginDoc) {
		t.Error("HasOrigin(SPEC-001, OriginDoc) = false, want true")
	}

	if set.HasOrigin("SPEC-001", OriginCode) {
		t.Error("HasOrigin(SPEC-001, OriginCode) = true, want false")
	}

	if set.HasOrigin("NONEXISTENT", OriginDoc) {
		t.Error("HasOrigin(NONEXISTENT, OriginDoc) = true, want false")
	}
}

// @test TEST-INT_MOD-017
func TestIdentifierSet_Get(t *testing.T) {
	set := NewIdentifierSet()
	spec := NewIdentifier("SPEC-001", TypeSpec, "", "", 0)
	set.Add(spec)

	got, ok := set.Get("SPEC-001")
	if !ok {
		t.Error("Get(SPEC-001) ok = false, want true")
	}
	if got.ID != "SPEC-001" {
		t.Errorf("Get(SPEC-001).ID = %q, want SPEC-001", got.ID)
	}

	_, ok = set.Get("NONEXISTENT")
	if ok {
		t.Error("Get(NONEXISTENT) ok = true, want false")
	}
}

// @test TEST-INT_MOD-018
func TestValidationResult_Sort_MultipleRules(t *testing.T) {
	result := NewValidationResult()
	result.AddError("zzz", "msg3", "", "", "")
	result.AddError("aaa", "msg1", "", "", "")
	result.AddError("mmm", "msg2", "", "", "")

	result.Sort()

	if result.Errors[0].Rule != "aaa" {
		t.Errorf("Errors[0].Rule = %q, want 'aaa'", result.Errors[0].Rule)
	}
	if result.Errors[1].Rule != "mmm" {
		t.Errorf("Errors[1].Rule = %q, want 'mmm'", result.Errors[1].Rule)
	}
	if result.Errors[2].Rule != "zzz" {
		t.Errorf("Errors[2].Rule = %q, want 'zzz'", result.Errors[2].Rule)
	}
}
