---
id: bank-transfer-matching
title: Automatic bank transfer matching
state: investigation
---

## Summary

Read the freelancer's bank feed and mark invoices as paid when a matching transfer arrives.

## Hypothesis

Most clients keep paying by bank transfer, so removing manual reconciliation saves more time than card payments.

## Evidence

- Support logs show 18 percent of tickets are about marking invoices as paid.

## Uncertainties

- Whether open banking access is available in all countries we serve.
- How often transfers carry a usable invoice reference.

## Decision

Decide whether bank matching is a separate bet or part of the payments bet.
