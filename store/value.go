package store

type Value struct {
	base Base
	key  []byte
}

func NewValue(base Base, key []byte) Value {
	return Value{
		base: base,
		key:  key,
	}
}

func (v Value) store(ctx Context) KVStore {
	return v.base.store(ctx)
}

func (v Value) Get(ctx Context, ptr interface{}) {
	v.base.MustUnmarshalBinaryBare(v.store(ctx).Get(v.key), ptr)
}

func (v Value) GetIfExists(ctx Context, ptr interface{}) {
	bz := v.store(ctx).Get(v.key)
	if bz != nil {
		v.base.MustUnmarshalBinaryBare(bz, ptr)
	}
}

func (v Value) Set(ctx Context, o interface{}) {
	v.store(ctx).Set(v.key, v.base.MustMarshalBinaryBare(o))
}

func (v Value) Exists(ctx Context) bool {
	return v.store(ctx).Has(v.key)
}

func (v Value) Delete(ctx Context) {
	v.store(ctx).Delete(v.key)
}
