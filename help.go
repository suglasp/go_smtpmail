
/*
 * Pieter De Ridder
 * 2017 
 *
 * Base64 smtp mail client
 * Can be cross compiled for Windows, Linux, Mac OS X and Solaris on Intel x86-64
 *
 * ------------------------------------------------------------
 * 
 * help smtpmail go file
 * contains function regarding user help
 *
 */


package main


// ----- application help ------ 
/* print help to console */
func printHelp() {
	verbose = true  // fixed enable verbose ouput
	
	// print the help line with only exe file name
	print("")
	print("smtpmail v" + PROGRAM_VERSION)
	print("Pieter De Ridder, 2017. Provided under the MIT license.")
	print("")
	print("Usage: ")
    print(stripExeFilename() + " -s <server> -p <port> -u <user> -pwd <passwd> -f <from@domain> -t <to@domain[,to-n@domain]> -b \"body msg\" -sub \"subject\" [-al|-ap] [-v] [-silent] [-c <configfile>] [-debug] [-nf]")
	print("")
	print("The -t parameter (= to mail address) can contain multiple addresses.")
    print("For example: -t john@europe.eu,fanny@australia.au,bender@mars-university.org")
	print("")
	print("Default SMTP TCP server ports are 25, 465 or 587.")
	print("")
	
	endApp(0)  // terminate
}

