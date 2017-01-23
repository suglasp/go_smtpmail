/*
 * Pieter De Ridder
 * 2017
 *
 * Base64 smtp mail client
 * Can be cross compiled for Windows, Linux, Mac OS X and Solaris on Intel x86-64
 *
 * original date : 18-01-2017
 * last updated  : 23-01-2017
 * 
 * ------------------------------------------------------------
 * 
 * main smtpmail go file
 *
 * ------------------------------------------------------------
 /
 * More info on auth
 * http://www.samlogic.net/articles/smtp-commands-reference-auth.htm
 *
 * golang smtp package doc:
 * https://golang.org/pkg/net/smtp/
 * https://golang.org/src/net/smtp/smtp_test.go
 * 
 */

package main

import (
	"bytes"           // bytes.Buffer
	"encoding/base64" // base64 encoding
	"fmt"
	//"log"             // logging
	"net/mail"
	"net/smtp"            
	"strconv"         // int to string
)


/* function to send a e-mail with authentication */
func doMail(a smtp.Auth, smtpSrv string, smtpP int, fromMA string, toMA string, s string, b string) {
	/*
	 * function variables:
	 * a       = auth
	 * smtpSrv = server
	 * smtpP   = port
	 * fromMA  = from Mail addresses
	 * toMA    = to Mail addresses
	 * s       = subject / title
	 * b       = body
	 */
	
	print("sending to " + toMA + ".")
	
	debugprint("")
	debugprint("doMail routine:")
	debugprint("smtp = " + smtpSrv + ":" + strconv.Itoa(smtpP)) 
	debugprint("from = " + fromMA)
	debugprint("to = " + toMA + " <-")
	debugprint("subject = " + s)
	debugprint("body = " + b)
	
	// from and to addresses
	from := mail.Address{"", fromMA}
	to := mail.Address{"", toMA}
	
	// build smtp mail header
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	// header["Subject"] = encodeRFC2047(title)
	header["Subject"] = s
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	//build message body
	var msg bytes.Buffer
	for k, v := range header {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(base64.StdEncoding.EncodeToString([]byte(b)))  //body message

	//concatenate server and port to "server:port""
	var hostAndPort bytes.Buffer
	hostAndPort.WriteString(smtpSrv)             // smtp server host
	hostAndPort.WriteString(":")                 // add colon ":" to separate "host:port"
	hostAndPort.WriteString(strconv.Itoa(smtpP)) // smtp server port
	
	// try sending mail
	if (a != nil) {
		err := smtp.SendMail( hostAndPort.String(), a, from.Address, []string{to.Address}, []byte( msg.String() ) )

		if err != nil {
			//log.Fatal(err)
			warning(err.Error(), ERROR)
			warning("abandoned loop", WARNING)
			endApp(-1)
		}
	}
}



// ----- main function ------ 
/* main function */
func main() {	
	InitApp()
		
	// setup smtp authentication
	var auth smtp.Auth
	
	if (useAuthLogin) {
		auth = LoginAuth(userName, userPasswd)
	} else {
		auth = smtp.PlainAuth("", userName, userPasswd, smtpServer)
	}
	
	debugprint("")
	debugprint("to addresses slice size = " + strconv.Itoa(len(toMailAddress)))
	
	// send mail to all mail recipes
	for mi := 0; mi < len(toMailAddress); mi++ {
		debugprint("")
		debugprint("fetch [" + strconv.Itoa(mi) + "] address = " + toMailAddress[mi] + " <-")
		
		doMail(auth, smtpServer, smtpPort, fromMailAddress, toMailAddress[mi], mailSubject, mailBody)
	}
	
	print("message(s) send.")
	
	//exit
	endApp(0)
}
