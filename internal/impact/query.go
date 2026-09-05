package impact

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jingxu9x/idd-cli/internal/collector"
	"github.com/jingxu9x/idd-cli/internal/graph"
)

const (
	ReportSchema       = "idd.impacted.v1"
	defaultImpactDepth = 2
)

// Location identifies the canonical declaration of an impacted entity.
type Location struct {
	Path string `json:"path"`
	Line int    `json:"line,omitempty"`
}

// Reason explains the deterministic path from a Git change to a queue item.
type Reason struct {
	Kind         string `json:"kind"`
	Path         string `json:"path,omitempty"`
	Line         int    `json:"line,omitempty"`
	SourceID     string `json:"source_id,omitempty"`
	Relation     string `json:"relation,omitempty"`
	RecordID     string `json:"record_id,omitempty"`
	Field        string `json:"field,omitempty"`
	Conservative bool   `json:"conservative,omitempty"`
}

// QueueItem is one stable unit for semantic review.
type QueueItem struct {
	ID         string    `json:"id"`
	Kind       string    `json:"kind"`
	Lifecycle  string    `json:"lifecycle,omitempty"`
	Resolution string    `json:"resolution"`
	Canonical  *Location `json:"canonical,omitempty"`
	Reasons    []Reason  `json:"reasons"`
}

// Warning requests human review without claiming semantic staleness.
type Warning struct {
	Rule    string   `json:"rule"`
	Entity  string   `json:"entity"`
	Paths   []string `json:"paths"`
	Message string   `json:"message"`
}

// Report is the deterministic change-aware review queue.
type Report struct {
	Schema   string       `json:"schema"`
	Base     string       `json:"base"`
	Depth    int          `json:"depth"`
	Changes  []FileChange `json:"changes"`
	Queue    []QueueItem  `json:"queue"`
	Warnings []Warning    `json:"warnings"`
}

type seed struct {
	id             string
	reason         Reason
	evidenceChange bool
}

// BuildReport maps Git changes to indexed entities and traverses the typed
// relation graph in both directions. It never infers relationships from text.
// @implement SPEC-CMD_IDD_CLI-012
func BuildReport(changes *ChangeSet, index *graph.TraceIndex) (*Report, error) {
	return BuildReportWithAnalyses(changes, index, nil)
}

