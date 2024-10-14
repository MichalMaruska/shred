// Package shred provides an API for shredding a file beyond possible recovery
//
// Creator: Pablo Martikian (pablomartikian@hotmail.com)
//
package shred_test

import (
	"bytes"
	"io/ioutil"
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
	file, err := ioutil.TempFile(".", "tmp_")
	if err != nil {
		t.Logf("Error creating temp file for testing.")
		t.FailNow()
	}
	// ...of size 6000
	os.Truncate(file.Name(), 6000)

	// shred test
	err = shred.Shred(file.Name())
	if err != nil {
		t.Errorf("Expected no error shredding a temporary file.")
		os.Remove(file.Name())
	}
}

func TestOverwriteRegularFile(t *testing.T) {
	// create temporary file...
	file, err := ioutil.TempFile(".", "tmp_")
	if err != nil {
		t.Logf("Error creating temp file for testing.")
		t.FailNow()
	}
	// ...of size 6000
	os.Truncate(file.Name(), 6000)
	f1, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Logf("Error reading file before randomization.")
		os.Remove(file.Name())
		t.FailNow()
	}

	// overwrite
	err = shred.Overwrite(file.Name())
	if err != nil {
		t.Errorf("Expected no error shredding a temporary file.")
	}

	f2, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Logf("Error reading file before randomization.")
		os.Remove(file.Name())
		t.FailNow()
	}

	// test file bytes before and after randomization
	if bytes.Equal(f1, f2) {
		t.Errorf("Expected different contect after randomization.")
	}
	os.Remove(file.Name())
}
