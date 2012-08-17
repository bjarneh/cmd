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
    "io/ioutil"
    "path/filepath"
    "bytes"
    "parse/gopt"
    "utilz/handy"
)


var(
    help    = false
    regx    = false
    posix   = false
    pattern = ""
    process = regexp.MustCompile("^[0-9]+$")
    // filters out desired processes
    matchOK func(string)bool
)


func main(){

    parseArgv()

    if help {
        printHelpAndExit()
    }

    if ! handy.IsDir("/proc") {
        fmt.Fprintln(os.Stderr, "/proc is not a directory on this machine")
        fmt.Fprintln(os.Stderr, "that makes me sad... goodbye")
        os.Exit(1)
    }

    addFilter()

    files, err := readDir("/proc")

    if err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }

    dirs  := processFilter(files)
    procs := matchFilter(dirs)

    printProcessInfo( procs )

}


func printProcessInfo(files []os.FileInfo) {

    for i := range files {
        fmt.Printf("pid: %5s\n", files[i].Name())
    }
}


func matchFilter(dirs []os.FileInfo) (procs []os.FileInfo) {

    root := "/proc"

    for i := range dirs {
        status  := filepath.Join(root, dirs[i].Name(), "status")
        cmdline := filepath.Join(root, dirs[i].Name(), "cmdline")
        if handy.IsFile(status) && handy.IsFile(cmdline) {
            if matchOK( cmdline ) {
                procs = append(procs, dirs[i])
            }
        }else{
            fmt.Printf("not a valid process: %s\n", dirs[i].Name())
        }
    }

    return procs
}

func processFilter(dirs []os.FileInfo) (procs []os.FileInfo) {

    for i := range dirs {
        if process.MatchString( dirs[i].Name() ) {
            procs = append(procs, dirs[i])
        }
    }

    return procs
}


func addFilter(){

    if pattern == "" {
        matchOK = func(s string)bool{ return true }
        return
    }

    if regx || posix {
        addRegexpFilter()
    }else{
        addSimpleFilter()
    }

}


func addRegexpFilter(){

    var reg *regexp.Regexp

    if posix {
        reg = regexp.MustCompilePOSIX(pattern)
    }else{
        reg = regexp.MustCompile(pattern)
    }

    matchOK = func(s string)bool{
        b := slurp(s)
        return reg.Match(b)
    }
}

func addSimpleFilter(){

    patternBytes := []byte(pattern)

    matchOK = func(s string)bool{
        b := slurp(s)
        return bytes.Index(b, patternBytes) != -1
    }
}

func slurp(strpath string) (b []byte) {
    b, err := ioutil.ReadFile(strpath)
    if err != nil {
        handy.Check(err)
    }
    return b
}


func readDir(dirname string) ([]os.FileInfo, error) {
    f, err := os.Open(dirname)
    if err != nil {
        return nil, err
    }
    list, err := f.Readdir(-1)
    f.Close()
    if err != nil {
        return nil, err
    }
    return list, nil
}


func printHelpAndExit(){

    var msg =`
  peg - pgrep alternative

  usage: peg [OPTIONS] pattern

  options:

    -h --help     print this menu and exit
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

    if len(rest) > 0 {
        pattern = rest[0]
    }

}
