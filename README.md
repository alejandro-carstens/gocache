## golavel-cache

[![Go Report Card](https://goreportcard.com/badge/github.com/alejandro-carstens/golavel-cache)](https://goreportcard.com/report/github.com/alejandro-carstens/golavel-cache)
[![Build Status](https://travis-ci.org/alejandro-carstens/golavel-cache.svg?branch=master)](https://travis-ci.org/alejandro-carstens/golavel-cache)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hyn/multi-tenant/2.x/license.md)

This package allows you to implement a store agnostic caching system via a common interface by providing an abtraction layer between the different store clients and your application. The latter allows for each store to be used interchangeably without any code changes other than the programmatic configuration of the desired store(s). For a more detailed documentation please refer to the [godoc](https://godoc.org/github.com/alejandro-carstens/golavel-cache).

To start using this package in your application simply run:`go get github.com/alejandro-carstens/golavel-cache`

## Usage

Set the params for the store you want

```go

params := make(map[string]interface{})
  
// Redis
params["password"] = ""
params["database"] = 0
params["address"] = "localhost:6379"
params["prefix"] = "golavel"
  
//Memcache (you can spacify multiple servers)
params["server"] = "127.0.0.1:11211"
params["prefix"] = "golavel:"
  
//Map
params["prefix"] = "golavel"

```

New up the cache by passing the store name and the appropiate params

```go

// Can be any of the following: "redis", "memcache" or "map"
store := "redis"

c, err := cache.New(store, params)
```

Example:

```go
package main

import (
	"fmt"
	"github.com/alejandro-carstens/golavel-cache"
)

func main() {
	params := make(map[string]interface{})

	params["password"] = ""
	params["database"] = 0
	params["address"] = "localhost:6379"
	params["prefix"] = "golavel"

	c, err := cache.New("redis", params)

	// Put a value in the cache for 10 mins.
	c.Put("foo", "bar", 10)

	// Retrieve a value from the cache,
	// may return a string, an int64, or a float64
	// depending on the value type)
	val, err := c.Get("foo")

	if err != nil {
		fmt.Print(val)
	}

	// Delete the k-v pair
	c.Forget("foo")

	// Remember the value forever
	c.Forever("baz", "buz")

	// Flush the cache
	c.Flush()
}



```

Example with Structs

```go
package main

import (
	"fmt"
	"github.com/alejandro-carstens/golavel-cache"
)

type Foo struct {
	Name        string
	Description string
}

func main() {
	params := make(map[string]interface{})

	params["server"] = "127.0.0.1:11211"
	params["prefix"] = "golavel"

	c, err := cache.New("memcache", params)

	var foo Foo

	foo.Name = "Alejandro"
	foo.Description = "Whatever"

	c.Put("foo", foo, 10)

	var bar Foo

	// Retrieve a struct from the cache
	val, err := c.GetStruct("foo", &bar)

	if err != nil {
		fmt.Print(bar.Name)        // Alejandro
		fmt.Print(bar.Description) // Whatever
	}
}
```

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

