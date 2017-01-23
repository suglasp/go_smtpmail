
/*
 * Pieter De Ridder
 * 2017 
 *
 * Base64 smtp mail client
 * Can be cross compiled for Windows, Linux, Mac OS X and Solaris on Intel x86-64
 *
 * ------------------------------------------------------------
 * 
 * encoding go file
 * extra function for Base64 encoding following RFC2047
 *
 */


package main

import (
	"strings"
	"net/mail"
)

// ----- Encoding ------ 
/* encode for base64 */
func encodeRFC2047(String string) string {
	/* use mail's rfc2047 to encode any string */
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}
