package main

import (
	"fmt"
	"os"
	"time"

	"gdff/filemgr"
    "gdff/format"
)


func main() {
    start := time.Now()
    if len(os.Args) < 2 {
        fmt.Println(`Usage: go run main.go <directory>`)
        return
    }

    dir := os.Args[1]
    if fileMap, err := filemgr.HashFilesInDir(dir); err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        printDuplicateFile(fileMap)
    }
    
    fmt.Printf("Scanning duration: %v\n", time.Since(start))
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
