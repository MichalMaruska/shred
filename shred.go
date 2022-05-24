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

	// use buffered io for performance
	wr = bufio.NewWriter(fd)

	randbuff := make([]byte, 4096)

	// write random data in blocks of 4K when possible
	for sizeleft = lstat.Size(); sizeleft > 0; sizeleft -= 4096 {
		if sizeleft < 4096 {
			// last block smaller than 4K sometimes
			randbuff = make([]byte, sizeleft)
		}
		rand.Read(randbuff)
		wr.Write(randbuff)
	}

	// flush & close & return
	wr.Flush()
	fd.Close()
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
