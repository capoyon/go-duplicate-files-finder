package filemgr

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

    "gdff/hash"
)

type FileInfo struct {
    Path string
    Size int64
}


// hashFilesInDir walks through the directory and hashes all files.
func HashFilesInDir(dir string) (map[string][]FileInfo, error) {
    fileHashMap := make(map[string][]FileInfo)
    err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() {

            hash, err := hash.HashFile(path)
            if err != nil {
                fmt.Errorf("%v\n", err)
            }else {
                size, err := getFileSize(path) 
                if err != nil {
                    return err
                }

                fileInfo := FileInfo{
                    Path: path,
                    Size: size,
                }

                if _, exist := fileHashMap[hash]; exist {
                    fileHashMap[hash] = append(fileHashMap[hash], fileInfo)
                } else {
                    fileHashMap[hash] = []FileInfo{fileInfo}
                }

            }

        }
        return nil
    })

    return fileHashMap, err
}



// getFileSize computes the human-readable size of a file.
func GetFileSize(path string) (int64, error) {
    // Open the file
    file, err := os.Open(path)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    // Get the file statistics
    fileInfo, err := file.Stat()
    if err != nil {
        return 0, err
    }

    // Get the file size and format it
    bytes := fileInfo.Size()

    return bytes, nil
}

// hashFilesInDir walks through the directory and hashes all files.
func hashFilesInDir(dir string) (map[string][]FileInfo, error) {
    fileHashMap := make(map[string][]FileInfo)
    err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() {

            hash, err := hash.HashFile(path)
            if err != nil {
                fmt.Errorf("%v\n", err)
            }else {
                size, err := getFileSize(path) 
                if err != nil {
                    return err
                }

                fileInfo := FileInfo{
                    Path: path,
                    Size: size,
                }

                if _, exist := fileHashMap[hash]; exist {
                    fileHashMap[hash] = append(fileHashMap[hash], fileInfo)
                } else {
                    fileHashMap[hash] = []FileInfo{fileInfo}
                }

            }

        }
        return nil
    })

    return fileHashMap, err
}



// getFileSize computes the human-readable size of a file.
func getFileSize(path string) (int64, error) {
    // Open the file
    file, err := os.Open(path)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    // Get the file statistics
    fileInfo, err := file.Stat()
    if err != nil {
        return 0, err
    }

    // Get the file size and format it
    bytes := fileInfo.Size()

    return bytes, nil
}
