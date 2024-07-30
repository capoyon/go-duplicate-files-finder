package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type FileInfo struct {
    Path string
    Size int64
}

func main() {
    start := time.Now()
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <directory>")
        return
    }

    dir := os.Args[1]
    if fileMap, err := hashFilesInDir(dir); err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        printDuplicateFile(fileMap)
    }
    
    fmt.Printf("Scanning duration: %v\n", time.Since(start))
}

func printDuplicateFile(m map[string][]FileInfo) {
    totalSize := int64(0)
    totalSave := int64(0)
   
    for _, values := range m {
        if len(values) > 1 {

            firstSize := values[0].Size
            allSameSize := true
            allsize := int64(0)
            saveSize := int64(0)
        
            for _, file := range values {
                allsize += file.Size
                
                // check for hash collision
                if file.Size != firstSize {
                    allSameSize = false
                }

            }

            totalSize += allsize
            saveSize = allsize - values[0].Size
            totalSave += saveSize

            fmt.Println("Size:", formatBytes(values[0].Size), "Save: ", formatBytes(saveSize))
            for _, file := range values {
                fmt.Println(" -", file.Path)

            }
            if !allSameSize {
                fmt.Println("Hash collision detected!")
            }
            fmt.Println()
        }
    }
    fmt.Printf("Total Duplicate size: %v\n", formatBytes(totalSize))
    fmt.Printf("Size save if only one file stay: %v\n", formatBytes(totalSave))
 
}

// hashFilesInDir walks through the directory and hashes all files.
func hashFilesInDir(dir string) (map[string][]FileInfo, error) {
    fileHashMap := make(map[string][]FileInfo)
    err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() {

            hash, err := hashFile(path)
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

// hashFile computes the MD5 hash of the file at the given path.
func hashFile(path string) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer file.Close()

    hasher := md5.New()
    if _, err := io.Copy(hasher, file); err != nil {
        return "", err
    }

    return hex.EncodeToString(hasher.Sum(nil)), nil
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

// formatBytes converts a byte size into a human-readable string.
func formatBytes(bytes int64) string {
    const (
        _       = iota
        KB = 1 << (10 * iota)
        MB
        GB
        TB
        PB
        EB
        ZB
        YB
    )

    var value float64 = float64(bytes)
    var unit string

    switch {
    case value >= PB:
        value /= PB
        unit = "PB"
    case value >= TB:
        value /= TB
        unit = "TB"
    case value >= GB:
        value /= GB
        unit = "GB"
    case value >= MB:
        value /= MB
        unit = "MB"
    case value >= KB:
        value /= KB
        unit = "KB"
    default:
        unit = "B"
    }

    return fmt.Sprintf("%.2f %s", value, unit)
}
