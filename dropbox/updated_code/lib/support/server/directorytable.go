package server

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

/**
 * Creates a root and shared directory for the given user
 *
 * Input:
 *	username string - the username to be checked
 * Output: if the username is unique then it returns nil, else error
 *
 * Precondition: username and corresponding session token have been validated and sanitized
 * Postcondition: username has a root and shared directory in the directories table
 */

func createRootandShared(username string) error {
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return err
	}

	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	//check if

	// insert root directory
	stmt, err := db.Prepare(`INSERT INTO directories 
		(username, path, name)
		VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// and execute the insert
	_, err = stmt.Exec(username, "", "/")
	if err != nil {
		return err
	}

	//insert shared folder
	stmt, err = db.Prepare(`INSERT INTO directories 
		(username, path, name)
		VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// and execute the insert
	_, err = stmt.Exec(username, "/", "shared/")
	if err != nil {
		return err
	}

	return nil
}

/**
 * Removes all directories from the directories table for the given user
 *
 * Input:
 *	username string - the username to be checked
 * Output: if the username is unique then it returns nil, else error
 *
 * Precondition: username and corresponding session token have been validated and sanitized
 * Postcondition: username does not have any directories in the directories table
 */

func removeAllDirectories(username string) error {
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return err
	}

	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`DELETE FROM directories
		WHERE username = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Adds directory to the directories table with the given username and file path and directory name
 *
 * Input:
 *	username string - the username to be checked
 *	curr_dir string - the path of the new directory
 *	dirname string - teh new directory name
 * Output: if the username is unique and the directory doesn't exist yet then it returns nil, else error
 *
 * Precondition: username and corresponding session token and directory path have been validated and sanitized
 * Postcondition: username has the new directory in its directories table
 */

func createDirectoryByUsername(username string, curr_dir string, dirname string) error {
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return err
	}

	dirname = dirname + "/"

	count_dirs, err := getNumDirectoriesByUsername(username, curr_dir, dirname)
	if err != nil {
		return err
	}
	if count_dirs != 0 {
		return errors.New("Directory already exists")
	}

	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`INSERT INTO directories 
		(username, path, name)
		VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// and execute the insert
	_, err = stmt.Exec(username, curr_dir, dirname)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Adds directory to the table named with the sharer' username with the path of /shared/ to the sharees directory tree
 *
 * Input:
 *	sharer_username string - the sharer's username
 *	sharee_username string - the sharee's username
 * Output: if both usernames are unique and the directory has not already been created return nil, else error
 *
 * Precondition: both usernames and corresponding session tokens have been validated and sanitized
 * Postcondition: username has the new directory in its shared folder table
 */

