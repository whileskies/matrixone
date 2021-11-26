// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package avg

import (
	"fmt"
	"matrixone/pkg/container/nulls"
	"matrixone/pkg/container/ring"
	"matrixone/pkg/container/types"
	"matrixone/pkg/container/vector"
	"matrixone/pkg/encoding"
	"matrixone/pkg/vm/mheap"
)

func NewAvg(typ types.Type) *AvgRing {
	return &AvgRing{Typ: typ}
}

func (r *AvgRing) String() string {
	return fmt.Sprintf("%v-%v", r.Vs, r.Ns)
}

func (r *AvgRing) Free(m *mheap.Mheap) {
	if r.Da != nil {
		mheap.Free(m, r.Da)
		r.Da = nil
		r.Vs = nil
		r.Ns = nil
	}
}

func (r *AvgRing) Count() int {
	return len(r.Vs)
}

func (r *AvgRing) Size() int {
	return cap(r.Da)
}

func (r *AvgRing) Dup() ring.Ring {
	return &AvgRing{
		Typ: r.Typ,
	}
}

func (r *AvgRing) Type() types.Type {
	return r.Typ
}

func (r *AvgRing) SetLength(n int) {
	r.Vs = r.Vs[:n]
	r.Ns = r.Ns[:n]
}

func (r *AvgRing) Shrink(sels []int64) {
	for i, sel := range sels {
		r.Vs[i] = r.Vs[sel]
		r.Ns[i] = r.Ns[sel]
	}
	r.Vs = r.Vs[:len(sels)]
	r.Ns = r.Ns[:len(sels)]
}

func (r *AvgRing) Shuffle(_ []int64, _ *mheap.Mheap) error {
	return nil
}

func (r *AvgRing) Grow(m *mheap.Mheap) error {
	n := len(r.Vs)
	if n == 0 {
		data, err := mheap.Alloc(m, 64)
		if err != nil {
			return err
		}
		r.Da = data
		r.Ns = make([]int64, 0, 8)
		r.Vs = encoding.DecodeFloat64Slice(data)
	} else if n+1 >= cap(r.Vs) {
		r.Da = r.Da[:n*8]
		data, err := mheap.Grow(m, r.Da, int64(n+1)*8)
		if err != nil {
			return err
		}
		mheap.Free(m, r.Da)
		r.Da = data
		r.Vs = encoding.DecodeFloat64Slice(data)
	}
	r.Vs = r.Vs[:n+1]
	r.Vs[n] = 0
	r.Ns = append(r.Ns, 0)
	return nil
}

func (r *AvgRing) Fill(i int64, sel, z int64, vec *vector.Vector) {
	switch vec.Typ.Oid {
	case types.T_int8:
		r.Vs[i] += float64(vec.Col.([]int8)[sel]) * float64(z)
	case types.T_int16:
		r.Vs[i] += float64(vec.Col.([]int16)[sel]) * float64(z)
	case types.T_int32:
		r.Vs[i] += float64(vec.Col.([]int32)[sel]) * float64(z)
	case types.T_int64:
		r.Vs[i] += float64(vec.Col.([]int64)[sel]) * float64(z)
	case types.T_uint8:
		r.Vs[i] += float64(vec.Col.([]uint8)[sel]) * float64(z)
	case types.T_uint16:
		r.Vs[i] += float64(vec.Col.([]uint16)[sel]) * float64(z)
	case types.T_uint32:
		r.Vs[i] += float64(vec.Col.([]uint32)[sel]) * float64(z)
	case types.T_uint64:
		r.Vs[i] += float64(vec.Col.([]uint64)[sel]) * float64(z)
	case types.T_float32:
		r.Vs[i] += float64(vec.Col.([]float32)[sel]) * float64(z)
	case types.T_float64:
		r.Vs[i] += float64(vec.Col.([]float64)[sel]) * float64(z)
	}
	if nulls.Contains(vec.Nsp, uint64(sel)) {
		r.Ns[i]++
	}
}

func (r *AvgRing) BulkFill(i int64, zs []int64, vec *vector.Vector) {
	switch vec.Typ.Oid {
	case types.T_int8:
		vs := vec.Col.([]uint8)
		for j, v := range vs {
			r.Vs[i] += float64(v) * float64(zs[j])
		}
	case types.T_int16:
		vs := vec.Col.([]uint16)
		for j, v := range vs {
			r.Vs[i] += float64(v) * float64(zs[j])
		}
	case types.T_int32:
		vs := vec.Col.([]uint32)
		for j, v := range vs {
			r.Vs[i] += float64(v) * float64(zs[j])
		}
	case types.T_int64:
		vs := vec.Col.([]uint64)
		for j, v := range vs {
			r.Vs[i] += float64(v) * float64(zs[j])
		}
	case types.T_uint8:
		vs := vec.Col.([]uint8)
		for j, v := range vs {
			r.Vs[i] += float64(v) * float64(zs[j])
		}
	case types.T_uint16:
		vs := vec.Col.([]uint16)
		for j, v := range vs {
			r.Vs[i] += float64(v) * float64(zs[j])
		}
	case types.T_uint32:
		vs := vec.Col.([]uint32)
		for j, v := range vs {
			r.Vs[i] += float64(v) * float64(zs[j])
		}
	case types.T_uint64:
		vs := vec.Col.([]uint64)
		for j, v := range vs {
			r.Vs[i] += float64(v) * float64(zs[j])
		}
	}
	r.Ns[i] += int64(nulls.Length(vec.Nsp))
}

func (r *AvgRing) Add(a interface{}, x, y int64) {
	ar := a.(*AvgRing)
	r.Vs[x] += ar.Vs[y]
	r.Ns[x] += ar.Ns[y]
}

func (r *AvgRing) Mul(x, z int64) {
	r.Ns[x] *= z
	r.Vs[x] *= float64(z)
}

func (r *AvgRing) Eval(zs []int64) *vector.Vector {
	defer func() {
		r.Da = nil
		r.Vs = nil
		r.Ns = nil
	}()
	nsp := new(nulls.Nulls)
	for i, z := range zs {
		if n := z - r.Ns[i]; n == 0 {
			nulls.Add(nsp, uint64(i))
		} else {
			r.Vs[i] /= float64(n)
		}
	}
	return &vector.Vector{
		Nsp:  nsp,
		Data: r.Da,
		Col:  r.Vs,
		Or:   false,
		Typ:  types.Type{Oid: types.T_float64, Size: 8},
	}
}