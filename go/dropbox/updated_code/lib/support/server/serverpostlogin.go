package server

import (
	"encoding/hex"
	"github.com/s17-ahorowi2-jtc2/updated_code/internal"
	"io/ioutil"
	"os"
	"strings"
)

/**
 * Get the working directory of the given user.
 *
 * Input:
 *   username string - the username given by the client
 *   session_token []byte - the session token from the client
 * Output: a pwd return type with the working directory for the user on success, or an error message on failure
 *
 * Precondition: The username and session token must be valid (checked at start).
 * Postcondition: N/A
 */
func getWorkingDirectoryHandler(username string, session_token []byte) internal.PwdReturn {
	// ensure valid user
	err := checkSessionToken(username, session_token)
	if err != nil {
		return internal.PwdReturn{Path: "", Err: err.Error(), Logout: true}
	}

	// get and return working directory
	curr_dir, err := getCurrDir(username)
	if err != nil {
		return internal.PwdReturn{Path: "", Err: err.Error(), Logout: false}
	}

	return internal.PwdReturn{Path: curr_dir, Err: "", Logout: false}
}

/**
 * Get the subdirectories and files of the current directory for the given user.
 *
 * Input:
 *   username string - the username given by the client
 *   session_token []byte - the session token from the client
 * Output: A list return type with the subdirectories and files of the current directory, or an error message
 *   on failure.
 *
 * Precondition: The username and session token must be valid (checked at start).
 * Postcondition: N/A
 */
func lsHandler(username string, session_token []byte) internal.ListReturn {
	err := checkSessionToken(username, session_token)
	if err != nil {
		return internal.ListReturn{Dirs: []string{}, Files: []string{}, Err: err.Error(), Logout: true}
	}

	// get current directory
	curr_dir, err := getCurrDir(username)
	if err != nil {
		return internal.ListReturn{Dirs: []string{}, Files: []string{}, Err: err.Error(), Logout: true}
	}

	// get subdirectories
	subdirectories, err := getSubdirectoriesByUsernameAndCurrDir(username, curr_dir)
	if err != nil {
		return internal.ListReturn{Dirs: []string{}, Files: []string{}, Err: err.Error(), Logout: false}
	}

	var files []string

	//Check which tables to get files from depending if they are in shared folder or not
	if !checkIfCurrDirIsInSharedFolder(curr_dir) {
		// get files from files table
		files, err = getFilesByUsernameAndCurrDir(username, curr_dir)
		if err != nil {
			return internal.ListReturn{Dirs: []string{}, Files: []string{}, Err: err.Error(), Logout: false}
		}
	} else {
		// get files from shared table
		owner := getSharedOwner(curr_dir)
		files, err = getSharedFilesByUsernameAndOwner(username, owner)
		if err != nil {
			return internal.ListReturn{Dirs: []string{}, Files: []string{}, Err: err.Error(), Logout: false}
		}
	}

	return internal.ListReturn{Dirs: subdirectories, Files: files, Err: "", Logout: false}
}

/**
 * Change the directory of the current user given the working path.
 *
 * Input:
 *   username string - the username given by the client
 *   session_token []byte - the session token from the client
 * Output: An error return type indicating if there was an error.
 *
 * Precondition: The username and session token must be valid (checked at start).
 * Postcondition: If successful, the working directory for the user is changed to the given path.
 */
func cdHandler(username string, session_token []byte, path string) internal.ErrorReturn {
	err := checkSessionToken(username, session_token)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: true}
	}

	// ensure valid path
	path, err = checkValidDirectoryPath(username, path)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	// update working directory as appropriate
	if checkIfAbsolutePath(path) {
		err = setCurrDir(username, path)
		if err != nil {
			return internal.ErrorReturn{Err: err.Error(), Logout: false}
		}
	} else {
		old_path, err := getCurrDir(username)
		if err != nil {
			return internal.ErrorReturn{Err: err.Error(), Logout: false}
		}

		err = setCurrDir(username, old_path+path)
		if err != nil {
			return internal.ErrorReturn{Err: err.Error(), Logout: false}
		}
	}

	return internal.ErrorReturn{Err: "", Logout: false}
}

