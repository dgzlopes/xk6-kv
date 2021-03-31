package kv

import (
	"context"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/js/modules"
)

func init() {
	modules.Register("k6/x/kv", new(KV))
}

// KV is the k6 key-value extension.
type KV struct{}

type Client struct {
	db *badger.DB
}

var check = false
var client *Client

// XClient represents the Client constructor (i.e. `new kv.Client()`) and
// returns a new Key Value client object.
func (r *KV) XClient(ctxPtr *context.Context, name string, memory bool) interface{} {
	rt := common.GetRuntime(*ctxPtr)
	if check != true {
		if name == "" {
			name = "/tmp/badger"
		}
		var db *badger.DB
		if memory {
			db, _ = badger.Open(badger.DefaultOptions("").WithLoggingLevel(badger.ERROR).WithInMemory(true))
		} else {
			db, _ = badger.Open(badger.DefaultOptions(name).WithLoggingLevel(badger.ERROR))
		}
		client = &Client{db: db}
		check = true
		return common.Bind(rt, client, ctxPtr)
	} else {
		return common.Bind(rt, client, ctxPtr)
	}

}

// Set the given key with the given value.
func (c *Client) Set(key string, value string) error {
	err := c.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(value))
		return err
	})
	return err
}

// Get returns the value for the given key.
func (c *Client) Get(key string) string {
	var valCopy []byte
	_ = c.db.View(func(txn *badger.Txn) error {
		item, _ := txn.Get([]byte(key))
		_ = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		return nil
	})
	return string(valCopy)
}

// ViewPrefix return all the key value pairs where the key starts with some prefix.
func (c *Client) ViewPrefix(prefix string) map[string]string {
	m := make(map[string]string)
	c.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(prefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				m[string(k)] = string(v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return m
}
