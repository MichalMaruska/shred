// Package shred provides an API for shredding a file beyond possible recovery
//
// Creator: Pablo Martikian (pablomartikian@hotmail.com)
//
package shred_test

import (
	"bytes"
	"io/fs"
	"os"
	"shred"
	"testing"
)

func TestShredDirectory(t *testing.T) {
	err := shred.Shred(".")
	if err == nil {
		t.Errorf("Expected error trying to shred a directory, but it passed.")
	}
}

func TestShredWrongFilename(t *testing.T) {
	err := shred.Shred("")
	if err == nil {
		t.Errorf("Expected error trying to shred nothing, but it passed.")
	}
}


func temporaryFile(t *testing.T) (*os.File, error) {
	file, err := os.CreateTemp(".", "tmp_")
	if err != nil {
		t.Logf("Error creating temp file for testing.\n\t%s", err.Error())
	}
	return file, err
}


func TestUnreadableFile(t *testing.T) {
	file, err := temporaryFile(t)
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())
	err = os.Chmod(file.Name(), fs.ModePerm & 0400)

	err = shred.Shred(file.Name())
	if err == nil {
		t.Fatalf("Expected error on read-only file.")
	}
}


// write something?
func generateContent(file *os.File, t *testing.T) error {
	// os.Truncate(file.Name(), 6000)

	fd, err := os.OpenFile(file.Name(), os.O_WRONLY, 0)
	if err != nil {
		t.Logf("Error opening temp file.")
		return err
	}
	defer fd.Close()

	_, err = fd.WriteString("Hello")

	if err != nil {
		t.Logf("Error opening temp file.")
	}
	return err
}

func TestShredCheckRemoval(t *testing.T) {
	file, err := temporaryFile(t)
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	err = generateContent(file, t)
	if err != nil {
		t.FailNow()
	}

	err = shred.Shred(file.Name())
	if err != nil {
		t.Fatalf("Expected no error shredding a temporary file.")
	}

	lstat, err := os.Lstat(file.Name())
	if (lstat != nil) {
		t.Errorf("File not removed after Shred.")
	}
}


func TestSingleOverwrite(t *testing.T) {
	file, err := temporaryFile(t)
	if err != nil {
		t.FailNow()
	}
	defer os.Remove(file.Name())

	err = generateContent(file, t)
	if err != nil {
		t.FailNow()
	}

	// it is all zeros!
	f1, err := os.ReadFile(file.Name())
	if err != nil {
		t.Logf("Error reading file before randomization.")
		t.FailNow()
	}

	// overwrite
	err = shred.Overwrite(file.Name())
	if err != nil {
		t.Fatalf("Expected no error shredding a temporary file.")
	}

	f2, err := os.ReadFile(file.Name())
	if err != nil {
		t.Logf("Error reading file before randomization.")
		t.FailNow()
	}

	// test file bytes before and after randomization
	if bytes.Equal(f1, f2) {
		// unless the "random" data happens to be all zeros :)
		t.Errorf("Expected different contect after randomization.")
	}
}
