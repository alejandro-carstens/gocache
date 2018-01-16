## Golavel Cache

[![Go Report Card](https://goreportcard.com/badge/github.com/alejandro-carstens/golavel-cache)](https://goreportcard.com/report/github.com/alejandro-carstens/golavel-cache)
[![Build Status](https://travis-ci.org/alejandro-carstens/golavel-cache.svg?branch=master)](https://travis-ci.org/alejandro-carstens/golavel-cache)

Inspired by the Laravel Cache System, Golavel Cache allows you to implement a store agnostic caching system 
via a common interface by providing an abtraction layer between the different store clients and your application so that each store can be use interchangeably without having to modify any code, but just the programmatic configuration of the desired store. 

## Supported Stores

- Redis
- Memcache
- Array

## Future Stores 

- Go-Cache 
- Propose a store you would like to see implemented (Ex: BlockDB, LMDB, Mongo, MySQL, etc.)
- Build the implementation yourself for the store you want and submit a PR

## TODO List:

- Documentation

