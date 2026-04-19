package engine

import (
	"context"
	"testing"

	"github.com/jingxu9x/idd-link-validator/internal/config"
	"github.com/jingxu9x/idd-link-validator/internal/model"
)

func TestEngine_ConsistencyCheck_Warning(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			ConsistencyCheck: config.ConsistencyCheck{
				Enabled:   true,
				Threshold: 0.3,
			},
		},
	}

	ids := model.NewIdentifierSet()

	docID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "Test Spec", "validates user authentication and JWT token issuance", "docs/spec.md", 1)
	docID.SetOrigin(model.OriginDoc)
	ids.Add(docID)

	codeID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "", "database connection pooling settings", "main.go", 10)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if len(result.Warnings) == 0 {
		t.Error("Expected consistency warning for low similarity")
	}

	found := false
	for _, w := range result.Warnings {
		if w.Rule == "consistency-check" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected consistency-check warning")
	}
}

func TestEngine_ConsistencyCheck_NoWarning(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			ConsistencyCheck: config.ConsistencyCheck{
				Enabled:   true,
				Threshold: 0.3,
			},
		},
	}

	ids := model.NewIdentifierSet()

	docID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "Test Spec", "validates user authentication and JWT token issuance", "docs/spec.md", 1)
	docID.SetOrigin(model.OriginDoc)
	ids.Add(docID)

	codeID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "", "validates user authentication and JWT token issuance for session", "main.go", 10)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, w := range result.Warnings {
		if w.Rule == "consistency-check" {
			t.Error("Expected no consistency warning for similar texts")
		}
	}
}

func TestEngine_ConsistencyCheck_Disabled(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			ConsistencyCheck: config.ConsistencyCheck{
				Enabled:   false,
				Threshold: 0.3,
			},
		},
	}

	ids := model.NewIdentifierSet()

	docID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "Test Spec", "validates user authentication", "docs/spec.md", 1)
	docID.SetOrigin(model.OriginDoc)
	ids.Add(docID)

	codeID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "", "completely unrelated content here", "main.go", 10)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, w := range result.Warnings {
		if w.Rule == "consistency-check" {
			t.Error("Expected no consistency warnings when disabled")
		}
	}
}

func TestEngine_ConsistencyCheck_MissingDescribe(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			ConsistencyCheck: config.ConsistencyCheck{
				Enabled:   true,
				Threshold: 0.3,
			},
		},
	}

	ids := model.NewIdentifierSet()

	docID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "Test Spec", "", "docs/spec.md", 1)
	docID.SetOrigin(model.OriginDoc)
	ids.Add(docID)

	codeID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "", "validates user authentication", "main.go", 10)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, w := range result.Warnings {
		if w.Rule == "consistency-check" {
			t.Error("Expected no consistency warning when doc describe is empty")
		}
	}
}

func TestEngine_ConsistencyCheck_CodeOnly(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			ConsistencyCheck: config.ConsistencyCheck{
				Enabled:   true,
				Threshold: 0.3,
			},
		},
	}

	ids := model.NewIdentifierSet()

	codeID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "", "validates user authentication", "main.go", 10)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	for _, w := range result.Warnings {
		if w.Rule == "consistency-check" {
			t.Error("Expected no consistency warning when only code exists")
		}
	}
}

func TestEngine_ConsistencyCheck_HighThreshold(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			ConsistencyCheck: config.ConsistencyCheck{
				Enabled:   true,
				Threshold: 0.95,
			},
		},
	}

	ids := model.NewIdentifierSet()

	docID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "Test Spec", "validates user authentication and JWT tokens", "docs/spec.md", 1)
	docID.SetOrigin(model.OriginDoc)
	ids.Add(docID)

	codeID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "", "validates user authentication and JWT tokens for sessions", "main.go", 10)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, w := range result.Warnings {
		if w.Rule == "consistency-check" {
			found = true
			break
		}
	}
	if found {
		t.Error("Expected no consistency warning for very similar texts with high threshold")
	}
}

func TestEngine_ConsistencyCheck_LowSimilarity(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			ConsistencyCheck: config.ConsistencyCheck{
				Enabled:   true,
				Threshold: 0.1,
			},
		},
	}

	ids := model.NewIdentifierSet()

	docID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "Test Spec", "user login and authentication", "docs/spec.md", 1)
	docID.SetOrigin(model.OriginDoc)
	ids.Add(docID)

	codeID := model.NewIdentifierWithDescribe("SPEC-001", model.TypeSpec, "", "database connection pool settings", "main.go", 10)
	codeID.SetOrigin(model.OriginCode)
	ids.Add(codeID)

	eng := New(cfg)
	result, err := eng.Run(context.Background(), ids)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	found := false
	for _, w := range result.Warnings {
		if w.Rule == "consistency-check" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected consistency warning for very dissimilar texts")
	}
}

func TestEngine_AddStructuralErrors(t *testing.T) {
	cfg := config.Default()
	eng := New(cfg)

	errors := []*model.ValidationError{
		{Rule: "test-rule", Message: "test error"},
	}
	eng.AddStructuralErrors(errors)

	if len(eng.result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(eng.result.Errors))
	}
}

func TestEngine_inferLinkType(t *testing.T) {
	cfg := config.Default()
	eng := New(cfg)

	tests := []struct {
		fromType model.IdentifierType
		wantLink model.LinkType
	}{
		{model.TypeSpec, model.LinkTests},
		{model.TypeTest, model.LinkImplements},
		{model.TypeContract, model.LinkReferences},
		{model.TypeDesign, model.LinkReferences},
	}

	for _, tt := range tests {
		link := eng.inferLinkType(tt.fromType, "TEST-BE-001")
		if link != tt.wantLink {
			t.Errorf("inferLinkType(%v) = %v, want %v", tt.fromType, link, tt.wantLink)
		}
	}
}

func TestEngine_validateBidirectional(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Validation: config.ValidationConfig{
			RequireBidirectional: false,
			AllowOrphans:         true,
		},
	}
	eng := New(cfg)

	ids := model.NewIdentifierSet()
	specID := model.NewIdentifier("SPEC-001", model.TypeSpec, "", "docs/spec.md", 1)
	specID.SetOrigin(model.OriginDoc)
	ids.Add(specID)

	eng.Run(context.Background(), ids)

	if len(eng.result.Errors) != 0 {
		t.Errorf("Expected 0 errors when RequireBidirectional=false, got %d", len(eng.result.Errors))
	}
}

func TestEngine_validateDocCodeCorrespondence(t *testing.T) {
	cfg := config.Default()
	eng := New(cfg)

	ids := model.NewIdentifierSet()

	codeOnlyID := model.NewIdentifier("SPEC-001", model.TypeSpec, "", "main.go", 1)
	codeOnlyID.SetOrigin(model.OriginCode)
	ids.Add(codeOnlyID)

	eng.Run(context.Background(), ids)

	found := false
	for _, e := range eng.result.Errors {
		if e.Rule == "doc-code-correspondence" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error for code without doc")
	}
}

func TestEngine_BuildReport(t *testing.T) {
	cfg := config.Default()
	eng := New(cfg)

	ids := model.NewIdentifierSet()
	eng.Run(context.Background(), ids)

	report := eng.BuildReport()
	if report.Tool != "idd-cli" {
		t.Errorf("Tool = %q, want idd-cli", report.Tool)
	}
	if report.Version != "1.0.0" {
		t.Errorf("Version = %q, want 1.0.0", report.Version)
	}
}
