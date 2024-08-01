package main

import (
    "flag"
	"fmt"
	"os"
	"time"

	"gdff/filemgr"
    "gdff/format"
)


func main() {
    var path string
    flag.StringVar(&path, "path", "", "path/directory to scan")
    flag.Parse()

    if path == "" {
        fmt.Printf("Error: --path flag is required\n\n")
        flag.Usage()
        os.Exit(1)
    }
    if !isValidDir(path) {
        fmt.Printf("%s is not a valid directory.\n", path)
        os.Exit(1)
    }


    start := time.Now()
    fileMap, err := filemgr.GroupBySize(path)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Print("\n\nHashing...\n")     
    hashToFileInfo, err := filemgr.GroupByHash(fileMap)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        printDuplicateFile(*hashToFileInfo)
    }
    fmt.Printf("Scanning duration: %v\n", time.Since(start))
}

func isValidDir(path string) bool {
    info, err := os.Stat(path)
    if err != nil {
        return false
    }
    if !info.IsDir() {
        return false
    }

    return true
}

func printSizeToPath(paths *[]string) {
     if paths == nil {
        fmt.Println("The arr is nil")
        return
    }
    
    for _, path := range *paths {
        fmt.Printf("Size: %s bytes\n", path)
    }
}

func printDuplicateFile(m map[string][]filemgr.FileInfo) {
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

            fmt.Println("Size:", format.FileSize(values[0].Size), "Save: ", format.FileSize(saveSize))
            for _, file := range values {
                fmt.Println(" -", file.Path)

            }
            if !allSameSize {
                fmt.Println("Hash collision detected!")
            }
            fmt.Println()
        }
    }
    fmt.Printf("Total Duplicate size: %v\n", format.FileSize(totalSize))
    fmt.Printf("Size save if only one file stay: %v\n", format.FileSize(totalSave))
 
}
