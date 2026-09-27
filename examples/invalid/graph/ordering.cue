// Invalid on purpose: a value cycle and a reference to an unknown record.
comparisons: value: [
	{more: "alpha", than: "beta"},
	{more: "beta", than: "gamma"},
	{more: "gamma", than: "alpha"},
	{more: "alpha", than: "delta"},
]
