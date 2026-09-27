// Pairwise comparisons. {more: a, than: b} reads "a ranks above b".
// Only stated and transitively implied comparisons are used; nothing is scored.
comparisons: {
	// Which opportunity is more strategically valuable?
	value: [
		{more: "alpha", than: "beta"},
		{more: "alpha", than: "gamma"},
	]

	// Which opportunity is more uncertain?
	uncertainty: [
		{more: "gamma", than: "beta"},
		{more: "beta", than: "alpha"},
	]

	// Which opportunity is more complex?
	complexity: [
		{more: "gamma", than: "alpha"},
		{more: "alpha", than: "beta"},
	]
}
