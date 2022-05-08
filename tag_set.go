package gocache

import (
	"fmt"
	"strings"

	"github.com/rs/xid"
)

// TagSet is the representation of a tag set for the caching stores
type TagSet struct {
	store store
	names []string
}

func (ts *TagSet) Reset() error {
	for i, name := range ts.names {
		id, err := ts.ResetTag(name)
		if err != nil {
			return err
		}

		ts.names[i] = id
	}

	return nil
}

func (ts *TagSet) Namespace() (string, error) {
	tagsIds, err := ts.TagIds()
	if err != nil {
		return "", err
	}

	return strings.Join(tagsIds, "|"), err
}

func (ts *TagSet) TagIds() ([]string, error) {
	tagIds := make([]string, len(ts.names))
	for i, name := range ts.names {
		val, err := ts.TagId(name)
		if err != nil {
			return tagIds, err
		}

		tagIds[i] = val
	}

	return tagIds, nil
}

func (ts *TagSet) TagId(name string) (string, error) {
	value, err := ts.store.GetString(ts.TagKey(name))
	if err != nil && !isErrNotFound(err) {
		return "", err
	}
	if len(value) == 0 {
		return ts.ResetTag(name)
	}

	return fmt.Sprint(value), nil
}

func (ts *TagSet) TagKey(name string) string {
	return "tag:" + name + ":key"
}

func (ts *TagSet) ResetTag(name string) (string, error) {
	id := xid.New().String()

	return id, ts.store.Forever(ts.TagKey(name), id)
}

func (ts *TagSet) Flush() error {
	for _, name := range ts.names {
		if err := ts.FlushTag(name); err != nil {
			return err
		}
	}

	return nil
}

func (ts *TagSet) FlushTag(name string) error {
	_, err := ts.store.Forget(ts.TagKey(name))

	return err
}

func (ts *TagSet) Names() []string {
	return ts.names
}
