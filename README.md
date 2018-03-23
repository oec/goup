# goup - An updater for go

Download and install a version of go under `/opt` as a symlink <tt>/opt/go → /opt/go<i>version</i></tt>

Usage: goup _version_

Parameters:

    -dst string
          directory to install go to (default "/opt")
    -url string
          download-url (default "https://dl.google.com/go/")

Output:

    > goup 1.10
    2018/03/23 18:58:42 Downloading https://dl.google.com/go/go1.10.linux-amd64.tar.gz
    2018/03/23 18:58:51 Unpacking go1.10.linux-amd64.tar.gz
    2018/03/23 18:58:56 Creating symlink go → go1.10
    2018/03/23 18:58:56 Ugprade to 1.10 done.
