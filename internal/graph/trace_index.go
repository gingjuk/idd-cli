package graph

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jingxu9x/idd-cli/internal/model"
)

// Resolution describes whether an entity has exactly one canonical declaration.
type Resolution string

const (
	ResolutionResolved   Resolution = "resolved"
	ResolutionUnresolved Resolution = "unresolved"
	ResolutionAmbiguous  Resolution = "ambiguous"
)

// OccurrenceKind classifies each place where an entity appears.
type OccurrenceKind string

const (
	OccurrenceDeclaration    OccurrenceKind = "declaration"
	OccurrenceReference      OccurrenceKind = "reference"
	OccurrenceImplementation OccurrenceKind = "implementation"
	OccurrenceTestEvidence   OccurrenceKind = "test_evidence"
	OccurrenceMention        OccurrenceKind = "mention"
)

// Occurrence preserves one physical appearance of a logical entity.
type Occurrence struct {
	ID           string         `json:"id"`
	EntityID     string         `json:"entity_id"`
	Origin       string         `json:"origin"`
	DocumentRole string         `json:"document_role,omitempty"`
	Kind         OccurrenceKind `json:"kind"`
	Path         string         `json:"path"`
	Line         int            `json:"line,omitempty"`
	Field        string         `json:"field,omitempty"`
	RecordID     string         `json:"record_id,omitempty"`
	Title        string         `json:"title,omitempty"`
	Description  string         `json:"description,omitempty"`
}

// Entity is the logical identity shared by all declarations, references, and evidence.
type Entity struct {
	ID          string       `json:"id"`
	Kind        string       `json:"kind"`
	Namespace   string       `json:"namespace,omitempty"`
	Lifecycle   string       `json:"lifecycle,omitempty"`
	Resolution  Resolution   `json:"resolution"`
	Owner       *Occurrence  `json:"owner,omitempty"`
	Owners      []Occurrence `json:"owners,omitempty"`
	Occurrences []Occurrence `json:"occurrences"`
}

// Provenance identifies the authored field or annotation that created a relation.
type Provenance struct {
	OccurrenceID string         `json:"occurrence_id"`
	Kind         OccurrenceKind `json:"kind"`
	Path         string         `json:"path"`
	Line         int            `json:"line,omitempty"`
	Field        string         `json:"field,omitempty"`
	RecordID     string         `json:"record_id,omitempty"`
}

// Relation is one directed, typed relationship with lossless provenance.
type Relation struct {
	From       string         `json:"from"`
	To         string         `json:"to"`
	Type       model.LinkType `json:"type"`
	Provenance Provenance     `json:"provenance"`
}

// TraceIndex stores logical entities separately from their physical occurrences
// and keeps both directions of every provenance-bearing relation.
type TraceIndex struct {
	entities      map[string]*Entity
	occurrences   map[string][]Occurrence
	byPath        map[string][]Occurrence
	occurrenceIDs map[string]int
	inbound       map[string][]Relation
	outbound      map[string][]Relation
	relations     []Relation
}

// NewTraceIndex creates an empty project traceability index.
func NewTraceIndex() *TraceIndex {
	return &TraceIndex{
		entities:      make(map[string]*Entity),
		occurrences:   make(map[string][]Occurrence),
		byPath:        make(map[string][]Occurrence),
		occurrenceIDs: make(map[string]int),
		inbound:       make(map[string][]Relation),
		outbound:      make(map[string][]Relation),
	}
}

// BuildTraceIndex projects the collected identifier stream into the canonical
// entity/occurrence and relation/provenance layers. It does not discard broken
// references: a relation target without an owner becomes an unresolved entity.
// @implement SPEC-INTERNAL_GRAPH-009
func BuildTraceIndex(ids *model.IdentifierSet) *TraceIndex {
	return BuildTraceIndexWithOccurrences(ids, nil)
}

// BuildTraceIndexWithOccurrences adds parser-owned occurrences such as plain
// Markdown mentions that are intentionally absent from IdentifierSet because
// they must not satisfy correspondence or coverage validation.
// @implement SPEC-INTERNAL_GRAPH-009
func BuildTraceIndexWithOccurrences(ids *model.IdentifierSet, occurrences []Occurrence) *TraceIndex {
	index := NewTraceIndex()
	if ids == nil {
		for _, occurrence := range occurrences {
			index.addExternalOccurrence(occurrence)
		}
		index.finalize()
		return index
	}
	for _, identifier := range ids.AllIdentifiers() {
		index.addIdentifier(identifier)
	}
	for _, occurrence := range occurrences {
		index.addExternalOccurrence(occurrence)
	}
	for _, identifier := range ids.AllIdentifiers() {
		for _, ref := range identifier.Links {
			index.addRelation(identifier, model.IdentifierLink{Ref: ref, Type: inferTraceLinkType(identifier.Type, ref)})
		}
		for _, link := range identifier.TypedLinks {
			index.addRelation(identifier, link)
		}
	}
	index.finalize()
	return index
}

