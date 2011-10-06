// © Knug Industries 2011 all rights reserved
// GNU GENERAL PUBLIC LICENSE VERSION 3.0
// Author bjarneh@ifi.uio.no

package counter

import (
    "os"
    "fmt"
    "io/ioutil"
)

func CountLines(verbose bool, files ...string) (total int) {

    var (
        content []byte
        err     os.Error
        current int
    )

    for i := 0; i < len(files); i++ {

        content, err = ioutil.ReadFile(files[i])

        current = 0

        if err == nil {
            for _, b := range content {
                if b == '\n' {
                    current++
                }
            }
        } else {
            fmt.Fprintf(os.Stderr, "%s\n", err)
        }

        if verbose {
            fmt.Printf("match: %6d %s\n", current, files[i])
        }

        total += current
    }

    fmt.Printf("lines: %6d\n", total)

    return
}
