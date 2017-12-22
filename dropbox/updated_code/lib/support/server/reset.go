package server

import (
	"fmt"
	"os"
)

/**
 * Reset the server to its default state. Called when the --reset option is given to the server.
 *
 * Input: N/A
 * Output: N/A
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func resetServer() {
	// get the checksums of all files stored in the server. Will not contain duplicates
	checksums, dberr := getAllChecksums()
	if dberr != nil {
		fmt.Printf("ERROR: %v\n", dberr)
		return
	}

	// remove every file from the basedir. There should be no files in basedir that are in not in the files table,
	// assuming basedir never changes
	for _, checksum := range checksums {
		if checksum != "" {
			dirpath := basedir + "/"
			file_to_delete := dirpath + checksum
			err := os.Remove(file_to_delete)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				return
			}
		}
	}

	// delete all the tables
	dberr = deleteTables()
	if dberr != nil {
		fmt.Printf("ERROR: %v\n", dberr)
		return
	}

	// create all the tables
	dberr = createTables()
	if dberr != nil {
		// if create tables failed, delete any tables that were created
		dberr2 := deleteTables()
		if dberr2 != nil {
			fmt.Printf("ERROR: %v\n", dberr2)
			return
		}

		fmt.Printf("ERROR: %v\n", dberr)
		return
	}

	fmt.Printf("Reset Completed, exiting...\n")
}