// BuildReportWithAnalyses uses the source models already collected by
// collector.TraceProject to map changed function bodies to their annotated
// declaration ranges without scanning source a second time.
// @implement SPEC-CMD_IDD_CLI-012
func BuildReportWithAnalyses(
	changes *ChangeSet,
	index *graph.TraceIndex,
	analyses []*collector.SourceAnalysis,
) (*Report, error) {
	if changes == nil {
		return nil, fmt.Errorf("git change set is nil")
	}
	if index == nil {
		return nil, fmt.Errorf("trace index is nil")
	}
	normalizedChanges := indexedAnnotationChanges(changes.Files, index)
	changedPaths := make(map[string]bool, len(changes.Files)*2)
	for _, change := range normalizedChanges {
		changedPaths[cleanPath(change.Path)] = true
		if change.OldPath != "" {
			changedPaths[cleanPath(change.OldPath)] = true
		}
	}

	seeds := make([]seed, 0)
	for _, change := range normalizedChanges {
		for _, id := range change.AnnotationIDs {
			if _, indexed := index.Entity(id); !indexed && change.Status != "deleted" {
				// Patch text may contain annotation examples inside strings. The
				// AST-backed project index is authoritative for current files;
				// raw unknown IDs are retained only when the whole source vanished.
				continue
			}
			seeds = append(seeds, seed{
				id:             id,
				reason:         Reason{Kind: "changed_annotation", Path: change.Path},
				evidenceChange: true,
			})
		}
	}
	seeds = append(seeds, occurrenceSeeds(normalizedChanges, index, analyses)...)
	sortSeeds(seeds)

	reasons := make(map[string][]Reason)
	distance := make(map[string]int)
	evidencePaths := make(map[string]bool)
	frontier := make([]string, 0)
	for _, value := range seeds {
		if _, found := distance[value.id]; !found {
			distance[value.id] = 0
			frontier = append(frontier, value.id)
		}
		reasons[value.id] = appendUniqueReason(reasons[value.id], value.reason)
		if value.evidenceChange && value.reason.Path != "" {
			evidencePaths[cleanPath(value.reason.Path)] = true
		}
	}
	sort.Strings(frontier)
	for level := 0; level < defaultImpactDepth && len(frontier) > 0; level++ {
		next := make([]string, 0)
		for _, id := range frontier {
			relations := append(index.Outbound(id), index.Inbound(id)...)
			sortRelations(relations)
			for _, relation := range relations {
				neighbor := relation.To
				if neighbor == id {
					neighbor = relation.From
				}
				if _, entityEndpoint := index.Entity(neighbor); !entityEndpoint {
					// Physical declaration/document endpoints participate in trace
					// provenance but are not logical IDD review queue entities.
					continue
				}
				reason := Reason{
					Kind:     "related_entity",
					Path:     relation.Provenance.Path,
					Line:     relation.Provenance.Line,
					SourceID: id,
					Relation: string(relation.Type),
					RecordID: relation.Provenance.RecordID,
					Field:    relation.Provenance.Field,
				}
				reasons[neighbor] = appendUniqueReason(reasons[neighbor], reason)
				if _, found := distance[neighbor]; !found {
					distance[neighbor] = level + 1
					next = append(next, neighbor)
				}
			}
		}
		frontier = uniqueSorted(next)
	}

	queue := make([]QueueItem, 0, len(distance))
	for id := range distance {
		entity, found := index.Entity(id)
		if !found {
			if isPublicIDDID(id) {
				queue = append(queue, QueueItem{ID: id, Kind: "unknown", Resolution: "unresolved", Reasons: sortedReasons(reasons[id])})
			}
			continue
		}
		if !reviewableKind(entity.Kind) {
			continue
		}
		item := QueueItem{
			ID: id, Kind: entity.Kind, Lifecycle: entity.Lifecycle,
			Resolution: string(entity.Resolution), Reasons: sortedReasons(reasons[id]),
		}
		if entity.Owner != nil {
			item.Canonical = &Location{Path: entity.Owner.Path, Line: entity.Owner.Line}
		}
		queue = append(queue, item)
	}
	sort.Slice(queue, func(i, j int) bool {
		if kindRank(queue[i].Kind) != kindRank(queue[j].Kind) {
			return kindRank(queue[i].Kind) < kindRank(queue[j].Kind)
		}
		return queue[i].ID < queue[j].ID
	})
	warnings := unchangedSpecWarnings(queue, changedPaths, evidencePaths)

	return &Report{
		Schema: ReportSchema, Base: changes.Base, Depth: defaultImpactDepth,
		Changes: normalizedChanges, Queue: queue, Warnings: warnings,
	}, nil
}

func indexedAnnotationChanges(changes []FileChange, index *graph.TraceIndex) []FileChange {
	result := make([]FileChange, 0, len(changes))
	for _, change := range changes {
		copy := change
		copy.Hunks = append([]Hunk(nil), change.Hunks...)
		copy.AnnotationIDs = make([]string, 0, len(change.AnnotationIDs))
		for _, id := range change.AnnotationIDs {
			if _, indexed := index.Entity(id); indexed || change.Status == "deleted" {
				copy.AnnotationIDs = append(copy.AnnotationIDs, id)
			}
		}
		result = append(result, copy)
	}
	return result
}

func occurrenceSeeds(
	changes []FileChange,
	index *graph.TraceIndex,
	analyses []*collector.SourceAnalysis,
) []seed {
	analysisByPath := make(map[string]*collector.SourceAnalysis, len(analyses))
	for _, analysis := range analyses {
		if analysis != nil {
			analysisByPath[cleanPath(analysis.Path)] = analysis
		}
	}
	result := make([]seed, 0)
	for _, change := range changes {
		path := cleanPath(change.Path)
		analysis := analysisByPath[path]
		matchedDeclaration := false
		if analysis != nil {
			for _, declaration := range analysis.Declarations {
				if !declarationIntersectsHunks(declaration.Line, declaration.EndLine, change.Hunks) {
					continue
				}
				for _, annotation := range declaration.Annotations {
					if annotation.Ignored || !annotation.Attached {
						continue
					}
					for _, id := range annotation.Refs {
						result = append(result, seed{
							id:             id,
							reason:         Reason{Kind: "changed_declaration", Path: path, Line: declaration.Line},
							evidenceChange: true,
						})
						matchedDeclaration = true
					}
				}
			}
		}
		// Documentation changes seed every structured occurrence in that
		// canonical file. A source file falls back to the same conservative
		// behavior only when no annotated declaration range was hit.
		if analysis != nil && matchedDeclaration {
			continue
		}
		for _, occurrence := range index.OccurrencesIn(path, 0, 0) {
			evidence := occurrence.Kind == graph.OccurrenceImplementation ||
				occurrence.Kind == graph.OccurrenceTestEvidence
			result = append(result, seed{
				id: occurrence.EntityID,
				reason: Reason{
					Kind:         "changed_occurrence",
					Path:         occurrence.Path,
					Line:         occurrence.Line,
					RecordID:     occurrence.RecordID,
					Field:        occurrence.Field,
					Conservative: analysis != nil && evidence,
				},
				evidenceChange: evidence,
			})
		}
	}
	return result
}

