// © Knug Industries 2011 all rights reserved
// GNU GENERAL PUBLIC LICENSE VERSION 3.0
// Author bjarneh@ifi.uio.no

package main

import (
    "os"
    "fmt"
    "log"
    "regexp"
    "parse/gopt"
    "utilz/walker"
    "utilz/handy"
    "utilz/counter"
)

// option parser
var getopt *gopt.GetOpt

// default values for root directory and files that are worth counting :-)
var (
    root    = "." // PWD
    match   = `(\.sh|\.php|\.c|\.h|\.st|\.scala|\.rb|\.html|\.css|\.java|\.py|\.cpp|\.pl|\.maude|\.go|\.rs)$`
)


func init() {

    getopt = gopt.New()

    getopt.BoolOption("-h -help --help help")
    getopt.BoolOption("-l -list --list")
    getopt.BoolOption("-v -verbose --verbose")
    getopt.StringOption("-m -m= -match -match= --match --match=")

    addMatchFunc(match)
}

func main() {

    args := getopt.Parse(os.Args[1:])

    switch {
    case getopt.IsSet("-list"):
        listing()
    case getopt.IsSet("-help"):
        usage()
    case getopt.IsSet("-match"):
        addMatchFunc(getopt.Get("-match"))
    }

    if len(args) > 0 {
        root = args[0]
    }

    if handy.IsDir( root ) {
        files := walker.PathWalk(root)
        counter.CountLines(getopt.IsSet("-verbose"), files...)
    }else if handy.IsFile( root ) {
        counter.CountLines(false, root)
    }else{
        log.Fatalf("[ERROR] '%s' neither file nor directory\n", root)
    }

}

func usage() {
    fmt.Printf("%s", helpMsg)
    os.Exit(0)
}

func listing() {
    fmt.Println("\noptions:\n")
    fmt.Printf("  -h --help    : %v\n", getopt.IsSet("-help"))
    fmt.Printf("  -v --verbose : %v\n", getopt.IsSet("-verbose"))
    fmt.Printf("  -l --list    : %v\n", getopt.IsSet("-list"))
    fmt.Printf("  -m --match   : %v\n\n", match)
    os.Exit(0)
}

func addMatchFunc(reg string) {

    m := regexp.MustCompile(reg)

    walker.IncludeFile = func(s string) bool {
        return m.MatchString(s)
    }
}

var helpMsg = `
  wcl - recursive: wc -l (line count)

  usage: wcl [OPTIONS] DIR

  DIR defaults to '.' ($PWD)

  options:

    -h --help    :  print this message and exit
    -v --verbose :  verbose output, i.e. info for each file
    -l --list    :  list options and values and exit
    -m --match   :  only count lines in files matching regexp


`
