# Integration lib

It aims to be a modular library for integration.
There is no strict limitation what to include or not.

Currently, includes (see [docs](https://godoc.org/github.com/reddec/integration-lib)):

- cloud
  - google sheets
- communication
  - email
  - telegram
- finance
  - exchange rates (fiat and crypto)
  - tinkoff investing

# Style guide

## Plain functions

- Should have docs
- Docs should contain API provider(s) if applicable
- Should have context and context-less (optional) version if uses external services
- Errors should be wrapped and tagged


## Structured

In addition to plain function

- Configuration structure should not contain state
- Should have default constructor that fills parameters from [environment](https://github.com/caarlos0/env)
- Environment variables should have prefix, related to the provider (package)
- Configuration structure could be used as container for functions if functions are not modifying state
- It's ok to panic in default constructor