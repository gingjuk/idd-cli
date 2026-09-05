package impact

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

// @test TEST-CMD_IDD_CLI-008
func TestParseNameStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []FileChange
		wantErr bool
	}{
		{
			name:  "ordinary statuses",
			input: "M\x00internal/a.go\x00A\x00docs/a.md\x00D\x00old.go\x00",
			want: []FileChange{
				{Path: "internal/a.go", Status: "modified", Hunks: []Hunk{}},
				{Path: "docs/a.md", Status: "added", Hunks: []Hunk{}},
				{Path: "old.go", Status: "deleted", Hunks: []Hunk{}},
			},
		},
		{
			name:  "rename preserves both paths",
			input: "R097\x00old name.go\x00new name.go\x00",
			want:  []FileChange{{Path: "new name.go", OldPath: "old name.go", Status: "renamed", Hunks: []Hunk{}}},
		},
		{name: "truncated ordinary entry", input: "M\x00", wantErr: true},
		{name: "truncated rename entry", input: "R100\x00old.go\x00", wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseNameStatus([]byte(test.input))
			if (err != nil) != test.wantErr {
				t.Fatalf("parseNameStatus() error = %v, wantErr %t", err, test.wantErr)
			}
			if !test.wantErr && !reflect.DeepEqual(got, test.want) {
				t.Errorf("parseNameStatus() = %#v, want %#v", got, test.want)
			}
		})
	}
}

// @test TEST-CMD_IDD_CLI-008
func TestParsePatch(t *testing.T) {
	tests := []struct {
		name      string
		patch     string
		wantHunks []Hunk
		wantIDs   []string
	}{
		{
			name:      "addition deletion and default counts",
			patch:     "@@ -8 +8,2 @@\n-// @implement SPEC-OLD_OWNERSHIP-001\n+// @implement SPEC-NEW_OWNERSHIP-002, SPEC-NEW_OWNERSHIP-001\n+func changed() {}\n",
			wantHunks: []Hunk{{OldStart: 8, OldLines: 1, NewStart: 8, NewLines: 2}},
			wantIDs:   []string{"SPEC-NEW_OWNERSHIP-001", "SPEC-NEW_OWNERSHIP-002", "SPEC-OLD_OWNERSHIP-001"},
		},
		{
			name:      "deleted test contract annotation",
			patch:     "@@ -4,2 +4,0 @@\n-// @test-contract TEST-AUTH-009\n-func removed() {}\n",
			wantHunks: []Hunk{{OldStart: 4, OldLines: 2, NewStart: 4, NewLines: 0}},
			wantIDs:   []string{"TEST-AUTH-009"},
		},
		{
			name:      "metadata is not an annotation",
			patch:     "--- a/spec.go\n+++ b/spec.go\n@@ -1 +1 @@\n-package old\n+package current\n",
			wantHunks: []Hunk{{OldStart: 1, OldLines: 1, NewStart: 1, NewLines: 1}},
			wantIDs:   []string{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hunks, ids := parsePatch([]byte(test.patch))
			if !reflect.DeepEqual(hunks, test.wantHunks) {
				t.Errorf("parsePatch() hunks = %#v, want %#v", hunks, test.wantHunks)
			}
			if !reflect.DeepEqual(ids, test.wantIDs) {
				t.Errorf("parsePatch() IDs = %#v, want %#v", ids, test.wantIDs)
			}
		})
	}
}

// @test TEST-CMD_IDD_CLI-008
func TestCollectGitChangesIncludesTrackedAndUntracked(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init", "-q")
	runGit(t, root, "config", "user.name", "IDD Test")
	runGit(t, root, "config", "user.email", "idd@example.invalid")
	writeFixture(t, filepath.Join(root, "service.go"), "package fixture\n\nfunc service() {}\n")
	runGit(t, root, "add", "service.go")
	runGit(t, root, "commit", "-q", "-m", "baseline")
	writeFixture(t, filepath.Join(root, "service.go"), "package fixture\n\n// @implement SPEC-FIXTURE-001\nfunc service() {}\n")
	writeFixture(t, filepath.Join(root, "service_test.go"), "package fixture\n\n// @test TEST-FIXTURE-001\nfunc example() {}\n")

	changes, err := CollectGitChanges(context.Background(), root, "HEAD")
	if err != nil {
		t.Fatalf("CollectGitChanges() error = %v", err)
	}
	if changes.Schema != ChangeSetSchema || len(changes.Files) != 2 {
		t.Fatalf("CollectGitChanges() = %#v", changes)
	}
	want := map[string][]string{
		"service.go":      {"SPEC-FIXTURE-001"},
		"service_test.go": {"TEST-FIXTURE-001"},
	}
	for _, change := range changes.Files {
		if !reflect.DeepEqual(change.AnnotationIDs, want[change.Path]) {
			t.Errorf("%s annotation IDs = %#v, want %#v", change.Path, change.AnnotationIDs, want[change.Path])
		}
		if len(change.Hunks) == 0 {
			t.Errorf("%s has no changed line range", change.Path)
		}
	}
}

func runGit(t *testing.T, root string, args ...string) {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = root
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}

func writeFixture(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", path, err)
	}
}
