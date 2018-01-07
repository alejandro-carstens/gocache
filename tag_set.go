package cache

import (
	"github.com/segmentio/ksuid"
	"strings"
)

type TagSet struct {
	Store StoreInterface
	Names []string
}

func (this *TagSet) GetNamespace() (string, error) {
	tagsIds, err := this.tagIds()

	if err != nil {
		return "", err
	}

	return strings.Join(tagsIds, "|"), err
}

func (this *TagSet) resetTag(name string) (string, error) {
	id := ksuid.New().String()

	err := this.Store.Forever(this.tagKey(name), id)

	return id, err
}

func (this *TagSet) Reset() error {
	for i, name := range this.Names {
		id, err := this.resetTag(name)

		if err != nil {
			return err
		}

		this.Names[i] = id
	}

	return nil
}

func (this *TagSet) tagId(name string) (string, error) {
	value, err := this.Store.Get(this.tagKey(name))

	if err != nil {
		return value.(string), err
	}

	if value == "" {
		return this.resetTag(name)
	}

	return value.(string), nil
}

func (this *TagSet) tagKey(name string) string {
	return "tag:" + name + ":key"
}

func (this *TagSet) tagIds() ([]string, error) {
	tagIds := make([]string, len(this.Names))

	for i, name := range this.Names {
		val, err := this.tagId(name)

		if err != nil {
			return tagIds, err
		}

		tagIds[i] = val
	}

	return tagIds, nil
}