func makeSharerDirectoryForSharee(sharer_username string, sharee_username string) error {
	err := validateUniqueUserByUsername(sharer_username)
	if err != nil {
		return err
	}

	err = validateUniqueUserByUsername(sharee_username)
	if err != nil {
		return err
	}

	if sharer_username == sharee_username {
		return errors.New("Sharer and sharee cannot be the same person")
	}

	dirname := sharer_username + "/"
	path := "/shared/" // with this function, it must be in the shared directory

	count_dirs, err := getNumDirectoriesByUsername(sharee_username, path, dirname)
	if err != nil {
		return err
	}
	if count_dirs == 1 {
		// because the directory already exists, it is not an error, so just return nil
		return nil
	} else if count_dirs >= 2 {
		// should never happen
		return errors.New("Shared directory exists more than once")
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`INSERT INTO directories
		(username, path, name)
		VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sharee_username, path, dirname)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Removes sharer's directory from the shared folder of the sharee
 *
 * Input:
 *	sharer_username string - the sharer's username
 *	sharee_username string - the sharee's username
 * Output: if both usernames are unique and the directory has not already been created return nil, else error
 *
 * Precondition: both usernames and corresponding session tokens have been validated and sanitized
 * Postcondition: the sharee no longer has the sharer's folder in their shared folder in the directories table
 */

func removeSharerDirectoryForSharee(sharer_username string, sharee_username string) error {
	err := validateUniqueUserByUsername(sharer_username)
	if err != nil {
		return err
	}

	err = validateUniqueUserByUsername(sharee_username)
	if err != nil {
		return err
	}

	if sharer_username == sharee_username {
		return errors.New("Sharer and sharee cannot be the same person")
	}

	dirname := sharer_username + "/"
	path := "/shared/" // with this function, it must be in the shared directory

	count_dirs, err := getNumDirectoriesByUsername(sharee_username, path, dirname)
	if err != nil {
		return err
	}
	if count_dirs != 1 {
		// trying to remove nonexistent directory should never happen
		// should never call this function without first verifying sharer has shared a file with sharee
		// and that the number of shared files just changed to 0
		return errors.New("Shared directory does not exist")
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`DELETE FROM directories
		WHERE username = ? AND path = ? AND name = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sharee_username, path, dirname)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Removes directory from the directories table for the given username and file path and directory name
 *
 * Input:
 *	username string - the username to be checked
 *	curr_dir string - the path of the new directory
 *	dirname string - teh new directory name
 * Output: if the username and directory exist it returns nil, else error
 *
 * Precondition: username and corresponding session token and directory path have been validated and sanitized
 * Postcondition: username no longer has the directory in its directories table
 */

func removeDirectoryByUsername(username string, curr_dir string, dirname string) error {
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return err
	}

	dirname = dirname + "/"

	count_dirs, err := getNumDirectoriesByUsername(username, curr_dir, dirname)
	if err != nil {
		return err
	}
	if count_dirs != 1 {
		// either 0, in which case does not exist, or 2, which should not happen
		return errors.New("Directory does not exist")
	}

	// the full path of the subdirectory
	subdir := curr_dir + dirname

	// ensure no files in folder being removed
	files, err := getFilesByUsernameAndCurrDir(username, subdir)
	if err != nil {
		return err
	}
	if len(files) != 0 {
		return errors.New("Cannot remove directory that has files.")
	}

	// ensure no subdirectories in directory being removed
	subdirectories, err := getSubdirectoriesByUsernameAndCurrDir(username, subdir)
	if err != nil {
		return err
	}
	if len(subdirectories) != 0 {
		return errors.New("Cannot remove directory that has subdirectories.")
	}

	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`DELETE FROM directories
		WHERE username = ? AND path = ? AND name = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, curr_dir, dirname)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Returns the number of directories given a username, path, and directory name
 *
 * Input:
 *	username string - the username to be checked
 *	parent string - the path of the directory
 *	child string - the directory name
 * Output: return the number of times the directory shows up in the table, error if db fails
 *
 * Precondition: username and corresponding session token and directory path and namehave been validated and sanitized
 * Postcondition: N/A
 */

func getNumDirectoriesByUsername(username string, parent string, child string) (int, error) {
	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return -1, err
	}
	defer db.Close()

	// prepare statement for determining number of users
	stmt, err := db.Prepare("SELECT COUNT(*) FROM directories WHERE username = ? AND path = ? AND name = ?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	// query the db to determine number of users
	var numUsers int
	err = stmt.QueryRow(username, parent, child).Scan(&numUsers) // it will only return 1 row cause count
	if err != nil {
		return -1, err
	}

	return numUsers, nil
}

/**
 * Checks whether the given directory exists exactly once given a username, path, and directory name
 *
 * Input:
 *	username string - the username to be checked
 *	parent string - the path of the directory
 *	child string - the directory name
 * Output: if the directory is in the table once return nil, else error
 *
 * Precondition: username and corresponding session token and directory path and name have been validated and sanitized
 * Postcondition: N/A
 */

func validateDirectoryExistsByUsername(username string, parent string, child string) error {
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	numDirectories, err := getNumDirectoriesByUsername(username, parent, child)
	if err != nil {
		return err
	}
	if numDirectories != 1 {
		return errors.New("Directory does not exist")
	}

	return nil
}

/**
 * Returns all subdirectories given a username and path
 *
 * Input:
 *	username string - the username to be checked
 *	curr_dir string - the path of the directory
 * Output: returns a string array containing all subdirectories of the given path, error if db fails
 *
 * Precondition: username and corresponding session token and directory path have been validated and sanitized
 * Postcondition: N/A
 */

func getSubdirectoriesByUsernameAndCurrDir(username string, curr_dir string) ([]string, error) {
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return []string{}, err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return []string{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT name FROM directories
		WHERE username = ? AND path = ?
		ORDER BY name ASC`)
	if err != nil {
		return []string{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, curr_dir)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	// read the name of all subdirectories
	var subdirectories []string
	for rows.Next() {
		var dirname string
		err = rows.Scan(&dirname)
		if err != nil {
			return []string{}, err
		}
		subdirectories = append(subdirectories, dirname)
	}

	err = rows.Err()
	if err != nil {
		return []string{}, err
	}

	return subdirectories, nil
}

/**
 * Checks if a path is valid for a given username. Adds backslash at the end of path if it doesn't have one.
 * Splits path and gets last directory and parent path. Then uses those to see if they are in the directory	table.
 *
 *
 * Input:
 *	username string - the username to be checked
 *	path string - the path of the directory
 * Output: returns an altered path which converts any .. into a valid path, error path is empty or doesn't exist
 *
 * Precondition: username and corresponding session token and directory path have been validated and sanitized
 * Postcondition: N/A
 */

