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
	return strings.Join(this.tagIds(), "|")
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

func (this *TagSet) tagId(name string) string {
	value := this.Store.Get(this.tagKey(name))

	if value == "" {
		return this.resetTag(name)
	}

	return value.(string)
}

func (this *TagSet) tagKey(name string) string {
	return "tag:" + name + ":key"
}

func (this *TagSet) tagIds() []string {
	tagIds := make([]string, len(this.Names))

	for i, name := range this.Names {
		tagIds[i] = this.tagId(name)
	}

	return tagIds
}
