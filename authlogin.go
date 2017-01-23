/*
 * Pieter De Ridder
 * 2017 
 *
 * Base64 smtp mail client
 * Can be cross compiled for Windows, Linux, Mac OS X and Solaris on Intel x86-64
 *
 * ------------------------------------------------------------
 * 
 * authlogin go file
 * specific smtp authentication for MS Exchange servers or similar.
 *
 */
 
package main

import (
	"net/smtp"
	"errors"
)
	
// ----- custom AUTH LOGIN structures ------ 
/* SMTP AUTH LOGIN handler (MS Exchange and Office 356) */
type loginAuth struct {
  username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

