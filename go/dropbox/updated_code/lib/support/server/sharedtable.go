package server

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

/**
 * Shares a file with a sharee by putting the shared file in the shared table with a given sharer username,
 * a filename, the sharer's filepath, the checksum of the file and the sharee's username
 *
 * Input:
 *  sharer_username string - the sharer's username
 *  filename string - the file name
 *  filepath string - the file path
 *  checksum string - the files checksum
 *  sharee_username string - the sharee's username
 *  write_perms bool - whether the sharee has write permissions or not
 * Output: if both usernames validate, the two usernames are not the same, and the the file has
 * not been shared with this sharee by this sharer, return nil, else error
 *
 * Precondition: both usernames and corresponding session tokens have been validated and sanitized
 * as well as the filename and file path
 * Postcondition: The shared file is added to the shared table with the above information
 */

func shareFileWithSharee(sharer_username string, filename string, filepath string, checksum string, sharee_username string, write_perms bool) error {
	err := validateUniqueUserByUsername(sharer_username)
	if err != nil {
		return err
	}

	err = validateUniqueUserByUsername(sharee_username)
	if err != nil {
		return err
	}

	if sharer_username == sharee_username {
		return errors.New("Cannot share file with yourself.")
	}

	// get count in shared table of files with the same info
	count_files, err := getNumSharedFilesByInfo(sharer_username, filename, filepath, sharee_username)
	if err != nil {
		return err
	}
	if count_files != 0 {
		return errors.New("You have already shared a file with the same name with the given user")
	}

	// it is safe to give the user the file, so update the shared table to have a new row
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`INSERT INTO shared
		(owner, filename, filepath, checksum, sharee, write_perms)
		VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sharer_username, filename, filepath, checksum, sharee_username, write_perms)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Changes permissions a sharee has on a file given the sharer of the file
 *
 * Input:
 *  sharer_username string - the sharer's username
 *  filename string - the file name
 *  filepath string - the file path
 *  checksum string - the files checksum
 *  sharee_username string - the sharee's username
 *  write_perms bool - whether the sharee has write permissions or not
 * Output: if both usernames validate, the two usernames are not the same, and the the file is
 * shared with this sharee by this sharer, return nil, else error
 *
 * Precondition: both usernames and corresponding session tokens have been validated and sanitized
 * as well as the filename and file path
 * Postcondition: The write permissions are changed for this file for this sharee
 */

func chmodFileWithSharee(sharer_username string, filename string, filepath string, checksum string, sharee_username string, write_perms bool) error {
	err := validateUniqueUserByUsername(sharer_username)
	if err != nil {
		return err
	}

	err = validateUniqueUserByUsername(sharee_username)
	if err != nil {
		return err
	}

	if sharer_username == sharee_username {
		return errors.New("Cannot share file with yourself.")
	}

	// get count in shared table of files with the same info
	count_files, err := getNumSharedFilesByInfo(sharer_username, filename, filepath, sharee_username)
	if err != nil {
		return err
	}
	if count_files != 1 {
		return errors.New("You have not shared this file with the given user.")
	}

	// it is safe to give the user the file, so update the shared table to have a new row
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`UPDATE shared
		SET write_perms = ? 
		WHERE owner = ? AND filename = ? AND filepath = ? AND checksum = ? AND sharee = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(write_perms, sharer_username, filename, filepath, checksum, sharee_username)
	if err != nil {
		return err
	}
	return nil

}

/**
 * Unshares the file from a sharee given a sharer
 *
 * Input:
 *  sharer_username string - the sharer's username
 *  filename string - the file name
 *  filepath string - the file path
 *  checksum string - the files checksum
 *  sharee_username string - the sharee's username
 * Output: if both usernames validate, the two usernames are not the same, and the the file is
 * shared with this sharee by this sharer, return nil, else error
 *
 * Precondition: both usernames and corresponding session tokens have been validated and sanitized
 * as well as the filename and file path
 * Postcondition: The file is removed from the shared table corresponding to the above inputs
 */