/**
 * Make a subdirectory with the given name in the current directory.
 *
 * Input:
 *   username string - the username given by the client
 *   session_token []byte - the session token from the client
 * Output: An error return type indicating if there was an error.
 *
 * Precondition: The username and session token must be valid (checked at start).
 * Postcondition: If successful, the current directory now has a subdirectory with the given name.
 */
func mkdirHandler(username string, session_token []byte, dirname string) internal.ErrorReturn {
	err := checkSessionToken(username, session_token)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: true}
	}

	// ensure valid directory name
	if !validateDirname(dirname) {
		return internal.ErrorReturn{Err: "Invalid directory name.", Logout: false}
	}

	// get the current directory
	curr_dir, err := getCurrDir(username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	// cannot create directory in shared folder
	if strings.HasPrefix(curr_dir, "/shared/") {
		// cannot remove shared directory
		return internal.ErrorReturn{Err: "Unable to create folder in shared directory.", Logout: false}
	}

	// create the directory. This will check if valid or not
	err = createDirectoryByUsername(username, curr_dir, dirname)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	return internal.ErrorReturn{Err: "", Logout: false}
}

/**
 * Remove the subdirectory with the given name in the current directory.
 *
 * Input:
 *   username string - the username given by the client
 *   session_token []byte - the session token from the client
 * Output: An error return type indicating if there was an error.
 *
 * Precondition: The username and session token must be valid (checked at start).
 * Postcondition: If successful, the current directory no longer has a directory with the given dirname. To be
 *   successful, that directory must have existed and been empty.
 */
func rmdirHandler(username string, session_token []byte, dirname string) internal.ErrorReturn {
	err := checkSessionToken(username, session_token)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: true}
	}

	// ensure valid dirname
	if !validateDirname(dirname) {
		return internal.ErrorReturn{Err: "Invalid directory name.", Logout: false}
	}

	// get current directory
	curr_dir, err := getCurrDir(username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	if curr_dir == "/" && dirname == "shared" {
		// cannot remove shared directory
		return internal.ErrorReturn{Err: "Unable to delete shared directory.", Logout: false}
	}

	if strings.HasPrefix(curr_dir, "/shared/") {
		// cannot remove shared file
		return internal.ErrorReturn{Err: "Unable to remove directory in shared directory.", Logout: false}
	}

	// remove the directory. This should check if directory is empty or not
	err = removeDirectoryByUsername(username, curr_dir, dirname)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	return internal.ErrorReturn{Err: "", Logout: false}
}

/**
 * Remove the file with the given name in the current directory.
 *
 * Input:
 *   username string - the username given by the client
 *   session_token []byte - the session token from the client
 * Output: An error return type indicating if there was an error.
 *
 * Precondition: The username and session token must be valid (checked at start).
 * Postcondition: If successful, the current directory no longer has a file with the given filename. To be
 *   successful, that file must not have been shared with anyone.
 */
func rmfileHandler(username string, session_token []byte, filename string) internal.ErrorReturn {
	err := checkSessionToken(username, session_token)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: true}
	}

	curr_dir, err := getCurrDir(username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	//check for absolute path
	if checkIfAbsolutePath(filename) {
		//check if valid and convert dots
		true_path, file, err := checkValidAbsolutePathWithFile(username, filename)
		if err != nil {
			return internal.ErrorReturn{Err: err.Error(), Logout: false}
		}

		//change variables
		curr_dir = true_path
		filename = file
	}

	// ensure valid filename
	if !validateFilename(filename) {
		return internal.ErrorReturn{Err: "Invalid file name.", Logout: false}
	}

	if strings.HasPrefix(curr_dir, "/shared/") {
		// cannot remove shared file
		return internal.ErrorReturn{Err: "Unable to remove file in shared directory.", Logout: false}
	}

	// Handle sharing inside removeFileByUsername
	file_checksum, err := removeFileByUsername(username, filename, curr_dir)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}
	// physically delete the file if a checksum was returned
	if file_checksum != "" {
		dirpath := basedir + "/"
		file_to_delete := dirpath + file_checksum
		err = os.Remove(file_to_delete)
		if err != nil {
			return internal.ErrorReturn{Err: err.Error(), Logout: false}
		}
	}

	return internal.ErrorReturn{Err: "", Logout: false}
}

