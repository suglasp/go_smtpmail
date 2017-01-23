/*
 * Pieter De Ridder
 * 2017 
 *
 * Base64 smtp mail client
 * Can be cross compiled for Windows, Linux, Mac OS X and Solaris on Intel x86-64
 *
 * ------------------------------------------------------------
 * 
 * global smtpmail go file
 * all "central" program data and structures are in here.
 *
 */
 
package main

import (
    "fmt"
	"os"              // os.Args
	"strings"
	"strconv"
	"runtime"
	"bufio"
)


// ----- constants ------
const (
	WARNING = iota  
	ERROR 
	)

	
const MAIL_ADDRESS_MINIMAL_LENGTH = 6
const PROGRAM_VERSION             = "1.2"




// ----- global vars ------ 
var (
	preinitPhase bool = true     // tells the parseArgs function to only lookup debug output --> processStdInPiped func

	useAuthLogin bool = false    // user AUTH LOGIN or AUTH PLAIN method
	verbose bool = false         // verbose output
	silent bool = false          // fully silent (no output at all)
	debug bool = false           // debugging
	
	userName, userPasswd string  // smtp authentication user + passwd
	 
	smtpServer string            // smtp server
	smtpPort int = 25            // smtp port
	
	fromMailAddress string       // single from mail address
	toMailAddress []string = make([]string, 0)   // multiple to mail adresses
	
	mailSubject, mailBody string  // mail subject and body
	 
	configFileOverride string     // config file user defined
	skipConfigFileLoad bool = false   // user want to skip config load?
)


// ----- Init ------
func InitApp() {
	parseArgs(preinitPhase)
	
	processStdInPiped()  // process piped input from cmdline
	
	parseArgs(preinitPhase)      // parse arguments second. overrides config file parameters
	readConfigFile() // first, read config file
	
	dumpstack()			//debugging purpose
	
    sanityArgs()     // sanity checking on provided user input
}


// ----- verbose output ------ 
/* verbose output to stdout */
func print(msg string) {
    // if debugging is enabled, override verbose flag to so we get "all" output
	if (verbose || debug) {
		fmt.Println( msg )
	}
}

/* debug output to stdout */
func debugprint(dbgmsg string) {
	if (debug) {
		fmt.Println( "DEBUG: " + dbgmsg )
	}
}


/* output warning or error messages */
func warning(msg string, code byte) {
	var opcode string = "warning, "
	
	switch (code) {
		case ERROR:
			opcode = "error, "	 // output: "error, {msg}"
		/* case WARNING:
			opcode = "warning, "   // output: "warning, {msg}" (default value)
		*/
	}
	
	print(opcode + msg)  //uses print() to output it's message.
}


/* exit application */
func endApp(exitCode int) {
	os.Exit(exitCode)
}


/* own routine to inspect variables */
func dumpstack() {
	debugprint("---")
	debugprint("stack dump global variables...")
	debugprint("useAuthLogin: " + strconv.FormatBool(useAuthLogin))
	debugprint("verbose: " + strconv.FormatBool(verbose))
	debugprint("silent: " + strconv.FormatBool(silent))
	debugprint("debug: " + strconv.FormatBool(debug))
	debugprint("userName: " + userName)
	
	// do not output password for security reasons (for the case of dumping to log file)
	//debugprint("userPasswd: " + userPasswd)
	if (len(userPasswd) > 0) {
		debugprint("userPasswd: [set, not shown]")
	} else {
		debugprint("userPasswd: [empty]")
	}
	
	debugprint("smtpServer: " + smtpServer)
	debugprint("smtpPort: " + strconv.Itoa(smtpPort))
	debugprint("mailSubject: \"" + mailSubject + "\"")
	debugprint("mailBody: \"" + mailBody + "\"")
	debugprint("fromMailAddress: " + fromMailAddress)
	
	for i := 0; i < len(toMailAddress); i++ {
		debugprint("toMailAddress: " + toMailAddress[i])
	}
	debugprint("---")
	debugprint("")
}


