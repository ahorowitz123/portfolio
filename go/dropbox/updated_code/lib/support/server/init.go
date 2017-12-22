package server

import (
	"fmt"
	//"io/ioutil
	// "golang.org/x/crypto/bcrypt"
	"math/rand"
	// "strings"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
	// "github.com/s17-ahorowi2-jtc2/updated_code/internal"
	"github.com/s17-ahorowi2-jtc2/updated_code/lib/support/rpc"
)

// the base directory for the server. Set on server startup
var basedir string

/**
 * Start the server. Ran when the server starts.
 *
 * Input: N/A
 * Output: N/A
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func InitServer() {
	// from http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
	// ensure password verification codes do not repeat when server restarted
	rand.Seed(time.Now().UnixNano())

	//make sure database tables are made
	err := createTables()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	// parse command line arguments
	var reset bool
	var listenAddr string
	switch {
	case len(os.Args) == 3 && os.Args[2] == "--reset":
		basedir = os.Args[1]
		reset = true
	case len(os.Args) == 3 && (len(os.Args[1]) == 0 || os.Args[1][0] != '-'):
		basedir = os.Args[1]
		listenAddr = os.Args[2]
	default:
		fmt.Fprintf(os.Stderr, "Usage: %v <base-dir> [--reset | <listen-address>]\n", os.Args[0])
		os.Exit(1)
	}

	initBaseDir(basedir)

	// if reset, then delete and recreate the database
	if reset {
		resetServer()
		return
	}

	// pre login handlers
	rpc.RegisterHandler("emailregister", emailRegisterHandler)
	rpc.RegisterHandler("validateregister", validateRegisterHandler)
	rpc.RegisterHandler("validatelogin", validateLoginHandler)

	// post login handlers
	rpc.RegisterHandler("getworkingdirectory", getWorkingDirectoryHandler)
	rpc.RegisterHandler("ls", lsHandler)
	rpc.RegisterHandler("cd", cdHandler)
	rpc.RegisterHandler("mkdir", mkdirHandler)
	rpc.RegisterHandler("rmdir", rmdirHandler)
	rpc.RegisterHandler("rmfile", rmfileHandler)
	rpc.RegisterHandler("logout", logoutHandler)
	rpc.RegisterHandler("upload", uploadHandler)
	rpc.RegisterHandler("cat", catHandler)
	rpc.RegisterHandler("sharefile", sharefileHandler)
	rpc.RegisterHandler("unsharefile", unsharefileHandler)
	rpc.RegisterHandler("listsharees", listshareesHandler)
	rpc.RegisterHandler("deleteacct", deleteacctHandler)
	rpc.RegisterHandler("chmodfile", chmodfileHandler)

	// register finalize
	rpc.RegisterFinalizer(finalizer)

	// run the server
	err = rpc.RunServer(listenAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not run server: %v\n", err)
		os.Exit(1)
	}

}

/**
 * Set the basedir global variable. Made by my partner. I have no clue why this is needed.
 * He doesn't know either.
 *
 * Input:
 *   base string - the basedir for the server, from the command line args
 * Output: N/A
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func initBaseDir(base string) {
	basedir = base
}

/**
 * Finalizer to be ran when server shuts down.
 *
 * Input: N/A
 * Output: N/A
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func finalizer() {
	fmt.Println("Shutting down...")
}
