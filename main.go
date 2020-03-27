package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"io"
	"log"
	"os"
	"path/filepath"
)

// ManifestEntry describes a file entry in the manifest
type ManifestEntry struct {
	FilePath  string `json:"file_path"`
	SizeBytes int64  `json:"size_bytes"`
	Sha256Sum string `json:"sha256sum"`
	CID       string `json:"CID"`
}

var manifest []*ManifestEntry
var startAt string

func main() {
	startAt, _ = os.Getwd()
	err := processDir(startAt)
	if err != nil {
		log.Fatalf("Failed to start: %s", err)
	}
	data, _ := json.MarshalIndent(manifest, "", "    ")
	oh, err := os.Create("manifest.json")
	if err != nil {
		log.Fatalf("Failed to create manifest file: %s", err)
	}
	if _, err = oh.Write(data); err != nil {
		log.Fatalf("Failed to write manifest file: %s", err)
	}
}

func processDir(dir string) error {
	fmt.Printf("Processing directory: %s\n", dir)
	dirInfo, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !dirInfo.IsDir() {
		return fmt.Errorf("Not a directory: %s", dir)
	}
	filepath.Walk(dir, func(path string, pathInfo os.FileInfo, err error) error {
		fmt.Printf("Checking entry: %s\n", path)
		if pathInfo.IsDir() {
			return nil
		}
		return processFile(path, pathInfo)
	})

	return nil
}

func processFile(filePath string, fileInfo os.FileInfo) error {
	fh, err := os.Open(filePath)
	defer fh.Close()
	if err != nil {
		return err
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, fh); err != nil {
		return err
	}
	fh.Seek(0, io.SeekStart)
	cID, err := fileCID(fh)
	if err != nil {
		return err
	}
	entry := &ManifestEntry{
		FilePath:  filePath[len(startAt):],
		SizeBytes: fileInfo.Size(),
		Sha256Sum: hex.EncodeToString(hash.Sum(nil)),
		CID:       cID,
	}
	manifest = append(manifest, entry)
	return nil
}

func fileCID(fh *os.File) (string, error) {
	sh := shell.NewShell("localhost:5001")
	cid, err := sh.Add(fh)
	if err != nil {
		return "", err
	}
	return cid, nil
}
