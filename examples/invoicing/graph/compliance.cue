// Tax and legal obligations that shape the invoice format.

"e-invoicing": {
	// The invoice data model changes for the mandate; a second currency
	// field should not be designed before that change is known.
	blocks: ["multi-currency"]
	relates: ["accountant-portal"]

	moreValuableThan: ["mobile-app", "tax-reports"]
	moreUncertainThan: ["tax-reports"]
	moreComplexThan: ["tax-reports", "multi-currency"]
}

"tax-reports": {
	// Accountants would consume the VAT summary through the portal.
	enables: ["accountant-portal"]

	moreValuableThan: ["accountant-portal"]
}
