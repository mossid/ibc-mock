package store

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func NewIndexer(cdc *codec.Codec, store sdk.KVStore, prefix []byte, enc IndexerKeyEncoding) Indexer {
	return Indexer{
		m:   NewMapping(cdc, store, prefix),
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

func (ix Indexer) Get(index uint64, ptr interface{}) {
	ix.m.Get(ElemKey(index, ix.enc), ptr)
}

func (ix Indexer) GetIfExists(index uint64, ptr interface{}) {
	ix.m.GetIfExists(ElemKey(index, ix.enc), ptr)
}

func (ix Indexer) Set(index uint64, o interface{}) {
	ix.m.Set(ElemKey(index, ix.enc), o)
}

func (ix Indexer) Has(index uint64) bool {
	return ix.m.Has(ElemKey(index, ix.enc))
}

func (ix Indexer) Delete(index uint64) {
	ix.m.Delete(ElemKey(index, ix.enc))
}

func (ix Indexer) IsEmpty() bool {
	return ix.m.IsEmpty()
}

func (ix Indexer) Iterate(ptr interface{}, fn func(uint64) bool) {
	ix.m.Iterate(ptr, func(bz []byte) bool {
		key, err := DecodeElemKey(bz, ix.enc)
		if err != nil {
			panic(err)
		}
		return fn(key)
	})
}

func (ix Indexer) First(ptr interface{}) (key uint64, ok bool) {
	keybz, ok := ix.m.First(ptr)
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

func (ix Indexer) Last(ptr interface{}) (key uint64, ok bool) {
	keybz, ok := ix.m.Last(ptr)
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
