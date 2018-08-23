DP Identity API
==============

### Getting started

#### MongoDB
* Run ```brew install mongodb```
* Run ```brew services start mongodb```

#### kafka
* Run ```brew install kafka```
* Run ```brew services start zookeeper```
* Run ```brew services start kafka```

### Configuration

| Environment variable        | Default                                   | Description
| --------------------------- | ----------------------------------------- | -----------
| BIND_ADDR                   |                                          | The host and port to bind to
| MONGODB_BIND_ADDR           | localhost:20111                          | The MongoDB bind address
| MONGODB_DATABASE            | identities                               | The MongoDB dataset database
| MONGODB_COLLECTION          | identities                               | MongoDB collection
| HEALTHCHECK_INTERVAL       | 30s                                       | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_TIMEOUT         | 2s                                     | The timeout that the healthcheck allows for checked subsystems

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2016-2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
