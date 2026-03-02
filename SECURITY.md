# Security Policy

## Reporting Security Issues

If you discover a security vulnerability in this SDK, please report it responsibly by opening a [GitHub issue](https://github.com/dominicnunez/opencode-sdk-go/issues) or contacting the maintainer directly.

For security issues related to the OpenCode platform itself, please follow [Anomaly Co's](https://anomaly.co) security reporting guidelines.

## Scope

This SDK communicates with a local or remote OpenCode server over HTTP. It can transmit authentication credentials through `Auth.Set` request payloads, while the OpenCode server manages credential validation and storage.

If you use `Auth.Set`, treat request bodies as sensitive data and avoid logging raw payloads that may include API keys or tokens. If you find issues related to request construction, response handling, or data exposure, please report them.

## Responsible Disclosure

Please allow reasonable time for investigation and remediation before public disclosure.