/**
 * Logout the current user.
 *
 * Input:
 *   username string - the username given by the client
 *   session_token []byte - the session token from the client
 * Output: An error return type indicating if there was an error.
 *
 * Precondition: The username and session token must be valid (checked at start).
 * Postcondition: The user's curr dir is updated to be reset to "/" in the users table. The user no longer has a
 *   session token in the user table, so they cannot make requests after logout.
 */
func logoutHandler(username string) internal.ErrorReturn {
	//set cur directory to empty string
	err := resetCurrDir(username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: true}
	}

	err = removeSessionTokenInfoByUsername(username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: true}
	}

	return internal.ErrorReturn{Err: "", Logout: false}
}

/**
 * Upload the given file bytes to the given filename. Modify indicates if it should delete and replace the file
 * that has the given filename on the server.
 *
 * Input:
 *   username string - the username given by the client
 *   session_token []byte - the session token from the client
 *   file_bytes []byte - the contents of the file being uploaded
 *   dest_filename - the filepath/filename of where to upload the file to on the server
 *   modify bool - true if must replace an existing file, false if should be a new file on the server
 *     in the directory (directory from dest_filename)
 * Output: An error return type indicating if there was an error.
 *
 * Precondition: The username and session token must be valid (checked at start). User is not out of space. File
 *   is wihtin filesize limits.
 * Postcondition:
 *   If modify true and successful, then the original file at dest_filename must have been replaced with the new
 *     contents.
 *   If modify false and successful, a new file exists at dest_filename with this contents. There was no file at
 *     dest_filename before.
 *   If unsucessful, nothing changed.
 */
