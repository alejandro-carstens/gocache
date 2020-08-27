package gocache

import (
	"fmt"
	"github.com/rs/xid"
	"strings"
)

// TagSet is the representation of a tag set for the caching stores
type TagSet struct {
	Store Store
	Names []string
}

// Reset resets the tag set
func (ts *TagSet) reset() error {
	for i, name := range ts.Names {
		id, err := ts.resetTag(name)
		if err != nil {
			return err
		}

		ts.Names[i] = id
	}

	return nil
}

// GetNamespace gets the current TagSet namespace
func (ts *TagSet) getNamespace() (string, error) {
	tagsIds, err := ts.tagIds()
	if err != nil {
		return "", err
	}

	return strings.Join(tagsIds, "|"), err
}

func (ts *TagSet) tagIds() ([]string, error) {
	tagIds := make([]string, len(ts.Names))

	for i, name := range ts.Names {
		val, err := ts.tagId(name)
		if err != nil {
			return tagIds, err
		}

		tagIds[i] = val
	}

	return tagIds, nil
}

func (ts *TagSet) tagId(name string) (string, error) {
	value, err := ts.Store.GetString(ts.tagKey(name))
	if err != nil && !isCacheMissedError(err) {
		return "", err
	}
	if len(value) == 0 {
		return ts.resetTag(name)
	}

	return fmt.Sprint(value), nil
}

func (ts *TagSet) tagKey(name string) string {
	return "tag:" + name + ":key"
}

func (ts *TagSet) resetTag(name string) (string, error) {
	id := xid.New().String()

	return id, ts.Store.Forever(ts.tagKey(name), id)
}
