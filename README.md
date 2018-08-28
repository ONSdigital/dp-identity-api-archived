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

### Usage

`make debug` to run locally.

### DEV NOTES

1.) Have added basic endpoints for GetIdentityByID and PostIdentity `('identities/{id}' and '/identities'` respectively), this was just to sanity check the database and auditing
functionality was working as intended. Delete/modify/purge as needed.

2.) The auditor will trigger an error on GetIdentityByID for any requests that lack the correct context fields. This is expected behaviour and not a bug.

### Configuration

| Environment variable        | Default                                   | Description
| --------------------------- | ----------------------------------------- | -----------
| BIND_ADDR                   |                                           | The host and port to bind to
| MONGODB_BIND_ADDR           | localhost:23800                           | The MongoDB bind address
| MONGODB_DATABASE            | identities                                | The MongoDB dataset database
| MONGODB_COLLECTION          | identities                                | MongoDB collection
| HEALTHCHECK_INTERVAL        | 30s                                       | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_TIMEOUT         | 2s                                        | The timeout that the healthcheck allows for checked subsystems

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2016-2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
