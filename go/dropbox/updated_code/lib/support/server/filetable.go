package server

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

/**
 * Get the number of files with the given checksum.
 *
 * Input:
 *   checksum string - the checksum to be used in the query
 * Output:
 *   int - the number of files with the checksum
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client must be valid. The checksum must have been generated via the MakeChecksum function.
 * Postcondition: N/A
 */

func countNumFileByChecksum(checksum string) (int, error) {
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return -1, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT COUNT(*) FROM files 
		WHERE checksum = ?`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(checksum).Scan(&count) // it will only return 1 row cause unique check
	if err != nil {
		return -1, err
	}

	return count, nil
}

/**
 * The number of files with the given username, filename, and filepath.
 *
 * Input:
 *   username string - the username to look up
 *   filename string - the filename to look up
 *   filepath string - the filepath to look up
 * Output:
 *   int - the number of files with the username, filename, and filepath
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client must be valid. The username, filename, and filepath must be valid.
 * Postcondition: N/A
 */
func countNumFilesByUsernameAndFileInfo(username string, filename string, filepath string) (int, error) {
	// validate only 1 user with this username
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return -1, err
	}

	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return -1, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT COUNT(*) FROM files
		WHERE username = ? AND filename = ? AND filepath = ?`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(username, filename, filepath).Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

/**
 * Insert the file with the given info into the files table.
 *
 * Input:
 *   username string - the username of the uploader of the file
 *   filename string - the name of the file
 *   filepath string - the path to the file
 *   checksum string - the sha512 checksum of the file
 * Output: The error, if an error occurred, else nil
 *
 * Precondition: The client, username, filename, and filepath are all valid. Checksum belongs to the uploaded file.
 *   The file is not already uploaded by this user in the same location.
 * Postcondition: The database will have info on the file.
 */
func uploadFile(username string, filename string, filepath string, checksum string) error {
	// validate only 1 user with this username
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

	stmt, err := db.Prepare(`INSERT INTO files
		(username, filename, filepath, checksum)
		VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filename, filepath, checksum)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Get the checksum of the file with the given info
 *
 * Input:
 *   username string - the username of the uploader of the file
 *   filename string - the name of the file
 *   filepath string - the path to the file
 * Output:
 *   checksum string - the checksum of the file
 *   err error - The error, if an error occurred, else nil
 *
 * Precondition: Client, username, filename, and filepath are all valid. Only one row exists with this combination
 *   of username, filename, and filepath.
 * Postcondition: N/A
 */
func getChecksum(username string, filename string, filepath string) (checksum string, err error) {
	// validate only 1 user with this username
	err = validateUniqueUserByUsername(username)
	if err != nil {
		return "", err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return "", err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT checksum FROM files
		WHERE username = ? AND filename = ? AND filepath = ?`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username, filename, filepath).Scan(&checksum)
	if err != nil {
		return "", err
	}

	return checksum, nil
}

/**
 * Get all unique checksums from the database. For reseting the server.
 *
 * Input: N/A
 * Output:
 *   []string - the string of all checksums in the files table
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: Files table exists.
 * Postcondition: N/A
 */
func getAllChecksums() ([]string, error) {
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return []string{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT DISTINCT checksum FROM files`)
	if err != nil {
		return []string{}, err
	}
	defer db.Close()

	rows, err := stmt.Query()
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	var checksums []string
	for rows.Next() {
		var checksum string
		err = rows.Scan(&checksum)
		if err != nil {
			return []string{}, err
		}
		checksums = append(checksums, checksum)
	}

	err = rows.Err()
	if err != nil {
		return []string{}, err
	}

	return checksums, nil
}

/**
 * Get all the files in the user's current directory. For ls.
 *
 * Input:
 *   username string - the username of the current user
 *   curr_dir string - the current directory for the user
 * Output:
 *   []string - the list of files in the current directory for the given user
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: Client is valid. The user's current directory in the user's table is curr_dir.
 * Postcondition: N/A
 */
func getFilesByUsernameAndCurrDir(username string, curr_dir string) ([]string, error) {
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

	stmt, err := db.Prepare(`SELECT filename FROM files
		WHERE username = ? AND filepath = ?
		ORDER BY filename ASC`)
	if err != nil {
		return []string{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, curr_dir)
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

/**
 * Remove the selected file for the selected user.
 *
 * Input:
 *   username string - the username of the user deleting the file
 *   filename string - the name of the file being removed
 *   curr_dir string - the directory containing the file to be removed
 * Output:
 *   string - the checksum if no one else has the file, else "". Indicates if file should be removed from
 *     deduplicated area
 *   error - the error, if an error occurred, else nil
 *
 * Precondition: The client, username, filename, and curr dir are valid. The user owns the selected file and
 *   can safely delete it.
 * Postcondition: If successful, the file was shared with no one and was deleted. The callee is responsible for
 *   handling deduplication. If not successful, nothing changes.
 */
func removeFileByUsername(username string, filename string, curr_dir string) (string, error) {
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return "", err
	}

	// get number of files, should only exist once
	count_files, err := countNumFilesByUsernameAndFileInfo(username, filename, curr_dir)
	if err != nil {
		return "", err
	}
	if count_files != 1 {
		// either 0, in which case does not exist, or 2, which should not happen
		return "", errors.New("File does not exist")
	}

	// get the checksum, used later to see if anyone else has the file
	checksum, err := getChecksum(username, filename, curr_dir)
	if err != nil {
		return "", err
	}

	// count number of users with which the given user has shared the file
	num_sharees, err := countNumShareesByOwnerAndFileInfo(username, filename, curr_dir)
	if err != nil {
		return "", err
	}
	if num_sharees != 0 {
		return "", errors.New("You must unshare a file with all users before deleting the file.")
	}

	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return "", err
	}
	defer db.Close()

	stmt, err := db.Prepare(`DELETE FROM files
		WHERE username = ? AND filename = ? AND filepath = ?`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filename, curr_dir)
	if err != nil {
		return "", err
	}

	// Check if that was only existence of files
	checksum_count, err := countNumFileByChecksum(checksum)
	if err != nil {
		return "", err
	}
	if checksum_count == 0 {
		// no one else has file, so remove it
		return checksum, nil
	}

	return "", nil
}

/**
 * Update the checksum of the file for owner. Called by modify function.
 *
 * Input:
 *   username string - the username of the user deleting the file
 *   filename string - the name of the file being removed
 *   curr_dir string - the directory containing the file to be removed
 *   new_checksum string - the new checksum for the file
 * Output: The error, if an error occurred, else nil
 *
 * Precondition: The client, username, filename, and filepath are valid. The checksum belongs to the new file.
 *   The user is the owner of the file.
 * Postcondition: If successful, the file checksum is updated to the new one. If not successful, nothing changes.
 */
func updateOwnerChecksum(username string, filename string, filepath string, new_checksum string) error {
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return err
	}

	count_files, err := countNumFilesByUsernameAndFileInfo(username, filename, filepath)
	if err != nil {
		return err
	}
	if count_files != 1 {
		// either 0, in which case does not exist, or 2, which should not happen
		return errors.New("File does not exist")
	}

	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	old_checksum, err := getChecksum(username, filename, filepath)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`UPDATE files SET checksum = ?
		WHERE username = ? AND filename = ? AND filepath = ? AND checksum = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(new_checksum, username, filename, filepath, old_checksum)
	if err != nil {
		return err
	}

	return nil

}
