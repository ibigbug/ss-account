# ss-account


[![CircleCI](https://circleci.com/gh/ibigbug/ss-account/tree/master.svg?style=svg)](https://circleci.com/gh/ibigbug/ss-account/tree/master)
[![codecov](https://codecov.io/gh/ibigbug/ss-account/branch/master/graph/badge.svg)](https://codecov.io/gh/ibigbug/ss-account)


## Introduction

ss-account is a TCP level proxy with user management, user port mapping, metrics and accounting features. More features like traffic control and powerful dashboard is still on the way.

## Dependency

- Redis: required.
- Prometheus: optional.
- Grafana: optional.
- Sentry: optional.

## Usage

### Run with docker

```shell
$ docker run -d --name ss-account --restart=unless-stopped --net=host ibigbug/ss-account:latest ./app
```

## Parameters

```
  -bind string
        management server listening address (default "0.0.0.0:9000")
  -port-range string
        accounting port range, e.g: 30000-40000 (default "30000-40000")
  -redis-db int
        redis database
  -redis-host string
        redis host (default "localhost")
  -redis-pass string
        redis password
  -redis-port string
        redis port (default "6379")
  -sentry-dsn string
        sentry DSN
```

## Roadmap

It's designed to be a wrapper for shadowsocks so that multipule users can share one shadowsocks server via indiviual ports, though it does work like a pure level 4 proxy (yet another HAProxy maybe).

- [ ] Management dashboard
- [ ] User register
- [ ] Traffic control/limit
- [ ] Speed limit