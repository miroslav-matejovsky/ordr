// Invalid on purpose: a value cycle and a reference to an unknown record.
alpha: moreValuableThan: ["beta", "delta"]
beta: moreValuableThan: ["gamma"]
gamma: moreValuableThan: ["alpha"]
