/*
 * Copyright 2021 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rt

import (
    `sync`
    `unsafe`
)

type Bitmap struct {
    N int
    B []byte
}

func (self *Bitmap) grow() {
    if self.N >= len(self.B) * 8 {
        self.B = append(self.B, 0)
    }
}

func (self *Bitmap) mark(i int, bv int) {
    if bv != 0 {
        self.B[i / 8] |= 1 << (i % 8)
    } else {
        self.B[i / 8] &^= 1 << (i % 8)
    }
}

func (self *Bitmap) Set(i int, bv int) {
    if i >= self.N {
        panic("bitmap: invalid bit position")
    } else {
        self.mark(i, bv)
    }
}

func (self *Bitmap) Append(bv int) {
    self.grow()
    self.mark(self.N, bv)
    self.N++
}

func (self *Bitmap) AppendMany(n int, bv int) {
    for i := 0; i < n; i++ {
        self.Append(bv)
    }
}

var (
    _stackMapLock  = sync.Mutex{}
    _stackMapCache = make(map[*StackMap]struct{})
)

type StackMap struct {
    N int32
    L int32
    B [1]byte
}

func (self *StackMap) Pin() uintptr {
    self.freeze()
    return uintptr(unsafe.Pointer(self))
}

func (self *StackMap) freeze() {
    _stackMapLock.Lock()
    _stackMapCache[self] = struct{}{}
    _stackMapLock.Unlock()
}

var (
    byteType = UnpackEface(byte(0)).Type
)

const (
    _StackMapSize = unsafe.Sizeof(StackMap{})
)

//go:linkname mallocgc runtime.mallocgc
//goland:noinspection GoUnusedParameter
func mallocgc(nb uintptr, vt *GoType, zero bool) unsafe.Pointer

type StackMapBuilder struct {
    b Bitmap
}

func (self *StackMapBuilder) Build() (p *StackMap) {
    nb := len(self.b.B)
    bm := mallocgc(_StackMapSize + uintptr(nb) - 1, byteType, false)

    /* initialize as 1 bitmap of N bits */
    p = (*StackMap)(bm)
    p.N, p.L = 1, int32(self.b.N)
    copy(BytesFrom(unsafe.Pointer(&p.B), nb, nb), self.b.B)
    return
}

func (self *StackMapBuilder) AddField(ptr bool) {
    if ptr {
        self.b.Append(1)
    } else {
        self.b.Append(0)
    }
}

func (self *StackMapBuilder) AddFields(n int, ptr bool) {
    if ptr {
        self.b.AppendMany(n, 1)
    } else {
        self.b.AppendMany(n, 0)
    }
}
