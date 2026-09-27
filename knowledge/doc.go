// Package knowledge represents and parses knowledge records.
//
// Knowledge is the source of truth in ORDR. A record captures what is known,
// why it is believed and which decision it prepares. Records never contain
// calculated rankings or scores; ordering lives in the graph package.
//
// The only record type in the proof of concept is [Opportunity], written as
// Markdown with a minimal front matter block. See README.md in the repository
// root for the file format.
//
// This package has no filesystem access. Callers pass file content and a
// source name used in error messages.
package knowledge