func unshareFileFromSharee(sharer_username string, filename string, filepath string, checksum string, sharee_username string) error {
	err := validateUniqueUserByUsername(sharer_username)
	if err != nil {
		return err
	}

	err = validateUniqueUserByUsername(sharee_username)
	if err != nil {
		return err
	}

	if sharer_username == sharee_username {
		return errors.New("Cannot share file with yourself.")
	}

	count_files, err := getNumSharedFilesByInfo(sharer_username, filename, filepath, sharee_username)
	if err != nil {
		return err
	}
	if count_files != 1 {
		return errors.New("You have not shared the selected file with the given user")
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`DELETE FROM shared
		WHERE owner = ? AND filename = ? AND filepath = ? AND checksum = ? AND sharee = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sharer_username, filename, filepath, checksum, sharee_username)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Get the number of files shared by the sharer to the sharee.
 *
 * Input:
 *   sharer_username string - the username of the sharer
 *   sharee_username string - the username of the sharee
 * Output:
 *   int - the number of files shared from sharer to sharee
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: Both sharer and sharee are valid users
 * Postcondition: N/A
 */
func getNumSharedFilesBySharerAndShareeUsername(sharer_username string, sharee_username string) (int, error) {
	err := validateUniqueUserByUsername(sharer_username)
	if err != nil {
		return -1, err
	}

	err = validateUniqueUserByUsername(sharee_username)
	if err != nil {
		return -1, err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return -1, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT COUNT(*) FROM shared
		WHERE owner = ? AND sharee = ?`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var num_files int
	err = stmt.QueryRow(sharer_username, sharee_username).Scan(&num_files)
	if err != nil {
		return -1, err
	}

	return num_files, err
}

/**
 * Get the number of files with the given file info shared by the owner to the sharee.
 *
 * Input:
 *   owner string - the username of the sharer
 *   filename string - the name of the file
 *   filepath string - the path to the file
 *   sharee_username string - the username of the sharee
 * Output:
 *   int - the number of files shared from owner to sharee with the given filename and filepath
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: Both sharer and sharee are valid users, filename and filepath are valid
 * Postcondition: N/A
 */
func getNumSharedFilesByInfo(owner string, filename string, filepath string, sharee string) (int, error) {
	err := validateUniqueUserByUsername(owner)
	if err != nil {
		return -1, err
	}

	err = validateUniqueUserByUsername(sharee)
	if err != nil {
		return -1, err
	}

	if owner == sharee {
		return -1, errors.New("Cannot share file with yourself.")
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return -1, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT COUNT(*) FROM shared
		WHERE owner = ? AND filename = ? AND filepath = ? AND sharee = ?`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var num_files int
	err = stmt.QueryRow(owner, filename, filepath, sharee).Scan(&num_files)
	if err != nil {
		return -1, err
	}
	return num_files, nil
}

/**
 * Update the checksums of all sharees based off owner and file info.
 *
 * Input:
 *   owner string - the username of the owner
 *   owner_path string - the file path in the owner's file structure
 *   filename string - the name of the file
 *   old_checksum string - the original checksum
 *   new_checksum string - the new checksum
 * Output:
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: Both sharer and sharee are valid users, path and filename are all valid
 * Postcondition: N/A
 */
func updateSharedChecksumsByOwner(owner string, owner_path string, filename string, old_checksum string, new_checksum string) error {
	err := validateUniqueUserByUsername(owner)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`UPDATE shared
		SET checksum = ? WHERE owner = ? AND filename = ? AND filepath = ? AND checksum = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(new_checksum, owner, filename, owner_path, old_checksum)
	if err != nil {
		return err
	}
	return nil
}

//update all sharee's (including itself) and the owner's checksum
func updateSharedChecksumsBySharee(owner string, owner_path string, filename string, old_checksum string, new_checksum string) error {
	err := validateUniqueUserByUsername(owner)
	if err != nil {
		return err
	}

	fmt.Println("OWNER:", owner, "OLDCHECKSUM:", old_checksum, "NEWCHECKSUM:", new_checksum)

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`UPDATE shared
		SET checksum = ? WHERE owner = ? AND filename = ? AND checksum = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(new_checksum, owner, filename, old_checksum)
	if err != nil {
		return err
	}

	stmt, err = db.Prepare(`UPDATE files
		SET checksum = ? WHERE username = ? AND filename = ? AND filepath = ? AND checksum = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(new_checksum, owner, filename, owner_path, old_checksum)
	if err != nil {
		return err
	}
	return nil
}

/**
 * Determine if sharee has write permissions for a file.
 *
 * Input:
 *   sharee string - the username of the sharee
 *   filename string - the name of the file
 *   filepath string - the path to the file
 *   checksum string - the checksum of the file
 * Output:
 *   has_write bool - true if has write permission, false otherwise
 *   err error - the error, if an error occurred, else nil
 *
 * Precondition: all arguments are valid, sharee is a valid user
 * Postcondition: N/A
 */
func hasWritePermissions(sharee string, filename string, filepath string, checksum string) (has_write bool, err error) {
	err = validateUniqueUserByUsername(sharee)
	if err != nil {
		return false, err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return false, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT write_perms FROM shared
		WHERE filename = ? AND filepath = ? AND sharee = ? AND checksum = ?`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(filename, filepath, sharee, checksum).Scan(&has_write)
	if err != nil {
		return false, err
	}

	return has_write, nil
}

/**
 * Get the list of all files shared by owner to username.
 *
 * Input:
 *   username string - the username of the sharee
 *   owner string - the username of the owner
 * Output:
 *   []string - the name of all files shared from owner to username
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: Both sharer and sharee are valid users
 * Postcondition: N/A
 */
func getSharedFilesByUsernameAndOwner(username string, owner string) ([]string, error) {
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return []string{}, err
	}

	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return []string{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT filename FROM shared
		WHERE sharee = ? AND owner = ?
		ORDER BY filename ASC`)
	if err != nil {
		return []string{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, owner)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	var files []string
	for rows.Next() {
		var filename string
		err = rows.Scan(&filename)
		if err != nil {
			return []string{}, err
		}
		files = append(files, filename)
	}

	err = rows.Err()
	if err != nil {
		return []string{}, err
	}

	return files, nil
}

// get the path of the owner for the shared user's file, for modify
func getSharedOwnerPath(sharee string, owner string, filename string) (owner_path string, err error) {
	err = validateUniqueUserByUsername(sharee)
	if err != nil {
		return "", err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return "", err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT filepath FROM shared
		WHERE filename = ? AND owner = ? AND sharee = ?`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	err = stmt.QueryRow(filename, owner, sharee).Scan(&owner_path)
	if err != nil {
		return "", err
	}

	return owner_path, nil
}

// delete all files shared with the given username, for delete account
func deleteAllFilesSharedWithUserByUsername(username string) error {
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`DELETE FROM shared
		WHERE sharee = ?`)
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

// get the number of sharees of owner for filename and filepath
func countNumShareesByOwnerAndFileInfo(owner string, filename string, filepath string) (int, error) {
	err := validateUniqueUserByUsername(owner)
	if err != nil {
		return -1, err
	}

	// Validate file exists
	count_files, err := countNumFilesByUsernameAndFileInfo(owner, filename, filepath)
	if err != nil {
		return -1, err
	}
	if count_files != 1 {
		return -1, errors.New("File does not exist")
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return -1, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT COUNT(*) FROM shared
		WHERE owner = ? AND filename = ? AND filepath = ?`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var count_sharees int
	err = stmt.QueryRow(owner, filename, filepath).Scan(&count_sharees) // it will only return 1 row cause unique check
	if err != nil {
		return -1, err
	}

	return count_sharees, nil
}

/**
 * Gets all sharees given an owner and a file name and path
 *
 * Input:
 *  owner string - the owner's username
 *  filename string - the file name
 *  filepath string - the file path
 * Output: if the owners username validates returns a string array of all sharees usernames for that file, error if db fails
 *
 * Precondition: owners username and corresponding session token have been validated and sanitized
 * as well as the filename and file path
 * Postcondition: N/A
 */
func getShareesByOwnerAndFileInfo(owner string, filename string, filepath string) ([]string, error) {
	err := validateUniqueUserByUsername(owner)
	if err != nil {
		return []string{}, err
	}

	// Validate file exists
	count_files, err := countNumFilesByUsernameAndFileInfo(owner, filename, filepath)
	if err != nil {
		return []string{}, err
	}
	if count_files != 1 {
		return []string{}, errors.New("File does not exist")
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return []string{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT sharee FROM shared
		WHERE owner = ? AND filename = ? AND filepath = ?`)
	if err != nil {
		return []string{}, err
	}
	defer db.Close()

	rows, err := stmt.Query(owner, filename, filepath)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	// read the name of all sharees
	var sharees []string
	for rows.Next() {
		var sharee string
		err = rows.Scan(&sharee)
		if err != nil {
			return []string{}, err
		}
		sharees = append(sharees, sharee)
	}

	err = rows.Err()
	if err != nil {
		return []string{}, err
	}

	return sharees, nil
}
