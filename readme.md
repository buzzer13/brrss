# BrRSS

[![latest release](https://gitlab.com/buzzer13/brrss/-/badges/release.svg)](https://gitlab.com/buzzer13/brrss/-/releases)

HTML-to-RSS bridge that can be used as a serverless function or run as a server.

# Getting started

[//]: # (## Releases)

[//]: # ()
[//]: # (1. Download a binary for your OS from the [GitLab releases]&#40;https://gitlab.com/buzzer13/healthpose/-/releases&#41; page.)

[//]: # (2. Prepare a [configuration file]&#40;#configuration&#41; and put it in the supported directory.)

## Container

1. Pull and run `registry.gitlab.com/buzzer13/brrss:latest` image (tag can be either `latest`, or a specific version like `v1.0.0`).
    - Command: `docker run --name=brrss -it "registry.gitlab.com/buzzer13/brrss:latest"`
    - Server doesn't enable an authentication by default, so you may want to check [configuration](#configuration) section.

# Configuration

## Environment variables

- `API_KEY` - when set - enables key auth for the API and requires `api-key` parameter to be present in every query.
- `API_USERNAME`, `API_PASSWORD` - when both are set - enables basic auth for the API.

Both `API_KEY` and `API_PASSWORD` can be either the plain string or bcrypt-compatible hash.
