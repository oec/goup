package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	// https://dl.google.com/go/go1.10.linux-amd64.tar.gz
	dlurl   = flag.String("url", "https://dl.google.com/go/", "download-url")
	dst     = flag.String("dst", "/opt", "directory to install go to")
	version = ""
)

func main() {
	flag.Parse()

	version = flag.Arg(0)
	if len(version) == 0 {
		log.Fatal("version needed")
	}

	filename := fmt.Sprintf("go%s.%s-%s.tar.gz", version, runtime.GOOS, runtime.GOARCH)

	var (
		err   error
		resp  *http.Response
		targz *os.File
	)

	err = os.Chdir(*dst)
	if err != nil {
		log.Fatal(err)
	}

	targz, err = os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Downloading", *dlurl+filename)
	resp, err = http.Get(*dlurl + filename)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(targz, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	_, err = targz.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Unpacking", filename)
	err = untar("go"+version, targz)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Creating symlink go →", "go"+version)
	if inf, err := os.Lstat("go"); err == nil && inf.Mode()&os.ModeSymlink != 0 {
		os.Remove("go") // just try to delete the existing link.
	}

	err = os.Symlink("go"+version, "go")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Ugprade to", version, "done.")
}

func untar(dst string, r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	defer gzr.Close()
	if err != nil {
		return err
	}

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

			// return any other error
		case err != nil:
			return err

			// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, strings.TrimPrefix(header.Name, "go"))

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

			// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			f.Close()
		}
	}
}
