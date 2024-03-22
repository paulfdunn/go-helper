// Package ziph provides helper functions for the zip package.
package ziph

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ZipStats is for getting statistics on a zip file; currently only
// supports the number of zip.File in an archive.
type ZipStats struct {
	FileCount int
}

// AsyncUnzip asynchronously unzips inputPath to outputPath; outputPath will be
// created if it does not exist. Directories are created with permDir permissions.
// Progress can be monitored via the returned channels, which return done, processed
// paths, and any errors. The done channel will fire once when work is done.
func AsyncUnzip(inputPath, outputPath string, permDir os.FileMode) (<-chan bool, <-chan string, <-chan error) {
	zs, err := GetZipStats(inputPath)
	if err != nil {
		// If GetZipStats failed, just set buffer lengths to 1 because the zip.OpenReader
		// will fail again below and return.
		zs.FileCount = 1
	}
	done := make(chan bool, 1)
	// Size channels so that they don't block if the caller is only checking done.
	processedPaths := make(chan string, zs.FileCount)
	errors := make(chan error, zs.FileCount)
	go func() {
		zr, err := zip.OpenReader(inputPath)
		if err != nil {
			done <- true
			errors <- err
			return
		}
		defer zr.Close()

		outputPath, err = filepath.Abs(outputPath)
		if err != nil {
			done <- true
			errors <- err
			return
		}

		for _, f := range zr.File {
			err := removeFromZip(f, outputPath, permDir)
			processedPaths <- outputPath
			if err != nil {
				done <- true
				errors <- err
				return
			}
		}
		done <- true
	}()

	return done, processedPaths, errors
}

// AsyncZip asynchronously creates a compressed ZIP file, of file/directories
// in paths, at zipPath. The paths are turned into absolute paths, then made
// relative by removing the leading filepath.Separator.
// Progress can be monitored via the returned channels, which return done, processed
// paths, and any errors. The done channel will fire once when work is done.
func AsyncZip(zipPath string, paths []string) (<-chan bool, <-chan string, <-chan error) {
	done := make(chan bool, 1)
	// Size channels so that they don't block if the caller is only checking done.
	processedPaths := make(chan string, len(paths))
	errors := make(chan error, len(paths))
	go func() {
		f, err := os.Create(zipPath)
		if err != nil {
			done <- true
			errors <- err
			return
		}
		defer f.Close()

		zipWriter := zip.NewWriter(f)
		defer zipWriter.Close()

		for _, path := range paths {
			err := filepath.WalkDir(path, addToZip(zipWriter))
			processedPaths <- path
			if err != nil {
				done <- true
				errors <- err
			}
		}

		done <- true
	}()

	return done, processedPaths, errors
}

// GetZipStats is for getting statistics on a zip file; currently only
// supports the number of zip.File in an archive.
func GetZipStats(inputPath string) (*ZipStats, error) {
	zr, err := zip.OpenReader(inputPath)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	zs := ZipStats{FileCount: 0}
	for range zr.File {
		zs.FileCount++
	}
	return &zs, nil
}

// addToZip is a closure called by Create to add a directory or file to the zipWriter;
// do not directly call this function.
// The paths are turned into absolute paths, then made relative by removing
// the leading filepath.Separator.
func addToZip(zipWriter *zip.Writer) func(string, fs.DirEntry, error) error {
	return func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return fs.SkipDir
		}

		info, err := dirEntry.Info()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name, err = filepath.Abs(path)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(header.Name, string(filepath.Separator))
		if info.IsDir() {
			header.Name += string(filepath.Separator)
		}
		header.Method = zip.Deflate

		headerWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	}
}

// removeFromZip removes a zipFile from its archive. The outputPath is checked
// for Zip Slip (https://github.com/golang/go/issues/40373) and an error is returned for
// inappropriate paths.
func removeFromZip(zipFile *zip.File, outputPath string, permDir os.FileMode) error {
	// zipFile.Name is a relative path and file name
	outputFilePath := filepath.Join(outputPath, zipFile.Name)
	// Reject paths that might Zip Slip; I.E. if zipFile.Name uses ../ to access
	// directories outside outputPath the file is rejected.
	if !strings.HasPrefix(outputFilePath, filepath.Clean(outputPath)+string(os.PathSeparator)) {
		return fmt.Errorf("removeFromZip invalid file path: %s", outputFilePath)
	}

	if zipFile.FileInfo().IsDir() {
		if err := os.MkdirAll(outputFilePath, permDir); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(outputFilePath), permDir); err != nil {
		return err
	}

	f, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
	if err != nil {
		return err
	}
	defer f.Close()

	irc, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer irc.Close()

	if _, err := io.Copy(f, irc); err != nil {
		return err
	}
	return nil
}
