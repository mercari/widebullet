# Widebullet

Widebullet is [JSON-RPC](http://www.jsonrpc.org/) base API gateway server. It implements [JSON-RPC batch](http://www.jsonrpc.org/specification#batch) endponts with extended format for HTTP REST request (see [SPEC](/SPEC.md)). For example, it receives one single JSON-RPC array which defines multiple HTTP requests and converts it into multiple concurrent HTTP requests. If you have multiple backend microservices and need to request them at same time for one transaction, Widebullet simplifies it.

# Status

Production ready.

# Requirement

Widebullet requires Go1.7.3 or later.

# Installation

Widebullet provides a executable named `wbt` to kick server. To install `wbt`, use `go get`,

```
$ go get -u github.com/mercari/widebullet/...
```

# Usage

To run `wbt`, you must provide configuration path via `-c` option (See [CONFIGURATION.md](/CONFIGURATION.md)) about details and [`config/example.toml`](/config/example.toml) for example usage.

```
$ wbt -c config/example.toml
```

Use `-help` to see more options.


# Configuration

See [CONFIGURATION.md](/CONFIGURATION.md) about details.

# Specification

See [SPEC.md](/SPEC.md) about details.

# Comitters

 * Tatsuhiko Kubo [@cubicdaiya](https://github.com/cubicdaiya)

# Contribution

Please read the CLA below carefully before submitting your contribution.

https://www.mercari.com/cla/

# License

Copyright 2016 Mercari, Inc.

Licensed under the MIT License.
