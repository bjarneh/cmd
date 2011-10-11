package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "utilz/handy"
    "utilz/walker"
    "parse/gopt"
)

const (
    UTF8 = iota
    LATIN1
    ASCII
    UNSURE
)

var getopt *gopt.GetOpt

func init() {
    getopt = gopt.New()
    getopt.BoolOption("-h -help --help")
    getopt.BoolOption("-r -recursive --recursive")
}

func main() {

    args := getopt.Parse(os.Args[1:])

    if getopt.IsSet("-help") {
        fmt.Println(info)
    }

    for i := 0; i < len(args); i++ {
        if getopt.IsSet("-recursive") && handy.IsDir(args[i]) {
            for f := range walker.ChanWalk(args[i]) {
                r := guessFile(f)
                report(f, r)
            }
        } else {
            r := guessFile(args[i])
            report(args[i], r)
        }
    }
}

func report(fname string, result int) {
    switch result {
    case UTF8:
        fmt.Printf(" UTF-8   : %s\n", fname)
    case LATIN1:
        fmt.Printf(" LATIN-1 : %s\n", fname)
    case ASCII:
        fmt.Printf(" ASCII   : %s\n", fname)
    default:
        fmt.Printf(" UNSURE  : %s\n", fname)
    }
}

func guessFile(fname string) int {

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

var info = `
  nog - Norwegian ISO-8859-1 / UTF-8 guess

  usage: nog [FILES]

  options:
    -h --help      :  print this menu and exit
    -r --recursive :  guess encoding of files in directories
`
