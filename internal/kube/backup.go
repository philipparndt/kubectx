package kube

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func timestamp() string {
	return time.Now().Format("20060102T150405")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// Cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}

	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func Backup() {
	original := FileName()
	backup := original + ".kubectx.backup." + timestamp()
	err := CopyFile(original, backup)
	if err != nil {
		panic(err)
	}

	fmt.Println("Backup saved to", backup)
}
