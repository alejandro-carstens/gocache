## golavel-cache

[![Go Report Card](https://goreportcard.com/badge/github.com/alejandro-carstens/golavel-cache)](https://goreportcard.com/report/github.com/alejandro-carstens/golavel-cache)
[![Build Status](https://travis-ci.org/alejandro-carstens/golavel-cache.svg?branch=master)](https://travis-ci.org/alejandro-carstens/golavel-cache)
[![GoDoc](https://godoc.org/github.com/alejandro-carstens/golavel-cache?status.svg)](https://godoc.org/github.com/alejandro-carstens/golavel-cache)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/alejandro-carstens/golavel-cache/blob/master/LICENSE)

This package allows you to implement a store agnostic caching system via a common interface by providing an abstraction layer between the different store clients and your application. The latter allows for each store to be used interchangeably without any code changes other than the programmatic configuration of the desired store(s). For a more detailed documentation please refer to the [godoc](https://godoc.org/github.com/alejandro-carstens/golavel-cache).


## Contributing

<b>By using this package you are already contributing, if you would like to go a bit further simply give the project a star and spread the word (it would be greatly appreciated)</b>. Otherwise, find an area you can help with and do it. Open source is about collaboration and open participation. Try to make your code look like what already exists or hopefully better and submit a pull request. Also, if you have any ideas on how to make the code better or on improving its scope and functionality please raise an issue and I will do my best to address it in a timely manner.

## Usage

To start using this package in your application simply run:`go get github.com/alejandro-carstens/golavel-cache`

Set the params for the store you want:

```go

params := make(map[string]interface{})
  
// Redis
params["password"] = ""
params["database"] = 0
params["address"] = "localhost:6379"
params["prefix"] = "golavel"
  
// Memcache (you can spacify multiple servers)
params["server"] = "127.0.0.1:11211"
params["prefix"] = "golavel:"
  
// Map
params["prefix"] = "golavel"

```

New up the cache by passing the store name and the appropiate params:

```go

// Can be any of the following: "redis", "memcache" or "map"
store := "redis"

c, err := cache.New(store, params)
```

Start using it:

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
	// may return a string, an int64, 
	// or a float64 depending on 
	// the value type
	val, err := c.Get("foo") 

	if err != nil {
		fmt.Print(val) // bar
	}

	// Delete a k-v pair
	c.Forget("foo")

	// Remember a value forever
	c.Forever("baz", "buz")

	// Flush the cache
	c.Flush()
}
```

Use it with structs:

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

Use it with tags:

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

	params["prefix"] = "golavel"

	c, err := cache.New("map", params)

	var foo Foo

	foo.Name = "Alejandro"
	foo.Description = "Whatever"
	
	c.Tags("tag").Forever("foo", foo, 10)

	var bar Foo

	// Retrieve a struct from the cache
	val, err := c.Tags("tag").GetStruct("foo", &bar)

	if err != nil {
		fmt.Print(bar.Name)        // Alejandro
		fmt.Print(bar.Description) // Whatever
	}
	
	tags := make([]string, 3)
	
	tags[0] = "tag1"
	tags[1] = "tag2"
	tags[2] = "tag3"
	
	// Put a value in the cache for 10 mins.
	c.Tags(tags...).Put("foo", "bar", 10)

	val, err := c.Tags(tags...).Get("foo") 

	if err != nil {
		fmt.Print(val) // bar
	}

	// Delete a k-v pair
	c.Tags(tags...).Forget("foo")
	
	c.Tags(tags...).Flush()
}
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
