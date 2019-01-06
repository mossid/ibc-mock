package store

import (
	"fmt"
)

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
	v.base.cdc.MustUnmarshalBinaryBare(v.store(ctx).Get(v.key), ptr)
}

func (v Value) GetIfExists(ctx Context, ptr interface{}) {
	bz := v.store(ctx).Get(v.key)
	if bz != nil {
		v.base.cdc.MustUnmarshalBinaryBare(bz, ptr)
	}
}

func (v Value) GetSafe(ctx Context, ptr interface{}) *GetSafeError {
	bz := v.store(ctx).Get(v.key)
	if bz == nil {
		return ErrEmptyValue()
	}
	err := v.base.cdc.UnmarshalBinaryBare(bz, ptr)
	if err != nil {
		return ErrUnmarshal(err)
	}
	return nil
}

func (v Value) Set(ctx Context, o interface{}) {
	v.store(ctx).Set(v.key, v.base.cdc.MustMarshalBinaryBare(o))
}

func (v Value) Exists(ctx Context) bool {
	return v.store(ctx).Has(v.key)
}

func (v Value) Delete(ctx Context) {
	v.store(ctx).Delete(v.key)
}

func (v Value) Key() (res []byte) {
	res = make([]byte, len(v.key))
	copy(res, v.key)
	return
}

/*
func (v Value) KeyPath() KeyPath {
	return v.base.KeyPath().AppendKey(v.key, KeyEncodingHex)
}
*/
type GetSafeErrorType byte

const (
	ErrTypeEmptyValue GetSafeErrorType = iota
	ErrTypeUnmarshal
)

func (ty GetSafeErrorType) Format(msg string) (res string) {
	switch ty {
	case ErrTypeEmptyValue:
		res = fmt.Sprintf("Empty value found")
	case ErrTypeUnmarshal:
		res = fmt.Sprintf("Error while unmarshal")
	default:
		panic("Unknown error type")
	}

	if msg != "" {
		res = fmt.Sprintf("%s: %s", res, msg)
	}

	return
}

type GetSafeError struct {
	ty    GetSafeErrorType
	inner error
}

var _ error = (*GetSafeError)(nil) // TODO: sdk.Error

func (err *GetSafeError) Error() string {
	if err.inner == nil {
		return err.ty.Format("")
	}
	return err.ty.Format(err.inner.Error())
}

func ErrEmptyValue() *GetSafeError {
	return &GetSafeError{
		ty: ErrTypeEmptyValue,
	}
}

func ErrUnmarshal(err error) *GetSafeError {
	return &GetSafeError{
		ty:    ErrTypeUnmarshal,
		inner: err,
	}
}
