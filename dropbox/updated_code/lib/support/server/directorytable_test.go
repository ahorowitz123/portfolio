package server

import (
	"os"
	"testing"
)

func TestgetParentPathAndChildofDirectory(t *testing.T) {
	_, _, err := getParentPathAndChildofDirectory("")
	if err == nil {
		t.Errorf("Should have gotten empty directory not allowed error")
	}

	_, _, err = getParentPathAndChildofDirectory("/hello/world")
	if err == nil {
		t.Errorf("Should have gotten not a directory error")
	}

	parent, child, err := getParentPathAndChildofDirectory("/hello/world/")
	if err != nil {
		t.Errorf("Should not return an error for valid path")
	}

	if parent != "/hello/" {
		t.Errorf("Parent: %v should be /hello/", parent)
	}

	if child != "world/" {
		t.Errorf("Child: %v should be world/", child)
	}
}

func TestcheckValidDirectoryPath(t *testing.T) {
	createTables()
	createNewUserPreEmail("ahorowi2@cs.brown.edu", "jabh1354", []byte{}, "", []byte{})
	createRootandShared("jabh1354")

	err := checkValidDirectoryPath("jabh1354", "")
	if err == nil {
		t.Errorf("Should have gotten empty path not allowed error")
	}

	err = checkValidDirectoryPath("jabh1354", "/")
	if err != nil {
		t.Errorf("Should not have gotten an error as root is valid")
	}

	//check valid absolute path
	err = checkValidDirectoryPath("jabh1354", "/shared/")
	if err != nil {
		t.Errorf("Should not return error because absolute path is valid")
	}

	err = checkValidDirectoryPath("jabh1354", "/blah/")
	if err == nil {
		t.Errorf("Should return error because absolute path does not exist")
	}

	//set cur_dir
	setCurrDir("jabh1354", "/")

	//check valid relative path
	err = checkValidDirectoryPath("jabh1354", "shared/")
	if err != nil {
		t.Errorf("Should not return error because realtive path is valid")
	}

	err = checkValidDirectoryPath("jabh1354", "blah/")
	if err == nil {
		t.Errorf("Should return error because relative path does not exist")
	}

	os.Remove("boxdrop.sqlite")
}
