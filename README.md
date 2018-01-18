
[![Go Report Card](https://goreportcard.com/badge/github.com/alejandro-carstens/golavel-cache)](https://goreportcard.com/report/github.com/alejandro-carstens/golavel-cache)
[![Build Status](https://travis-ci.org/alejandro-carstens/golavel-cache.svg?branch=master)](https://travis-ci.org/alejandro-carstens/golavel-cache)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hyn/multi-tenant/2.x/license.md)

This package allows you to implement a store agnostic caching system via a common interface by providing an abtraction layer between the different store clients and your application. The latter allows for each store to be used interchangeably without any code changes other than the programmatic configuration of the desired store(s). For a more detailed documentation please refer to the [godoc](https://godoc.org/github.com/alejandro-carstens/golavel-cache).

To start using this package in your application simply run:`go get github.com/alejandro-carstens/cache`

## Supported Stores

- Redis
- Memcache
- Map

## Future Stores 

- Go-Cache 
- Propose a store you would like to see implemented (Ex: BlockDB, LMDB, Mongo, MySQL, etc.)
- Build the implementation for the store you want and submit a PR (preferred)

## Contributing

Find an area you can help with and do it. Open source is about collaboration and open participation. Try to make your code look like what already exists or hopefully better and submit a pull request. Also, if you have any ideas on how to make the code better or on improving its scope and functionality please raise an issue and I will do my best to address it in a timely manner.

## TODO List:

- Documentation

