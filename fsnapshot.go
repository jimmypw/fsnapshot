package fsnapshot

import (
	"os"
	"path/filepath"
)

func readFileData(initialdirectory, subdirectory, f string) (FileData, error) {
	var fd FileData
	t, err := os.Open(filepath.Join(initialdirectory, subdirectory, f))
	defer t.Close()
	if err != nil {
		return fd, err
	}
	cs, err := checksumfile(t)
	if err != nil {
		return fd, err
	}

	fd.Filename = f
	fd.FilenameChecksum = filenameChecksum([]byte(filepath.Join(subdirectory, f)))
	fd.Subdirectory = subdirectory
	fd.Checksum = cs
	fd.FileStatus = fUnmodified
	fd.seen = true

	return fd, nil
}

func scandirectory(initialdirectory string, subdirectory string, files *[]FileData) error {
	var directories []string = make([]string, 0)

	handle, err := os.Open(filepath.Join(initialdirectory, subdirectory))
	defer handle.Close()

	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	listing, err := handle.Readdir(0)
	if err != nil {
		return err
	}

	for i := range listing {
		if listing[i].IsDir() {
			directories = append(directories, filepath.Join(subdirectory, listing[i].Name()))
		} else {
			t, err := readFileData(initialdirectory, subdirectory, listing[i].Name())
			if err != nil {
				return err
			}

			*files = append(*files, t)
		}
	}

	for i := range directories {
		scandirectory(initialdirectory, directories[i], files)
	}

	return nil
}

// Snapshot will checksum all files in the specified directory
func Snapshot(initialdirectory string) (FileManifest, error) {
	var fm FileManifest
	fm.Files = make([]FileData, 0)
	scandirectory(initialdirectory, "", &fm.Files)
	return fm, nil
}
