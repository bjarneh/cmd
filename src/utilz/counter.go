// © Knug Industries 2011 all rights reserved
// GNU GENERAL PUBLIC LICENSE VERSION 3.0
// Author bjarneh@ifi.uio.no

package counter

import (
    "os"
    "fmt"
    "io/ioutil"
)

func NewLineFiles(verbose bool, ch chan string) (total int) {

    var current int

    for file := range ch {

        current = NewLineFile(file)

        if verbose {
            fmt.Printf("match: %6d %s\n", current, file)
        }

        total += current
    }

    return
}

func NewLineFile(fname string) (nl int) {

    var (
        content []byte
        err     os.Error
    )

    content, err = ioutil.ReadFile(fname)

    if err == nil {
        for _, b := range content {
            if b == '\n' {
                nl++
            }
        }
    } else {
        fmt.Fprintf(os.Stderr, "%s\n", err)
    }

    return
}
