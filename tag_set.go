package gocache

import (
	"fmt"
	"strings"

	"github.com/rs/xid"
)

// TagSet is the representation of a set of tags to be used to interact with the caching stores
type TagSet struct {
	store store
	names []string
}

// Reset will reinitialize the value of all tags in the TagSet
func (ts *TagSet) Reset() error {
	for i, name := range ts.names {
		id, err := ts.resetTag(name)
		if err != nil {
			return err
		}

		ts.names[i] = id
	}

	return nil
}

// Flush clears removes all tags associated with the TagSet from the cache
func (ts *TagSet) Flush() error {
	for _, name := range ts.names {
		if err := ts.flushTag(name); err != nil {
			return err
		}
	}

	return nil
}

// TagIds returns all the ids associated to the TagSet tags
func (ts *TagSet) TagIds() ([]string, error) {
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

func (ts *TagSet) namespace() (string, error) {
	tagsIds, err := ts.TagIds()
	if err != nil {
		return "", err
	}

	return strings.Join(tagsIds, "|"), err
}

func (ts *TagSet) tagId(name string) (string, error) {
	value, err := ts.store.GetString(ts.tagKey(name))
	if err != nil && !isErrNotFound(err) {
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

	return id, ts.store.Forever(ts.tagKey(name), id)
}

func (ts *TagSet) flushTag(name string) error {
	_, err := ts.store.Forget(ts.tagKey(name))

	return err
}

func (ts *TagSet) Names() []string {
	return ts.names
}
