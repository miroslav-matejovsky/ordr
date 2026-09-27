// Typed relationships. {from: a, to: b} under a kind reads "a <kind> b".
relations: {
	// Repair meetup attendance is evidence for tool borrowing demand.
	supports: [{from: "beta", to: "alpha"}]

	// A shared inventory service needs at least one running library first.
	enables: [{from: "alpha", to: "gamma"}]

	relates: [{from: "beta", to: "gamma"}]
}
