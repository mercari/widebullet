# Configuration for Widebullet

A configuration file format for Widebullet is [TOML](https://github.com/toml-lang/toml).

A configuration for Widebullet has some sections. A example is [here](config/example.toml).

 * [Global Section](#core-section)
 * [Endpoints Section](#endpoints-section)

## Global Section

|name               |type  |description                                 |default         |note                                                              |
|-------------------|------|--------------------------------------------|----------------|------------------------------------------------------------------|
|Port               |string|port number or unix socket path             |29300           |e.g.)29300, unix:/tmp/wbt.sock <br/> `-p` option can overwrite    |
|LogLevel           |string|log-level                                   |error           |                                                                  |
|Timeout            |int   |timeout for proxing request                 |5               |unit is second                                                    |
|MaxIdleConnsPerHost|int   |maximum idle to keep per-host               |100             |                                                                  |
|DisableCompression |bool  |delete `Accept-Encoding: gzip` in header    |false           |                                                                  |

## Endpoints Section

|name           |type          |description                        |default|note|
|---------------|--------------|-----------------------------------|-------|----|
|Name           |string        |Endpoint name                      |       |    |
|Ep             |string        |Endpoint URL                       |       |    |
|ProxySetHeaders|array of array|Headers appended on proxing request|       |    |

See [SPEC.md](SPEC.md) about details for APIs.
