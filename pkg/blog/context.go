package blog

import (
	"context"
)

type ctxMarker struct{}

var (
	ctxMarkerKey = &ctxMarker{}
	NoopTags     = &noopTags{}
)

type Tags interface {
	// Set sets the given key in the metadata tags.
	Set(key string, value interface{}) Tags
	// Has checks if the given key exists.
	Has(key string) bool
	// Values returns a map of key to values.
	// Do not modify the underlying map, please use Set instead.
	Values() map[string]interface{}
}

type mapTags struct {
	values map[string]interface{}
}

func (t *mapTags) Set(key string, value interface{}) Tags {
	t.values[key] = value
	return t
}

func (t *mapTags) Has(key string) bool {
	_, ok := t.values[key]
	return ok
}

func (t *mapTags) Values() map[string]interface{} {
	return t.values
}

type noopTags struct{}

func (t *noopTags) Set(key string, value interface{}) Tags {
	return t
}

func (t *noopTags) Has(key string) bool {
	return false
}

func (t *noopTags) Values() map[string]interface{} {
	return nil
}

func Extract(ctx context.Context) Tags {
	t, ok := ctx.Value(ctxMarkerKey).(Tags)
	if !ok {
		return NoopTags
	}

	return t
}

func SetInContext(ctx context.Context, tags Tags) context.Context {
	return context.WithValue(ctx, ctxMarkerKey, tags)
}

func NewTags() Tags {
	return &mapTags{values: make(map[string]interface{})}
}
