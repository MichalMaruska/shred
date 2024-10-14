// Package shred provides an API for shredding a file beyond possible recovery
//
// Creator: Pablo Martikian (pablomartikian@hotmail.com)
//
package shred

import (
	"bufio"
	"crypto/rand"
	"errors"
	"io/fs"
	"os"
)

func Overwrite(path string) error {
	var BLOCKSIZE int64 = 4096
	// variable declaration
	var err error
	var lstat fs.FileInfo
	var fd *os.File
	var sizeleft int64
	var wr *bufio.Writer

	// 1st Lstat() to verify it is a regular file
	lstat, err = os.Lstat(path)
	if err != nil {
		return err
	}
	if !lstat.Mode().IsRegular() {
		return errors.New("only regular files can be shredded")
	}

	// 2nd let's open the file
	fd, err = os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer fd.Close()

	// use buffered io for performance
	wr = bufio.NewWriter(fd)

	buffer := make([]byte, BLOCKSIZE)

	// write random data in blocks
	// lseek(0)
	for sizeleft = lstat.Size(); sizeleft > 0; sizeleft -= BLOCKSIZE {
		// avoid generating useless random data, so shorten:
		if sizeleft < 4096 {
			buffer = make([]byte, sizeleft)
		}
		rand.Read(buffer)
		wr.Write(buffer)
	}

	wr.Flush()
	fd.Sync()
	return nil
}

func Shred(path string) error {
	var err error
	for i := 0; i < 3; i++ {
		err = Overwrite(path)
		if err != nil {
			return err
		}
	}
	err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
