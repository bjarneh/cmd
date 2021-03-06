//  Copyright © 2012 bjarneh
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

package main

import (
    "fmt"
    "os"
    "parse/gopt"
    "regexp"
    "strings"
    "utilz/walker"
)

var (
    help    = false
    regx    = false
    posix   = false
    lower   = false
    include = false
    topdir  = "."
    pattern = ""
)

func main() {

    parseArgv()

    if help {
        printHelpAndExit()
    }

    addFileFilter()
    addDirFilter()

    files := walker.ChanWalk(topdir)

    for f := range files {
        fmt.Printf("%s\n", f)
    }

}

func isVersionControlDir(s string) bool {

    if strings.HasSuffix(s, ".git") ||
        strings.HasSuffix(s, ".hg") ||
        strings.HasSuffix(s, ".svn") ||
        strings.HasSuffix(s, ".bzr") ||
        strings.HasSuffix(s, ".cvs") {
        return true
    }

    return false
}

func addDirFilter() {

    if include {
        return
    }

    walker.IncludeDir = func(s string) bool {
        if isVersionControlDir(s) {
            return false
        }
        return true
    }
}

func addFileFilter() {

    if pattern == "" {
        return
    }

    if regx || posix {
        addRegexFilter()
    } else {
        addSimpleFilter()
    }
}

func addSimpleFilter() {

    pattern_low := strings.ToLower(pattern)

    if !lower {
        walker.IncludeFile = func(s string) bool {
            return strings.Index(s, pattern) != -1
        }
    } else {
        walker.IncludeFile = func(s string) bool {
            s_low := strings.ToLower(s)
            return strings.Index(s_low, pattern_low) != -1
        }
    }

}

func addRegexFilter() {

    var reg *regexp.Regexp

    if posix {
        reg = regexp.MustCompilePOSIX(pattern)
    } else {
        reg = regexp.MustCompile(pattern)
    }

    walker.IncludeFile = func(s string) bool {
        return reg.MatchString(s)
    }
}

func printHelpAndExit() {

    var msg = `
  fnd - fast and simple find tool

  usage: fnd [OPTIONS] pattern

  options:

    -h --help     print this menu and exit
    -l --lower    match case insensitive
    -d --dir      directory to search (default: PWD)
    -r --regex    treat pattern as regular expression
    -p --posix    treat pattern as POSIX regular expression
    -i --include  include version control dirs (.svn/.cvs/.hg/.git/.bzr)
    `

    fmt.Println(msg)
    os.Exit(0)
}

func parseArgv() {

    getopt := gopt.New()

    getopt.BoolOption("-h -help --help")
    getopt.BoolOption("-r -regex --regex")
    getopt.BoolOption("-p -posix --posix")
    getopt.BoolOption("-l -lower --lower")
    getopt.BoolOption("-i -include --include")
    getopt.StringOptionFancy("-d --dir")

    rest := getopt.Parse(os.Args[1:])

    if getopt.IsSet("-help") {
        help = true
    }

    if getopt.IsSet("-regex") {
        regx = true
    }

    if getopt.IsSet("-dir") {
        topdir = getopt.Get("-dir");
    }

    if getopt.IsSet("-posix") {
        posix = true
    }

    if getopt.IsSet("-lower") {
        lower = true
    }

    if getopt.IsSet("-include") {
        include = true
    }

    if len(rest) > 0 {
        pattern = rest[0]
    }

}
