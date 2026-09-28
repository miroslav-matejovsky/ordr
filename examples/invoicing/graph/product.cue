// Product reach: mobile and expenses.

"receipt-capture": {
	// Interviews show receipts are forwarded by email, not photographed,
	// which undermines the need for a native app.
	invalidates: ["mobile-app"]
	supports: ["tax-reports"]
}

"mobile-app": {
	relates: ["receipt-capture"]

	moreUncertainThan: ["receipt-capture", "accountant-portal"]
	moreComplexThan: ["receipt-capture"]
}
