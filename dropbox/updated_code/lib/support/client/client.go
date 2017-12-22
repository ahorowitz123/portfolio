// Author: jliebowf
// Date: Spring 2016

// Package client provides support code for implementing a command-line client.
// Its two primary components are a command-line interaction wrapper that provides
// a usable interface around the logic of the client, and an auto-test framework
// that tests basic correctness properties about the client implementation.
package client

// Client represents an authenticated client. All methods should be carried out
// as whatever user the current client is authenticated as. This package is
// agnostic to how this authentication is implemented (it could even consist
// of the same login credentials being sent with every request).

// NOTE: The header comments for these functions are all in client/client.go
type Client interface {
	Help()
	Register() error
	Login() error
	Logout() error
	GetSessionToken() []byte
	SetSessionToken(session_token []byte)
	Pwd() (string, error)
	Ls() ([]string, []string, error)
	CD(path string) error
	Mkdir(dirname string) error
	Rmdir(dirname string) error
	Upload(src_filepath string, dest_filename string, modify bool) error
	Rm(filename string) error
	Cat(filename string) (string, error)
	Download(src_filename string, dest_filepath string) error
	Share(filename string, sharee_username string, write_perms bool) error
	Unshare(filename string, sharee_username string) error
	DeleteAcct() error
	LsSharees(filename string) ([]string, error)
	Chmod(filename string, sharee_username string, write_perms bool) error
}
