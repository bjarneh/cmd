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


import(
    "os"
    "fmt"
    "regexp"
    "strings"
    "parse/gopt"
    "utilz/walker"
)


var(
    help    = false
    regx    = false
    posix   = false
    lower    = false
    pattern = ""
)


func main(){

    parseArgv()

    if help {
        printHelpAndExit()
    }

    addFilter()

    files := walker.ChanWalk(".")

    for f := range files {
        fmt.Printf("%s\n", f)
    }

}

func addFilter(){

    if pattern == "" {
        return
    }

    if regx || posix {
        addRegexFilter()
    }else{
        addSimpleFilter()
    }
}

func addSimpleFilter(){

    pattern_low := strings.ToLower(pattern)

    if ! lower {
        walker.IncludeFile = func(s string) bool {
            return strings.Index(s, pattern) != -1
        }
    }else{
        walker.IncludeFile = func(s string) bool {
            s_low := strings.ToLower(s)
            return strings.Index(s_low, pattern_low) != -1
        }
    }

}

func addRegexFilter(){

    var reg *regexp.Regexp

    if posix {
        reg = regexp.MustCompile(pattern)
    }else{
        reg = regexp.MustCompilePOSIX(pattern)
    }

    walker.IncludeFile = func(s string) bool {
        return reg.MatchString(s)
    }
}

func printHelpAndExit(){

    var msg =`
  fnd - fast and simple find tool

  usage: fnd [OPTIONS] pattern

  options:

    -h --help     print this menu and exit
    -l --lower    match case insensitive
    -r --regex    treat pattern as regular expression
    -p --posix    treat pattern as POSIX regular expression
    `

    fmt.Println(msg)
    os.Exit(0)
}

func parseArgv(){

    getopt := gopt.New()

    getopt.BoolOption("-h -help --help")
    getopt.BoolOption("-r -regex --regex")
    getopt.BoolOption("-p -posix --posix")
    getopt.BoolOption("-l -lower --lower")

    rest   := getopt.Parse(os.Args[1:])

    if getopt.IsSet("-help") {
        help = true
    }

    if getopt.IsSet("-regex") {
        regx = true
    }

    if getopt.IsSet("-posix") {
        posix = true
    }

    if getopt.IsSet("-lower") {
        lower = true
    }

    if len(rest) > 0 {
        pattern = rest[0]
    }

}
