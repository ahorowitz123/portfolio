// Author: jliebowf
// Date: Spring 2016

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/s17-ahorowi2-jtc2/updated_code/lib/support/client"
	"github.com/s17-ahorowi2-jtc2/updated_code/lib/support/rpc"

	"github.com/fatih/color"
)

/**
 * The cmdInfo struct is used to store the info for each command listed in Help.
 *
 * cmd string - the name of the command
 * info string - the description of the command
 */
type cmdInfo struct {
	cmd  string
	info string
}

/**
 * The main function of the client. This is ran automatically when the client is started.
 *
 * Input: N/A
 * Output: N/A
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func main() {
	// Usage: ./bin/server <server>
	// error and exit on invalid number of arguments
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %v <server>\n", os.Args[0])
		os.Exit(1)
	}

	// create a new connection to the server
	server := rpc.NewServerRemote(os.Args[1])

	// create a new client struct. It starts with no session token and uses the connection to the
	// server created above
	c := Client{server: server, session_token: []byte{}, username: ""}

	// and run the client REPL
	err := client.RunCLI(&c)
	if err != nil {
		fmt.Printf("fatal error: %v\n", err)
		os.Exit(1)
	}
}

/**
 * The client struct is used to store the info for the running client.
 *
 * server *rpc.ServerRemote - The connection to the server, instantiated in main
 * session_token []byte - The session token. Is empty if not logged in. Is set to the return value of
 *    the server after logging in
 * username string - The username of the client logging in. Is the empty string if not logged in.
 */
type Client struct {
	server        *rpc.ServerRemote
	session_token []byte
	username      string
}

/**
 * Called when the user types help into the terminal with no arguments. This function will print out the
 * commands the user can enter at the given time. This function does not interact with the server.
 *
 * Input: N/A
 * Output: N/A
 *
 * Precondition: If the session token is not set, only shows pre login commands. If the session token is set,
 *   only shows the post login commands. Does not validate the session token, just checks it is set.
 * Postcondition: N/A.
 */
func (c *Client) Help() {
	fmt.Println("Available commands:")
	var cmds []cmdInfo
	if len(c.session_token) == 0 {
		cmds = []cmdInfo{
			cmdInfo{cmd: "help", info: "List all possible commands"},
			cmdInfo{cmd: "register", info: "Create a new account"},
			cmdInfo{cmd: "login", info: "Log into an existing account"},
			cmdInfo{cmd: "exit", info: "Exit dropbox"},
		}
	} else {
		cmds = []cmdInfo{
			cmdInfo{cmd: "help", info: "List all possible commands"},
			cmdInfo{cmd: "pwd", info: "Print current directory path"},
			cmdInfo{cmd: "ls", info: "List all contents of the current directory"},
			cmdInfo{cmd: "cd <directory>", info: "Navigate to directory"},
			cmdInfo{cmd: "mkdir <directory>", info: "Create new directory in current directory"},
			cmdInfo{cmd: "rmdir <directory>", info: "Remove the subdirectory with the given name"},
			cmdInfo{cmd: "upload <src_filepath> <dest_filename>", info: "Upload the given file to the current directory with the given filename"},
			cmdInfo{cmd: "modify <src_filepath> <dest_filename>", info: "Modify the given file and overwrite the given file in the current directory"},
			cmdInfo{cmd: "download <src_filename> <dest_filepath>", info: "Download the given file to given directory name on your local system"},
			cmdInfo{cmd: "cat <file>", info: "View contents of file"},
			cmdInfo{cmd: "rm <file>", info: "Remove the given file from the current directory."},
			cmdInfo{cmd: "share_r <file> <username>", info: "Share the selected file with the given user (read-only)"},
			cmdInfo{cmd: "share_rw <file> <username>", info: "Shared the selected file with the given user (read-write)"},
			cmdInfo{cmd: "unshare <file> <username>", info: "Unshare the selected file from the given user"},
			cmdInfo{cmd: "delete_acct", info: "Delete the current user's account. Must have no owned files or directories to call this function"},
			cmdInfo{cmd: "logout", info: "Log out of the server"},
			cmdInfo{cmd: "chmod_r <file> <username>", info: "Change share permissions on file with username to be readonly. File must previously have been shared with the given user"},
			cmdInfo{cmd: "chmod_r <file> <username>", info: "Change share permissions on file with username to be read-write. File must previously have been shared with the given user"},
		}
	}
	for _, c := range cmds {
		fmt.Println(c.cmd)
		fmt.Println("\t" + c.info)
	}
}

