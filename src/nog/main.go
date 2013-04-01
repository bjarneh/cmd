// © Knug Industries 2011 all rights reserved
// GNU GENERAL PUBLIC LICENSE VERSION 3.0
// Author bjarneh@ifi.uio.no

package main

import (
    "fmt"
    "os"
    "parse/gopt"
    "utilz/guess"
    "utilz/handy"
    "utilz/walker"
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
                r := guess.NorwayFile(f)
                report(f, r)
            }
        } else {
            r := guess.NorwayFile(args[i])
            report(args[i], r)
        }
    }
}

func report(fname string, result int) {
    switch result {
    case guess.UTF8:
        fmt.Printf(" UTF-8   : %s\n", fname)
    case guess.LATIN1:
        fmt.Printf(" LATIN-1 : %s\n", fname)
    case guess.ASCII:
        fmt.Printf(" ASCII   : %s\n", fname)
    default:
        fmt.Printf(" UNSURE  : %s\n", fname)
    }
}

var info = `
  nog - Norwegian ISO-8859-1 / UTF-8 guess

  usage: nog [OPTIONS] [FILES]

  options:
    -h --help      :  print this menu and exit
    -r --recursive :  guess encoding recursively
`
