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
// Set the bufSize to the number of files (from GetZipStats) to prevent this function
// from being blocked output channels not being read fast enough.
// Progress can be monitored via the returned channels, which return cancel, processed
// paths, and any errors. The cancel channel can be used to cancel an operation.
// The operation is complete when both processed paths and errors channels are closed.
func AsyncUnzip(inputPath, outputPath string, bufSize int, permDir os.FileMode) (chan<- bool, <-chan string, <-chan error) {
	cancel := make(chan bool, 1)
	// Size channels so that they don't block if the caller is only checking done.
	processedPaths := make(chan string, bufSize)
	errors := make(chan error, bufSize)
	go func() {
		zr, err := zip.OpenReader(inputPath)
		if err != nil {
			errors <- err
			close(processedPaths)
			close(errors)
			return
		}
		defer func() {
			if err := zr.Close(); err != nil {
				fmt.Printf("defer zr.Close() error:%+v\n", err)
			}
		}()

		outputPath, err = filepath.Abs(outputPath)
		if err != nil {
			errors <- err
			close(processedPaths)
			close(errors)
			return
		}

		for _, f := range zr.File {
			select {
			case <-cancel:
				errors <- fmt.Errorf("AsynUnzip canceled")
				close(processedPaths)
				close(errors)
				return
			default:
			}
			err := removeFromZip(f, outputPath, permDir)
			processedPaths <- filepath.Join(outputPath, f.Name)
			if err != nil {
				errors <- err
			}
		}

		close(processedPaths)
		close(errors)
	}()

	return cancel, processedPaths, errors
}

// AsyncZip asynchronously creates a compressed ZIP file, of file/directories in paths, at zipPath.
// The paths are turned into absolute paths, then made relative by removing the leading
// filepath.Separator. When zipped, if trimFilepath !=nil, all strings in trimFilepath are left
// trimmed from paths to create relative paths. Progress can be monitored via the returned channels,
// which return cancel, processed paths, and any errors. The cancel channel can be used to cancel an
// operation. The operation is complete when both processed paths and errors channels are closed.
func AsyncZip(zipPath string, paths []string, trimFilepath []string) (chan<- bool, <-chan string, <-chan error) {
	cancel := make(chan bool, 1)
	// Size channels so that they don't block if the caller is only checking done.
	processedPaths := make(chan string, len(paths))
	errors := make(chan error, len(paths))
	go func() {
		f, err := os.Create(zipPath)
		if err != nil {
			errors <- err
			close(processedPaths)
			close(errors)
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Printf("f.Close() error:%+v\n", err)
			}
		}()

		zipWriter := zip.NewWriter(f)
		defer func() {
			//nolint:errcheck
			// The error will always be "zip: writer closed twice", which is not typed and not useful to log.
			zipWriter.Close()
		}()

		for _, path := range paths {
			select {
			case <-cancel:
				errors <- fmt.Errorf("AsyncZip canceled")
				close(processedPaths)
				close(errors)
				return
			default:
			}
			err := filepath.WalkDir(path, addToZip(zipWriter, trimFilepath))
			processedPaths <- path
			if err != nil {
				errors <- err
			}
		}

		if err := zipWriter.Close(); err != nil {
			fmt.Printf("zipWriter.Close error:%+v\n", err)
		}
		close(processedPaths)
		close(errors)
	}()

	return cancel, processedPaths, errors
}

// GetZipStats is for getting statistics on a zip file; currently only
// supports the number of zip.File in an archive.
func GetZipStats(inputPath string) (*ZipStats, error) {
	zr, err := zip.OpenReader(inputPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := zr.Close(); err != nil {
			fmt.Printf("defer zr.Close() error:%+v\n", err)
		}
	}()

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
func addToZip(zipWriter *zip.Writer, trimFilepath []string) func(string, fs.DirEntry, error) error {
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
		zipFilepath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		if trimFilepath != nil {
			for _, trm := range trimFilepath {
				if strings.HasPrefix(zipFilepath, trm) {
					header.Name = strings.TrimPrefix(zipFilepath, trm)
					break
				}
			}
		} else {
			header.Name = zipFilepath
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
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Printf("defer f.Close() error:%+v\n", err)
			}
		}()

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
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("defer f.Close() error:%+v\n", err)
		}
	}()

	irc, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := irc.Close(); err != nil {
			fmt.Printf("defer irc.Close() error:%+v\n", err)
		}
	}()

	if _, err := io.Copy(f, irc); err != nil {
		return err
	}
	return nil
}
