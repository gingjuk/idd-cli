---
idd:
  version: "1.0"
  package: pkg/pattern
  document: testing
---

# Testing: pkg/pattern

## TEST-PKG_PATTERN-001: Strict identifier grammar cases

- **Kind:** `contract`
- **Covers:** `SPEC-PKG_PATTERN-001`, `SPEC-PKG_PATTERN-009`
- **Contracts:** `IdentifierSyntax`

**Purpose:**

Prove accepted TYPE-MODULE-NUMBER forms and rejected segment, type, module, and
number forms with a table whose expected validity is the oracle. Internal
PATTERN/WALK section markers are explicitly rejected.

**Oracle:** The test passes only when its assertions confirm accepted
TYPE-MODULE-NUMBER forms and rejected segment, type, module, and number forms with a
table whose expected validity is the oracle. Internal PATTERN/WALK section markers are
explicitly rejected.

## TEST-PKG_PATTERN-002: Broad pattern recognition cases

- **Kind:** `test`
- **Covers:** `SPEC-PKG_PATTERN-001`, `SPEC-PKG_PATTERN-007`

**Purpose:**

Prove that built-in IDD patterns are recognized, including CONTRACT and DESIGN,
while an unknown type returns an error. The test intentionally exercises broad
recognition rather than strict whole-string acceptance.

**Oracle:** The test passes only when its assertions confirm
built-in IDD patterns are recognized, including CONTRACT and DESIGN, while an unknown
type returns an error. The test intentionally exercises broad recognition rather than
strict whole-string acceptance.

## TEST-PKG_PATTERN-003: Identifier type detection

- **Kind:** `contract`
- **Covers:** `SPEC-PKG_PATTERN-006`
- **Contracts:** `DocumentReferenceSyntax`

**Purpose:**

Prove type detection for every registered identifier family, case-insensitive
recognition, unknown text, and empty input.

**Oracle:** The test passes only when its assertions confirm type
detection for every registered identifier family, case-insensitive recognition, unknown
text, and empty input.

## TEST-PKG_PATTERN-004: Annotation prefix mapping

- **Kind:** `test`
- **Covers:** `SPEC-PKG_PATTERN-002`, `SPEC-PKG_PATTERN-003`

**Purpose:**

Prove exact mapping of the three supported prefixes and rejection of differently
cased, unknown, and empty prefixes.

**Oracle:** The test passes only when its assertions confirm exact
mapping of the three supported prefixes and rejection of differently cased, unknown, and
empty prefixes.

## TEST-PKG_PATTERN-005: Quoted reference extraction

- **Kind:** `test`
- **Covers:** `SPEC-PKG_PATTERN-004`

**Purpose:**

Prove that plain identifier-shaped prose is ignored while backtick- and
quote-adjacent identifiers are returned. Cases also distinguish three-part IDs
from internal two-part section markers and verify empty input.

**Oracle:** The test passes only when its assertions confirm plain
identifier-shaped prose is ignored while backtick- and quote-adjacent identifiers are
returned. Cases also distinguish three-part IDs from internal two-part section markers
and verify empty input.

## TEST-PKG_PATTERN-006: Annotation list splitting

- **Kind:** `test`
- **Covers:** `SPEC-PKG_PATTERN-005`

**Purpose:**

Prove comma splitting, whitespace trimming, trailing-comma handling, and empty
input behavior using exact returned slices as the oracle.

**Oracle:** The test passes only when its assertions confirm comma
splitting, whitespace trimming, trailing-comma handling, and empty input behavior using
exact returned slices as the oracle.

## TEST-PKG_PATTERN-007: Mixed annotation extraction

- **Kind:** `test`
- **Covers:** `SPEC-PKG_PATTERN-002`, `SPEC-PKG_PATTERN-008`

**Purpose:**

Prove extraction from one annotation, multiple annotation kinds, and a
comma-separated target list, while unrelated content returns no identifiers.

**Oracle:** The test passes only when its assertions confirm extraction
from one annotation, multiple annotation kinds, and a comma-separated target list, while
unrelated content returns no identifiers.

## TEST-PKG_PATTERN-008: Quote adjacency helper

- **Kind:** `contract`
- **Covers:** `SPEC-PKG_PATTERN-004`
- **Contracts:** `AnnotationSyntax`

**Purpose:**

Prove the current lexical quote rule for backticks, double quotes, unquoted
text, and partial identifier matches. This test documents adjacency behavior;
it does not claim full Markdown delimiter parsing.

**Oracle:** The test passes only when its assertions confirm the current
lexical quote rule for backticks, double quotes, unquoted text, and partial identifier
matches. This test documents adjacency behavior; it does not claim full Markdown
delimiter parsing.

## Strategy

All cases are pure and table-driven, so fixtures consist only of input strings
and expected values or error presence. No regular expression internals are
asserted; tests observe public lexical behavior.

Important exclusions are deterministic ordering across the map-backed pattern
registry, duplicate removal, quoted lowercase extraction, empty module/number
segments, and balanced Markdown quoting. Those gaps are recorded so a future
grammar tightening starts with explicit failing cases rather than silently
changing collector behavior.
