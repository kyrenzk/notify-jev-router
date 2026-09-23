# Security policy

## Reporting a vulnerability

Please report security issues privately to the repository maintainers via GitHub **Security Advisories** on this repository, or by opening a confidential issue if advisories are unavailable.

Do not file public issues for exploitable vulnerabilities.

## Response expectations

- We aim to acknowledge reports within 5 business days.
- We will coordinate disclosure and patches before public detail when appropriate.

## Data handling expectations

- **Do not put PII** (email, phone, names, message bodies, tokens) into `RoutingContext` or Jev state. The library validates and rejects common PII fields.
- Host services remain responsible for redaction before calling `Resolve`.
- Production logs should not dump full Jev state; use `PlanMetadata` and destination lists only.

## Supported versions

| Version | Supported |
|---------|-----------|
| 0.1.x   | Yes       |
