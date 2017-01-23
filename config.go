
/*
 * Pieter De Ridder
 * 2017 
 *
 * Base64 smtp mail client
 * Can be cross compiled for Windows, Linux, Mac OS X and Solaris on Intel x86-64
 *
 * ------------------------------------------------------------
 * 
 * encoding smtpmail go file
 * extra function for Base64 encoding following RFC2047
 *
 */


package main

import (
	"strings"
	"bufio"      // files
	"os"         // os.Args
	"strconv"
)


// ----- file I/O ------ 
/* read a text file line by line to a slice */
func readLines(path string) ([]string, error) {
	var lines []string
	
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	
	return lines, scanner.Err()
}


/* file exists check */
func fileExists(path string) bool {
	var found bool = false

    _, err := os.Stat(path)
	
    if err == nil { 
		found = true
	}
	
    /*if os.IsNotExist(err) {
		found false
	}*/
	
    return found
}

/* search of the config file */
func fetchConfigFile() string {
	var configFile string
		

	// used defined config file with cmdline parameter "-config <configfile>"
	if (len(configFileOverride) > 0) {
		if (fileExists(configFileOverride)) {
			configFile = configFileOverride
			debugprint("user config " + configFileOverride + " set.")		
		} else {
			debugprint("user config " + configFileOverride + " not found.")
			warning("user defined config not present!", WARNING)
			warning("given file: " + configFileOverride, WARNING)
		}	
	} else {
		// application auto search
		var exeFilename string = stripExeFilenameExt()
		var currentPath string = stripExeFilepath()
		
		switch getOSFlavor() {
			case "win32":	
				var win32SearchPath string = currentPath + exeFilename + ".conf"  
				
				// windows has a static folder patch "<current folder>\<exefilename>.conf"
				// filename : "<current folder>\<exefilename>.conf"
				if (fileExists(win32SearchPath)) {
					configFile = win32SearchPath 
					debugprint("win config " + win32SearchPath + " set.")				
				} else {
					debugprint("win config " + win32SearchPath + " not found.")
					configFile = "" // return empty
				}
			
			case "bsd":
				fallthrough
			case "darwin":
				fallthrough
			case "unix":
			
				var unixSearchPath string = "/etc/" + exeFilename + ".conf"
				
				if (fileExists(unixSearchPath)) {
					// try /etc folder
					// filename :  "/etc/<exefilename>.conf"
					configFile = unixSearchPath 
					debugprint("unix config " + unixSearchPath + " set.")
				} else {
					debugprint("unix config " + unixSearchPath + " not found.")
					
					// fallback to "current" folder
					// filename :  "<current folder>\<exefilename>.conf"
					unixSearchPath = currentPath + exeFilename + ".conf"
				
					if (fileExists(unixSearchPath)) {
						configFile = unixSearchPath 
						debugprint("unix config " + unixSearchPath + " set.")
					} else {
						configFile = "" // return empty
						debugprint("unix config " + unixSearchPath + " not found.")
					}
				}
			default:
				configFile = ""  // return empty
		}

	}
	return configFile
}


/* read and parse config file */
func readConfigFile() {

	if !(skipConfigFileLoad) {
	
	

		var configFile string = fetchConfigFile()
		
		if (len(configFile) > 0) {
			lines, err := readLines(configFile)
			
			debugprint("reading config file")
			
			if err == nil {
				print("config file present")
			
				for l := 0; l < len(lines); l++ {
					if (len(lines[l]) >= 2 ) {	
						line := strings.TrimSpace(lines[l])   // rim the whole line first and store it in a var
								
						if !(strings.HasPrefix(line, "//")) {   //exclude comments starting with a double slash
							codes := strings.Split(line, "=")   // try splitting on a '=' char
							
							opcode := strings.TrimSpace(codes[0])   // the "opcode", or command so to say
							opvalue := strings.TrimSpace(codes[1])  // the value of the opcode
					
							// parse opcodes
							switch (opcode) {
								case "server", "smtpserver":
									smtpServer = opvalue
								case "port", "smtpport":
									smtpPort, _ = strconv.Atoi(opvalue)
								case "authusername", "authuser", "username":
									userName = opvalue
								case "authpassword", "authpwd", "password":
									userPasswd = opvalue
								case "from", "frommail":
									fromMailAddress = opvalue
								case "body":
									mailBody = opvalue
								case "subject", "title":
									mailSubject = opvalue
								case "authlogin", "authmethod", "method":
									if (opvalue == "1") {
										useAuthLogin = true
									} else {
										useAuthLogin = false
									}
								case "verbose":
									if (opvalue == "1") {
										verbose = true
									} else {
										verbose = false
									}
								case "silent":
									if (opvalue == "1") {
										silent = true
									} else {
										silent = false
									}
								/*case "debug":
									if (opvalue == "1") {
										debug = true
									} else {
										debug = false
									}*/
							}
						}
					}
				}
			} else {
				warning("problem while reading the config file!", WARNING)  // throw as warning instead of error, because the application should still work without config file
			}
		} else {
			print("no config file present.")
		}
	} else {
		print("user skipped config file.")
	}
}
