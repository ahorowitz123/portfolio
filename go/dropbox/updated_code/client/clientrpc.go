package main

import (
	"fmt"

	"github.com/s17-ahorowi2-jtc2/updated_code/internal"
)

/**
 * Make the request to the server for validating the first part of the registration process.
 *
 * Input:
 *   email string - the email given in registration
 *   username string - the username given in registration
 *   password string - the password given in registration
 * Output:
 *   bool - true if first part of registration is valid, else false
 *   error - the error, if an error occurred, else false
 *
 * Precondition: Client is not logged in. All arguments validated client side.
 * Postcondition: Server verified the arguments, generated and sent an email verification code to the given email,
 *   providing no errors.
 */
func (c *Client) EmailRegister(email string, username string, password string) (bool, error) {
	var ret internal.RegisterReturn
	err := c.server.Call("emailregister", &ret, email, username, password)
	if err != nil {
		return false, err
	}
	if ret.Err != "" {
		return false, fmt.Errorf(ret.Err)
	}
	return ret.Success, nil
}

/**
 * Make the request to the server for validating the second part of the registration process, which is checking
 * the verification code.
 *
 * Input:
 *   email string - the email given in registration
 *   username string - the username given in registration
 *   password string - the password given in registration
 *   ver_code string - the verification code given in registration
 * Output:
 *   bool - true if ver_code correct and everything validates, else false
 *   error - the error, if an error occurred, else false
 *
 * Precondition: Client is not logged in. First part of registration (EmailRegister)
 * Postcondition: Server validated verification code, and user is a valid user that can now log in, if no error.
 */
func (c *Client) ValidateRegister(email string, username string, password string, ver_code string) (bool, error) {
	var ret internal.RegisterReturn
	err := c.server.Call("validateregister", &ret, email, username, password, ver_code)
	if err != nil {
		return false, err
	}
	if ret.Err != "" {
		return false, fmt.Errorf(ret.Err)
	}
	return ret.Success, nil
}

/**
 * Make the request to the server for validating the login.
 *
 * Input:
 *   email string - the email given in login
 *   password string - the password given in login
 * Output:
 *   []byte - the session token for the user, if error not nil
 *   error - the error, if an error occurred, else false
 *
 * Precondition: Client is not logged in. All arguments validated client side.
 * Postcondition: Client is logged in. The server generated a session token for the client.
 */
func (c *Client) ValidateLogin(username string, password string) ([]byte, error) {
	var ret internal.LoginReturn
	err := c.server.Call("validatelogin", &ret, username, password)
	if err != nil {
		return []byte{}, err
	}
	if ret.Err != "" {
		return []byte{}, fmt.Errorf(ret.Err)
	}
	return ret.SessionToken, nil
}

/**
 * Make the request to the server for handling logout.
 *
 * Input:
 *   username string - the username of the client
 * Output: The error, if an error occurred, else nil.
 *
 * Precondition: Client is logged in on both client and server side.
 * Postcondition: Client is not logged in on server. Session token and username are reset.
 */
func (c *Client) HandleLogout(username string) error {
	var ret internal.ErrorReturn
	err := c.server.Call("logout", &ret, username)
	if err != nil {
		return err
	}
	if ret.Logout {
		c.SetSessionToken([]byte{})
	}
	if ret.Err != "" {
		return fmt.Errorf(ret.Err)
	}
	return nil
}

/**
 * Make a request to get the current working directory of the user.
 *
 * Input:
 *   username string - the username stored client side for the user
 *   session_token []byte - the session token stored client side for the user
 * Output:
 *   string - the current working directory, if error not nil
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in.
 * Postcondition: If uesrname and/or session token is invalid, client is logged out. Else, nothing.
 */
func (c *Client) GetWorkingDirectory(username string, session_token []byte) (string, error) {
	var ret internal.PwdReturn
	err := c.server.Call("getworkingdirectory", &ret, username, session_token)
	if err != nil {
		return "", err
	}
	if ret.Logout {
		c.SetSessionToken([]byte{})
	}
	if ret.Err != "" {
		return "", fmt.Errorf(ret.Err)
	}
	return ret.Path, nil
}

/**
 * Make a request to get the files and subdirectories of the user's current working directory.
 *
 * Input:
 *   username string - the username stored client side for the user
 *   session_token []byte - the session token stored client side for the user
 * Output:
 *   []string - the names of the subdirectories. empty if error not nil
 *   []string - the names of the files. empty if error not nil
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in.
 * Postcondition: If uesrname and/or session token is invalid, client is logged out. Else, nothing.
 */
func (c *Client) GetSubdirectoriesAndFiles(username string, session_token []byte) ([]string, []string, error) {
	var ret internal.ListReturn
	err := c.server.Call("ls", &ret, username, session_token)
	if err != nil {
		return []string{}, []string{}, err
	}
	if ret.Logout {
		c.SetSessionToken([]byte{})
	}
	if ret.Err != "" {
		return []string{}, []string{}, fmt.Errorf(ret.Err)
	}
	return ret.Dirs, ret.Files, nil
}

