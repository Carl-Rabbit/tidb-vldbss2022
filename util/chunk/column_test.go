// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package chunk

import (
	"github.com/pingcap/check"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/tidb/types"
)

func equalColumn(c1, c2 *Column) bool {
	if c1.length != c2.length ||
		c1.nullCount != c2.nullCount {
		return false
	}
	if len(c1.nullBitmap) != len(c2.nullBitmap) ||
		len(c1.offsets) != len(c2.offsets) ||
		len(c1.data) != len(c2.data) ||
		len(c1.elemBuf) != len(c2.elemBuf) {
		return false
	}
	for i := range c1.nullBitmap {
		if c1.nullBitmap[i] != c2.nullBitmap[i] {
			return false
		}
	}
	for i := range c1.offsets {
		if c1.offsets[i] != c2.offsets[i] {
			return false
		}
	}
	for i := range c1.data {
		if c1.data[i] != c2.data[i] {
			return false
		}
	}
	for i := range c1.elemBuf {
		if c1.elemBuf[i] != c2.elemBuf[i] {
			return false
		}
	}
	return true
}

func (s *testChunkSuite) TestColumnCopy(c *check.C) {
	col := newFixedLenColumn(8, 10)
	for i := 0; i < 10; i++ {
		col.AppendInt64(int64(i))
	}

	c1 := col.copyConstruct()
	c.Check(equalColumn(col, c1), check.IsTrue)
}

func (s *testChunkSuite) TestLargeStringColumnOffset(c *check.C) {
	numRows := 1
	col := newVarLenColumn(numRows, nil)
	col.offsets[0] = 6 << 30
	c.Check(col.offsets[0], check.Equals, int64(6<<30)) // test no overflow.
}

func (s *testChunkSuite) TestAppendInt64s(c *check.C) {
	col1 := NewColumn(*types.NewFieldType(mysql.TypeLonglong), 1024, 1024)
	col2 := NewColumn(*types.NewFieldType(mysql.TypeLonglong), 0, 1024)

	i64s := col1.Int64s()
	for i := 0; i < 1024; i++ {
		i64s[i] = int64(i)
		col2.AppendInt64(int64(i))
	}

	for i := 0; i < 1024; i++ {
		c.Assert(col1.GetInt64(i), check.Equals, int64(i))
		c.Assert(col2.GetInt64(i), check.Equals, int64(i))
	}
}
