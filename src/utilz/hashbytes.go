//  Copyright © 2013 bjarneh
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

package hashbytes

import (
    "crypto/sha1"
    "hash"
    "io"
)

type HashBytes struct {
    passwd   string
    pos      int
    bytes    []byte
    sha1hash hash.Hash
}

func New(password string) *HashBytes {

    s := sha1.New()
    io.WriteString(s, password)
    b := s.Sum(nil)

    return &HashBytes{password, 0, b, s}
}

func (h *HashBytes) Xor(b []byte) {

    for i := 0; i < len(b); i++ {

        b[i] = b[i] ^ h.bytes[h.pos]
        h.pos++

        if h.pos >= len(h.bytes) {
            h.next()
        }
    }

}

func (h *HashBytes) next() {
    io.WriteString(h.sha1hash, h.passwd)
    h.bytes = h.sha1hash.Sum(nil)
    h.pos = 0
}
