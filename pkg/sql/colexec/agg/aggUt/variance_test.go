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

package aggut

import (
	"github.com/matrixorigin/matrixone/pkg/sql/colexec/agg"
	"testing"

	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/container/vector"
	"github.com/matrixorigin/matrixone/pkg/testutil"
	"github.com/matrixorigin/matrixone/pkg/vm/mheap"
	"github.com/matrixorigin/matrixone/pkg/vm/mmu/guest"
	"github.com/matrixorigin/matrixone/pkg/vm/mmu/host"
	"github.com/stretchr/testify/require"
)

func TestVariance(t *testing.T) {
	inputType := types.New(types.T_int8, 0, 0, 0)
	variance1 := agg.NewVariance[int8]()
	variance2 := agg.NewVariance[int8]()
	variance3 := agg.NewVariance[int8]()
	m := mheap.New(guest.New(1<<30, host.New(1<<30)))
	vec := testutil.NewVector(Rows, inputType, m, false, nil)
	{
		agg := agg.NewUnaryAgg(variance1, true, inputType, types.New(types.T_float64, 0, 0, 0), variance1.Grows, variance1.Eval, variance1.Merge, variance1.Fill, nil)
		err := agg.Grows(1, m)
		require.NoError(t, err)
		for i := 0; i < Rows; i++ {
			agg.Fill(0, int64(i), 1, []*vector.Vector{vec})
		}
		v, err := agg.Eval(m)
		require.NoError(t, err)
		require.Equal(t, []float64{8.25}, vector.GetColumn[float64](v))
		v.Free(m)
	}
	{
		agg0 := agg.NewUnaryAgg(variance2, true, inputType, types.New(types.T_float64, 0, 0, 0), variance2.Grows, variance2.Eval, variance2.Merge, variance2.Fill, nil)
		err := agg0.Grows(1, m)
		require.NoError(t, err)
		for i := 0; i < Rows; i++ {
			agg0.Fill(0, int64(i), 1, []*vector.Vector{vec})
		}
		agg1 := agg.NewUnaryAgg(variance3, true, inputType, types.New(types.T_float64, 0, 0, 0), variance3.Grows, variance3.Eval, variance3.Merge, variance3.Fill, nil)
		err = agg1.Grows(1, m)
		require.NoError(t, err)
		for i := 0; i < Rows; i++ {
			agg1.Fill(0, int64(i), 1, []*vector.Vector{vec})
		}
		agg0.Merge(agg1, 0, 0)
		{
			v, err := agg0.Eval(m)
			require.NoError(t, err)
			require.Equal(t, []float64{8.25}, vector.GetColumn[float64](v))
			v.Free(m)
		}
		{
			v, err := agg1.Eval(m)
			require.NoError(t, err)
			require.Equal(t, []float64{8.25}, vector.GetColumn[float64](v))
			v.Free(m)
		}
	}
	vec.Free(m)
	require.Equal(t, int64(0), m.Size())
}