/**
 * Make a request to change the current working directory of the user.
 *
 * Input:
 *   username string - the username stored client side for the user
 *   session_token []byte - the session token stored client side for the user
 *   path string - the new working directory for the user, relative or absolute
 * Output:
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in.
 * Postcondition: If uesrname and/or session token is invalid, client is logged out. Else, if no error, the
 *   working directory is changed to the given path.
 */
func (c *Client) ChangeDirectory(username string, session_token []byte, path string) error {
	var ret internal.ErrorReturn
	err := c.server.Call("cd", &ret, username, session_token, path)
	if err != nil {
		return err
	}
	if ret.Logout {
		c.SetSessionToken([]byte{})
	}
	if ret.Err != "" {
		return fmt.Errorf(ret.Err)
	}
	return nil
}

/**
 * Make a request to make a directory named dirname in the current working directory of the user.
 *
 * Input:
 *   username string - the username stored client side for the user
 *   session_token []byte - the session token stored client side for the user
 *   dirname string - the name of the new directory
 * Output:
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in.
 * Postcondition: If uesrname and/or session token is invalid, client is logged out. Else, if no error, new
 *   subdirectory called dirname.
 */
func (c *Client) MakeDirectory(username string, session_token []byte, dirname string) error {
	var ret internal.ErrorReturn
	err := c.server.Call("mkdir", &ret, username, session_token, dirname)
	if err != nil {
		return err
	}
	if ret.Logout {
		c.SetSessionToken([]byte{})
	}
	if ret.Err != "" {
		return fmt.Errorf(ret.Err)
	}
	return nil
}

/**
 * Make a request to remove the directory named dirname in the current working directory of the user.
 *
 * Input:
 *   username string - the username stored client side for the user
 *   session_token []byte - the session token stored client side for the user
 *   dirname string - the name of the directory being removed
 * Output:
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in.
 * Postcondition: If uesrname and/or session token is invalid, client is logged out. Else, if no error,
 *   subdirectory called dirname is removed.
 */
func (c *Client) RemoveDirectory(username string, session_token []byte, dirname string) error {
	var ret internal.ErrorReturn
	err := c.server.Call("rmdir", &ret, username, session_token, dirname)
	if err != nil {
		return err
	}
	if ret.Logout {
		c.SetSessionToken([]byte{})
	}
	if ret.Err != "" {
		return fmt.Errorf(ret.Err)
	}
	return nil
}

/**
 * Make a request to remove the file named filename in the current working directory of the user.
 *
 * Input:
 *   username string - the username stored client side for the user
 *   session_token []byte - the session token stored client side for the user
 *   filename string - the name of the file being removed
 * Output:
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in.
 * Postcondition: If uesrname and/or session token is invalid, client is logged out. Else, if no error,
 *   file named filename is removed.
 */
func (c *Client) RemoveFile(username string, session_token []byte, filename string) error {
	var ret internal.ErrorReturn
	err := c.server.Call("rmfile", &ret, username, session_token, filename)
	if err != nil {
		return err
	}
	if ret.Logout {
		c.SetSessionToken([]byte{})
	}
	if ret.Err != "" {
		return fmt.Errorf(ret.Err)
	}
	return nil
}

/**
 * Make a request to upload the file_bytes to the filepath/filename dest_filename. Modify indicates if it should
 * (and must) replace the existing file.
 *
 * Input:
 *   username string - the username stored client side for the user
 *   session_token []byte - the session token stored client side for the user
 *   file_bytes []byte - the contents to be uploaded
 *   dest_filename string - the filepath/filename on the server to store the byte
 *   modify bool - true if replacing a file, false if making a new file
 * Output:
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in.
 * Postcondition: If uesrname and/or session token is invalid, client is logged out. Else, if no error,
 *   file is uploaded/modified as appropriate.
 */
func (c *Client) uploadFile(username string, session_token []byte, file_bytes []byte, dest_filename string, modify bool) error {
	var ret internal.ErrorReturn
	err := c.server.Call("upload", &ret, username, session_token, file_bytes, dest_filename, modify)
	if err != nil {
		return err
	}
	if ret.Logout {
		c.SetSessionToken([]byte{})
	}
	if ret.Err != "" {
		return fmt.Errorf(ret.Err)
	}
	return nil
}

/**
 * Make a request to get the file contents of the given file from the server.
 *
 * Input:
 *   username string - the username stored client side for the user
 *   session_token []byte - the session token stored client side for the user
 *   file string - the fliepath/filename of the file which contents are being retrieved
 * Output:
 *   []byte - the contents of the file, empty if an error
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in.
 * Postcondition: If uesrname and/or session token is invalid, client is logged out. Else, if no error,
 *   nothing changes.
 */
