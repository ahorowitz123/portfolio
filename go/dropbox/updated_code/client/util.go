package main

import (
	"errors"
	"regexp"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

/**
 * Checks if the email is a valid email string.
 *
 * Input:
 *   email string - the email address to be checked
 * Output: true if a valid email address, else false
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func validateEmail(email string) bool {
	// Regex from https://www.socketloop.com/tutorials/golang-validate-email-address-with-regular-expression
	emailRegexp := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegexp.MatchString(email)
}

/**
 * Checks if the username is a valid username string.
 *
 * Input:
 *   username string - the username to be checked
 * Output: true if a valid username, else false
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func validateUsername(username string) bool {
	// Regex from https://www.socketloop.com/tutorials/golang-regular-expression-alphanumeric-underscore
	usernameRegexp := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return len(username) <= 16 && usernameRegexp.MatchString(username)
}

/**
 * Check if the passwords are equal and are a valid password.
 *
 * Input:
 *   password1 string - the first password enetered during registration
 *   password2 string - the second password entered during registration
 * Output:
 *   bool - true if passwords equal and valid, else false
 *   error - error describing why passwords failed to validate, nil if valid
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func validatePassword(password1 string, password2 string) (bool, error) {
	if password1 != password2 {
		return false, errors.New("Passwords do not match")
	}

	if len(password1) < 8 || len(password1) > 24 {
		return false, errors.New("Passwords must be 8-24 characters")
	}

	whitespaceRegexp := regexp.MustCompile(`[\s]+`)
	lowercaseRegexp := regexp.MustCompile(`[a-z]+`)
	uppercaseRegexp := regexp.MustCompile(`[A-Z]+`)
	numericRegexp := regexp.MustCompile(`[0-9]+`)
	punctuationRegexp := regexp.MustCompile(`[(!#$&*+,-.:;?@^_~)]+`)
	invalidRegexp := regexp.MustCompile(`[^(!#$&*+,-.:;?@^_~)0-9A-Za-z]+`)

	if whitespaceRegexp.MatchString(password1) || invalidRegexp.MatchString(password1) {
		return false, errors.New("Passwords have invalid characters")
	}

	counter := 0
	if lowercaseRegexp.MatchString(password1) {
		counter++
	}
	if uppercaseRegexp.MatchString(password1) {
		counter++
	}
	if numericRegexp.MatchString(password1) {
		counter++
	}
	if punctuationRegexp.MatchString(password1) {
		counter++
	}

	if counter >= 3 {
		return true, nil
	} else {
		return false, errors.New("Password did not meet constraints")
	}
}

/**
 * Check that the given dirname is a valid name for a directory.
 *
 * Input:
 *   dirname string - the name of the directory to be checked
 * Output: true if a valid dirname, else false
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func validateDirname(dirname string) bool {
	dirnameRegexp := regexp.MustCompile("^[a-zA-Z0-9]+$")
	dirnameRegexp2 := regexp.MustCompile("^\\.{1,2}$")
	return (len(dirname) <= 32 && dirnameRegexp.MatchString(dirname)) || dirnameRegexp2.MatchString(dirname)
}

/**
 * Check that the given filename is a valid name for a file.
 *
 * Input:
 *   filename string - the name of the file to be checked
 * Output: true if a valid filename, else false
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func validateFilename(filename string) bool {
	filenameRegexp := regexp.MustCompile("^[a-zA-Z0-9_\\-.]+$")
	return len(filename) <= 18 && filenameRegexp.MatchString(filename)
}

/**
 * Get the password from standard in. Does not show password on standard in.
 *
 * Input: N/A
 * Output:
 *   string - the password typed into stdin. "" if error not nil
 *   error - the error, if there was one, else nil
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func getPassword() (string, error) {
	// read the password as bytes
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	// return the password as a string
	return string(password), nil

}

/**
 * Check that the given path is a valid path.
 *
 * Input:
 *   path string - the name of the path to be checked
 * Output: true if a valid path, else false
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func validatePath(path string) bool {
	path = strings.TrimSpace(path)
	//make sure string is not empty

	if path == "" {
		return false
	}

	//check if path is root
	if path == "/" {
		return true
	}

	//check if path is .
	if path == "." {
		return true
	}

	//check if path is ..
	if path == ".." {
		return true
	}

	//remove first slash
	if string(path[0]) == "/" {
		path = path[1:len(path)]
	}

	//remove last slash
	if string(path[len(path)-1]) == "/" {
		path = path[:len(path)-1]
	}

	//get all directories in path
	dirs := strings.Split(path, "/")

	//check all dirs
	for _, dir := range dirs {
		valid := validateDirname(dir)
		if !valid {
			return false
		}
	}

	return true
}

/**
 * Get the filename and path from the given filepath.
 *
 * Input:
 *   path string - the path to the file, including filename
 * Output:
 *   filepath string - the path to the file
 *   file string - the name of the file
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func GetFileNameAndPath(path string) (filepath string, file string) {
	last_slash_index := strings.LastIndex(path, "/")
	filepath = path[:last_slash_index+1]
	file = path[last_slash_index+1:]

	return filepath, file
}

/**
 * Check if path is absolute.
 *
 * Input:
 *   path string - the path to be checked
 * Output: true if absolute, else false
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func checkIfAbsolutePath(path string) bool {
	if path == "" {
		return false
	}

	if string(path[0]) == "/" {
		return true
	}

	return false
}
