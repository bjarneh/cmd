// © Knug Industries 2009 all rights reserved
// GNU GENERAL PUBLIC LICENSE VERSION 3.0
// Author bjarneh@ifi.uio.no

package walker /* texas ranger */

import (
	"os"
	"path/filepath"
)

// This package does something along the lines of: find PATH -type f
// Filters can be added on both directory and filenames in order to filter
// the resulting slice of pathnames.

// reassign to filter pathwalk
var IncludeDir  = func(p string) bool { return true }
var IncludeFile = func(p string) bool { return true }

func PathWalk(root string) (files []string) {

    fn := func(p string, d *os.FileInfo, e error) error{

        if d.IsDirectory() && ! IncludeDir(p) {
            return filepath.SkipDir
        }

        if d.IsRegular() && IncludeFile(p) {
            files = append(files, p)
        }

        return e
    }

    filepath.Walk(root, fn)

    return files
}

func helper(root string, ch chan string) {

    fn := func(p string, d *os.FileInfo, e error) error{

        if d.IsDirectory() && ! IncludeDir(p) {
            return filepath.SkipDir
        }

        if d.IsRegular() && IncludeFile(p) {
            ch <- p
        }

        return e
    }

    filepath.Walk(root, fn)

    close(ch)
}

// Same as PathWalk part from returning path names in a channel,
// note that this function returns immediatlely, most likely this is
// what you want unless you need all path names at once..
func ChanWalk(root string) chan string {
	fchan := make(chan string)
	go helper(root, fchan)
	return fchan
}