func (c *Client) Concatenate(username string, session_token []byte, file string) ([]byte, error) {
	var ret internal.DownloadReturn
	err := c.server.Call("cat", &ret, username, session_token, file)
	if err != nil {
		return []byte{}, err
	}
	if ret.Logout {
		c.SetSessionToken([]byte{})
	}
	if ret.Err != "" {
		return []byte{}, fmt.Errorf(ret.Err)
	}
	return ret.Body, nil
}

/**
 * Make a request to share the given filename with the sharee. Write_perms indicate if read/write or read only.
 *
 * Input:
 *   username string - the username stored client side for the user
 *   session_token []byte - the session token stored client side for the user
 *   filename string - the path/name of the file to be shared
 *   sharee_username - the name of the user with whom the file is being shared
 *   write_perms bool - true if shared with read/write, false if readonly
 * Output:
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in.
 * Postcondition: If uesrname and/or session token is invalid, client is logged out. Else, if no error,
 *   file is shared with sharee.
 */
func (c *Client) ShareFile(username string, session_token []byte, filename string, sharee_username string, write_perms bool) error {
	var ret internal.ErrorReturn
	err := c.server.Call("sharefile", &ret, username, session_token, filename, sharee_username, write_perms)
	if err != nil {
		return err
	}
	if ret.Logout {
		c.SetSessionToken([]byte{})
	}
	if ret.Err != "" {
		return fmt.Errorf(ret.Err)
	}
	return nil
}

/**
 * Make a request to unshare the given filename with the sharee.
 *
 * Input:
 *   username string - the username stored client side for the user
 *   session_token []byte - the session token stored client side for the user
 *   filename string - the path/name of the file to be shared
 *   sharee_username - the name of the user with whom the file is being shared
 * Output:
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in.
 * Postcondition: If uesrname and/or session token is invalid, client is logged out. Else, if no error,
 *   file is no longer shared with sharee.
 */
func (c *Client) UnshareFile(username string, session_token []byte, filename string, sharee_username string) error {
	var ret internal.ErrorReturn
	err := c.server.Call("unsharefile", &ret, username, session_token, filename, sharee_username)
	if err != nil {
		return err
	}
	if ret.Logout {
		c.SetSessionToken([]byte{})
	}
	if ret.Err != "" {
		return fmt.Errorf(ret.Err)
	}
	return nil
}

/**
 * Delete the account of the current user.
 *
 * Input:
 *   username string - the username stored client side for the user
 *   session_token []byte - the session token stored client side for the user
 * Output:
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in. The working directory is "/". The user has no owned file or
 *   directory.
 * Postcondition: If uesrname and/or session token is invalid, client is logged out. Else, the user no
 *   longer exists anywhere on the server, and the session_token is reset on the client side.
 */
func (c *Client) DeleteAccount(username string, session_token []byte) error {
	var ret internal.ErrorReturn
	err := c.server.Call("deleteacct", &ret, username, session_token)
	if err != nil {
		return err
	}
	if ret.Logout {
		// do it like this because on success need to reset session token
		c.SetSessionToken([]byte{})
	}
	if ret.Err != "" {
		return fmt.Errorf(ret.Err)
	}
	return nil
}

/**
 * Make a request to get the sharees of the file with the given filename.
 *
 * Input:
 *   username string - the username stored client side for the user
 *   session_token []byte - the session token stored client side for the user
 *   filename string - the filename/filepath of the file
 * Output:
 *   []string - the names of the sharees. empty if error not nil
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in.
 * Postcondition: If uesrname and/or session token is invalid, client is logged out. Else, nothing.
 */
func (c *Client) ListSharees(username string, session_token []byte, filename string) ([]string, error) {
	var ret internal.ShareesReturn
	err := c.server.Call("listsharees", &ret, username, session_token, filename)
	if err != nil {
		return []string{}, err
	}
	if ret.Logout {
		c.SetSessionToken([]byte{})
	}
	if ret.Err != "" {
		return []string{}, fmt.Errorf(ret.Err)
	}
	return ret.Sharees, nil

}

/**
 * Make a request to change the permissions of the given filename with the sharee. Write_perms
 * indicate if read/write or read only.
 *
 * Input:
 *   username string - the username stored client side for the user
 *   session_token []byte - the session token stored client side for the user
 *   filename string - the path/name of the file to be shared
 *   sharee_username - the name of the user with whom the file is being shared
 *   write_perms bool - true if change perms to read/write, false if readonly
 * Output:
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in.
 * Postcondition: If uesrname and/or session token is invalid, client is logged out. Else, if no error,
 *   file permissions for sharee are modified to read/write if write_perms is true, else read only. If
 *   permissions already the appropriate value, nothing changes.
 */
func (c *Client) ChmodFile(username string, session_token []byte, filename string, sharee_username string, write_perms bool) error {
	var ret internal.ErrorReturn
	err := c.server.Call("chmodfile", &ret, username, session_token, filename, sharee_username, write_perms)
	if err != nil {
		return err
	}
	if ret.Logout {
		c.SetSessionToken([]byte{})
	}
	if ret.Err != "" {
		return fmt.Errorf(ret.Err)
	}
	return nil

}
