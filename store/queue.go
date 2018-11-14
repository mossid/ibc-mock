package store

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Queue struct {
	Indexer
}

func NewQueue(cdc *codec.Codec, store sdk.KVStore, prefix []byte, enc IndexerKeyEncoding) Queue {
	return Queue{
		Indexer: NewIndexer(cdc, store, prefix, enc),
	}
}

func (q Queue) Peek(ptr interface{}) (ok bool) {
	_, ok = q.First(ptr)
	return
}

func (q Queue) Pop(ptr interface{}) (ok bool) {
	key, ok := q.First(ptr)
	if !ok {
		return
	}

	q.Delete(key)
	return
}

func (q Queue) Push(o interface{}) {
	key, ok := q.Last(nil)
	if !ok {
		key = 0
	} else {
		key = key + 1
	}

	q.Set(key, o)
}
