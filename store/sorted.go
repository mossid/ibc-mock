package store

type Sorted struct {
	m Mapping

	enc IndexEncoding
}

func NewSorted(base Base, prefix []byte, enc IndexEncoding) Sorted {
	return Sorted{
		m: NewMapping(base, prefix),

		enc: enc,
	}
}

func (s Sorted) key(power uint64, key []byte) []byte {
	return append(EncodeIndex(power, s.enc), key...)
}

func (s Sorted) Get(ctx Context, power uint64, key []byte, ptr interface{}) {
	s.m.Get(ctx, s.key(power, key), ptr)
}

func (s Sorted) GetIfExists(ctx Context, power uint64, key []byte, ptr interface{}) {
	s.m.GetIfExists(ctx, s.key(power, key), ptr)
}

func (s Sorted) Set(ctx Context, power uint64, key []byte, o interface{}) {
	s.m.Set(ctx, s.key(power, key), o)
}

func (s Sorted) Has(ctx Context, power uint64, key []byte) bool {
	return s.m.Has(ctx, s.key(power, key))
}

func (s Sorted) Delete(ctx Context, power uint64, key []byte) {
	s.m.Delete(ctx, s.key(power, key))
}

func (s Sorted) IsEmpty(ctx Context) bool {
	return s.m.IsEmpty(ctx)
}

func (s Sorted) Prefix(prefix []byte) Sorted {
	return Sorted{
		m:   s.m.Prefix(prefix),
		enc: s.enc,
	}
}

func (s Sorted) IterateAscending(ctx Context, ptr interface{}, fn func([]byte) bool) {
	s.m.Iterate(ctx, ptr, fn)
}

func (s Sorted) IterateDescending(ctx Context, ptr interface{}, fn func([]byte) bool) {
	s.m.ReverseIterate(ctx, ptr, fn)
}

func (s Sorted) Clear(ctx Context) {
	s.m.Clear(ctx)
}
