// Getting paid: the payments bet and its parts.
// Under record a, "supports: [b]" reads "a supports b" and
// "moreValuableThan: [b]" reads "a ranks above b on value".

payments: {
	// Card payments and bank matching are the two ways to deliver the bet.
	contains: ["card-payments", "bank-transfer-matching"]

	moreValuableThan: ["multi-currency", "mobile-app", "accountant-portal"]
}

"card-payments": {
	// The payment provider converts currencies, so foreign clients can pay.
	enables: ["multi-currency"]

	moreValuableThan: ["bank-transfer-matching"]
	lessUncertainThan: ["bank-transfer-matching"]
	lessComplexThan: ["bank-transfer-matching"]
	moreComplexThan: ["payment-reminders"]
}

"bank-transfer-matching": {
	// Reconciled payments are the input for a correct VAT summary.
	enables: ["tax-reports"]
}

"payment-reminders": {
	// The reminder beta shows late payment is the pain worth solving.
	supports: ["payments"]

	moreValuableThan: ["accountant-portal"]
	lessUncertainThan: ["accountant-portal", "receipt-capture"]
	lessComplexThan: ["multi-currency"]
}
