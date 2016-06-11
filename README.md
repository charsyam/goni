# goniplus
Application metric collector for [Goni](https://github.com/goniapm/goni)

[![Go Report Card](https://goreportcard.com/badge/github.com/goniapm/goniplus)](https://goreportcard.com/report/github.com/goniapm/goniplus) [![Build Status](https://travis-ci.org/goniapm/goniplus.svg?branch=develop)](https://travis-ci.org/goniapm/goniplus) [![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/goniapm/goniplus) [![codecov](https://codecov.io/gh/goniapm/goniplus/branch/develop/graph/badge.svg)](https://codecov.io/gh/goniapm/goniplus)

## Install

``` go
import "github.com/goniapm/goniplus"
```

### Get API Key
* If you use [hosted service](https://dashboard.goniapm.io),
    * Signup and create new project.
* If you use self hosted(private) service,
    * Go to hosted website dashboard, signup, and create new project.
    * For more information, please check [Goni wiki](https://github.com/goniapm/goni/wiki/Install-Guide).

## Metrics

Goni+ can collect
* CPU Usage (linux only)
* Expvar Metrics
* HTTP Metrics
* Runtime Metrics

Some metrics can be collected manually.
* Breadcrumb (HTTP request track)
* Error

## Middleware Support
Goni is specialized to web application monitoring, so Goni+ supports multiple web frameworks. Supported web frameworks are listed below.

* [net/http](https://github.com/goniapm/goni/wiki/goniplus-nethttp)
* [Gin](https://github.com/goniapm/goniplus-gin)

## Example
* [net/http](https://github.com/goniapm/goni/wiki/goniplus-nethttp)
* [Gin](https://github.com/goniapm/goniplus-gin)

## Links
* [Goni](https://github.com/goniapm/goni)
