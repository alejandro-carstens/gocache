## Golavel Cache

[![Go Report Card](https://goreportcard.com/badge/github.com/alejandro-carstens/golavel-cache)](https://goreportcard.com/report/github.com/alejandro-carstens/golavel-cache)
[![Build Status](https://travis-ci.org/alejandro-carstens/golavel-cache.svg?branch=master)](https://travis-ci.org/alejandro-carstens/golavel-cache)

Golavel Cache allows you to implement a store agnostic caching system via a common interface by providing an abtraction layer between the store clients and your application. The latter allows for each store to be use interchangeably without having to change code other than the programmatic configuration of the desired store(s). For a more detailed documentation please refer to the [godoc](https://godoc.org/github.com/alejandro-carstens/golavel-cache).

To start using Golavel Cache in your application simply run:`go get github.com/alejandro-carstens/golavel-cache`

## Supported Stores

- Redis
- Memcache
- Map

## Future Stores 

- Go-Cache 
- Propose a store you would like to see implemented (Ex: BlockDB, LMDB, Mongo, MySQL, etc.)
- Build the implementation for the store you want and submit a PR (preferred)

## TODO List:

- Documentation

