// Graph statements per opportunity. Under record a, "supports: [b]" reads
// "a supports b" and "moreValuableThan: [b]" reads "a ranks above b on value".
// Only stated and transitively implied comparisons are used; nothing is scored.

alpha: {
	// A shared inventory service needs at least one running library first.
	enables: ["gamma"]

	moreValuableThan: ["beta", "gamma"]
	moreComplexThan: ["beta"]
}

beta: {
	// Repair meetup attendance is evidence for tool borrowing demand.
	supports: ["alpha"]
	relates: ["gamma"]

	moreUncertainThan: ["alpha"]
}

gamma: {
	moreUncertainThan: ["beta"]
	moreComplexThan: ["alpha"]
}
