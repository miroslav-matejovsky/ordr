// Package graph models explicit relationships between knowledge records.
//
// Two kinds of statements exist:
//
//   - [Comparison]: a pairwise relative judgment on one [Dimension], for
//     example "alpha has more strategic value than beta". No numeric scores exist.
//   - [Relation]: a typed link such as supports, blocks, enables, invalidates
//     or relates.
//
// [New] validates statements against the known identifiers and rejects unknown
// references, self-references, duplicates, contradictory direct comparisons
// and ordering cycles. A valid [Graph] derives a partial order per dimension
// with [Graph.Order]; it never invents an order the statements do not imply.
//
// The model is independent of any notation. Package cuegraph reads it from CUE.
package graph
