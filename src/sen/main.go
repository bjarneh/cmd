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

package main

import (
    "fmt"
    "io"
    "os"
    "parse/gopt"
    "utilz/hashbytes"
    "strings"
    "utilz/handy"
    "utilz/walker"
)

var (
    help       = false
    safe       = false
    recursive  = false
    quiet      = false
    password   = ""
    bufferSize = 4096
)

func main() {

    files := parseArgv()

    if help {
        fmt.Printf(helpMessage, bufferSize)
        os.Exit(0)
    }

    if len(files) == 0 {
        die("No input files given\n")
    }

    if password == "" {
        fmt.Printf("password: ")
        _, err := fmt.Scanln(&password)
        handy.Check(err)
    }

    for i := range files {
        convert(files[i])
    }

}

func convert(fname string) {

    if handy.IsFile(fname) {

        convertFile(fname)

    } else if handy.IsDir(fname) {

        convertDir(fname)

    } else {

        die("[ERROR] File not found: %s\n", fname)

    }

}

func convertDir(dirname string) {

    if !recursive {

        errStr := "[WARNING] Ignoring: %s (see -recursive)\n"
        fmt.Fprintf(os.Stderr, errStr, dirname)

    } else {

        files := walker.ChanWalk(dirname)

        for file := range files {
            convertFile(file)
        }

    }
}

func convertFile(fname string) {

    var (
        outFname string
        action   string
    )

    infile, err := os.Open(fname)
    handy.Check(err)
    defer infile.Close()

    if strings.HasSuffix(fname, ".sen") {
        outFname = fname[:len(fname)-4]
        action = "decrypt"
    } else {
        outFname = fname + ".sen"
        action = "encrypt"
    }

    outfile, err := os.Create(outFname)
    handy.Check(err)
    defer outfile.Close()

    cryptr := hashbytes.New(password)
    buffer := make([]byte, bufferSize)

    if !quiet {
        fmt.Printf("%s: %s\n", action, fname)
    }

    for {

        n, err := infile.Read(buffer)

        if n == 0 && err == io.EOF {
            break
        } else {
            handy.Check(err)
        }

        // we read some bytes encrypt them
        cryptr.Xor(buffer[:n])

        _, err = outfile.Write(buffer[:n])
        if err != nil {
            die("Failed to write buffer to file\n")
        }

    }

    if !safe {
        defer handy.Delete(fname, true) // dies on error
    }

}

func parseArgv() []string {

    var err error

    getopt := gopt.New()

    getopt.BoolOption("-h -help --help")
    getopt.BoolOption("-s -safe --safe")
    getopt.BoolOption("-q -quiet --quiet")
    getopt.BoolOption("-r -recursive --recursive")
    getopt.StringOptionFancy("-p --password")
    getopt.StringOptionFancy("-b --buffer")

    rest := getopt.Parse(os.Args[1:])

    if getopt.IsSet("-help") {
        help = true
    }

    if getopt.IsSet("-quiet") {
        quiet = true
    }

    if getopt.IsSet("-password") {
        password = getopt.Get("-password")
    }

    if getopt.IsSet("-recursive") {
        recursive = true
    }

    if getopt.IsSet("-safe") {
        safe = true
    }

    if getopt.IsSet("-buffer") {
        bufferSize, err = getopt.GetInt("-buffer")
        if err != nil {
            die("Invalid int: '%s'\n", getopt.Get("-b"))
        }
    }

    return rest

}

func die(fm string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, fm, args...)
    os.Exit(1)
}

var helpMessage = `
  sen - (en|de)crypt files with SHA1-hash-xor-scheme

  usage: sen [OPTIONS] FILE(S)

  options:

    -h --help       print this menu and exit
    -s --safe       do not delete or overwrite files
    -q --quiet      do not report what's going on
    -r --recursive  recursively traverse directories
    -p --password   use this password to encrypt/decrypt
    -b --buffer     file reading buffer-size (default: %d)

`
