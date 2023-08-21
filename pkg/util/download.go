package util

import (
	"errors"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
)

func DownloadAndMoveToLocation(url, dir, fileName string) error {
	if string(dir[len(dir)-1]) == "/" {
		dir = dir[:len(dir)-1]
	}
	path := fmt.Sprintf("%s/%s", dir, fileName)

	tmpFile, err := os.CreateTemp("", "bpmtempfile.*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	if err = DownloadFile(url, tmpFile, fileName); err != nil {
		deleteFile(tmpFile)
		return fmt.Errorf("failed to download file: %w", err)
	}
	if err = tmpFile.Close(); err != nil {
		deleteFile(tmpFile)
		return fmt.Errorf("unable to close temp file: %w", err)
	}

	if err = MoveFile(tmpFile.Name(), path); err != nil {
		deleteFile(tmpFile)
		return fmt.Errorf("unable to copy downloaded file in place: %w", err)
	}

	return nil
}

func DownloadFile(url string, file *os.File, fileName string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer handleErrResponse(resp.Body.Close, "unable to close response body")

	f, err := os.OpenFile(file.Name(), os.O_CREATE|os.O_WRONLY, 0755)
	defer handleErrResponse(f.Close, "unable to close temporary file")

	bar := progressbar.DefaultBytes(resp.ContentLength, fileName)
	if _, err = io.Copy(io.MultiWriter(f, bar), resp.Body); err != nil {
		return errors.Join(errors.New(fmt.Sprintf("unable to write file %s", f.Name())), err)
	}

	return nil
}

func MoveFile(source, destination string) error {
	if !fileExists(source) {
		return errors.New(fmt.Sprintf("source file %s does not exist", source))
	}
	if fileExists(destination) {
		return errors.New(fmt.Sprintf("destination file %s already exists", destination))
	}

	src, err := os.OpenFile(source, os.O_RDWR, 0755)
	if err != nil {
		return fmt.Errorf("failed to open file for copying: %w", err)
	}
	defer deleteFile(src)

	dst, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer handleErrResponse(dst.Close, "failed to close destination file")

	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file: %w")
	}

	return nil
}

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func deleteFile(file *os.File) {
	_ = file.Close()

	if err := os.Remove(file.Name()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to delete file %s: %s\n", file.Name(), err.Error())
	}
}

func handleErrResponse(f func() error, errMsg string) {
	if err := f(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", errMsg, err.Error())
	}
}
