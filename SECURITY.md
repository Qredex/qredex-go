# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| 0.1.x   | ✅ Yes     |

Older versions are not supported.  Please upgrade to the latest release before reporting a vulnerability.

---

## Reporting a Vulnerability

**Do not open a public GitHub issue for security vulnerabilities.**

Please report security issues by emailing **security@qredex.com** with:

1. A description of the vulnerability and its potential impact
2. Steps to reproduce or a minimal proof-of-concept
3. Affected versions
4. Any suggested mitigations you are aware of

We will acknowledge your report within **2 business days** and aim to provide an initial assessment within **5 business days**.

---

## Responsible Disclosure

We ask that you:

- Give us reasonable time to investigate and release a fix before any public disclosure
- Avoid accessing, modifying, or deleting data that does not belong to you during research
- Do not disrupt or degrade production services

We will credit researchers in release notes (with your permission) once a fix is published.

---

## Security Design Notes

The Qredex Go SDK is designed with the following security properties:

- **No secrets in logs** — Client secrets, access tokens, IITs, and PITs are never written to any logger by default.
- **HTTPS enforced** — The `Production` and `Staging` environments use HTTPS-only base URLs.  The `Development` environment is intentionally plaintext for local use only.
- **Short-lived tokens** — OAuth access tokens are cached with a 30-second safety buffer before their expiry.  Tokens are never persisted to disk.
- **Integrations API only** — The SDK does not expose Merchant API or Internal API endpoints, reducing the blast radius of a compromised credential.
- **No plaintext credential logging** — The `Bootstrap()` and `New()` constructors validate credentials but never echo them back.

---

## Contact

- **Security email**: security@qredex.com
- **General contact**: os@qredex.com
- **Website**: https://qredex.com
