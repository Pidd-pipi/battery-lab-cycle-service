package main

import (
	"context"
	"sync"
)

// OpsCache 是记录查询缓存：写入深拷贝、读取返回独立副本，调用方改写不会污染缓存。
type OpsCache struct {
	mu    sync.RWMutex
	items map[string]OpsRecord
}

func newOpsCache() *OpsCache { return &OpsCache{items: map[string]OpsRecord{}} }

func (c *OpsCache) Get(ctx context.Context, id string) (OpsRecord, bool) {
	select {
	case <-ctx.Done():
		return OpsRecord{}, false
	default:
	}
	item, ok := c.items[id]
	if !ok {
		return OpsRecord{}, false
	}
	return item, true
}

func (c *OpsCache) Store(id string, item OpsRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[id] = item
}

func (c *OpsCache) Drop(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, id)
}

func (c *OpsCache) Size() int {
	return len(c.items)
}