// ----- the 'current' path and filename ------
/* get current exe filename */
func stripExeFilename() string {
	var exeFilename = os.Args[0]
	
	// strip exe-filename from paths
	if (strings.Contains(os.Args[0], `\`)) {
		winPaths := strings.Split(exeFilename, `\`)
		
		exeFilename = winPaths[len(winPaths) -1]
	}
	
	if (strings.Contains(os.Args[0], `/`)) {
		uxPaths := strings.Split(exeFilename, `/`)
		
		exeFilename = uxPaths[len(uxPaths) -1]
	}
	
	return exeFilename
}

/* extended exe file name strip. Strips also the ".exe" from the filename. */
func stripExeFilenameExt() string {
	var exeFilename = stripExeFilename()
	
	// strip exe-filename from paths
	if (strings.Contains(exeFilename, ".")) {
		filenameOnly := strings.Split(exeFilename, ".")
		
		exeFilename = filenameOnly[0]
	}
	
	return exeFilename
}

/* get current exe filepath */
func stripExeFilepath() string {
	var exeFilepath string = ""
	
	// strip exe-filename from paths
	if (strings.Contains(os.Args[0], `\`)) {
		winPaths := strings.Split(os.Args[0], `\`)
		
		for p := 0; p < (len(winPaths) -1); p++ {
			exeFilepath += winPaths[p] + `\`
		}
		
	}
	
	if (strings.Contains(os.Args[0], `/`)) {
		uxPaths := strings.Split(os.Args[0], `/`)
		
		for p := 0; p < (len(uxPaths) -1); p++ {
			exeFilepath += uxPaths[p] + `/`
		}
	}
	
	return exeFilepath
}


// ----- OS version ---------
/* returns a flavor instead of a typical OS name */
func getOSFlavor() string {
	var os string = "unknown"
	
	switch runtime.GOOS {
	    case "windows":
			os = "win32"
		case "solaris":
			fallthrough
		case "linux":
			os = "unix"
		case "dragonfly":
			fallthrough
		case "openbsd":
			fallthrough
		case "netbsd":
			fallthrough
		case "freebsd":
			os = "bsd"
		case "macosx":
			os = "darwin"
		default:
			os = "unknown"
	}
	
	return os
}



// ----- stdin input from a pipe -----
func processStdInPiped() {
	debugprint("")
	debugprint("parsing pipe stdin")

	in := bufio.NewReader(os.Stdin)  // open reader
	stat, err := os.Stdin.Stat()     // open stats on StdIn
	
	if err == nil {
		// specific, filter out StdIn piped data from cmdline
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			debugprint( "something fishy on the pipe!" )			
			
			in, _, err := in.ReadLine()  //read the line
			if err == nil {
				pipemsg := string(in)   // convert bytes to string
				debugprint("pipe says to me \"" + pipemsg + "\"")
				
				mailBody = pipemsg  // attach the string to our mail body
				
				// filter out max 30 characters for the Subject. Can be overwritten with user args!
				if (len(pipemsg) >= 30) {
					mailSubject = pipemsg[:30]
				} else {
					mailSubject = pipemsg
				}
				
				//debugprint("set body to: " + mailBody)          // commented out, already shown with the dumpstack() func
				//debugprint("set subject to: " + mailSubject)    // commented out, already shown with the dumpstack() func
				debugprint("")
			} 			
		}
	} else {
		debugprint("nothing on the pipe, all quiet.")
	}
	
	preinitPhase = false   // switch off pre-init args phase. switch to lazy-args init phase.
}	
	


// ----- cmdline Arguments and sanity checking ------ 
/* 
 * read and parse cmdline vars, do not use the 'flags' package provided with golang.
 * yes, i do know how to make use of the 'flags' package. But i want more flexibility.
 */
func parseArgs(preinit bool) {
	if ( len(os.Args) > 1 ) {
	
		if preinit {
			// pre-init args. "-debug" is fetched early on.
			for i := 1; i < len(os.Args); i++ {		
				// for easyness sake, use a switch to scan multiple options for debug args
				switch (strings.ToLower(os.Args[i])) {
					case "-debug", "-dump", "-walk":
						debug = true     // enable debug on pre-init
						debugprint("debug output enabled")
				}
			}
		} else {
			// lazy-init args. all other flags are parsed.
			for i := 1; i < len(os.Args); i++ {
				switch (strings.ToLower(os.Args[i])) {
					case "-server", "-s", "-smtpserver":
						if ( (i+1) < len(os.Args) ) {
							smtpServer = strings.TrimSpace(os.Args[i+1])		
							print( "set server to " + smtpServer )   // smtp server
						}
					case "-port", "-p", "-smtpport":
						if ( (i+1) < len(os.Args) ) {
							if smtpPort, err := strconv.Atoi(os.Args[i+1]); err == nil {    
								print("set port to " + strconv.Itoa(smtpPort) )    // smtp port
							}
						}
					case "-user", "-u":
						if ( (i+1) < len(os.Args) ) {
							userName = strings.TrimSpace(os.Args[i+1])
							print( "set auth user to " + userName )   // smtp auth user
						}
					case "-password", "-pwd":
						if ( (i+1) < len(os.Args) ) {
							userPasswd = os.Args[i+1]		    // smtp auth password
							print( "set auth password to ***")
						}		
					case "-authlogin", "-al":
						if (len(os.Args) == 2) {
							// seems like only -al is provided as a argument! so print help again.
							printHelp()
						} else {
							useAuthLogin = true    // AUTH LOGIN method
							print("set auth type to LOGIN")
						}					
					case "-authplain", "-ap":
						if (len(os.Args) == 2) {
							// seems like only -ap is provided as a argument! so print help again.
							printHelp()
						} else {
							useAuthLogin = false     // AUTH PLAIN method
							print("set auth type to PLAIN")
						}
					/*case "-debug", "-dump", "-walk":
						debug = true     // debug
						debugprint("debug output enabled") */
					case "-verbose", "-v":					
						if (len(os.Args) == 2) {
							// seems like only -v is provided as a argument! so print help again.
							printHelp()
						} else {
							verbose = true      // verbose
						}
					case "-from", "-f":
						if ( (i+1) < len(os.Args) ) {
							fromMailAddress = strings.TrimSpace(os.Args[i+1])   // from mail address
							print("set from address to " + fromMailAddress)
						}
					case "-to", "-t":
						if ( (i+1) < len(os.Args) ) {
							toMailAddress = strings.Split(os.Args[i+1], ",")     // to mail addresses
							
							//append(toMailAddress, os.Args[i+1])
							print("set to address")						
							for m := 0; m < len(toMailAddress); m++ {
								toMailAddress[m] = strings.TrimSpace(toMailAddress[m])
								print(toMailAddress[m])
							}						
						}	
					case "-title", "-ti", "-subject", "-sub":
						if ( (i+1) < len(os.Args) ) {
							mailSubject = os.Args[i+1]    // subject
						}
					case "-body", "-b":
						if ( (i+1) < len(os.Args) ) {
							mailBody = os.Args[i+1]       // body
						}
					case "-silent", "-q", "-quiet":
						if (len(os.Args) == 2) {
							// seems like only -silent is provided as a argument! so print help again.
							printHelp()
						} else {
							silent = true
						}
					case "-conf", "-config", "-cfg", "-file", "-c":
						if ( (i+1) < len(os.Args) ) {
							configFileOverride = os.Args[i+1]   // custom config file, user defined
						}
					case "-skipconfig", "-skip", "-nofile", "-nf":
						skipConfigFileLoad = true;
					/*case "-attach", "-a":
						*/
				}
			}
		}
	} else {
		printHelp()
	}
}


/* sanity check of argument parameters */
func sanityArgs() {
	debugprint("sanity args")

	var oldVerbose bool = verbose  // stock and load current verbose state 

	// temporary enable verbose output
	if !(silent) {
		verbose = true
	}

	//smtp host
	if (len(smtpServer) <= 0) {
		warning("smtp server is empty.", ERROR)
		endApp(0)
	}
	
	
	//sanity check To Mail adresses
	if (toMailAddress != nil) {
		if (len(toMailAddress) > 0) {
			for m := 0; m < len(toMailAddress); m++ {						
				if (!validateMailAddress(&toMailAddress[m])) {
					warning("to mail address is invalid -> " + toMailAddress[m], ERROR)
					endApp(0)
				}
			}						
		} else {
			warning("no to mail address(es) given!", ERROR)
			endApp(0)
		}
	} 
	
	//sanity check from Mail adress	
	if (len(fromMailAddress) > 0) {
		if (!validateMailAddress(&fromMailAddress)) {
			warning("from mail address is invalid -> " + fromMailAddress, ERROR)
			endApp(0)
		}								
	} else {
		warning("no from mail address(es) given!", ERROR)
		endApp(0)
	}
	 
	//smtp port
	if (smtpPort <= 0) {
		smtpPort = 25
		warning("smtp port was set to zero, changed to default.", WARNING)
	}
	
	//title or subject
	if (len(mailSubject) <= 0) {
		warning("subject is empty.", WARNING)
	}
	
	//body
	if (len(mailBody) <= 0) {
		warning("body is empty.", WARNING)
	}
	
	
	// switch back to previous state
	verbose = oldVerbose
}





// ----- mail address checking ------

/* very simple mail address checking */
func validateMailAddress(mailAddress *string) bool {
	var valid bool = true
			
	// minimal length
	if (len(*mailAddress) <= MAIL_ADDRESS_MINIMAL_LENGTH) {
		valid = false
	}
	
				
	// needs to contain a point
	if !(strings.Contains(*mailAddress, ".")) {
		valid = false
	}
	
	// needs to contain a @
	if !(strings.Contains(*mailAddress, "@")) {
		valid = false
	}
	
	
	
	
	// may not contain ,
	if (strings.Contains(*mailAddress, ",")) {
		valid = false
	}
	
	// may not contain :
	if (strings.Contains(*mailAddress, ":")) {
		valid = false
	}
	
	// may not contain ?
	if (strings.Contains(*mailAddress, "?")) {
		valid = false
	}
	
	// may not contain !
	if (strings.Contains(*mailAddress, "!")) {
		valid = false
	}
	
	// may not contain -
	if (strings.Contains(*mailAddress, "-")) {
		valid = false
	}
	
	// may not contain +
	if (strings.Contains(*mailAddress, "+")) {
		valid = false
	}
	
	// may not contain |
	if (strings.Contains(*mailAddress, "|")) {
		valid = false
	}
	
	// may not contain ^
	if (strings.Contains(*mailAddress, "^")) {
		valid = false
	}
	
	// may not contain é
	if (strings.Contains(*mailAddress, "é")) {
		valid = false
	}
	
	// may not contain '
	if (strings.Contains(*mailAddress, "'")) {
		valid = false
	}
	
	// may not contain §
	if (strings.Contains(*mailAddress, "§")) {
		valid = false
	}
	
	// may not contain è
	if (strings.Contains(*mailAddress, "è")) {
		valid = false
	}
	
	// may not contain ç
	if (strings.Contains(*mailAddress, "ç")) {
		valid = false
	}
	
	// may not contain à
	if (strings.Contains(*mailAddress, "à")) {
		valid = false
	}
	
	// may not contain °
	if (strings.Contains(*mailAddress, "°")) {
		valid = false
	}
	
	// may not contain =
	if (strings.Contains(*mailAddress, "=")) {
		valid = false
	}
	
	
	
	
	// may not begin with @
	if (strings.HasPrefix(*mailAddress, "@")) {
		valid = false
	}
	
	// may not end with @
	if (strings.HasSuffix(*mailAddress, "@")) {
		valid = false
	}
	
	// may not begin with .
	if (strings.HasPrefix(*mailAddress, ".")) {
		valid = false
	}

	// may not end with .
	if (strings.HasSuffix(*mailAddress, ".")) {
		valid = false
	}
	

	
	return valid
}

