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

type Region struct {
	Snippet   Snippet `json:"snippet,omitempty"`
	StartLine int     `json:"startLine,omitempty"`
}

type ArtifactLocation struct {
	URI string `json:"uri,omitempty"`
}

type ContextRegion struct {
	Snippet   Snippet `json:"snippet,omitempty"`
	EndLine   int     `json:"endLine,omitempty"`
	StartLine int     `json:"startLine,omitempty"`
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
