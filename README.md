## gocache

[![Go Report Card](https://goreportcard.com/badge/github.com/alejandro-carstens/gocache)](https://goreportcard.com/report/github.com/alejandro-carstens/gocache)
[![Build Status](https://travis-ci.org/alejandro-carstens/gocache.svg?branch=master)](https://travis-ci.org/alejandro-carstens/gocache)
[![GoDoc](https://godoc.org/github.com/alejandro-carstens/golavel-cache?status.svg)](https://godoc.org/github.com/alejandro-carstens/gocache)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/alejandro-carstens/golavel-cache/blob/master/LICENSE)

Some data retrieval performed by your application could be CPU intensive or take several seconds to complete. For these cases, it is common to cache the retrieved data for a period of time so that it can be retrieved quickly on subsequent requests for the same data. The cached data is usually stored in a very fast data store such as [Memcached](https://memcached.org) or [Redis](https://redis.io).

This package allows you to implement a store agnostic caching system via an expressive and unified interface by providing an abstraction layer between different data store drivers and your application. This allows for each store to be used interchangeably without any code changes other than the programmatic configuration of the desired store(s).

## Table Of Contents
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
    - [Obtaining A Cache Instance](#obtaining-a-cache-instance)
    - [Retrieving Items From The Cache](#retrieving-items-from-the-cache)
    - [Storing Items In The Cache](#storing-items-in-the-cache)
    - [Removing Items From The Cache](#removing-items-from-the-cache)
- [Cache Tags](#cache-tags)
    - [Storing Cache Tagged Items](#storing-cache-tagged-items)
    - [Accessing Cache Tagged Items](#accessing-cache-tagged-items)
    - [Removing Tagged Cache Items](#removing-tagged-cache-items)
- [Atomic Locks](#atomic-locks)
- [Contributing](#contributing)
- [Liscense](#liscense)

## Installation
To start using this package in your application simply run:`go get github.com/alejandro-carstens/gocache`

## Configuration

This package supports 3 backends out of the box: [Redis](https://redis.io), [Memcached](https://memcached.org) and Local (via [go-cache](https://github.com/patrickmn/go-cache)). Each store has a specific configuration whose parameters can be easily referenced in the following [GoDoc](https://pkg.go.dev/github.com/alejandro-carstens/gocache) sections:
- [RedisConfig](https://pkg.go.dev/github.com/alejandro-carstens/gocache#RedisConfig)
- [MemcacheConfig](https://pkg.go.dev/github.com/alejandro-carstens/gocache#MemcacheConfig)
- [LocalConfig](https://pkg.go.dev/github.com/alejandro-carstens/gocache#LocalConfig)

## Usage

### Obtaining A Cache Instance
In order to new up a cache implementation simply call ```gocache.New``` with the desired configuration: 
```go
// Redis
cache, err := gocache.New(&RedisConfig{
    Prefix: "gocache:",
    Addr:   "localhost:6379",
})
// handle err

// Memcache
cache, err := gocache.New(&MemcacheConfig{
    Prefix:  "gocache:",
    Servers: []string{"127.0.0.1:11211"},
})
// handle err

// Local
cache, err := gocache.New(&LocalConfig{
    Prefix:          "gocache:",
    DefaultInterval: time.Second,
})
// handle err
```

### Retrieving Items From The Cache

All methods including the prefix `Get` are used to retrieve items from the cache. In an item does not exist in the cache for the given key an error of type ```gocache.ErrNotFound``` will be raised. Please see the following examples:

```go

v, err := cache.GetFloat32("temperature")
// handle err

v, err := cache.GetFloat64("height")
// handle err

v, err := cache.GetInt("score")
// handle err

v, err := cache.GetInt64("counter")
// handle err

v, err := cache.GetString("username")
// handle err

v, err := cache.GetUint64("id")
// handle err

// Get any type e.g. Movie{Name string, Views int64}
var m Movie
err := cache.Get("e.t.", &m)
// handle err

// Handle missed entry for key
v, err := cache.GetString("entry-not-found-key")
if errors.Is(gocache.ErrNotFound, err) {
    // handle err
}
```
The method ```Many``` is also exposed in order to retrieve multiple cache records with one call. The results of the ```Many``` invocation will be returned in a map of [gocache.Item](https://pkg.go.dev/github.com/alejandro-carstens/gocache#Item) instances keyed by the retrieved cached entries keys. Please see the example below:

```go
items, err := cache.Many("string", "uint64", "int", "int64", "float64", "float32", "any")
// handle err

for key, item := range items {
    switch key:
    case "string":
        v, err := item.String()
        // handle err
    case "uint64":
        v, err := item.Uint64()
        // handle err
    case "int":
        v, err := item.GetInt()
        // handle err
    case "int64":
        v, err := item.GetInt64()
        // handle err
    case "float64":
        v, err := item.GetFloat64()
        // handle err
    case "float32":
        v, err := item.GetFloat32()
        // handle err
    case "any":
        var m Movie
        err := item.Get(&m)
        // handle err
    }
}

```
### Storing Items In The Cache
You can use the ```Put``` method to store items in the cache with a specified time to live:
```go
err := cache.Put("key", "value", 10 * time.Second)
// handle err

// You can store any value
err := cache.Put("most_watched_movie", &Movie{
    Name:  "Avatar",
    Views: 100,
}, 60 * time.Minute)
// handle err
```

To store a value indefinitely (without expiration) simply use the method ```Forever```:
```go
err := cache.Forever("key", "value")
// handle err
```

To store many values at once you can use ```PutMany```:
```go
var entries = []gocache.Entry{
    {
        Key:   "string",
        Value: "whatever",
        Duration: time.Minute,
    },
    {
        Key:   "any",
        Value: Movie{
          Name:  "Star Wars",
          Views: 10,
        },
        Duration: time.Minute,
    },
}
err := cache.PutMany(entries...)
// handle err
```
To increment and decrement values (for now you can only increment using ```int64``` values) simply use ```Increment``` & ```Decrement```. Please note that if there is not entroy for the key being incremented the initial value will be 0 and the entry will set to not expire:
```go
val, err := cache.Increment("a", 1) // a = 1
// handle err

val, err := cache.Increment("a", 10) // a = 11
// handle err

val, err := cache.Decrement("a", 2) // a = 9
// handle err

val, err := cache.Decrement("b", 5) // b = -5
// handle err
```

### Removing Items From The Cache
You may remove items from the cache using the ```Forget``` method:
```go
// Note that res will be true if the cache entry was removed and false 
// if no entry was for the given key
res, err := cache.Forget("key") 
// handle err
```
If you want to clear all entries from the cache you can use the ```Flush``` method:
```go
err := cache.Flush()
// handle err
```
## Cache Tags

### Storing Cache Tagged Items
Cache tags allow you to tag related items in the cache and then flush all cached values that have been assigned a given tag. You may access a tagged cache by passing in an ordered sliced of tag names. For example, let's access a tagged cache and put a value into the cache:
```go
err := cache.Tags("person", "artist").Put("John", "Doe", time.Minute)
// handle err

err := cache.Tags("person", "accountant").Put("Jane", "Doe", time.Minute)
// handle err
```
### Accessing Cache Tagged Items
To retrieve a tagged cache item, pass the same ordered list of tags to the tags method and then call the any of the methods shown in the [Retrieving Items From The Cache](#retrieving-items-from-the-cache) section above:
```go
v, err := cache.Tags("person", "artist").GetString("John")
// handle err

v, err := cache.Tags("person", "accountant").GetString("Jane")
// handle err
```
### Removing Tagged Cache Items
You may flush all items that are assigned a tag or list of tags. For example, this statement would remove all caches tagged with either people, authors, or both. So, both Jane and John would be removed from the cache:
```go
err := cache.Tags("person", "accountant").Flush()
// handle err
```
In contrast, this statement would remove only cached values tagged with authors, so Anne would be removed, but not John:
```go
err := cache.Tags("accountant").Flush()
// handle err
```

In addition you can also call ```Forget```:
```go
res, err := cache.Tags("person", "accountant").Forget("Jane")
// handle err
```

<b>Important Note:</b> with the exception of the [Redis](https://redis.io) driver, when calling `Flush` with tags, the underlying entries won't be deleted so please make sure to set expiration values when using tags and flushing.

## Atomic Locks

Atomic locks allow for the manipulation of distributed locks without worrying about race conditions. An example of this would be that you only want one process to work on one object at a time, such as the same file should only be uploaded one at a time (we do not want the same file to be uploaded more than once at any given time). You may create and manage locks using the via the ```Lock``` method:
```go
var (
    lock          = cache.Lock("merchant_1", "pid_1", 30 * time.Second)
    acquired, err = lock.Acquire()
)
if err != nil {
    // handle err
}
if acquired {
    defer func() {
        // released will be true if the lock was released before expiration
        released, err := lock.Release()
        // handle err
    }
    // do something here
}
```

## Contributing

Find an area you can help with and do it. Open source is about collaboration and open participation. Try to make your code look like what already exists or hopefully better and submit a pull request. Also, if you have any ideas on how to make the code better or on improving its scope and functionality please raise an issue and I will do my best to address it in a timely manner.

## Liscense

MIT.
