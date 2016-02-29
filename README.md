[![Build Status](https://travis-ci.org/yageek/apiversion.svg?branch=master)](https://travis-ci.org/yageek/apiversion)
[![GoDoc](https://godoc.org/github.com/yageek/apiversion?status.png)](https://godoc.org/github.com/yageek/apiversion)  [![Report Cart](http://goreportcard.com/badge/yageek/apiversion)](http://goreportcard.com/report/yageek/apiversion)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

# Versioning API with Go

Simple API versioning with Go with or without any middleware stack.

# Installation

```
go get -v github.com/yageek/apiversion
```

# How it works

## Vendor identifier

Before tagging a version, an API is identified by a vendor identifier with
the form `application/vnd.mybusiness.com`.
If it matches the value of the `Accept` header, the handler will next try
find a version that suits your needs.

## Version

AS you can registered several versions for your API, if you don't specify one,
the last registered will be selected by default.

To select, one particular version, the handler will match the version specified
after within the `Accept` header after the vendor identifier.

For example, to match the `v2` version, the request `Accept` header will be:

```
Accept: application/vnd.mybusiness.com-v1
```

# Usage

The package works out of the box as a simple HTTP handler or can be integrated
easily with any middleware stack.

For usage, see examples folder.