func declarationIntersectsHunks(start, end int, hunks []Hunk) bool {
	for _, hunk := range hunks {
		if hunk.NewLines <= 0 {
			continue
		}
		hunkEnd := hunk.NewStart + hunk.NewLines - 1
		if start <= hunkEnd && end >= hunk.NewStart {
			return true
		}
	}
	return false
}

func isPublicIDDID(id string) bool {
	return strings.HasPrefix(id, "SPEC-") || strings.HasPrefix(id, "TEST-") ||
		strings.HasPrefix(id, "contract:") || strings.HasPrefix(id, "component:")
}

func unchangedSpecWarnings(queue []QueueItem, changedPaths, evidencePaths map[string]bool) []Warning {
	if len(evidencePaths) == 0 {
		return []Warning{}
	}
	paths := make([]string, 0, len(evidencePaths))
	for path := range evidencePaths {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	warnings := make([]Warning, 0)
	for _, item := range queue {
		if item.Kind != "spec" || item.Canonical == nil || changedPaths[cleanPath(item.Canonical.Path)] {
			continue
		}
		warnings = append(warnings, Warning{
			Rule: "related-spec-unchanged", Entity: item.ID, Paths: append([]string(nil), paths...),
			Message: "implementation or test evidence changed while the related canonical SPEC file did not; review intent and evidence alignment",
		})
	}
	return warnings
}

func reviewableKind(kind string) bool {
	switch strings.ToLower(kind) {
	case "spec", "test", "contract", "component":
		return true
	default:
		return false
	}
}

func kindRank(kind string) int {
	switch kind {
	case "spec":
		return 0
	case "test":
		return 1
	case "contract":
		return 2
	case "component":
		return 3
	default:
		return 4
	}
}

func cleanPath(path string) string {
	path = filepath.ToSlash(filepath.Clean(path))
	return strings.TrimPrefix(path, "./")
}

func sortSeeds(values []seed) {
	sort.Slice(values, func(i, j int) bool {
		if values[i].id != values[j].id {
			return values[i].id < values[j].id
		}
		return reasonKey(values[i].reason) < reasonKey(values[j].reason)
	})
}

func sortRelations(values []graph.Relation) {
	sort.Slice(values, func(i, j int) bool {
		left, right := values[i], values[j]
		return fmt.Sprintf("%s\x00%s\x00%s\x00%s\x00%09d", left.Type, left.From, left.To, left.Provenance.Path, left.Provenance.Line) <
			fmt.Sprintf("%s\x00%s\x00%s\x00%s\x00%09d", right.Type, right.From, right.To, right.Provenance.Path, right.Provenance.Line)
	})
}

func appendUniqueReason(values []Reason, value Reason) []Reason {
	key := reasonKey(value)
	for _, current := range values {
		if reasonKey(current) == key {
			return values
		}
	}
	return append(values, value)
}

func sortedReasons(values []Reason) []Reason {
	result := append([]Reason(nil), values...)
	sort.Slice(result, func(i, j int) bool { return reasonKey(result[i]) < reasonKey(result[j]) })
	return result
}

func reasonKey(value Reason) string {
	return fmt.Sprintf("%s\x00%s\x00%09d\x00%s\x00%s\x00%s\x00%s\x00%t", value.Kind, value.Path, value.Line, value.SourceID, value.Relation, value.RecordID, value.Field, value.Conservative)
}

func uniqueSorted(values []string) []string {
	sort.Strings(values)
	result := values[:0]
	for _, value := range values {
		if len(result) == 0 || result[len(result)-1] != value {
			result = append(result, value)
		}
	}
	return result
}
