## smtpmail

simple smtp mail cmdline client written in GO (GOLANG)


### os versions supported (using gc): 
windows (x86 and x64), linux, solaris 11 (intel x64), freebsd, openbsd, dragonfly, darwin (mac os) 


### os versions tested:
windows x64, solaris 11.3 (intel x64), linux x64 (Linux Mint)


### cmdline parameters:
-s <smtp server>
-p <smtp port> (default 25)
-f <from>
-t <to@domain.ext,to@domain.ext, ...>
-b "message"
-sub "subject"
-u <authuser>      [ note: for exchange you should use user@domain for login ]
-pwd <password>
-al or -ap
-v
-q
-debug
-c <config file>
-nf

### expanded cmdline parameters:
-s          =>  -smtpserver or -server
-p          =>  -smtpport or -port
-u          =>  -user
-pwd        =>  -password
-al         =>  -authlogin
-ap         =>  -authplain
-b          =>  -body
-sub        =>  -subject or -ti or -title
-f          =>  -from
-t          =>  -to
-v          =>  -verbose
-q          =>  -silent or -quiet
-c          =>  -config or -cfg or -file or -conf
-nf         =>  -nofile, -skip, -skipconfig


### verbose and debugging
 Parameter '-verbose' outputs basic operational information.

 Parameter '-debug' has no short versions. Instead you can use '-walk' or '-dump', but those serve the same purpose as -debug. 
 Debug can be used as a "very verbose" function.


### config file naming convention and location:
 The smtpmail.conf.example config file can be used. Simply copy file and rename it to smtpmail.conf. (*)

 On Windows, the config file should be placed in the working directory.
 On unix/bsd/mac flavors, the same rule applies then for Windows, but the file may also reside under folder /etc.

 (*) Naming convention of the config file should be <executable name>.conf. On windows, the ".exe" extension should be stripped off.


 The config file may contain parameters, but not all are required. You could easly make a combination of config file parameters (persistent ones)
 and variable ones through the cmdline.

 config file structure follows the following rules/parameters:
   // double slash is a comment line
   server=smtp.domain.ext
   port=25
   authuser=authuser@domain
   authpwd=authuserpwd
   frommail=fromuser@domain
   subject=a subject
   body=a body message
   authlogin=0|1
   verbose=0|1
   silent=0|1

 The order does not matter.

 
### Authentication method:
 Auth login and auth plain, serve the purpose of the smtp server authentication method.
 I've noticed that MS Exchange servers from MS Exchange 2010 or upward, need to use auth login instead of auth plain.


### TLS encryption:
 smtpmail wil automatically try using plain or TLS encryption. This is a buildin feature of the GOLANG net/smtp package.


### Encoding:
 smtpmail always uses base64 encoding in the mail header it sends.
