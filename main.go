package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
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

type VersionInfo struct {
	Stable  bool   `json:"stable"`
	Version string `json:"version"`
	Files   []struct {
		Arch     string `json:"arch"`
		Filename string `json:"filename"`
		Kind     string `json:"kind"`
		Os       string `json:"os"`
		Sha256   string `json:"sha256"`
		Size     int    `json:"size"`
		Version  string `json:"version"`
	} `json:"files"`
}

const VersionsURL = "https://golang.org/dl/?mode=json"

var (
	// https://dl.google.com/go/go1.10.linux-amd64.tar.gz

	dlurl = flag.String("url", "https://dl.google.com/go/", "download-url")
	dst   = flag.String("dst", "/opt", "directory to install go to")
	zos   = flag.String("os", runtime.GOOS, "OS to install")
	arch  = flag.String("arch", runtime.GOARCH, "architecture to install")
	dry   = flag.Bool("n", false, "dry run, don't install")
)

func main() {
	flag.Parse()

	r, e := http.Get(VersionsURL)
	if e != nil {
		log.Fatal(e)
	}
	defer r.Body.Close()

	var (
		index    int
		version  string
		filename string
		hash     string
		found    bool
		versions = []VersionInfo{}
	)

	dec := json.NewDecoder(r.Body)
	if e = dec.Decode(&versions); e != nil {
		log.Fatal(e)
	}

	if version = flag.Arg(0); len(version) != 0 {
		for i := range versions {
			if versions[i].Version == version {
				index = i
				found = true
				break
			}
		}
		if !found {
			fmt.Println("No such version:", version)
			fmt.Println("Only the following version are available:")
			for _, v := range versions {
				fmt.Println("\t" + v.Version)
			}
			os.Exit(1)
		}
	} else {
		version = versions[index].Version
	}

	if version == runtime.Version() {
		fmt.Println("Version", version, "is already the current version")
		return
	}

	fmt.Println("Using version", version)

	found = false
	for _, f := range versions[index].Files {
		if f.Arch == *arch && f.Os == *zos {
			filename = f.Filename
			hash = f.Sha256
			found = true
			break
		}
	}
	if !found {
		fmt.Println("No such architecture+os:", *arch, "+", *zos)
		os.Exit(1)
	}

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

	fmt.Println("Downloading", *dlurl+filename)
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

	fmt.Println("Checking Signature", filename)
	hasher := sha256.New()
	_, err = io.Copy(hasher, targz)
	if err != nil {
		log.Fatal(err)
	} else if hval := fmt.Sprintf("%x", hasher.Sum(nil)); hval != hash {
		log.Fatal("sha256 mismatch! Expected ", hash, ", but got ", hval)
	}

	_, err = targz.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Unpacking", filename)
	err = untar(version, targz)
	if err != nil {
		log.Fatal(err)
	}

	if *dry {
		fmt.Println("Not going any further")
		return
	}

	fmt.Println("Creating symlink go →", version)
	if inf, err := os.Lstat("go"); err == nil && inf.Mode()&os.ModeSymlink != 0 {
		os.Remove("go") // just try to delete the existing link.
	}

	err = os.Symlink(version, "go")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ugprade to", version, "done.")
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
