package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
)

// ManifestEntry describes a file entry in the manifest
type ManifestEntry struct {
	FilePath  string `json:"file_path"`
	SizeBytes int64  `json:"size_bytes"`
	Sha256Sum string `json:"sha256sum"`
	CID       string `json:"CID"`
}

// TemplateData is all the data we'll send to the template
type TemplateData struct {
	Manifest  []*ManifestEntry `json:"manifest"`
	CreatedAt time.Time        `json:"created_at"`
}

var manifest []*ManifestEntry
var startAt string

func main() {
	versionFlag := flag.Bool("version", false, "print version information")
	flag.Parse()
	if *versionFlag {
		fmt.Printf("(version=%s, branch=%s, gitcommit=%s)\n", Version, GitBranch, GitCommit)
		fmt.Printf("(go=%s, user=%s, date=%s)\n", GoVersion, BuildUser, BuildDate)
		os.Exit(0)
	}
	startAt, _ = os.Getwd()
	err := processDir(startAt)
	if err != nil {
		log.Fatalf("Failed to start: %s", err)
	}
	// sort entries by filePath
	sort.Slice(manifest, func(i, j int) bool { return manifest[i].FilePath > manifest[j].FilePath })
	if err = writeJSON(manifest); err != nil {
		log.Fatalf("Failed to write manifest file: %s", err)
	}
	if err = writeHTML(manifest); err != nil {
		log.Fatalf("Failed to write manifest file: %s", err)
	}
}

func writeHTML(manifest []*ManifestEntry) error {
	tData := &TemplateData{
		Manifest:  manifest,
		CreatedAt: time.Now(),
	}

	t, err := template.ParseFiles("index.tmpl")
	if err != nil {
		return err
	}
	f, err := os.Create("index.html")
	if err != nil {
		return err
	}
	err = t.Execute(f, tData)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func writeJSON(manifest []*ManifestEntry) error {
	data, _ := json.MarshalIndent(manifest, "", "    ")
	oh, err := os.Create("manifest.json")
	if err != nil {
		return err
	}
	if _, err = oh.Write(data); err != nil {
		return err
	}
	return nil
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
