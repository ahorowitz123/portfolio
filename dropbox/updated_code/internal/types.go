// Author: jliebowf
// Date: Spring 2016

package internal

// contains the return info for ls requests
type ListReturn struct {
	Dirs   []string // The subdirectories
	Files  []string // The files
	Err    string   // If no error was encountered, this will be empty
	Logout bool     // If the client should log itself out after
}

// contains the return info for pwd requests
type PwdReturn struct {
	Path   string // The working directory
	Err    string // If no error was encountered, this will be empty
	Logout bool   // If the client should log itself out after
}

// contains the return info for download/cat requests
type DownloadReturn struct {
	Body   []byte // The contents of the file
	Err    string // If no error was encountered, this will be empty
	Logout bool   // If the client should log itself out after
}

// contains the return info for registration requests
type RegisterReturn struct {
	Success bool   // True if registration successful, false otherwise
	Err     string // If no error was encountered, this will be empty
	// does not need logout because this is before log in
}

// contains the return info for login requests
type LoginReturn struct {
	SessionToken []byte // The session token to be used by the client
	Err          string // If no error was encountered, this will be empty
	// does not need logout because this is before log in
}

// contains the return info for requests that only return errors
type ErrorReturn struct {
	Err    string // If no error was encountered, this will be empty
	Logout bool   // If the client should log itself out after
}

// contains the return info for ls_sharees requests
type ShareesReturn struct {
	Sharees []string // The usernames of all sharees
	Err     string   // If no error was encountered, this will be empty
	Logout  bool     // If the client should log itself out after
}
