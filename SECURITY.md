# Security Policy

## Status of this library

This library is unaudited. No claim of an audit or production proof is made.

## Supported Versions

Use this table to determine which versions are currently being supported with
security fixes.

| Version          | Supported          |
| ---------------- | ------------------ |
| v0.1.0           | :white_check_mark: |
| Earlier versions | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in this library, please report it via
GitHub Security Advisories.

### What to Include

- **Title**: Clear, concise description of the issue
- **Context**: How the vulnerability can be exploited
- **Impact**: What an attacker could achieve
- **Reproduction Steps**: Minimal steps to reproduce the issue
- **Expected vs Actual Behavior**: What should happen vs what happens
- **Test Keys/Entries**: If applicable, the test vectors that demonstrate the
  issue

### What Is Out of Scope

- Documentation improvements or typos
- Feature requests
- Performance optimizations
- Use of the library in unsafe ways (this library is designed for correct
  signing; misuse is out of scope)

### Response Timeline

- **Initial response**: Within 5 business days
- **Triage**: Within 10 business days — maintainers will confirm, categorize,
  and assign severity
- **Fix Etch**: Target within 30 business days for critical issues, 90 days for
  medium
- **Public Disclosure**: Coordinated disclosure with the reporter, typically 90
  days after a fix is released

### Reporting Format

Use the GitHub Security Advisory form when creating the issue. The form includes
sections for:

1. **Advisory Link**: Will be auto-generated
2. **What to include**: As described above
3. **Out of scope**: As described above
4. **Supported versions**: As described above
5. **Acknowledgements**: Reporters will be acknowledged (unless requesting
   anonymity)

### Key Management

This library does not manage cryptographic keys. Signers are an interface, and
the only shipped implementation is an in-memory ed25519 signer. Key management,
keystores, KMS integration, and wallet connections are out of scope for this
library and should be handled by the caller.

## Security Scopes

The following types of issues are considered **critical** and should be reported
immediately:

- **Signature correctness bugs**: Any issue where signatures produced by the
  library do not verify on-chain, potentially moving money to unintended
  recipients
- **Wire format deviations**: Any change to the XDR wire format that breaks
  compatibility with the Stellar network or golden vectors
- **Replay vulnerabilities**: Issues that allow a valid signature to be replayed
  in a different context
- **Credential arm mismanagement**: Issues where the wrong credential arm is
  used or written

The following are **not** considered security issues:

- The library being "un audited" or "unproven" (this is explicitly stated in the
  software)
- Missing features (passkey support, secp256r1, etc.) - these are contributor
  issues
- Documentation typos or formatting issues
- CLI usability improvements
- Benchmark results

## Security Contact

- **GitHub Security Advisories**: Use the form at
  `https://github.com/soroauth/soroauth-go/security/advisories/new`
- **Email**: Not applicable - use the GitHub form only

## Verifying a Fix

After a fix is released, the following should pass:

1. `go test ./...` (unit tests, no network)
2. `go test -run TestReadmeSnippetsMatchTheirSource -v` (README snippets match
   source)
3. `make vectors-check` (golden vectors byte-for-byte match)
4. `go vet ./...` (no vet issues)

All three must pass before a security fix is considered verified.