func TestVarianceDist(t *testing.T) {
	inputType := types.New(types.T_int8, 0, 0, 0)
	variance1 := agg.NewVariance[int8]()
	variance2 := agg.NewVariance[int8]()
	variance3 := agg.NewVariance[int8]()
	m := mheap.New(guest.New(1<<30, host.New(1<<30)))
	vec := testutil.NewVector(Rows, inputType, m, false, nil)
	{
		agg := agg.NewUnaryDistAgg(true, inputType, types.New(types.T_float64, 0, 0, 0), variance1.Grows, variance1.Eval, variance1.Merge, variance1.Fill)
		err := agg.Grows(1, m)
		require.NoError(t, err)
		for i := 0; i < Rows; i++ {
			agg.Fill(0, int64(i), 1, []*vector.Vector{vec})
		}
		v, err := agg.Eval(m)
		require.NoError(t, err)
		require.Equal(t, []float64{8.25}, vector.GetColumn[float64](v))
		v.Free(m)
	}
	{
		agg0 := agg.NewUnaryDistAgg(true, inputType, types.New(types.T_float64, 0, 0, 0), variance2.Grows, variance2.Eval, variance2.Merge, variance2.Fill)
		err := agg0.Grows(1, m)
		require.NoError(t, err)
		for i := 0; i < Rows; i++ {
			agg0.Fill(0, int64(i), 1, []*vector.Vector{vec})
		}
		agg1 := agg.NewUnaryDistAgg(true, inputType, types.New(types.T_float64, 0, 0, 0), variance3.Grows, variance3.Eval, variance3.Merge, variance3.Fill)
		err = agg1.Grows(1, m)
		require.NoError(t, err)
		for i := 0; i < Rows; i++ {
			agg1.Fill(0, int64(i), 1, []*vector.Vector{vec})
		}
		agg0.Merge(agg1, 0, 0)
		{
			v, err := agg0.Eval(m)
			require.NoError(t, err)
			require.Equal(t, []float64{8.25}, vector.GetColumn[float64](v))
			v.Free(m)
		}
		{
			v, err := agg1.Eval(m)
			require.NoError(t, err)
			require.Equal(t, []float64{8.25}, vector.GetColumn[float64](v))
			v.Free(m)
		}
	}
	vec.Free(m)
	require.Equal(t, int64(0), m.Size())
}

func TestDecimalDist(t *testing.T) {
	inputType := types.New(types.T_decimal64, 0, 0, 0)
	outputType := types.New(types.T_decimal128, 0, 0, 0)
	variance1 := agg.NewVD64()
	variance2 := agg.NewVD64()
	variance3 := agg.NewVD64()
	m := mheap.New(guest.New(1<<30, host.New(1<<30)))
	vec := testutil.NewVector(Rows, inputType, m, false, nil)
	{
		agg := agg.NewUnaryDistAgg(true, inputType, outputType, variance1.Grows, variance1.Eval, variance1.Merge, variance1.Fill)
		err := agg.Grows(1, m)
		require.NoError(t, err)
		for i := 0; i < Rows; i++ {
			agg.Fill(0, int64(i), 1, []*vector.Vector{vec})
		}
		v, err := agg.Eval(m)
		require.NoError(t, err)
		d, _ := types.Decimal128_FromFloat64(8.25, 64, 4)
		require.Equal(t, []types.Decimal128{d}, vector.GetColumn[types.Decimal128](v))
		v.Free(m)
	}
	{
		agg0 := agg.NewUnaryDistAgg(true, inputType, outputType, variance2.Grows, variance2.Eval, variance2.Merge, variance2.Fill)
		err := agg0.Grows(1, m)
		require.NoError(t, err)
		for i := 0; i < Rows; i++ {
			agg0.Fill(0, int64(i), 1, []*vector.Vector{vec})
		}
		agg1 := agg.NewUnaryDistAgg(true, inputType, outputType, variance3.Grows, variance3.Eval, variance3.Merge, variance3.Fill)
		err = agg1.Grows(1, m)
		require.NoError(t, err)
		for i := 0; i < Rows; i++ {
			agg1.Fill(0, int64(i), 1, []*vector.Vector{vec})
		}
		agg0.Merge(agg1, 0, 0)
		{
			v, err := agg0.Eval(m)
			require.NoError(t, err)
			d, _ := types.Decimal128_FromFloat64(8.25, 64, 4)
			require.Equal(t, []types.Decimal128{d}, vector.GetColumn[types.Decimal128](v))
			v.Free(m)
		}
		{
			v, err := agg1.Eval(m)
			require.NoError(t, err)
			d, _ := types.Decimal128_FromFloat64(8.25, 64, 4)
			require.Equal(t, []types.Decimal128{d}, vector.GetColumn[types.Decimal128](v))
			v.Free(m)
		}
	}
	vec.Free(m)
	require.Equal(t, int64(0), m.Size())
}
