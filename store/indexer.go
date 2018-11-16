package store

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

type IndexerKeyEncoding byte

const (
	DecIndexerEnc IndexerKeyEncoding = iota
	HexIndexerEnc
	BinIndexerEnc
)

type Indexer struct {
	m Mapping

	enc IndexerKeyEncoding
}

func NewIndexer(base Base, prefix []byte, enc IndexerKeyEncoding) Indexer {
	return Indexer{
		m:   NewMapping(base, prefix),
		enc: enc,
	}
}

func ElemKey(index uint64, enc IndexerKeyEncoding) (res []byte) {
	switch enc {
	case DecIndexerEnc:
		return []byte(fmt.Sprintf("%020d", index))
	case HexIndexerEnc:
		return []byte(fmt.Sprintf("%020x", index))
	case BinIndexerEnc:
		res = make([]byte, 8)
		binary.BigEndian.PutUint64(res, index)
		return
	default:
		panic("invalid IndexerKeyEncoding")
	}
}

func DecodeElemKey(bz []byte, enc IndexerKeyEncoding) (res uint64, err error) {
	switch enc {
	case DecIndexerEnc:
		return strconv.ParseUint(string(bz), 10, 64)
	case HexIndexerEnc:
		return strconv.ParseUint(string(bz), 16, 64)
	case BinIndexerEnc:
		return binary.BigEndian.Uint64(bz), nil
	default:
		panic("invalid IndexerKeyEncoding")
	}
}

func (ix Indexer) Get(ctx Context, index uint64, ptr interface{}) {
	ix.m.Get(ctx, ElemKey(index, ix.enc), ptr)
}

func (ix Indexer) GetIfExists(ctx Context, index uint64, ptr interface{}) {
	ix.m.GetIfExists(ctx, ElemKey(index, ix.enc), ptr)
}

func (ix Indexer) Set(ctx Context, index uint64, o interface{}) {
	ix.m.Set(ctx, ElemKey(index, ix.enc), o)
}

func (ix Indexer) Has(ctx Context, index uint64) bool {
	return ix.m.Has(ctx, ElemKey(index, ix.enc))
}

func (ix Indexer) Delete(ctx Context, index uint64) {
	ix.m.Delete(ctx, ElemKey(index, ix.enc))
}

func (ix Indexer) IsEmpty(ctx Context) bool {
	return ix.m.IsEmpty(ctx)
}

func (ix Indexer) Prefix(prefix []byte) Indexer {
	return Indexer{
		m: ix.m.Prefix(prefix),

		enc: ix.enc,
	}
}

func (ix Indexer) Iterate(ctx Context, ptr interface{}, fn func(uint64) bool) {
	ix.m.Iterate(ctx, ptr, func(bz []byte) bool {
		key, err := DecodeElemKey(bz, ix.enc)
		if err != nil {
			panic(err)
		}
		return fn(key)
	})
}

func (ix Indexer) First(ctx Context, ptr interface{}) (key uint64, ok bool) {
	keybz, ok := ix.m.First(ctx, ptr)
	if !ok {
		return
	}
	if len(keybz) != 0 {
		key, err := DecodeElemKey(keybz, ix.enc)
		if err != nil {
			return key, false
		}
	}
	return
}

func (ix Indexer) Last(ctx Context, ptr interface{}) (key uint64, ok bool) {
	keybz, ok := ix.m.Last(ctx, ptr)
	if !ok {
		return
	}
	if len(keybz) != 0 {
		key, err := DecodeElemKey(keybz, ix.enc)
		if err != nil {
			return key, false
		}
	}
	return
}
