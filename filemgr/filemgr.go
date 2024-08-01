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


// returns the path of the file that has the same size
func GroupBySize(dir string) (*[]string, error) {
    sizeToPath := make(map[int64][]string)

    err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() {
            fileSize, err := getFileSize(path)
            if err != nil {
                return fmt.Errorf("error getting file size for %s: %v", path, err)
            }
            sizeToPath[fileSize] = append(sizeToPath[fileSize], path)
        }
        return nil
    })

    if err != nil {
        return nil, err
    }

    var result []string

    //only include the path the has the same size
    for _, arr := range sizeToPath {
        if len(arr) > 1 {
            result = append(result, arr...)
        }
    }

    return &result, nil
}


// hashFilesInDir walks through the directory and hashes all files.
func GroupByHash(paths *[]string) (*map[string][]FileInfo, error) {
    hashToFileInfo := make(map[string][]FileInfo)

    for _, path := range *paths{
       
        _, err := os.Stat(path)
        if os.IsNotExist(err) { continue }

        hash, err := hash.HashFile(path)
        if err != nil {
            fmt.Println(err)
            continue
        }else {
            size, err := getFileSize(path) 
            if err != nil {
                fmt.Println(err)
                continue
            }

            fileInfo := FileInfo{
                Path: path,
                Size: size,
            }
            hashToFileInfo[hash] = append(hashToFileInfo[hash], fileInfo)
        }
    }


    return &hashToFileInfo, nil
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
                fmt.Println(err)
                return err
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
