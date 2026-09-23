package sarif

import (
	"time"
)

type Sarif struct {
	Runs    []Run  `json:"runs,omitempty"`
	Version string `json:"version,omitempty"`
	Schema  string `json:"$schema,omitempty"`
}

type Rules struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	HelpURI string `json:"helpUri,omitempty"`
}

type Driver struct {
	Name  string  `json:"name,omitempty"`
	Rules []Rules `json:"rules,omitempty"`
}

type Tool struct {
	Driver Driver `json:"driver,omitempty"`
}

type WorkingDirectory struct {
	URI string `json:"uri,omitempty"`
}

type Invocation struct {
	Arguments           []string         `json:"arguments,omitempty"`
	ExecutionSuccessful bool             `json:"executionSuccessful"`
	CommandLine         string           `json:"commandLine,omitempty"`
	EndTimeUtc          time.Time        `json:"endTimeUtc,omitempty"`
	WorkingDirectory    WorkingDirectory `json:"workingDirectory,omitempty"`
}

type Conversion struct {
	Tool       Tool       `json:"tool,omitempty"`
	Invocation Invocation `json:"invocation,omitempty"`
}

type Invocations struct {
	ExecutionSuccessful bool             `json:"executionSuccessful,omitempty"`
	EndTimeUtc          time.Time        `json:"endTimeUtc,omitempty"`
	WorkingDirectory    WorkingDirectory `json:"workingDirectory,omitempty"`
}

type Properties struct {
}

type Message struct {
	Text string `json:"text,omitempty"`
}

type Snippet struct {
	Text string `json:"text,omitempty"`
}

// Region.Snippet and ContextRegion.Snippet are *Snippet (not Snippet):
// encoding/json's "omitempty" only recognizes false, 0, a nil
// pointer/interface, and an empty array/slice/map/string as "empty" -- a
// struct value is never considered empty, so a plain Snippet field would
// always be serialized as "snippet":{} even when unset. Per the SARIF
// 2.1.0 schema, "snippet" is an artifactContent object that must carry
// "text" or "binary" when present at all; an empty object fails
// validation and causes strict consumers (GitLab's SARIF ingestion,
// GitHub's upload-sarif) to reject the whole report. Using a pointer
// lets omitempty actually omit the key until real snippet text exists.
type Region struct {
	Snippet   *Snippet `json:"snippet,omitempty"`
	StartLine int      `json:"startLine,omitempty"`
}

type ArtifactLocation struct {
	URI string `json:"uri,omitempty"`
}

type ContextRegion struct {
	Snippet   *Snippet `json:"snippet,omitempty"`
	EndLine   int      `json:"endLine,omitempty"`
	StartLine int      `json:"startLine,omitempty"`
}

type PhysicalLocation struct {
	Region           Region           `json:"region,omitempty"`
	ArtifactLocation ArtifactLocation `json:"artifactLocation,omitempty"`
	ContextRegion    ContextRegion    `json:"contextRegion,omitempty"`
}

type Locations struct {
	PhysicalLocation PhysicalLocation `json:"physicalLocation,omitempty"`
}

type ResultsProperties struct {
	IssueConfidence string `json:"issue_confidence,omitempty"`
	IssueSeverity   string `json:"issue_severity,omitempty"`
}

type Results struct {
	Message    Message           `json:"message,omitempty"`
	Level      string            `json:"level,omitempty"`
	Locations  []Locations       `json:"locations,omitempty"`
	Properties ResultsProperties `json:"properties,omitempty"`
	RuleID     string            `json:"ruleId,omitempty"`
	RuleIndex  int               `json:"ruleIndex,omitempty"`
}

type Run struct {
	Tool        Tool          `json:"tool,omitempty"`
	Conversion  Conversion    `json:"conversion,omitempty"`
	Invocations []Invocations `json:"invocations,omitempty"`
	Properties  Properties    `json:"properties,omitempty"`
	Results     []Results     `json:"results,omitempty"`
}
