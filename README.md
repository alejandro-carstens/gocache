## gocache

[![Go Report Card](https://goreportcard.com/badge/github.com/alejandro-carstens/gocache)](https://goreportcard.com/report/github.com/alejandro-carstens/gocache)
[![Build Status](https://travis-ci.org/alejandro-carstens/gocache.svg?branch=master)](https://travis-ci.org/alejandro-carstens/gocache)
[![GoDoc](https://godoc.org/github.com/alejandro-carstens/golavel-cache?status.svg)](https://godoc.org/github.com/alejandro-carstens/gocache)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/alejandro-carstens/golavel-cache/blob/master/LICENSE)

This package allows you to implement a store agnostic caching system via a common interface by providing an abstraction layer between different store clients and your application. The latter allows for each store to be used interchangeably without any code changes other than the programmatic configuration of the desired store(s). For a more detailed documentation please refer to the [godoc](https://godoc.org/github.com/alejandro-carstens/gocache).


## Contributing

Find an area you can help with and do it. Open source is about collaboration and open participation. Try to make your code look like what already exists or hopefully better and submit a pull request. Also, if you have any ideas on how to make the code better or on improving its scope and functionality please raise an issue and I will do my best to address it in a timely manner.

## Usage

To start using this package in your application simply run:`go get github.com/alejandro-carstens/gocache`

Set the params for the store you want:

```go
```

New up the cache by passing the store name and the appropiate params:

```go
```

Start using it:

```go
```

Use it with structs:

```go
```

Use it with tags:

```go
```

For more examples please refer to the tests.

## Supported Stores

- Redis
- Memcache
- Map

## Future Stores 

- Go-Cache 
- Propose a store you would like to see implemented (Ex: BlockDB, LMDB, Mongo, MySQL, etc.)
- Build the implementation for the store you want and submit a PR (preferred) 

## Testing

Run ```go test -v```

It is important to note that one must install the required stores or comment out the ones you do not want to test. Since this is an abstraction layer, <b>WHEN CONTRIBUTING YOU SHOULD NOT ADD OR MODIFY TESTS</b> just make your implementation conform to what is already there. However if you can make the tests better please do modify them.

## Liscense

MIT.
