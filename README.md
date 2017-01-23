# smtpmail
simple smtp mail cmdline client written in GO (GOLANG)

### os versions supported (using the default GOLANG gc compiler): 
windows (x86 and x64), linux, solaris 11 (intel x64), freebsd, openbsd, dragonfly, darwin (mac os) 

### os versions I have tested:
windows x64, solaris 11.3 (intel x64), linux x64 (Linux Mint)

### cmdline parameters:
```sh
./smtpmail -s smtp.domain.ext -p 25 -f from@domain.ext -t to1@domain.ext,to2@domain.ext -b "mail message" -sub "mail subject" -u <authuser> -pwd <password> [-al|-ap] [-q] [-c <config file>|-nf] [-v|-debug]
```

### expanded cmdline parameters:
-s = -smtpserver or -server<br />
-p = -smtpport or -port<br />
-u = -user<br />
-pwd = -password<br />
-al = -authlogin<br />
-ap = -authplain<br />
-b = -body<br />
-sub = -subject or -ti or -title<br />
-f = -from<br />
-t = -to<br />
-v = -verbose<br />
-q = -silent or -quiet<br />
-c = -config or -cfg or -file or -conf<br />
-nf = -nofile, -skip, -skipconfig<br />


### verbose and debugging
Parameter `-verbose` outputs basic operational information.

Parameter `-debug`  has no short versions. Instead you can use `-walk` or `-dump`, but those serve the same purpose as `-debug`. Debug can be seen as a "very verbose" parameter.


### config file naming convention and location:
The `smtpmail.conf.example` config file can be used. Simply copy file and rename it to smtpmail.conf. \(*\)

On Windows, the config file should be placed in the working directory (!= current directory).
On unix/bsd/mac flavors, the same rule applies then for Windows, but the file may also reside under folder `/etc`.

\(*\) Naming convention of the config file should be `<executable name>.conf`. On windows, the ".exe" extension should be stripped off.


The config file may contain parameters, but not all are required. You could easly make a combination or mix of config file parameters (acting as persistent settings) and variable ones provided through the cmdline. 

config file structure follows the following rules/parameters:
`// double slash is a comment line`<br />
`// you can make more the one comment line`<br />
`server=smtp.domain.ext`<br />
`port=25`<br />
`authuser=authuser@domain`<br />
`authpwd=authuserpwd`<br />
`frommail=frommail@domain.ext`<br />
`subject=some subject`<br />
`body=a body message`<br />
`authlogin=0|1`<br />
`verbose=0|1`<br />
`silent=0|1`<br />

The order of paramter appearance in the config file, does not matter.

 
### Authentication method:
AUTH LOGIN and AUTH PLAIN, serve the purpose of the smtp server authentication method.
I\'ve noticed that MS Exchange servers from MS Exchange 2010 or upward, prefer to use `AUTH LOGIN` instead of `AUTH PLAIN`.
The authentication method can be changed with the parameters ´-ap´ or ´al´ from the cmdline or by setting `authlogin=` to 0 or 1 in the config file.  

### TLS encryption:
smtpmail wil automatically try using plain or TLS encryption. This is a buildin feature of the GOLANG `net/smtp` package.

### Encoding:
smtpmail always uses base64 encoding in the mail header it sends.
