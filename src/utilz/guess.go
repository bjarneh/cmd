// © Knug Industries 2011 all rights reserved
// GNU GENERAL PUBLIC LICENSE VERSION 3.0
// Author bjarneh@ifi.uio.no

package guess

import (
    "os"
    "io/ioutil"
    "utilz/handy"
)

const (
    UTF8 = iota
    LATIN1
    ASCII
    UNSURE
)

func NorwayFile(fname string) int {

    var (
        ascii  byte = 1
        utf8   byte
        latin1 byte
        prev   byte
        b      []byte
        e      os.Error
    )

    b, e = ioutil.ReadFile(fname)
    handy.Check(e)

    for i := 0; i < len(b); i++ {
        if b[i] > 127 {
            ascii = 0
        }
        if looksUnicode(prev, b[i]) {
            utf8 = 1
        }
        if looksLatin(b[i]) {
            latin1 = 1
        }
        prev = b[i]

        if !conflict(ascii, utf8, latin1) && ascii == 0 {
            i = len(b) // breaks
        }
    }

    if !conflict(ascii, utf8, latin1) {
        switch {
        case latin1 == 1:
            return LATIN1
        case utf8 == 1:
            return UTF8
        case ascii == 1:
            return ASCII
        }
    }

    return UNSURE
}

func conflict(a, b, c byte) bool {
    return (a + b + c) != 1
}

func looksLatin(b byte) bool {
    switch b {
    case 230, 248, 229, 198, 216, 197:
        return true
    }
    return false
}

func looksUnicode(first, second byte) bool {
    if first == 195 {
        switch second {
        case 166, 184, 165, 134, 152, 133:
            return true
        }
    }
    return false
}