func uploadHandler(username string, session_token []byte, file_bytes []byte, dest_filename string, modify bool) internal.ErrorReturn {
	// ensure file can be uploaded
	if len(file_bytes) >= 5600 {
		return internal.ErrorReturn{Err: "This file is too big", Logout: false}
	}

	// ensure valid user
	err := checkSessionToken(username, session_token)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: true}
	}

	//Check size cap on user
	err = overUserSizeLimit(username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	// get the current directory
	curr_dir, err := getCurrDir(username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	//check for absolute path
	if checkIfAbsolutePath(dest_filename) {
		//check if valid and convert dots
		true_path, file_name, err := checkValidAbsolutePathWithFile(username, dest_filename)
		if err != nil {
			return internal.ErrorReturn{Err: err.Error(), Logout: false}
		}

		//change variables
		curr_dir = true_path
		dest_filename = file_name
	}

	//validate filename
	if !validateFilename(dest_filename) {
		return internal.ErrorReturn{Err: "Invalid file name.", Logout: false}
	}

	// cannot upload/modify in shared folder
	if checkIfCurrDirIsInSharedFolder(curr_dir) && !modify {
		return internal.ErrorReturn{Err: "Unable to create file in shared directory", Logout: false}
	}

	var countFiles int
	//check if file exists if not in the shared file
	if !checkIfCurrDirIsInSharedFolder(curr_dir) {
		countFiles, err = countNumFilesByUsernameAndFileInfo(username, dest_filename, curr_dir)
		if err != nil {
			return internal.ErrorReturn{Err: "This file does not exist", Logout: false}
		}
	} else { //check if file exists by owner if in the shared file
		owner := getSharedOwner(curr_dir)
		owner_path, err := getSharedOwnerPath(username, owner, dest_filename)
		if err != nil {
			return internal.ErrorReturn{Err: "This file does not exist", Logout: false}
		}

		countFiles, err = countNumFilesByUsernameAndFileInfo(owner, dest_filename, owner_path)
		if err != nil {
			return internal.ErrorReturn{Err: "This file does not exist", Logout: false}
		}

	}

	//If uploading the file should not already exist
	if countFiles != 0 && !modify {
		return internal.ErrorReturn{Err: "Unable to create file, file already exists", Logout: false}
	}

	//If modifying, the file should already exist
	if countFiles != 1 && modify {
		return internal.ErrorReturn{Err: "File to modify does not exist on server", Logout: false}
	}

	// get the checksum of the file
	checksum := makeFileChecksum(file_bytes)
	hexFilename := hex.EncodeToString(checksum)

	//see if any files in the table exist with that name
	numInstances, err := countNumFileByChecksum(hexFilename)
	if err != nil {
		return internal.ErrorReturn{Err: "This file does not exist", Logout: false}
	}

	//If modifying, if you are possibly sharing file as owner, update all sharees checksum and possible remove file
	if modify && !checkIfCurrDirIsInSharedFolder(curr_dir) {
		//get old checksum
		old_checksum, err := getChecksum(username, dest_filename, curr_dir)
		//check for modifications
		if old_checksum == hexFilename {
			return internal.ErrorReturn{Err: "No modifications detected", Logout: false}
		}
		if err != nil {
			return internal.ErrorReturn{Err: "This file does not exist", Logout: false}
		}

		//update users checksum
		err = updateOwnerChecksum(username, dest_filename, curr_dir, hexFilename)
		if err != nil {
			return internal.ErrorReturn{Err: "This file does not exist", Logout: false}
		}

		//update all sharees with old checksum with new checksum
		err = updateSharedChecksumsByOwner(username, curr_dir, dest_filename, old_checksum, hexFilename)
		if err != nil {
			return internal.ErrorReturn{Err: "There was an error", Logout: false}
		}

		//if modifying, remove old version of file if no one else is using it
		//get amount of people using that checksum excluding yourself
		numInstances, err := countNumFileByChecksum(old_checksum)
		if err != nil {
			return internal.ErrorReturn{Err: "There was an error", Logout: false}
		}

		//make sure no one has it
		if numInstances == 0 {
			dirpath := basedir + "/"
			file_to_delete := dirpath + old_checksum
			err = os.Remove(file_to_delete)
			if err != nil {
				return internal.ErrorReturn{Err: "There was an error", Logout: false}
			}
		}
	}

	//If modifying, if you are a sharee, check if you have write permissions,
	//Then change checksums of owner and other sharees
	if modify && checkIfCurrDirIsInSharedFolder(curr_dir) {
		//get old checksum
		owner := getSharedOwner(curr_dir)
		owner_path, err := getSharedOwnerPath(username, owner, dest_filename)
		if err != nil {
			return internal.ErrorReturn{Err: "This file does not exist", Logout: false}
		}

		old_checksum, err := getChecksum(owner, dest_filename, owner_path)
		if err != nil {
			return internal.ErrorReturn{Err: "This file does not exist", Logout: false}
		}
		//check for modifications
		if old_checksum == hexFilename {
			return internal.ErrorReturn{Err: "No modifications detected", Logout: false}
		}

		has_write, err := hasWritePermissions(username, dest_filename, owner_path, old_checksum)

		//if they do not have write permissions, return
		if !has_write {
			return internal.ErrorReturn{Err: "You do not have write permissions to this shared file.", Logout: false}
		}

		//update all sharee's (including itself) and the owner's checksum
		err = updateSharedChecksumsBySharee(owner, owner_path, dest_filename, old_checksum, hexFilename)
		if err != nil {
			return internal.ErrorReturn{Err: "This file does not exist", Logout: false}
		}

		//if modifying, remove old version of file if no one else is using it
		//get amount of people using that use checksum excluding owner
		numInstances, err := countNumFileByChecksum(old_checksum)
		if err != nil {
			return internal.ErrorReturn{Err: "This file does not exist", Logout: false}
		}

		//make sure no one has it
		if numInstances == 0 {
			dirpath := basedir + "/"
			file_to_delete := dirpath + old_checksum
			err = os.Remove(file_to_delete)
			if err != nil {
				return internal.ErrorReturn{Err: "There was an error", Logout: false}
			}
		}

	}

	if numInstances == 0 {
		// create the file
		dirpath := basedir + "/"
		full_filepath := dirpath + hexFilename
		ioutil.WriteFile(full_filepath, file_bytes, 0600)
	}

	// then update database regardless, if they are not modifying
	if !modify {
		err = uploadFile(username, dest_filename, curr_dir, hexFilename)
		if err != nil {
			return internal.ErrorReturn{Err: "There was an error", Logout: false}
		}
	}

	return internal.ErrorReturn{Err: "", Logout: false}
}

/**
 * Get the contents of the file stored at the given file.
 *
 * Input:
 *   username string - the username given by the client
 *   session_token []byte - the session token from the client
 *   file string - the path/name of the file whose contents should be retrieved
 * Output: An download return type with the contents of the file if there was no error, else indicating
 *   the error.
 *
 * Precondition: The username and session token must be valid (checked at start).
 * Postcondition: If successful, the file existed. Nothing about the file changed.
 */
func catHandler(username string, session_token []byte, file string) internal.DownloadReturn {
	err := checkSessionToken(username, session_token)
	if err != nil {
		return internal.DownloadReturn{Body: []byte{}, Err: err.Error(), Logout: true}
	}

	// get the current directory
	curr_dir, err := getCurrDir(username)
	if err != nil {
		return internal.DownloadReturn{Body: []byte{}, Err: err.Error(), Logout: false}
	}

	//check for absolute path
	if checkIfAbsolutePath(file) {
		//check if valid and convert dots
		true_path, file_name, err := checkValidAbsolutePathWithFile(username, file)
		if err != nil {
			return internal.DownloadReturn{Body: []byte{}, Err: err.Error(), Logout: false}
		}

		//change variables
		curr_dir = true_path
		file = file_name
	}

	if !validateFilename(file) {
		return internal.DownloadReturn{Body: []byte{}, Err: "Invalid file name.", Logout: false}
	}

	var checksum string
	//if not shared
	if !checkIfCurrDirIsInSharedFolder(curr_dir) {
		//validate file exists and get checksum
		checksum, err = checkIfFileExistsAndgetChecksum(username, file, curr_dir)
		if err != nil {
			return internal.DownloadReturn{Body: []byte{}, Err: err.Error(), Logout: false}
		}
	} else { //if shared
		owner := getSharedOwner(curr_dir)
		owner_path, err := getSharedOwnerPath(username, owner, file)
		checksum, err = checkIfFileExistsAndgetChecksum(owner, file, owner_path)
		if err != nil {
			return internal.DownloadReturn{Body: []byte{}, Err: err.Error(), Logout: false}
		}
	}

	dirpath := basedir + "/"
	file_bytes, err := ioutil.ReadFile(dirpath + checksum)
	if err != nil {
		return internal.DownloadReturn{Body: []byte{}, Err: err.Error(), Logout: false}
	}

	return internal.DownloadReturn{Body: file_bytes, Err: "", Logout: false}
}

/**
 * Share the file stored at filename with the given user. Write_perms determines the permissions on the file.
 *
 * Input:
 *   username string - the username given by the client
 *   session_token []byte - the session token from the client
 *   filename string - the filepath/filename of the file to be shared
 *   sharee_username string - the username of the person with whom the file is being shared
 *   write_perms bool - true if read/write, false if readonly
 * Output: An error return type indicating if there was an error.
 *
 * Precondition: The username and session token must be valid (checked at start). Must be a valid filename and
 *   an actual user for sharee.
 * Postcondition: On success, the file at filename must have existed and been shared with the sharee_username.
 *   They will have write permissions if write_perms is true, else will have read only permissions. The shared
 *   table should be updated to reflect these changes. The sharee will now see the file shared in the location
 *   /shared/<sharer_username>/<filename>. ls_sharees will no longer list the sharee as an actual sharee. If not
 *   successful, nothing changes.
 */
func sharefileHandler(sharer_username string, session_token []byte, filename string, sharee_username string, write_perms bool) internal.ErrorReturn {
	// write_perms should be true if read and write, false if read only
	err := checkSessionToken(sharer_username, session_token)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: true}
	}

	// ensure valid filename
	if !validateFilename(filename) {
		return internal.ErrorReturn{Err: "Invalid filename", Logout: false}
	}

	// ensure valid sharer_username
	if !validateUsername(sharer_username) {
		return internal.ErrorReturn{Err: "Invalid username", Logout: false}
	}

	// ensure valid username
	if !validateUsername(sharee_username) {
		return internal.ErrorReturn{Err: "Invalid sharee username", Logout: false}
	}

	if sharer_username == sharee_username {
		return internal.ErrorReturn{Err: "Cannot share file with yourself", Logout: false}
	}

	// get the current directory
	curr_dir, err := getCurrDir(sharer_username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	// cannot share file that has been shared with you
	// This check may be unnecessary
	if strings.HasPrefix(curr_dir, "/shared/") {
		return internal.ErrorReturn{Err: "Unable to share file that has been shared with you", Logout: false}
	}

	// check that file exists in the current directory in the files table - if so, you have permission to share it
	count_files, err := countNumFilesByUsernameAndFileInfo(sharer_username, filename, curr_dir)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}
	if count_files != 1 {
		return internal.ErrorReturn{Err: "File does not exist", Logout: false}
	}

	// get the checksum of the file
	checksum, err := getChecksum(sharer_username, filename, curr_dir)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	// check that the sharee exists
	num_sharees, err := getNumUsersWithUsername(sharee_username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}
	if num_sharees != 1 {
		return internal.ErrorReturn{Err: "Sharee does not exist", Logout: false}
	}

	// make the shared directory for the sharee. This will check if it already exists or not
	err = makeSharerDirectoryForSharee(sharer_username, sharee_username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	// actually share the file with the other user, meaning update shared file table
	err = shareFileWithSharee(sharer_username, filename, curr_dir, checksum, sharee_username, write_perms)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	return internal.ErrorReturn{Err: "", Logout: false}
}