func checkValidDirectoryPath(username string, path string) (altered_path string, err error) {
	//check for empty string
	if path == "" {
		return "", errors.New("Cannot give empty string as path")
	}

	var parent string
	var child string

	//add backslash at end of path if they don't have it
	if path[len(path)-1:] != "/" {
		path = path + "/"
	}

	if string(path[0]) == "/" { //check absolute path
		path = handleDots(path)
		parent, child, err = getParentPathAndChildofDirectory(path)
		if err != nil {
			return "", err
		}
	} else { //check relative path
		curr_dir, err := getCurrDir(username)
		if err != nil {
			return "", err
		}

		path = handleDots(curr_dir + path)
		parent, child, err = getParentPathAndChildofDirectory(path)
		if err != nil {
			return "", err
		}
	}

	fmt.Println("USERNAME:", username, "PARENT", parent, "CHILD", child)
	err = validateDirectoryExistsByUsername(username, parent, child)
	if err != nil {
		return "", err
	}

	return path, nil

}

/**
 * Checks if absolute path to a file is valid for a given username.
 * Splits path and gets last directory and parent path. Then uses those to see if they are in the directory	table.
 *
 *
 * Input:
 *	username string - the username to be checked
 *	path string - the path of the directory
 * Output: returns an altered path which converts any .. into a valid path, error path is empty or doesn't exist
 *
 * Precondition: username and corresponding session token and directory path have been validated and sanitized
 * Postcondition: N/A
 */

func checkValidAbsolutePathWithFile(username string, path string) (altered_path string, file_name string, err error) {
	//check for empty string
	if path == "" {
		return "", "", errors.New("Cannot give empty string as path")
	}

	var parent string
	var child string

	if path[len(path)-1:] == "/" {
		return "", "", errors.New("This is not a valid file path")
	}

	if string(path[0]) != "/" {
		return "", "", errors.New("This is not an absolute path")
	}

	new_path, file_name, err := getParentPathAndChildofFile(path)
	path = handleDots(new_path)
	parent, child, err = getParentPathAndChildofDirectory(path)
	if err != nil {
		return "", "", err
	}

	fmt.Println("USERNAME:", username, "PARENT", parent, "CHILD", child)
	err = validateDirectoryExistsByUsername(username, parent, child)
	if err != nil {
		return "", "", err
	}

	return new_path, file_name, nil

}

/**
 * Gets last directory and parent path to that directory from a given absolute path. If given root then return no parent
 *
 * Input:
 *	path string - the path of the directory
 * Output: returns the parent path and the child direcotry, error if path is empty or doesn't exist
 *
 * Precondition: directory path has been validated and sanitized
 * Postcondition: N/A
 */

func getParentPathAndChildofDirectory(path string) (parent string, child string, err error) {
	//check for empty string
	if path == "" {
		return "", "", errors.New("Cannot give empty string as path")
	}

	//check to make sure it is directory
	if path[len(path)-1:len(path)] != "/" {
		return "", "", errors.New("This path is not a directory")
	}

	//base case of root directory
	if path == "/" {
		return "", path, nil
	}

	//remove last slash
	no_last_slash := path[:len(path)-1]
	//get last slash index
	last_slash := strings.LastIndex(no_last_slash, "/")
	parent = path[:last_slash+1]
	child = path[last_slash+1 : len(path)]

	return parent, child, nil
}

/**
 * Gets last directory and parent path of a file of a given absolute path. Checks if it is a path to a file
 *
 * Input:
 *	path string - the absolute path of the file
 * Output: returns the parent path and the child direcotry, error if path is empty or doesn't exist
 *
 * Precondition: directory path has been validated and sanitized
 * Postcondition: N/A
 */

func getParentPathAndChildofFile(path string) (parent string, child string, err error) {
	//check for empty string
	if path == "" {
		return "", "", errors.New("Cannot give empty string as path")
	}

	//check to make sure it is not a directory
	if path[len(path)-1:len(path)] == "/" {
		return "", "", errors.New("This path does not direct to a file")
	}

	//get last slash index
	last_slash := strings.LastIndex(path, "/")
	parent = path[:last_slash+1]
	child = path[last_slash+1 : len(path)]

	return parent, child, nil
}

/**
 * Converts dots in path to valid path. Does not go outside root
 *
 * Input:
 *	path string - the given path
 * Output: returns the path with text instead of dots
 *
 * Precondition: directory path has been validated and sanitized
 * Postcondition: N/A
 */

func handleDots(path string) string {
	//fmt.Println("PATH:", path)
	dirs := strings.Split(path, "/")
	out_path := "/"
	for _, dir := range dirs {
		//handle ..
		if dir == ".." {
			//don't go out of root directory
			if out_path != "/" {
				remove_last_slash := out_path[:len(out_path)-1]
				last_slash_index := strings.LastIndex(remove_last_slash, "/")
				out_path = out_path[:last_slash_index+1]
			}
		} else if dir != "." && dir != "" {
			out_path += dir + "/"
		}
		//fmt.Println("DIR:", dir)
		//fmt.Println("OUT_PATH:", out_path)
	}

	//fmt.Println("FINAL OUT_PATH:", out_path)
	return out_path
}
