---
id: card-payments
title: Card payments via payment provider
state: evidence
---

## Summary

Accept card payments on the invoice page through an external payment provider.

## Hypothesis

Clients pay card-payable invoices at least a week earlier than invoices paid by bank transfer.

## Evidence

- A concierge test with 12 accounts sending manual payment links cut median time to paid to 9 days.
- The payment provider supports hosted checkout, so card data never touches our servers.

## Uncertainties

- Who carries the card fee: the freelancer or the client.
- Whether the provider's onboarding rejects many sole traders.

## Decision

Choose the payment provider and fee model, or stop after the concierge test.
