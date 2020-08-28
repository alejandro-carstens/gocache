package gocache

import "testing"

func TestPutGetInt64WithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		tags := tag()
		if err := cache.Tags(tags).Put("key", 100, 1); err != nil {
			t.Fatal(err)
		}

		got, err := cache.Tags(tags).GetInt64("key")
		if err != nil {
			t.Error(err.Error())
		}
		if got != int64(100) {
			t.Error("Expected 100, got ", got)
		}
		if _, err := cache.Tags(tags).Forget("key"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPutGetFloat64WithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		var expected float64

		expected = 9.99

		tags := tag()
		if err := cache.Tags(tags).Put("key", expected, 1); err != nil {
			t.Fatal(err)
		}

		got, err := cache.Tags(tags).GetFloat64("key")
		if err != nil {
			t.Error(err.Error())
		}
		if got != expected {
			t.Error("Expected 9.99, got ", got)
		}
		if _, err := cache.Tags(tags).Forget("key"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestIncrementWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		tags := tag()
		if _, err := cache.Tags(tags).Increment("increment_key", 1); err != nil {
			t.Fatal(err)
		}
		if _, err := cache.Tags(tags).Increment("increment_key", 1); err != nil {
			t.Fatal(err)
		}

		got, err := cache.Tags(tags).GetInt64("increment_key")
		if err != nil {
			t.Error(err.Error())
		}

		var expected int64 = 2
		if got != expected {
			t.Error("Expected 2, got ", got)
		}
		if _, err := cache.Tags(tags).Forget("increment_key"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDecrementWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		tags := tag()
		if _, err := cache.Tags(tags).Increment("decrement_key", 2); err != nil {
			t.Fatal(err)
		}
		if _, err := cache.Tags(tags).Decrement("decrement_key", 1); err != nil {
			t.Fatal(err)
		}

		var expected int64 = 1

		got, err := cache.Tags(tags).GetInt64("decrement_key")
		if err != nil {
			t.Error(err.Error())
		}
		if got != expected {
			t.Error("Expected "+string(expected)+", got ", got)
		}
		if _, err := cache.Tags(tags).Forget("decrement_key"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestForeverWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		expected := "value"

		tags := tag()
		if err := cache.Tags(tags).Forever("key", expected); err != nil {
			t.Fatal(err)
		}

		got, err := cache.Tags(tags).GetString("key")
		if err != nil {
			t.Error(err.Error())
		}
		if got != expected {
			t.Error("Expected "+expected+", got ", got)
		}
		if _, err := cache.Tags(tags).Forget("key"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPutGetManyWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		tags := tag()

		keys := make(map[string]string)

		keys["key_1"] = "value"
		keys["key_2"] = "100"
		keys["key_3"] = "9.99"

		if err := cache.Tags(tags).PutMany(keys, 10); err != nil {
			t.Fatal(err)
		}

		resultKeys := make([]string, 3)

		resultKeys[0] = "key_1"
		resultKeys[1] = "key_2"
		resultKeys[2] = "key_3"

		results, err := cache.Tags(tags).Many(resultKeys)
		if err != nil {
			t.Error(err.Error())
		}

		for i := range results {
			if results[i] != keys[i] {
				t.Error(i, results[i])
			}
		}

		if _, err := cache.Tags(tags).Flush(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPutGetWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		tags := make([]string, 3)

		tags[0] = "tag1"
		tags[1] = "tag2"
		tags[2] = "tag3"

		var firstExample example
		firstExample.Name = "Alejandro"
		firstExample.Description = "Whatever"

		if err := cache.Tags(tags...).Put("key", firstExample, 10); err != nil {
			t.Fatal(err)
		}

		var newExample example

		if err := cache.Tags(tags...).Get("key", &newExample); err != nil {
			t.Error(err.Error())
		}
		if newExample != firstExample {
			t.Error("The structs are not the same", newExample)
		}
		if _, err := cache.Forget("key"); err != nil && err.Error() != MemcacheNilErrorResponse {
			t.Fatal(err)
		}
	}
}

func TestTagSet(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		tagSet := cache.Tags("Alejandro").GetTags()

		namespace, err := tagSet.getNamespace()
		if err != nil {
			t.Error(err.Error())
		}
		if len([]rune(namespace)) != 20 {
			t.Error("The namespace is not 20 chars long.", namespace)
		}
		if got := tagSet.reset(); got != nil {
			t.Error("Reset did not return nil.", got)
		}
	}
}

func tag() string {
	return "tag"
}
