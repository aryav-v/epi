// Copyright (c) 2015, Peter Mrekaj. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.txt file.

package lists

import (
	"math/big"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestMedianOfSorted(t *testing.T) {
	for _, test := range []struct {
		l    []interface{}
		i    int
		want *big.Rat
		err  error
	}{
		// Test for errors.
		{[]interface{}{"x"}, 0, nil, ErrType},
		{[]interface{}{0, 2, 1, 3}, 0, nil, ErrSort},
		{[]interface{}{0, 1, 2}, -1, nil, ErrNode},

		// Test lists with cycle.
		{[]interface{}{1}, 0, big.NewRat(1, 1), nil},
		{[]interface{}{0, 1}, 0, big.NewRat(1, 2), nil},
		{[]interface{}{0, 1, 2}, 0, big.NewRat(1, 1), nil},
		{[]interface{}{0, 1, 2, 3}, 0, big.NewRat(1+2, 2), nil},
		{[]interface{}{10, 20, 30}, 1, big.NewRat(20, 1), nil},
		{[]interface{}{10, 20, 30, 40}, 2, big.NewRat(20+30, 2), nil},
		{[]interface{}{10, 20, 30, 40, 50}, 4, big.NewRat(30, 1), nil},

		// Test lists without cycle.
		{[]interface{}{}, -1, nil, nil},
		{[]interface{}{0, 1, 2, 3, 4, 5}, -1, big.NewRat(2+3, 2), nil},
		{[]interface{}{10, 20, 30, 40, 50}, -1, big.NewRat(30, 1), nil},
	} {
		l, n := createCycleList(test.l, test.i)
		if test.err == ErrNode { // Setup an unknown node for an ErrNode test.
			n = &Node{}
		}
		if got, err := MedianOfSorted(l, n); !reflect.DeepEqual(got, test.want) || err != test.err {
			t.Errorf("MedianOfSorted(%v, %v) = %v, %v; want %v, %v", test.l, n, got, err, test.want, test.err)
		}
	}
}

func benchMedianOfSorted(b *testing.B, size int) {
	b.StopTimer()
	if size%2 == 0 { // Decrement size to pickup the median from the middle.
		size--
	}
	for i := 0; i < b.N; i++ {
		ints := rand.New(rand.NewSource(int64(i))).Perm(size)
		sort.Ints(ints)
		var data []interface{}
		for _, n := range ints {
			data = append(data, n)
		}
		l, n := createCycleList(data, 0)
		b.StartTimer()
		m, err := MedianOfSorted(l, n)
		b.StopTimer()
		if int(m.Num().Int64()) != data[size/2] || err != nil {
			b.Error("MedianOfSorted did not find the median")
		}
	}
}

func BenchmarkMedianOfSorted1e1(b *testing.B) { benchMedianOfSorted(b, 1e1) }
func BenchmarkMedianOfSorted1e2(b *testing.B) { benchMedianOfSorted(b, 1e2) }
func BenchmarkMedianOfSorted1e3(b *testing.B) { benchMedianOfSorted(b, 1e3) }
