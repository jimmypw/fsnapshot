package fsnapshot

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
)

// FileManifest is the top level object
type FileManifest struct {
	Files []FileData
}

// Save will write the FileManifest to disk
func (fm *FileManifest) Save(loc string) error {
	var buf bytes.Buffer
	g := gob.NewEncoder(&buf)
	err := g.Encode(fm)
	if err != nil {
		return err
	}

	fd, err := os.Create(loc)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.Write(buf.Bytes())
	if err != nil {
		return err
	}

	fd.Close()

	return nil
}

// Load will read the FileManifest from disk
func (fm *FileManifest) Load(loc string) error {

	fd, err := os.Open(loc)
	if err != nil {
		return err
	}
	defer fd.Close()

	buf := bufio.NewReader(fd)
	g := gob.NewDecoder(buf)

	err = g.Decode(fm)
	if err != nil {
		return err
	}

	fd.Close()

	return nil
}

func (fm *FileManifest) setallseen(b bool) {
	for i := range fm.Files {
		fm.Files[i].seen = b
	}
}

// Compare will compare two manifests and produce a third with the differences
func Compare(oldfm, newfm *FileManifest) (*FileManifest, error) {
	var returndata FileManifest
	oldfm.setallseen(false)
	newfm.setallseen(false)

	for i := range newfm.Files {
		for x := range oldfm.Files {
			res := bytes.Compare(oldfm.Files[x].FilenameChecksum, newfm.Files[i].FilenameChecksum)
			if res == 0 {
				oldfm.Files[x].seen = true
				newfm.Files[i].seen = true
			}
		}
	}

	// At this point the seen flag all files that exist both arrays have been marked as "seen"
	// Files that are marked as not seen in oldfm have been deleted.
	// Files that are marked as not seen in newfm have been created.

	// Mark all unseen files in oldfm as deleted
	for i := range oldfm.Files {
		if !oldfm.Files[i].seen {
			oldfm.Files[i].FileStatus = fDeleted
		}
	}

	// Mark all unseen files in newfm as new
	for i := range newfm.Files {
		if !newfm.Files[i].seen {
			newfm.Files[i].FileStatus = fNew
		}
	}

	// Compare the checksums of all "seen" files to establish if the file has changed
	for i := range newfm.Files {
		for x := range oldfm.Files {
			res := bytes.Compare(oldfm.Files[x].FilenameChecksum, newfm.Files[i].FilenameChecksum)
			if res == 0 && oldfm.Files[x].seen && newfm.Files[i].seen {
				res := bytes.Compare(oldfm.Files[x].Checksum, newfm.Files[i].Checksum)
				if res == 0 {
					oldfm.Files[x].FileStatus = fUnmodified
					newfm.Files[i].FileStatus = fUnmodified
				} else {
					oldfm.Files[x].FileStatus = fChanged
					newfm.Files[i].FileStatus = fChanged
				}
			}
		}
	}

	// Finally we merge the two manifests to produce a delta
	returndata.Files = oldfm.Files
	for i := range newfm.Files {
		if !newfm.Files[i].seen {
			returndata.Files = append(returndata.Files, newfm.Files[i])
		}
	}

	return &returndata, nil
}

// Report will provide a textual report based on a provided file manifest
func (fm FileManifest) Report() {
	for i := range fm.Files {
		switch fm.Files[i].FileStatus {
		case fChanged:
			fmt.Printf("C %s\n", filepath.Join(fm.Files[i].Subdirectory, fm.Files[i].Filename))
		case fNew:
			fmt.Printf("N %s\n", filepath.Join(fm.Files[i].Subdirectory, fm.Files[i].Filename))
		case fDeleted:
			fmt.Printf("D %s\n", filepath.Join(fm.Files[i].Subdirectory, fm.Files[i].Filename))
		case fUnmodified:
			fmt.Printf("U %s\n", filepath.Join(fm.Files[i].Subdirectory, fm.Files[i].Filename))
		default:
			fmt.Println("Invalid file status")
		}
	}
}
