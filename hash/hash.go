package hash

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// hashFile computes the MD5 hash of the file at the given path.
func HashFile(path string) (string, error) {
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