/**
 * Registers a new account. Prompts the user for the following, in order:
 *   email
 *     -checks is valid email client side
 *   username
 *     -checks is valid email server side
 *   password
 *   password again
 *     -checks both are valid passwords and are equal server side
 *
 * If above parts pass client side validation, it will make a request to the server, which revalidates
 * everything and enforces unique email and username. If this passes, the server sends an email with
 * the verification code to the client. The client then has a preset amount of time to enter the verification
 * code. Upon entry, another request is made to the server, which will validate the verification code and if valid,
 * create the user account, else report an error.
 *
 * Input: N/A
 * Output: An error if an error occurs, else nil.
 *
 * Precondition: The client must not be logged in/must not have a session token.
 * Postcondition: Nothing changes (client remains logged out)
 */
func (c *Client) Register() error {
	// ensure not logged in
	if len(c.session_token) != 0 {
		return errors.New("Unable to register account.")
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Registering new account. Please enter the required information.")

	// get user email
	fmt.Print("E-mail: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return errors.New("Invalid email entry, canceling registration.")
	}
	email = strings.TrimSpace(email)
	// client side validate user email
	if !validateEmail(email) {
		return errors.New("Invalid email address, canceling registration.")
	}

	// get username
	fmt.Print("Username (max 16 alphanumeric characters): ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return errors.New("Invalid username entry, canceling registration.")
	}
	username = strings.TrimSpace(username)
	// client side validate username
	if !validateUsername(username) {
		return errors.New("Invalid username, canceling registration.")
	}

	// get password
	fmt.Println(`Passwords must be 8-24 characters. They must have no whitespace and must 
		contain at least 3 of the following: lowercase letters, 
		uppercase letters, numbers, punctuation (!#$&*+,-.:;?@^_~)`)
	fmt.Print("Password: ")
	password1, err := getPassword()
	if err != nil {
		return err
	}

	// get password again
	fmt.Print("\nConfirm Password: ")
	password2, err := getPassword()
	if err != nil {
		return err
	}
	fmt.Print("\n")

	// verify passwords match
	passwordValidated, err := validatePassword(password1, password2)
	if !passwordValidated {
		return err
	}

	// register the email on the server
	succ, err := c.EmailRegister(email, username, password1)
	if err != nil {
		return err
	}
	if !succ {
		return errors.New("Unable to register email.")
	}

	// get the verification code
	fmt.Print("You should have just received an email containing a verification code. Please type that code here to verify your account: ")
	ver_code, err := reader.ReadString('\n')
	if err != nil {
		return errors.New("Unable to read verification code, canceling registration.")
	}
	ver_code = strings.TrimSpace(ver_code)

	// validate the verification code
	succ, err = c.ValidateRegister(email, username, password1, ver_code)
	if err != nil {
		return err
	}
	if !succ {
		return errors.New("Unable to validate account.")
	}

	// if successfully registered
	color.Green("Account created!\n")
	return nil
}

/**
 * Prompt the users for their username and password. Validates both of them client side, then makes a request
 * to the server, which revalidates them and ensures they are correct. If a correct and valid combination,
 * the server creates a session token for the user and returns it.
 *
 * Input: N/A
 * Output: An error if an error occurred, else nil.
 *
 * Precondition: The client is not logged in.
 * Postcondition: On successful log in, the client's username is set to what was entered, and the session token
 *   is set to the server's return value. On failure, neither is set (both should be empty), and the client
 *   remains logged out.
 */
func (c *Client) Login() error {
	// ensure not logged in
	if len(c.session_token) != 0 {
		return errors.New("Unable to login account.")
	}

	reader := bufio.NewReader(os.Stdin)

	// get username
	fmt.Print("Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return errors.New("Invalid username entry, canceling login.")
	}
	c.username = strings.TrimSpace(username)
	// client side validate username
	if !validateUsername(c.username) {
		return errors.New("Invalid username, canceling login.")
	}

	// get password
	fmt.Print("Password: ")
	password, err := getPassword()
	fmt.Print("\n")
	if err != nil {
		return err
	}

	session_token, err := c.ValidateLogin(c.username, password)
	c.SetSessionToken(session_token)
	if err != nil {
		return err
	}

	return nil
}

/**
 * If the user is logged in, logout the user.
 *
 * Input: N/A
 * Output: If an error occured, an error, else nil.
 *
 * Precondition: The user is currently logged in.
 * Postcondition: The client's username and session token are reset to empty.
 */
func (c *Client) Logout() error {
	// ensure not logged in
	if len(c.session_token) == 0 {
		return errors.New("Unable to logout account.")
	}

	c.SetSessionToken([]byte{})
	err := c.HandleLogout(c.username)
	if err != nil {
		return err
	}
	c.username = ""
	return nil
}

/**
 * Get the session token from the client struct. Does not interact with the server.
 *
 * Input: N/A
 * Output: The session token stored in the client.
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func (c *Client) GetSessionToken() []byte {
	return c.session_token
}

/**
 * Set the session token in the client struct. Does not interact with the server.
 *
 * Input: The new session token.
 * Output: N/A
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func (c *Client) SetSessionToken(session_token []byte) {
	c.session_token = session_token
}

/**
 * Gets the current working directory for the user from the server.
 *
 * Input: N/A
 * Output:
 *   string - the current working directory. Should be "" if error is not nil
 *   error - the error, if an error returned, else nil
 *
 * Precondition: Client is logged in
 * Postcondition: N/A
 */
func (c *Client) Pwd() (string, error) {
	cwd, err := c.GetWorkingDirectory(c.username, c.session_token)
	if err != nil {
		return "", err
	}
	return cwd, nil
}

/**
 * Return the lists of subdirectories and files in the current directory.
 *
 * Input: N/A
 * Output:
 *   []string - the names of all subdirectories in the current directory, sorted alphabetically
 *   []string - the names of all files in the current directory, sorted alphabetically
 *   error - the error, if an error returned, else nil. If nil, the other arguments should be empty slices
 *
 * Precondition: Client is logged in
 * Postcondition: N/A
 */
func (c *Client) Ls() ([]string, []string, error) {
	// get subdirectories and sort them alphabetically
	subdirectories, files, err := c.GetSubdirectoriesAndFiles(c.username, c.session_token)
	if err != nil {
		return []string{}, []string{}, err
	}

	return subdirectories, files, err
}

/**
 * Change the directory to the given path. Will vaildate the path to ensure is valid.
 *
 * Input:
 *    path string - the path to cd into. Can be relative or absolute
 * Output: the error, if an error occurred, else nil.
 *
 * Precondition: Client is logged in
 * Postcondition: Working directory changed to path both client and server side, providing no error
 */
func (c *Client) CD(path string) error {
	if !validatePath(path) {
		return errors.New("Invalid path")
	}
	err := c.ChangeDirectory(c.username, c.session_token, path)
	if err != nil {
		return err
	}
	return nil
}

/**
 * Make a directory with the given name in the current directory. Validates the dirname is a valid directory name.
 *
 * Input:
 *   dirname string - the name of the new directory
 * Output: the error, if an error occurred, else nil
 *
 * Precondition: Client is logged in
 * Postcondition: Subdirectory exists with given dirname in current directory. An error if failure to make
 *   directory, or it already existed.
 */
func (c *Client) Mkdir(dirname string) error {
	if !validateDirname(dirname) {
		return errors.New("Invalid directory name")
	}
	err := c.MakeDirectory(c.username, c.session_token, dirname)
	if err != nil {
		return err
	}
	return nil
}

/**
 * Remove a directory with the given name in the current directory. Validates the dirname is a valid directory name.
 *
 * Input:
 *   dirname string - the name of the directory to be removed
 * Output: the error, if an error occurred, else nil
 *
 * Precondition: Client is logged in.
 * Postcondition: Subdirectory with given dirname no longer exists. An error if failure to remove directory, or it
 *   does not exist.
 */
func (c *Client) Rmdir(dirname string) error {
	if !validateDirname(dirname) {
		return errors.New("Invalid directory name")
	}
	err := c.RemoveDirectory(c.username, c.session_token, dirname)
	if err != nil {
		return err
	}
	return nil
}

/**
 * Remove the file with the given name from the current directory. Validates the filename is a valid filename.
 *
 * Input:
 *    filename string, the name of the file being removed
 * Output: the error, if an error occurred, else nil
 *
 * Precondition: Client is logged in.
 * Postcondition: File with filename no longer exists in current directory. An error if failure to remove file or
 *    it did not exist.
 */
func (c *Client) Rm(filename string) error {
	//allow absolute path
	if checkIfAbsolutePath(filename) {
		path, file := GetFileNameAndPath(filename)
		if !validatePath(path) {
			return errors.New("Invalid dest path")
		}

		if !validateFilename(file) {
			return errors.New("Invalid dest filename")
		}
	} else { // not absolute path
		if !validateFilename(filename) {
			return errors.New("Invalid dest filename")
		}
	}

	err := c.RemoveFile(c.username, c.session_token, filename)
	if err != nil {
		return err
	}
	return nil
}

/**
 * Upload the file stored locally in the given filepath to the given filename (which may or may not have an
 * absolute path). Modify indicates if the file ought to already exist and be overwritten.
 *
 * Input:
 *   src_filepath string - the path to the local file stored on the client computer
 *   dest_filename string - the filepath and filename of where to store the file on the server
 *   modify bool - true if delete and overwrite an existing file (it must exist), false if must be new file
 * Output: the error, if an error occurred, else nil
 *
 * Precondition: Client is logged in.
 * Postcondition: The file is uploaded to the server. An error if failure to upload file.
 */
func (c *Client) Upload(src_filepath string, dest_filename string, modify bool) error {
	_, src_filename := GetFileNameAndPath(src_filepath)
	if !validateFilename(src_filename) {
		return errors.New("Invalid src filename")
	}

	//allow absolute path
	if checkIfAbsolutePath(dest_filename) {
		dest_path, filename := GetFileNameAndPath(dest_filename)
		if !validatePath(dest_path) {
			return errors.New("Invalid dest path")
		}

		if !validateFilename(filename) {
			return errors.New("Invalid dest filename")
		}
	} else { // not absolute path
		if !validateFilename(dest_filename) {
			return errors.New("Invalid dest filename")
		}
	}
	file_bytes, err := ioutil.ReadFile(src_filepath)
	if err != nil {
		return err
	}

	err = c.uploadFile(c.username, c.session_token, file_bytes, dest_filename, modify)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Download the file at the given dest_filepath to the given src_filename.
 *
 * Input:
 *   src_filename - the local filepath/filename to store the downloaded file
 *   dest_filepath - the location of the file on the server to be downloaded
 * Output: the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in.
 * Postcondition: The client has the file downloaded to src_filename, providing there is no error on download.
 */
func (c *Client) Download(src_filename string, dest_filepath string) error {
	//check for absolute path
	if checkIfAbsolutePath(src_filename) { //validate absolute path
		src_path, filename := GetFileNameAndPath(src_filename)
		if !validatePath(src_path) {
			return errors.New("Invalid dest path")
		}

		if !validateFilename(filename) {
			return errors.New("Invalid dest filename")
		}
	} else { //validate relative path
		if !validateFilename(src_filename) {
			return errors.New("Invalid file name")
		}
	}

	file_contents, err := c.Concatenate(c.username, c.session_token, src_filename)
	if err != nil {
		return err
	}

	//write to file
	err = ioutil.WriteFile(dest_filepath, file_contents, 0600)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Print out the file at the given dest_filepath to standard out.
 *
 * Input:
 *   dest_filepath - the location of the file on the server to be downloaded
 * Output: the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in.
 * Postcondition: N/A.
 */
func (c *Client) Cat(file string) (string, error) {
	//check for absolute path
	if checkIfAbsolutePath(file) { //validate absolute path
		path, filename := GetFileNameAndPath(file)
		if !validatePath(path) {
			return "", errors.New("Invalid path")
		}

		if !validateFilename(filename) {
			return "", errors.New("Invalid file name")
		}
	} else { //validate relative path
		if !validateFilename(file) {
			return "", errors.New("Invalid file name")
		}
	}

	file_contents, err := c.Concatenate(c.username, c.session_token, file)
	if err != nil {
		return "", err
	}

	//convert bytes to string
	file_contents_string := string(file_contents)

	return file_contents_string, nil
}

/**
 * Share the file with name filename on the server with sharee_username. Write_perms indicates read only or
 * read/write.
 *
 * Input:
 *   filename - the name of the file in the directory to be shared
 *   sharee_username - the name of the user with whom the file is being shared
 *   write_perms - true if read/write, false if read-only
 * Output: the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in. The client must be owner of the file.
 * Postcondition: The file is shared with the sharee, providing there is no error on download.
 */
func (c *Client) Share(filename string, sharee_username string, write_perms bool) error {
	if !validateFilename(filename) {
		return errors.New("Invalid filename")
	}

	if !validateUsername(sharee_username) {
		return errors.New("Invalid username")
	}

	err := c.ShareFile(c.username, c.session_token, filename, sharee_username, write_perms)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Unshare the file with name filename on the server with sharee_username.
 *
 * Input:
 *   filename - the name of the file in the directory to be shared
 *   sharee_username - the name of the user with whom the file is being shared
 * Output: the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in. The client must be owner of the file.
 * Postcondition: The file is no longer shared with the sharee, providing there is no error on download.
 */
func (c *Client) Unshare(filename string, sharee_username string) error {

	if !validateFilename(filename) {
		return errors.New("Invalid filename")
	}

	if !validateUsername(sharee_username) {
		return errors.New("Invalid username")
	}

	err := c.UnshareFile(c.username, c.session_token, filename, sharee_username)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Delete the account of the current user.
 *
 * Input: N/A
 * Output: The error, if an error occurred, else nil.
 *
 * Precondition: The client is logged in and owns no directories or files. The client's working directory is "/".
 * Postcondition: The client is logged out, and no longer stores its session token or username. The server
 *   maintains no records of the client in users, directories, or shared table.
 */
func (c *Client) DeleteAcct() error {
	err := c.DeleteAccount(c.username, c.session_token)
	if err != nil {
		return err
	}
	c.SetSessionToken([]byte{})
	c.username = ""
	return nil
}

/**
 * List all people with whom the client has shared the given file.
 *
 * Input:
 *   filename string - the name of the file whose sharees are to be listed
 * Output:
 *   []string - the names of all sharees. Empty if error is not nil
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: client is logged in
 * Postcondition: N/A (nothing changes)
 */
func (c *Client) LsSharees(filename string) ([]string, error) {
	if !validateFilename(filename) {
		return []string{}, errors.New("Invalid filename")
	}

	sharees, err := c.ListSharees(c.username, c.session_token, filename)
	if err != nil {
		return []string{}, err
	}

	return sharees, nil
}

/**
 * Change the permissions on the given filename shared with sharee_username. Write_perms indicates read only or
 * read/write.
 *
 * Input:
 *   filename - the name of the file in the directory to have share permissions for sharee modified
 *   sharee_username - the name of the user with whom the file's permissions is being changed
 *   write_perms - true if read/write, false if read-only
 * Output: the error, if an error occurred, else nil
 *
 * Precondition: The client is logged in. The client must be owner of the file. The client has already shared the
 *   file with sharee_username.
 * Postcondition: The permissions for the sharee_username have been updated, providing there is no error on download.
 */
func (c *Client) Chmod(filename string, sharee_username string, write_perms bool) error {
	if !validateFilename(filename) {
		return errors.New("Invalid filename")
	}

	if !validateUsername(sharee_username) {
		return errors.New("Invalid username")
	}

	err := c.ChmodFile(c.username, c.session_token, filename, sharee_username, write_perms)
	if err != nil {
		return err
	}
	return nil

}
