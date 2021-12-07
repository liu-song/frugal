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

package encoder

import (
    `github.com/cloudwego/frugal/internal/atm`
)

type Linker interface {
    Link(p atm.Program) Encoder
}

var (
    linker   Linker
    F_encode atm.CallHandle
)

func init() {
    F_encode = atm.RegisterGCall(encode, emu_gcall_encode)
}

func Link(p atm.Program) Encoder {
    if linker == nil {
        return link_emu(p)
    } else {
        return linker.Link(p)
    }
}

func SetLinker(v Linker) {
    linker = v
}
