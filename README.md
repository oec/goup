# goup - An updater for go

Download and install a version of go under `/opt` as a symlink <tt>/opt/go → /opt/go<i>version</i></tt>

Usage: goup [version]

Parameters:

  -arch string
    	architecture to install (default "amd64")
  -dst string
    	directory to install go to (default "/opt")
  -n	dry run, don't install
  -os string
    	OS to install (default "linux")
  -url string
    	download-url (default "https://dl.google.com/go/")

Output:

     % ./goup
    2019/05/17 12:04:06 using version go1.12.5
    2019/05/17 12:04:06 Downloading https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz
    2019/05/17 12:04:43 Checking Signature go1.12.5.linux-amd64.tar.gz
    2019/05/17 12:04:44 Unpacking go1.12.5.linux-amd64.tar.gz
    2019/05/17 12:04:49 Creating symlink go → go1.12.5
    2019/05/17 12:04:49 Ugprade to go1.12.5 done.

