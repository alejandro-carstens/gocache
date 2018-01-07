package cache

import (
	"github.com/segmentio/ksuid"
	"strings"
)

type TagSet struct {
	Store StoreInterface
	Names []string
}

func (this *TagSet) GetNamespace() string {
	tagsIds, err := this.tagIds()

	if err != nil {
		panic(err)
	}

	return strings.Join(tagsIds, "|")
}

func (this *TagSet) resetTag(name string) string {
	id := ksuid.New().String()

	this.Store.Forever(this.tagKey(name), id)

	return id
}

func (this *TagSet) Reset() {
	for i, name := range this.Names {
		this.Names[i] = this.resetTag(name)
	}
}

func (this *TagSet) tagId(name string) (string, error) {
	value, err := this.Store.Get(this.tagKey(name))

	if err != nil {
		return value.(string), err
	}

	if value == "" {
		return this.resetTag(name), nil
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
