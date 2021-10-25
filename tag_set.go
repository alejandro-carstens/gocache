package gocache

import (
	"fmt"
	"strings"

	"github.com/rs/xid"
)

// tagSet is the representation of a tag set for the caching stores
type tagSet struct {
	store store
	names []string
}

// Reset resets the tag set
func (ts *tagSet) reset() error {
	for i, name := range ts.names {
		id, err := ts.resetTag(name)
		if err != nil {
			return err
		}

		ts.names[i] = id
	}

	return nil
}

// GetNamespace gets the current tagSet namespace
func (ts *tagSet) getNamespace() (string, error) {
	tagsIds, err := ts.tagIds()
	if err != nil {
		return "", err
	}

	return strings.Join(tagsIds, "|"), err
}

func (ts *tagSet) tagIds() ([]string, error) {
	tagIds := make([]string, len(ts.names))
	for i, name := range ts.names {
		val, err := ts.tagId(name)
		if err != nil {
			return tagIds, err
		}

		tagIds[i] = val
	}

	return tagIds, nil
}

func (ts *tagSet) tagId(name string) (string, error) {
	value, err := ts.store.GetString(ts.tagKey(name))
	if err != nil && !isErrNotFound(err) {
		return "", err
	}
	if len(value) == 0 {
		return ts.resetTag(name)
	}

	return fmt.Sprint(value), nil
}

func (ts *tagSet) tagKey(name string) string {
	return "tag:" + name + ":key"
}

func (ts *tagSet) resetTag(name string) (string, error) {
	id := xid.New().String()

	return id, ts.store.Forever(ts.tagKey(name), id)
}