func (i *TraceIndex) addExternalOccurrence(occurrence Occurrence) {
	if occurrence.EntityID == "" {
		return
	}
	if occurrence.ID == "" {
		occurrence = newOccurrence(
			occurrence.EntityID,
			occurrence.Origin,
			occurrence.DocumentRole,
			occurrence.Kind,
			occurrence.Path,
			occurrence.Line,
			occurrence.Field,
			occurrence.RecordID,
		)
	}
	entity := i.ensureEntity(occurrence.EntityID, inferredIdentifierType(occurrence.EntityID), "")
	occurrence = i.addOccurrence(entity, occurrence)
	if occurrence.Kind == OccurrenceMention {
		relation := Relation{
			From: "document:" + occurrence.Path,
			To:   occurrence.EntityID,
			Type: model.LinkMentions,
			Provenance: Provenance{
				OccurrenceID: occurrence.ID,
				Kind:         occurrence.Kind,
				Path:         occurrence.Path,
				Line:         occurrence.Line,
				Field:        occurrence.Field,
				RecordID:     occurrence.RecordID,
			},
		}
		i.relations = append(i.relations, relation)
		i.outbound[relation.From] = append(i.outbound[relation.From], relation)
		i.inbound[relation.To] = append(i.inbound[relation.To], relation)
	}
}

func (i *TraceIndex) addIdentifier(identifier *model.Identifier) {
	if identifier == nil || identifier.ID == "" {
		return
	}
	entity := i.ensureEntity(identifier.ID, identifier.Type, identifier.Kind)
	if entity.Namespace == "" {
		entity.Namespace = identifier.Namespace
		if entity.Namespace == "" {
			entity.Namespace = inferredNamespace(identifier.ID)
		}
	}
	if identifier.Status != "" {
		entity.Lifecycle = identifier.Status
	}
	kind := OccurrenceDeclaration
	origin := string(identifier.Origin)
	if identifier.Origin == model.OriginCode {
		kind = OccurrenceImplementation
		origin = "source"
		if identifier.Type == model.TypeTest {
			kind = OccurrenceTestEvidence
			origin = "test"
		}
	}
	role := identifier.DocumentRole
	if role == "" && identifier.Origin == model.OriginDoc {
		role = strings.TrimSuffix(filepath.Base(identifier.Source), filepath.Ext(identifier.Source))
	}
	occurrence := newOccurrence(identifier.ID, origin, role, kind, identifier.Source, identifier.Line, "", identifier.ID)
	occurrence.Title = identifier.Title
	occurrence.Description = identifier.Describe
	occurrence = i.addOccurrence(entity, occurrence)
	relationType := model.LinkImplements
	if kind == OccurrenceTestEvidence {
		relationType = model.LinkExecutes
	}
	relation := Relation{
		From: occurrence.ID,
		To:   identifier.ID,
		Type: relationType,
		Provenance: Provenance{
			OccurrenceID: occurrence.ID,
			Kind:         occurrence.Kind,
			Path:         occurrence.Path,
			Line:         occurrence.Line,
			Field:        "annotation",
			RecordID:     occurrence.RecordID,
		},
	}
	if kind == OccurrenceImplementation || kind == OccurrenceTestEvidence {
		i.relations = append(i.relations, relation)
		i.outbound[relation.From] = append(i.outbound[relation.From], relation)
		i.inbound[relation.To] = append(i.inbound[relation.To], relation)
	}
}

func (i *TraceIndex) addRelation(identifier *model.Identifier, link model.IdentifierLink) {
	if identifier == nil || link.Ref == "" {
		return
	}
	if link.Type == "" {
		link.Type = inferTraceLinkType(identifier.Type, link.Ref)
	}
	link.Type = normalizeTraceRelationType(identifier, link)
	i.ensureEntity(identifier.ID, identifier.Type, identifier.Kind)
	target := i.ensureEntity(link.Ref, inferredIdentifierType(link.Ref), "")
	path, line := link.Source, link.Line
	if path == "" {
		path = identifier.Source
	}
	if line == 0 {
		line = identifier.Line
	}
	field := link.Field
	if field == "" {
		field = inferredRelationField(link.Type)
	}
	recordID := link.RecordID
	if recordID == "" {
		recordID = identifier.ID
	}
	role := identifier.DocumentRole
	if role == "" && identifier.Origin == model.OriginDoc {
		role = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}
	occurrence := newOccurrence(link.Ref, string(identifier.Origin), role, OccurrenceReference, path, line, field, recordID)
	occurrence = i.addOccurrence(target, occurrence)
	relation := Relation{
		From: identifier.ID,
		To:   link.Ref,
		Type: link.Type,
		Provenance: Provenance{
			OccurrenceID: occurrence.ID,
			Kind:         occurrence.Kind,
			Path:         occurrence.Path,
			Line:         occurrence.Line,
			Field:        occurrence.Field,
			RecordID:     occurrence.RecordID,
		},
	}
	i.relations = append(i.relations, relation)
	i.outbound[relation.From] = append(i.outbound[relation.From], relation)
	i.inbound[relation.To] = append(i.inbound[relation.To], relation)
}

