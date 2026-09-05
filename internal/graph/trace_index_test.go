package graph

import (
	"reflect"
	"testing"

	"github.com/jingxu9x/idd-cli/internal/model"
)

// @test TEST-INTERNAL_GRAPH-016
func TestBuildTraceIndexResolutionAndProvenance(t *testing.T) {
	tests := []struct {
		name       string
		build      func() *model.IdentifierSet
		id         string
		resolution Resolution
		owners     int
	}{
		{
			name: "one declaration with implementation is resolved",
			build: func() *model.IdentifierSet {
				set := model.NewIdentifierSet()
				doc := model.NewIdentifier("SPEC-AUTH-001", model.TypeSpec, "Login", "docs/auth/spec.md", 12)
				doc.Namespace, doc.Status, doc.DocumentRole = "AUTH", "active", "spec"
				doc.AddTypedLinkAt("component:auth#Service", model.LinkReferences, doc.Source, 14, "components")
				set.Add(doc)
				code := model.NewIdentifier("SPEC-AUTH-001", model.TypeSpec, "", "internal/auth/service.go", 42)
				code.SetOrigin(model.OriginCode)
				set.Add(code)
				return set
			},
			id: "SPEC-AUTH-001", resolution: ResolutionResolved, owners: 1,
		},
		{
			name: "reference without declaration is unresolved",
			build: func() *model.IdentifierSet {
				set := model.NewIdentifierSet()
				test := model.NewIdentifier("TEST-AUTH-001", model.TypeTest, "", "docs/auth/testing.md", 20)
				test.AddTypedLinkAt("SPEC-AUTH-404", model.LinkImplements, test.Source, 22, "covers")
				set.Add(test)
				return set
			},
			id: "SPEC-AUTH-404", resolution: ResolutionUnresolved, owners: 0,
		},
		{
			name: "two declarations are ambiguous",
			build: func() *model.IdentifierSet {
				set := model.NewIdentifierSet()
				set.Add(model.NewIdentifier("SPEC-AUTH-001", model.TypeSpec, "", "docs/a/spec.md", 8))
				set.Add(model.NewIdentifier("SPEC-AUTH-001", model.TypeSpec, "", "docs/b/spec.md", 9))
				return set
			},
			id: "SPEC-AUTH-001", resolution: ResolutionAmbiguous, owners: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			index := BuildTraceIndex(test.build())
			entity, ok := index.Entity(test.id)
			if !ok {
				t.Fatalf("Entity(%q) missing", test.id)
			}
			if entity.Resolution != test.resolution || len(entity.Owners) != test.owners {
				t.Fatalf("entity resolution/owners = %s/%d, want %s/%d", entity.Resolution, len(entity.Owners), test.resolution, test.owners)
			}
		})
	}
}

func TestTraceIndexStableRelationsAndDetachedResults(t *testing.T) {
	set := model.NewIdentifierSet()
	spec := model.NewIdentifier("SPEC-AUTH-001", model.TypeSpec, "", "docs/auth/spec.md", 10)
	spec.AddTypedLinkAt("contract:auth#Authenticator", model.LinkContract, spec.Source, 13, "contracts")
	spec.AddTypedLinkAt("component:auth#Service", model.LinkReferences, spec.Source, 12, "components")
	set.Add(spec)
	index := BuildTraceIndex(set)

	got := index.Outbound(spec.ID)
	wantTypes := []model.LinkType{model.LinkConstrainedBy, model.LinkDesignedBy}
	if len(got) != 2 || !reflect.DeepEqual([]model.LinkType{got[0].Type, got[1].Type}, wantTypes) {
		t.Fatalf("outbound relation order = %#v", got)
	}
	if got[0].Provenance.Path != "docs/auth/spec.md" || got[0].Provenance.Line != 13 || got[0].Provenance.Field != "contracts" {
		t.Fatalf("contract provenance = %#v", got[0].Provenance)
	}
	got[0].Provenance.Path = "mutated"
	if index.Outbound(spec.ID)[0].Provenance.Path == "mutated" {
		t.Fatal("Outbound returned mutable backing storage")
	}
}

func TestTraceIndexEvidenceAndMentionOccurrences(t *testing.T) {
	set := model.NewIdentifierSet()
	testRecord := model.NewIdentifier("TEST-AUTH-001", model.TypeTest, "", "docs/auth/testing.md", 20)
	set.Add(testRecord)
	testCode := model.NewIdentifier("TEST-AUTH-001", model.TypeTest, "", "internal/auth/service_test.go", 31)
	testCode.SetOrigin(model.OriginCode)
	set.Add(testCode)
	mention := Occurrence{
		EntityID:     "TEST-AUTH-001",
		Origin:       "doc",
		DocumentRole: "design",
		Kind:         OccurrenceMention,
		Path:         "docs/auth/design.md",
		Line:         44,
	}

	index := BuildTraceIndexWithOccurrences(set, []Occurrence{mention})
	entity, ok := index.Entity("TEST-AUTH-001")
	if !ok || entity.Resolution != ResolutionResolved {
		t.Fatalf("entity = %#v, %v", entity, ok)
	}
	kinds := make(map[OccurrenceKind]bool)
	for _, occurrence := range entity.Occurrences {
		kinds[occurrence.Kind] = true
	}
	for _, kind := range []OccurrenceKind{OccurrenceDeclaration, OccurrenceTestEvidence, OccurrenceMention} {
		if !kinds[kind] {
			t.Errorf("missing occurrence kind %q in %#v", kind, entity.Occurrences)
		}
	}
	var hasExecutes, hasMention bool
	for _, relation := range index.Inbound("TEST-AUTH-001") {
		hasExecutes = hasExecutes || relation.Type == model.LinkExecutes
		hasMention = hasMention || relation.Type == model.LinkMentions
	}
	if !hasExecutes || !hasMention {
		t.Fatalf("inbound evidence relations = %#v", index.Inbound("TEST-AUTH-001"))
	}
	if got := index.OccurrencesIn("docs/auth/design.md", 40, 50); len(got) != 1 || got[0].Kind != OccurrenceMention {
		t.Fatalf("OccurrencesIn() = %#v", got)
	}
}

func TestTraceIndexOccurrenceIDsRemainUniqueForRepeatedReferences(t *testing.T) {
	set := model.NewIdentifierSet()
	testRecord := model.NewIdentifier("TEST-AUTH-001", model.TypeTest, "", "docs/auth/testing.md", 8)
	testRecord.AddTypedLinkAt("SPEC-AUTH-001", model.LinkImplements, testRecord.Source, 10, "Covers")
	testRecord.AddTypedLinkAt("SPEC-AUTH-001", model.LinkImplements, testRecord.Source, 10, "Covers")
	set.Add(testRecord)

	occurrences := BuildTraceIndex(set).Occurrences("SPEC-AUTH-001")
	if len(occurrences) != 2 {
		t.Fatalf("occurrences = %#v", occurrences)
	}
	if occurrences[0].ID == occurrences[1].ID {
		t.Fatalf("repeated references reused occurrence ID %q", occurrences[0].ID)
	}
}
