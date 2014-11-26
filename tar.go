package dirtar

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

// tar the given directory, paths inside the archive
// are relative to the directory
func Tar(dir string, w io.Writer) error {
	tw := tar.NewWriter(w)

	//walk and tar all files in a dir
	//@from http://stackoverflow.com/questions/13611100/how-to-write-a-directory-not-just-the-files-in-it-to-a-tar-gz-file-in-golang
	visit := func(fpath string, fi os.FileInfo, err error) error {

		//cancel walk if something went wrong
		if err != nil {
			return err
		}

		//skip root
		if fpath == dir {
			return nil
		}

		//dont 'add' dirs to archive
		if fi.IsDir() {
			return nil
		}

		f, err := os.Open(fpath)
		if err != nil {
			return err
		}
		defer f.Close()

		//use relative path inside archive
		rel, err := filepath.Rel(dir, fpath)
		if err != nil {
			return err
		}

		//create header from file info struct
		hdr, err := tar.FileInfoHeader(fi, rel)
		if err != nil {
			return err
		}

		//write header to archive
		// hdr.Name = rel?
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}

		//copy content into archive
		if _, err = io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	}

	//walk the context and create archive
	if err := filepath.Walk(dir, visit); err != nil {
		return err
	}

	return nil
}

// untar the archive into a given directory
func Untar(dir string, r io.Reader) error {

	//@todo check if dir is empty?

	//create a new reader
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return err
		}

		//create and open files
		//@todo this assumes the archives dir seperators are the same?
		path := filepath.Join(dir, hdr.Name)
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		//copy tar content into file, effectively untarring
		if _, err := io.Copy(f, tr); err != nil {
			return err
		}
	}

	return nil
}