func (i *TraceIndex) ensureEntity(id string, identifierType model.IdentifierType, identifierKind string) *Entity {
	if entity := i.entities[id]; entity != nil {
		return entity
	}
	kind := strings.ToLower(string(identifierType))
	if identifierKind == "component" || strings.HasPrefix(id, "component:") {
		kind = "component"
	} else if strings.HasPrefix(id, "contract:") {
		kind = "contract"
	}
	entity := &Entity{ID: id, Kind: kind, Namespace: inferredNamespace(id), Resolution: ResolutionUnresolved, Occurrences: []Occurrence{}}
	i.entities[id] = entity
	return entity
}

func (i *TraceIndex) addOccurrence(entity *Entity, occurrence Occurrence) Occurrence {
	baseID := occurrence.ID
	ordinal := i.occurrenceIDs[baseID]
	i.occurrenceIDs[baseID] = ordinal + 1
	if ordinal > 0 {
		occurrence.ID = fmt.Sprintf("%s#%d", baseID, ordinal+1)
	}
	i.occurrences[entity.ID] = append(i.occurrences[entity.ID], occurrence)
	i.byPath[occurrence.Path] = append(i.byPath[occurrence.Path], occurrence)
	entity.Occurrences = append(entity.Occurrences, occurrence)
	return occurrence
}

func (i *TraceIndex) finalize() {
	for _, entity := range i.entities {
		sortOccurrences(entity.Occurrences)
		i.occurrences[entity.ID] = append([]Occurrence(nil), entity.Occurrences...)
		for _, occurrence := range entity.Occurrences {
			if occurrence.Kind == OccurrenceDeclaration {
				entity.Owners = append(entity.Owners, occurrence)
			}
		}
		switch len(entity.Owners) {
		case 0:
			entity.Resolution = ResolutionUnresolved
		case 1:
			entity.Resolution = ResolutionResolved
			owner := entity.Owners[0]
			entity.Owner = &owner
			if entity.Lifecycle == "" {
				entity.Lifecycle = "active"
			}
		default:
			entity.Resolution = ResolutionAmbiguous
		}
	}
	sortRelations(i.relations)
	for id := range i.inbound {
		sortRelations(i.inbound[id])
	}
	for id := range i.outbound {
		sortRelations(i.outbound[id])
	}
	for path := range i.byPath {
		sortOccurrences(i.byPath[path])
	}
}

// Entity returns a detached entity snapshot and whether the ID appeared.
func (i *TraceIndex) Entity(id string) (Entity, bool) {
	entity, ok := i.entities[id]
	if !ok {
		return Entity{}, false
	}
	return cloneEntity(entity), true
}

// Entities returns all entity snapshots sorted by ID.
func (i *TraceIndex) Entities() []Entity {
	result := make([]Entity, 0, len(i.entities))
	for _, entity := range i.entities {
		result = append(result, cloneEntity(entity))
	}
	sort.Slice(result, func(a, b int) bool { return result[a].ID < result[b].ID })
	return result
}

// Occurrences returns every occurrence for an ID in stable location order.
func (i *TraceIndex) Occurrences(id string) []Occurrence {
	return append([]Occurrence(nil), i.occurrences[id]...)
}

// OccurrencesIn returns occurrences in a project-relative path whose line is
// within the inclusive range. A non-positive boundary leaves that side open.
func (i *TraceIndex) OccurrencesIn(path string, fromLine, toLine int) []Occurrence {
	path = filepath.ToSlash(filepath.Clean(path))
	result := make([]Occurrence, 0)
	for _, occurrence := range i.byPath[path] {
		if fromLine > 0 && occurrence.Line < fromLine {
			continue
		}
		if toLine > 0 && occurrence.Line > toLine {
			continue
		}
		result = append(result, occurrence)
	}
	return result
}

// Inbound returns relations targeting id in stable order.
func (i *TraceIndex) Inbound(id string) []Relation {
	return append([]Relation(nil), i.inbound[id]...)
}

// Outbound returns relations originating at id in stable order.
func (i *TraceIndex) Outbound(id string) []Relation {
	return append([]Relation(nil), i.outbound[id]...)
}

