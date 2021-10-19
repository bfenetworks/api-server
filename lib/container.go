// Copyright (c) 2021 The BFE Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lib

import "sort"

func Int64BoolMap2Slice(m map[int64]bool) (s []int64) {
	for k, ok := range m {
		if ok {
			s = append(s, k)
		}
	}
	return
}

func StringBoolMap2Slice(m map[string]bool) (s []string) {
	for k, ok := range m {
		if ok {
			s = append(s, k)
		}
	}
	return
}

func StringMap2Slice(m map[string]bool) []string {
	var rst []string
	for k := range m {
		rst = append(rst, k)
	}

	return rst
}

func StringSlice2Map(ss []string) map[string]bool {
	m := map[string]bool{}
	for _, s := range ss {
		m[s] = true
	}

	return m
}

func Int64Map2Slice(m map[int64]bool) []int64 {
	var rst []int64
	for k := range m {
		rst = append(rst, k)
	}

	return rst
}

func SortMapInt642String(m map[int64]string) []int64 {
	ss := []int64{}
	for k := range m {
		ss = append(ss, k)
	}

	sort.Slice(ss, func(i int, j int) bool {
		return ss[i] < ss[j]
	})

	return ss
}

func StringSliceHasElement(s []string, e string) bool {
	for _, one := range s {
		if one == e {
			return true
		}
	}

	return false
}

// StringSliceSubtract a - b
func StringSliceSubtract(a, b []string) (r []string) {
	bm := map[string]bool{}
	for _, one := range b {
		bm[one] = true
	}

	for _, one := range a {
		if !bm[one] {
			r = append(r, one)
		}
	}

	return r
}

func StringSliceSub(a, b []string) (diff []string) {
	bm := map[string]struct{}{}
	for _, one := range b {
		bm[one] = zero
	}

	for _, one := range a {
		if _, ok := bm[one]; !ok {
			diff = append(diff, one)
		}
	}

	return diff
}

var zero = struct{}{}