/**
 * Unshare the file stored at filename with the given user.
 *
 * Input:
 *   username string - the username given by the client
 *   session_token []byte - the session token from the client
 *   filename string - the filepath/filename of the file to be shared
 *   sharee_username string - the username of the person with whom the file is being shared
 * Output: An error return type indicating if there was an error.
 *
 * Precondition: The username and session token must be valid (checked at start). Must be a valid filename and
 *   an actual user for sharee. The file must have previously been sahred with the sharee.
 * Postcondition: On success, the file is no longer shared with the sharee. The shared table is updated to
 *   reflect this, and the sharee will no longer see the file. ls_sharees will no longer list the sharee as an
 *   actual sharee. If not successful, nothing changes.
 */

func unsharefileHandler(sharer_username string, session_token []byte, filename string, sharee_username string) internal.ErrorReturn {
	err := checkSessionToken(sharer_username, session_token)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: true}
	}

	// ensure valid filename
	if !validateFilename(filename) {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	// ensure valid sharer_username
	if !validateUsername(sharer_username) {
		return internal.ErrorReturn{Err: "Invalid username", Logout: false}
	}

	// ensure valid username
	if !validateUsername(sharee_username) {
		return internal.ErrorReturn{Err: "Invalid sharee username", Logout: false}
	}

	if sharer_username == sharee_username {
		return internal.ErrorReturn{Err: "Cannot share file with yourself", Logout: false}
	}

	// get the current directory
	curr_dir, err := getCurrDir(sharer_username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	// cannot unshare file that another user shared with you
	if strings.HasPrefix(curr_dir, "/shared/") {
		return internal.ErrorReturn{Err: "Unable to share file that has been shared with you", Logout: false}
	}

	// check that file exists in the current directory in the files table - if so, user has permission to unshare it
	count_files, err := countNumFilesByUsernameAndFileInfo(sharer_username, filename, curr_dir)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}
	if count_files != 1 {
		return internal.ErrorReturn{Err: "File does not exist", Logout: false}
	}

	// get the checksum of the file
	checksum, err := getChecksum(sharer_username, filename, curr_dir)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	// check that the sharee exists
	num_sharees, err := getNumUsersWithUsername(sharee_username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}
	if num_sharees != 1 {
		return internal.ErrorReturn{Err: "Sharee does not exist", Logout: false}
	}

	// Remove the shared file, which requires ensuring it exists
	err = unshareFileFromSharee(sharer_username, filename, curr_dir, checksum, sharee_username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	// Get the number of remaining files in that directory
	// if 0, delete the directory
	num_files, err := getNumSharedFilesBySharerAndShareeUsername(sharer_username, sharee_username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}
	if num_files == 0 {
		err = removeSharerDirectoryForSharee(sharer_username, sharee_username)
		if err != nil {
			return internal.ErrorReturn{Err: err.Error(), Logout: false}
		}
	}

	return internal.ErrorReturn{Err: "", Logout: false}
}

/**
 * Delete the account of the given user.
 *
 * Input:
 *   username string - the username given by the client
 *   session_token []byte - the session token from the client
 * Output: An error return type indicating if there was an error.
 *
 * Precondition: The username and session token must be valid (checked at start). The user must have no owned
 *   files or directories.
 * Postcondition: If successful, the user account is deleted from the database. Any files previously shared
 *   with the user are no longer shared with the user. The user's directories of / and /shared/ are deleted from
 *   the directories table. If not successful, nothing changes.
 */
func deleteacctHandler(username string, session_token []byte) internal.ErrorReturn {
	err := checkSessionToken(username, session_token)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: true}
	}

	// ensure valid username
	if !validateUsername(username) {
		return internal.ErrorReturn{Err: "Invalid username", Logout: false}
	}

	// get the current directory
	curr_dir, err := getCurrDir(username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}
	if curr_dir != "/" {
		return internal.ErrorReturn{Err: "Can only delete account form home (/) directory", Logout: false}
	}

	// ensure no files in folder being removed
	files, err := getFilesByUsernameAndCurrDir(username, curr_dir)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}
	if len(files) != 0 {
		return internal.ErrorReturn{Err: "To be deleted, account must have no owned files", Logout: false}
	}

	// ensure no subdirectories in directory being removed
	subdirectories, err := getSubdirectoriesByUsernameAndCurrDir(username, curr_dir)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}
	if len(subdirectories) == 0 {
		return internal.ErrorReturn{Err: "Unable to delete account", Logout: false}
	}
	if len(subdirectories) != 1 {
		return internal.ErrorReturn{Err: "To be deleted, account must have no owned directories other than /shared/", Logout: false}
	}
	if subdirectories[0] != "shared/" {
		return internal.ErrorReturn{Err: "Unable to delete account, account has directory other than shared/", Logout: false}
	}

	// Delete all files shared with user by username
	err = deleteAllFilesSharedWithUserByUsername(username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	// Delete all directories, should only be root and everything in shared but this works to be safe
	err = removeAllDirectories(username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	// Delete the user's account
	err = removeUserRowByUsername(username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	return internal.ErrorReturn{Err: "", Logout: true} // This should be true so user is logged out after deleting account

}

/**
 * Get a list of all sharees for the given file.
 *
 * Input:
 *   username string - the username given by the client
 *   session_token []byte - the session token from the client
 *   filename string - the name of the file whose sharees are to be retrieved
 * Output: A sharees return type with the list of sharees if no error, else indicating if there was an error.
 *
 * Precondition: The username and session token must be valid (checked at start).
 * Postcondition: N/A
 */
func listshareesHandler(username string, session_token []byte, filename string) internal.ShareesReturn {
	err := checkSessionToken(username, session_token)
	if err != nil {
		return internal.ShareesReturn{Sharees: []string{}, Err: err.Error(), Logout: true}
	}

	// ensure valid username
	if !validateUsername(username) {
		return internal.ShareesReturn{Sharees: []string{}, Err: "Invalid username", Logout: false}
	}

	// ensure valid filename
	if !validateFilename(filename) {
		return internal.ShareesReturn{Sharees: []string{}, Err: "Invalid filename", Logout: false}
	}

	// get the current directory
	curr_dir, err := getCurrDir(username)
	if err != nil {
		return internal.ShareesReturn{Sharees: []string{}, Err: err.Error(), Logout: false}
	}

	sharees, err := getShareesByOwnerAndFileInfo(username, filename, curr_dir)
	if err != nil {
		return internal.ShareesReturn{Sharees: []string{}, Err: err.Error(), Logout: false}
	}

	return internal.ShareesReturn{Sharees: sharees, Err: "", Logout: false}
}

/**
 * Change the permissions on the file stored at filename with the given sharee.
 * Write_perms determines the permissions on the file.
 *
 * Input:
 *   username string - the username given by the client
 *   session_token []byte - the session token from the client
 *   filename string - the filepath/filename of the file to be shared
 *   sharee_username string - the username of the person with whom the file is being shared
 *   write_perms bool - true if read/write, false if readonly
 * Output: An error return type indicating if there was an error.
 *
 * Precondition: The username and session token must be valid (checked at start). Must be a valid filename and
 *   an actual user for sharee. The file must have been previously shared with the sharee.
 * Postcondition: On success, the file at filename must have existed. The permissions will be updated for the user
 *   in the shared table, indicating if read/write or readonly. If permissions are changed to be the same as
 *   before, nothing actually changes. Nothing changes about the location of the shared file for the sharee.
 *   If not successful, nothing changes.
 */
func chmodfileHandler(sharer_username string, session_token []byte, filename string, sharee_username string, write_perms bool) internal.ErrorReturn {
	// write_perms should be true if read and write, false if read only
	err := checkSessionToken(sharer_username, session_token)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: true}
	}

	// ensure valid filename
	if !validateFilename(filename) {
		return internal.ErrorReturn{Err: "Invalid filename", Logout: false}
	}

	// ensure valid sharer_username
	if !validateUsername(sharer_username) {
		return internal.ErrorReturn{Err: "Invalid username", Logout: false}
	}

	// ensure valid username
	if !validateUsername(sharee_username) {
		return internal.ErrorReturn{Err: "Invalid sharee username", Logout: false}
	}

	if sharer_username == sharee_username {
		return internal.ErrorReturn{Err: "Cannot chmod file with yourself", Logout: false}
	}

	// get the current directory
	curr_dir, err := getCurrDir(sharer_username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	// cannot share file that has been shared with you
	// This check may be unnecessary
	if strings.HasPrefix(curr_dir, "/shared/") {
		return internal.ErrorReturn{Err: "Unable to chmod file that has been shared with you", Logout: false}
	}

	// check that file exists in the current directory in the files table - if so, you have permission to share it
	count_files, err := countNumFilesByUsernameAndFileInfo(sharer_username, filename, curr_dir)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}
	if count_files != 1 {
		return internal.ErrorReturn{Err: "File does not exist", Logout: false}
	}

	// get the checksum of the file
	checksum, err := getChecksum(sharer_username, filename, curr_dir)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	// check that the sharee exists
	num_sharees, err := getNumUsersWithUsername(sharee_username)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}
	if num_sharees != 1 {
		return internal.ErrorReturn{Err: "Sharee does not exist", Logout: false}
	}

	// try to update write perms
	err = chmodFileWithSharee(sharer_username, filename, curr_dir, checksum, sharee_username, write_perms)
	if err != nil {
		return internal.ErrorReturn{Err: err.Error(), Logout: false}
	}

	return internal.ErrorReturn{Err: "", Logout: false}
}
