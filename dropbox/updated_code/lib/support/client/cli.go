// Author: jliebowf
// Date: Spring 2016

package client

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

/**
 * RunCLI will the run the REPL for the given client. After this function returns, the caller should
 * not call this function again, and the client should exit cleanly.
 *
 * Input:
 *   c client - The client whose REPL is being ran.
 * Output:
 *   An error if an error occurs, else nil.
 *
 * Precondition: The client is already connected to the server.
 * Postcondition: N/A
 */
func RunCLI(c Client) error {
	fmt.Println("Welcome to the Dropbox client.")
	c.Help()

	s := bufio.NewScanner(os.Stdin)

Repl:
	for {
		// if there is a session token, print both working directory and prompt, else just print path
		if len(c.GetSessionToken()) != 0 {
			pwd, err := c.Pwd()
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)

			}
			fmt.Printf("%v> ", pwd)
		} else {
			fmt.Print("> ")
		}

		// handle CTRL-D
		if !s.Scan() {
			// if logged in, first logout
			if len(c.GetSessionToken()) != 0 {
				err := c.Logout()
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
				}
			}
			break Repl
		}

		// split string by whitespace
		parts := strings.Fields(s.Text())
		if len(parts) == 0 {
			continue
		}
		args := parts[1:]

		if len(c.GetSessionToken()) == 0 {
			// before log in options
			switch parts[0] {
			case "help":
				if len(args) != 0 {
					fmt.Printf("Usage: %v\n", parts[0])
					break
				}
				c.Help()
				break
			case "register":
				if len(args) != 0 {
					fmt.Printf("Usage: %v\n", parts[0])
					break
				}
				err := c.Register()
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
				}
				break
			case "login":
				if len(args) != 0 {
					fmt.Printf("Usage: %v\n", parts[0])
					break
				}
				err := c.Login()
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
				} else {
					color.Green("Connected.\n")
				}
				break
			case "exit":
				break Repl // break out of switch and and for loop
			default:
				fmt.Println("Unknown command; try \"help\"")
				break
			}
		} else {
			// after log in options
			switch parts[0] {
			case "help":
				if len(args) != 0 {
					fmt.Printf("Usage: %v\n", parts[0])
					break
				}
				c.Help()
				break
			case "logout":
				if len(args) != 0 {
					fmt.Printf("Usage: %v\n", parts[0])
					break
				}
				err := c.Logout()
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
				} else {
					fmt.Printf("Logged out.\n")
				}
				break
			case "pwd":
				if len(args) != 0 {
					fmt.Printf("Usage: %v\n", parts[0])
					break
				}
				pwd, err := c.Pwd()
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				fmt.Printf("%v\n", pwd)
				break
			case "ls":
				if len(args) != 0 {
					fmt.Printf("Usage: %v\n", parts[0])
					break
				}
				dirs, files, err := c.Ls()
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				// print returned directories
				for _, dirname := range dirs {
					color.Cyan(dirname)
				}
				// print returned files
				for _, filename := range files {
					fmt.Printf("%v\n", filename)
				}
				break
			case "cd":
				if len(args) != 1 {
					fmt.Printf("Usage: %v <directory>\n", parts[0])
					break
				}
				err := c.CD(args[0])
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				break
			case "mkdir":
				if len(args) != 1 {
					fmt.Printf("Usage: %v <directory>\n", parts[0])
					break
				}
				err := c.Mkdir(args[0])
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				break
			case "rmdir":
				if len(args) != 1 {
					fmt.Printf("Usage: %v <directory>\n", parts[0])
					break
				}
				err := c.Rmdir(args[0])
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				break
			case "rm":
				if len(args) != 1 {
					fmt.Printf("Usage: %v <filename>\n", parts[0])
					break
				}
				err := c.Rm(args[0])
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				break
			case "upload":
				if len(args) != 2 {
					fmt.Printf("Usage: %v <src_filepath> <dest_filename>\n", parts[0])
					break
				}
				err := c.Upload(args[0], args[1], false)
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				break
			case "download":
				if len(args) != 2 {
					fmt.Printf("Usage: %v <src_filename> <dest_filepath>\n", parts[0])
					break
				}
				err := c.Download(args[0], args[1])
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				break
			case "cat":
				if len(args) != 1 {
					fmt.Printf("Usage: %v <file>\n", parts[0])
					break
				}
				file_contents, err := c.Cat(args[0])
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				fmt.Printf(file_contents)
				break
			case "share_r":
				if len(args) != 2 {
					fmt.Printf("Usage: %v <file> <sharee_username>\n", parts[0])
					break
				}
				err := c.Share(args[0], args[1], false)
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				break
			case "share_rw":
				if len(args) != 2 {
					fmt.Printf("Usage: %v <file> <sharee_username>\n", parts[0])
					break
				}
				err := c.Share(args[0], args[1], true)
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				break
			case "unshare":
				if len(args) != 2 {
					fmt.Printf("Usage: %v <file> <sharee_username>\n", parts[0])
					break
				}
				err := c.Unshare(args[0], args[1])
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				break
			case "modify":
				if len(args) != 2 {
					fmt.Printf("Usage: %v <src_filepath> <dest_filename>\n", parts[0])
					break
				}
				err := c.Upload(args[0], args[1], true)
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				break
			case "delete_acct":
				if len(args) != 0 {
					fmt.Printf("Usage: %v\n", parts[0])
					break
				}
				err := c.DeleteAcct()
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				color.Red("Account deleted.")
				break
			case "ls_sharees":
				if len(args) != 1 {
					fmt.Printf("Usage: %v <file>\n", parts[0])
					break
				}
				sharees, err := c.LsSharees(args[0])
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				for _, sharee := range sharees {
					fmt.Printf("%v\n", sharee)
				}
				break
			case "chmod_r":
				if len(args) != 2 {
					fmt.Printf("Usage: %v <file> <username>\n", parts[0])
					break
				}
				err := c.Chmod(args[0], args[1], false)
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				break
			case "chmod_rw":
				if len(args) != 2 {
					fmt.Printf("Usage: %v <file> <username>\n", parts[0])
					break
				}
				err := c.Chmod(args[0], args[1], true)
				if err != nil {
					fmt.Printf("ERROR: %v\n", err)
					break
				}
				break
			default:
				fmt.Println("Unknown command; try \"help\"")
				break
			}
		}
	}

	// handle errors in scanning
	if err := s.Err(); err != nil {
		fmt.Printf("error scanning stdin: %v\n", err)
		return err
	}
	return nil
}