// Relations returns every project relation in stable order.
func (i *TraceIndex) Relations() []Relation {
	return append([]Relation(nil), i.relations...)
}

func cloneEntity(entity *Entity) Entity {
	result := *entity
	result.Owners = append([]Occurrence(nil), entity.Owners...)
	result.Occurrences = append([]Occurrence(nil), entity.Occurrences...)
	if entity.Owner != nil {
		owner := *entity.Owner
		result.Owner = &owner
	}
	return result
}

func newOccurrence(entityID, origin, role string, kind OccurrenceKind, path string, line int, field, recordID string) Occurrence {
	path = filepath.ToSlash(filepath.Clean(path))
	if path == "." {
		path = ""
	}
	key := fmt.Sprintf("%s|%s|%d|%s|%s|%s", entityID, path, line, kind, field, recordID)
	return Occurrence{ID: key, EntityID: entityID, Origin: origin, DocumentRole: role, Kind: kind, Path: path, Line: line, Field: field, RecordID: recordID}
}

func sortOccurrences(values []Occurrence) {
	sort.SliceStable(values, func(a, b int) bool {
		left, right := values[a], values[b]
		if left.Kind != right.Kind {
			return left.Kind < right.Kind
		}
		if left.Path != right.Path {
			return left.Path < right.Path
		}
		if left.Line != right.Line {
			return left.Line < right.Line
		}
		if left.Field != right.Field {
			return left.Field < right.Field
		}
		return left.ID < right.ID
	})
}

func sortRelations(values []Relation) {
	sort.SliceStable(values, func(a, b int) bool {
		left, right := values[a], values[b]
		if left.Type != right.Type {
			return left.Type < right.Type
		}
		if left.From != right.From {
			return left.From < right.From
		}
		if left.To != right.To {
			return left.To < right.To
		}
		if left.Provenance.Path != right.Provenance.Path {
			return left.Provenance.Path < right.Provenance.Path
		}
		if left.Provenance.Line != right.Provenance.Line {
			return left.Provenance.Line < right.Provenance.Line
		}
		return left.Provenance.Field < right.Provenance.Field
	})
}

func inferredNamespace(id string) string {
	if index := strings.IndexByte(id, ':'); index >= 0 {
		scoped := id[index+1:]
		if hash := strings.IndexByte(scoped, '#'); hash >= 0 {
			scoped = scoped[:hash]
		}
		return scoped
	}
	first := strings.IndexByte(id, '-')
	last := strings.LastIndexByte(id, '-')
	if first >= 0 && last > first {
		return id[first+1 : last]
	}
	return ""
}

func inferredIdentifierType(id string) model.IdentifierType {
	switch {
	case strings.HasPrefix(id, "SPEC-"):
		return model.TypeSpec
	case strings.HasPrefix(id, "TEST-"):
		return model.TypeTest
	case strings.HasPrefix(id, "CONTRACT-"), strings.HasPrefix(id, "contract:"):
		return model.TypeContract
	case strings.HasPrefix(id, "DESIGN-"), strings.HasPrefix(id, "component:"):
		return model.TypeDesign
	default:
		return ""
	}
}

func inferTraceLinkType(fromType model.IdentifierType, ref string) model.LinkType {
	toType := inferredIdentifierType(ref)
	switch fromType {
	case model.TypeSpec:
		if toType == model.TypeTest {
			return model.LinkTests
		}
		if toType == model.TypeContract {
			return model.LinkContract
		}
	case model.TypeTest:
		if toType == model.TypeContract {
			return model.LinkContractTests
		}
		return model.LinkImplements
	case model.TypeContract:
		return model.LinkContractImplements
	}
	return model.LinkReferences
}

func inferredRelationField(linkType model.LinkType) string {
	switch linkType {
	case model.LinkTests, model.LinkImplements, model.LinkVerifies:
		return "covers"
	case model.LinkContract, model.LinkContractTests, model.LinkContractImplements, model.LinkConstrainedBy, model.LinkProves:
		return "contracts"
	case model.LinkDesignedBy:
		return "components"
	case model.LinkDependsOn, model.LinkDependedBy:
		return "depends-on"
	case model.LinkSupersedes:
		return "supersedes"
	case model.LinkDeprecatedBy:
		return "deprecated-by"
	default:
		return "reference"
	}
}

func normalizeTraceRelationType(identifier *model.Identifier, link model.IdentifierLink) model.LinkType {
	switch link.Type {
	case model.LinkReferences:
		if strings.HasPrefix(link.Ref, "component:") {
			return model.LinkDesignedBy
		}
	case model.LinkContract:
		return model.LinkConstrainedBy
	case model.LinkImplements:
		if identifier.Type == model.TypeTest {
			return model.LinkVerifies
		}
	case model.LinkContractTests:
		return model.LinkProves
	}
	return link.Type
}
