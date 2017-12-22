package server

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

/**
 * Create the tables in the database, if they do not exist. Will make the database used, boxdrop.sqlite, if
 * it does not exist. Called every time the server is ran.
 *
 * Input: N/A
 * Output: the error, if an error occurred
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func createTables() error {
	err := createUsersTable()
	if err != nil {
		return err
	}

	err = createDirectoryTable()
	if err != nil {
		return err
	}

	err = createFilesTable()
	if err != nil {
		return err
	}

	err = createSharedTable()
	if err != nil {
		return err
	}

	return nil
}

/**
 * Create the users table if it does not exist.
 *
 * Input: N/A
 * Output: the error, if an error occurred
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func createUsersTable() error {
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	// TODO: Insert not nulls, primary keys, etc.
	createStmt := `CREATE TABLE IF NOT EXISTS users 
		(username text NOT NULL PRIMARY KEY, 
		 email text NOT NULL, 
		 ver_code blob, 
		 ver_code_creation_time int,
		 password blob NOT NULL, 
		 salt text NOT NULL, 
		 session_token blob,
		 session_token_last_access_time int,
		 curr_dir text)`

	_, err = db.Exec(createStmt) // exec returns nothing
	if err != nil {
		return err
	}

	return nil
}

/**
 * Create the directories table if it does not exist.
 *
 * Input: N/A
 * Output: the error, if an error occurred
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func createDirectoryTable() error {
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	createStmt := `CREATE TABLE IF NOT EXISTS directories 
		(id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		 username text NOT NULL, 
		 path text NOT NULL, 
		 name text NOT NULL)`

	_, err = db.Exec(createStmt) // exec returns nothing
	if err != nil {
		return err
	}

	return nil
}

/**
 * Create the files table if it does not exist.
 *
 * Input: N/A
 * Output: the error, if an error occurred
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func createFilesTable() error {
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	createStmt := `CREATE TABLE IF NOT EXISTS files
		(id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		 username string NOT NULL,
		 filename string NOT NULL,
		 filepath string NOT NULL,
		 checksum string NOT NULL)`

	_, err = db.Exec(createStmt)
	if err != nil {
		return err
	}
	return nil
}

/**
 * Create the shared table if it does not exist.
 *
 * Input: N/A
 * Output: the error, if an error occurred
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func createSharedTable() error {
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	createStmt := `CREATE TABLE IF NOT EXISTS shared
	(id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	 owner string NOT NULL,
	 filename string NOT NULL,
	 filepath string NOT NULL,
	 checksum string NOT NULL,
	 sharee string NOT NULL,
	 write_perms boolean NOT NULL)`
	_, err = db.Exec(createStmt)
	if err != nil {
		return err
	}
	return nil
}
