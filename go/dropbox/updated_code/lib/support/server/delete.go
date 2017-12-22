package server

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

/**
 * Create the tables in the database, if they exist. Called whenever the server is explicitly reset.
 *
 * Input: N/A
 * Output: the error, if an error occurred
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func deleteTables() error {
	err := deleteUsersTable()
	if err != nil {
		return err
	}

	err = deleteDirectoriesTable()
	if err != nil {
		return err
	}

	err = deleteFilesTable()
	if err != nil {
		return err
	}

	err = deleteSharedTable()
	if err != nil {
		return err
	}

	return nil
}

/**
 * Delete the users table if it exists.
 *
 * Input: N/A
 * Output: the error, if an error occurred
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func deleteUsersTable() error {
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	deleteStmt := `DROP TABLE IF EXISTS users`

	_, err = db.Exec(deleteStmt)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Delete the directories table if it exists.
 *
 * Input: N/A
 * Output: the error, if an error occurred
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func deleteDirectoriesTable() error {
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	deleteStmt := `DROP TABLE IF EXISTS directories`

	_, err = db.Exec(deleteStmt)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Delete the files table if it exists.
 *
 * Input: N/A
 * Output: the error, if an error occurred
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func deleteFilesTable() error {
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	deleteStmt := `DROP TABLE IF EXISTS files`

	_, err = db.Exec(deleteStmt)
	if err != nil {
		return err
	}
	return nil
}

/**
 * Delete the shared table if it exists.
 *
 * Input: N/A
 * Output: the error, if an error occurred
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func deleteSharedTable() error {
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	deleteStmt := `DROP TABLE IF EXISTS shared`

	_, err = db.Exec(deleteStmt)
	if err != nil {
		return err
	}
	return nil
}
