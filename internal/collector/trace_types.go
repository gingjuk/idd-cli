package collector

// TraceSchema is the stable schema emitted by the top-level trace command.
const TraceSchema = "idd.trace.v1"

// TraceOptions bounds relationship traversal and optional occurrence classes.
type TraceOptions struct {
	Depth           int  `json:"depth"`
	IncludeMentions bool `json:"include_mentions"`
}

// TraceQuery records the normalized query which produced a dossier.
type TraceQuery struct {
	ID              string `json:"id"`
	Depth           int    `json:"depth"`
	IncludeMentions bool   `json:"include_mentions"`
}

// TraceDossier is the stable, ID-centered projection shared by trace and
// higher-level review workflows. Status is unresolved or ambiguous when an ID
// has zero or several canonical owners; a planned entity uses planned.
type TraceDossier struct {
	Schema            string                     `json:"schema"`
	Query             TraceQuery                 `json:"query"`
	Status            string                     `json:"status"`
	Entity            *TraceEntity               `json:"entity,omitempty"`
	Canonical         *TraceCanonical            `json:"canonical,omitempty"`
	Owners            []TraceOccurrence          `json:"owners"`
	Occurrences       []TraceOccurrence          `json:"occurrences"`
	OutboundRelations []TraceRelation            `json:"outbound_relations"`
	InboundRelations  []TraceRelation            `json:"inbound_relations"`
	RelatedEntities   []TraceRelatedEntity       `json:"related_entities"`
	Components        []TraceComponent           `json:"components"`
	Contracts         []TraceContract            `json:"contracts"`
	Implementations   []ReviewContextDeclaration `json:"implementations"`
	Tests             []TraceTest                `json:"tests"`
	Mentions          []TraceOccurrence          `json:"mentions"`
	Findings          []ReviewContextIssue       `json:"findings"`
	Truncated         bool                       `json:"truncated"`
}

// TraceEntity describes the logical identity independently from its physical
// occurrences.
type TraceEntity struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace,omitempty"`
	Lifecycle  string `json:"lifecycle,omitempty"`
	Resolution string `json:"resolution"`
}

// TraceLocation identifies an authored record or declaration.
type TraceLocation struct {
	Path string `json:"path"`
	Line int    `json:"line"`
}

// TraceCanonical is the single authoritative record and its authored
// semantics. Fields irrelevant to the entity kind remain absent.
type TraceCanonical struct {
	TraceLocation
	Title       string `json:"title,omitempty"`
	Package     string `json:"package,omitempty"`
	Purpose     string `json:"purpose,omitempty"`
	Ownership   string `json:"ownership,omitempty"`
	Boundary    string `json:"boundary,omitempty"`
	Decisions   string `json:"decisions,omitempty"`
	Guarantees  string `json:"guarantees,omitempty"`
	Requirement string `json:"requirement,omitempty"`
	Acceptance  string `json:"acceptance,omitempty"`
	Kind        string `json:"kind,omitempty"`
	Oracle      string `json:"oracle,omitempty"`
	Markdown    string `json:"markdown,omitempty"`
	Truncated   bool   `json:"truncated"`
}

// TraceOccurrence preserves every typed appearance of an entity.
type TraceOccurrence struct {
	ID           string `json:"id"`
	EntityID     string `json:"entity_id"`
	Kind         string `json:"kind"`
	Origin       string `json:"origin"`
	DocumentRole string `json:"document_role,omitempty"`
	Path         string `json:"path"`
	Line         int    `json:"line"`
	Field        string `json:"field,omitempty"`
	RecordID     string `json:"record_id,omitempty"`
}

// TraceProvenance identifies the exact authored field or annotation which
// produced a relation.
type TraceProvenance struct {
	OccurrenceID   string `json:"occurrence_id"`
	Path           string `json:"path"`
	Line           int    `json:"line"`
	Field          string `json:"field,omitempty"`
	RecordID       string `json:"record_id,omitempty"`
	OccurrenceKind string `json:"occurrence_kind,omitempty"`
}

// TraceRelation is a directed, typed edge. Depth is the shortest traversal
// distance from the queried entity.
type TraceRelation struct {
	From       string          `json:"from"`
	To         string          `json:"to"`
	Type       string          `json:"type"`
	Depth      int             `json:"depth"`
	Provenance TraceProvenance `json:"provenance"`
}

// TraceRelatedEntity is a bounded summary discovered while traversing.
type TraceRelatedEntity struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	Lifecycle  string `json:"lifecycle,omitempty"`
	Resolution string `json:"resolution"`
	Depth      int    `json:"depth"`
}

// TraceComponent preserves the architecture context of a related Component.
type TraceComponent struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Purpose   string `json:"purpose"`
	Ownership string `json:"ownership,omitempty"`
	Boundary  string `json:"boundary,omitempty"`
	Decisions string `json:"decisions,omitempty"`
	File      string `json:"file"`
	Line      int    `json:"line"`
	Markdown  string `json:"markdown"`
	Truncated bool   `json:"truncated"`
}

// TraceContract preserves the guarantees of a related Contract.
type TraceContract struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Guarantees string `json:"guarantees"`
	File       string `json:"file"`
	Line       int    `json:"line"`
	Markdown   string `json:"markdown"`
	Truncated  bool   `json:"truncated"`
}

// TraceTest preserves one authored TEST plus its executable declarations.
type TraceTest struct {
	ID           string                     `json:"id"`
	Title        string                     `json:"title"`
	Kind         string                     `json:"kind"`
	Purpose      string                     `json:"purpose"`
	Oracle       string                     `json:"oracle"`
	Covers       []string                   `json:"covers"`
	Contracts    []string                   `json:"contracts"`
	File         string                     `json:"file"`
	Line         int                        `json:"line"`
	Markdown     string                     `json:"markdown"`
	Declarations []ReviewContextDeclaration `json:"declarations"`
	Truncated    bool                       `json:"truncated"`
}
