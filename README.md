# goup - An updater for go

Download and install a version of go under `/opt` as a symlink <tt>/opt/go → /opt/<i>version</i></tt>

Usage: goup [version]

Parameters:

```
   -arch string
     	architecture to install (default "amd64")
   -dst string
     	directory to install go to (default "/opt")
   -n	dry run, don't install
   -os string
     	OS to install (default "linux")
   -url string
     	download-url (default "https://dl.google.com/go/")
```

Output:

	 % goup        
	Using version go1.12.6
	Downloading https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz
	Checking Signature go1.12.6.linux-amd64.tar.gz
	Unpacking go1.12.6.linux-amd64.tar.gz
	Creating symlink go → go1.12.6
	Ugprade to go1.12.6 done.

