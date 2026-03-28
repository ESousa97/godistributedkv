# Security Policy

## Supported Versions

We only provide security updates for the following versions of **godistributedkv**:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of **godistributedkv** seriously. If you believe you have found a security vulnerability, please report it to us responsibly.

**Please do not report security vulnerabilities through public GitHub issues.**

### Responsible Disclosure Process

1.  **Email**: Send an email to [security@example.com](mailto:security@example.com) with the details of the vulnerability.
2.  **Details**: Please include as much information as possible, including:
    *   Type of vulnerability (e.g., buffer overflow, SQL injection, etc.).
    *   Steps to reproduce the issue.
    *   Potential impact.
    *   Any suggested fixes or mitigations.
3.  **Acknowledgement**: We will acknowledge receipt of your report within 48 hours.
4.  **Verification**: We will investigate and verify the vulnerability.
5.  **Fix**: We will work on a fix and coordinate a release.
6.  **Public Disclosure**: Once a fix is released, we will publicly disclose the vulnerability with proper attribution to you (if desired).

## Security Best Practices

To keep your installation secure:

- Always run the latest supported version.
- Use TLS/SSL for all gRPC communication in production (where applicable).
- Limit access to the administrative ports.

Thank you for helping keep our community safe!
