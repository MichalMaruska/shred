// Package shred provides an API for shredding a file beyond possible recovery
//
// Creator: Pablo Martikian (pablomartikian@hotmail.com)
//
package shred_test

import (
	"bytes"
	"os"
	"shred"
	"testing"
)

func TestShredDir(t *testing.T) {
	err := shred.Shred(".")
	if err == nil {
		t.Errorf("Expected error trying to shred a directory, but it passed.")
	}
}

func TestShredNothing(t *testing.T) {
	err := shred.Shred("")
	if err == nil {
		t.Errorf("Expected error trying to shred nothing, but it passed.")
	}
}

func TestShredRegularFile(t *testing.T) {
	// create temporary file...
	file, err := os.CreateTemp(".", "tmp_")
	if err != nil {
		t.Logf("Error creating temp file for testing.")
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// ...of size 6000
	os.Truncate(file.Name(), 6000)

	// shred test
	err = shred.Shred(file.Name())
	if err != nil {
		os.Remove(file.Name())
		t.Fatalf("Expected no error shredding a temporary file.")
	}

	// verify
	lstat, err := os.Lstat(file.Name())
	if (lstat != nil) {
		t.Errorf("File not removed after Shred.")
	}
}

func TestOverwriteRegularFile(t *testing.T) {
	// create temporary file...
	file, err := os.CreateTemp(".", "tmp_")
	if err != nil {
		t.Logf("Error creating temp file for testing.")
		t.FailNow()
	}
	defer os.Remove(file.Name())

	// ...of size 6000
	os.Truncate(file.Name(), 6000)
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
	os.Remove(file.Name())
}
