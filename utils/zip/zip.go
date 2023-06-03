package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func zipDir(dirPath string, zipWriter *zip.Writer, basePath string) error {
	fileInfos, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		path := filepath.Join(dirPath, fileInfo.Name())
		if fileInfo.IsDir() {
			err := zipDir(path, zipWriter, basePath)
			if err != nil {
				return err
			}
		} else {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			relPath, err := filepath.Rel(basePath, path)
			if err != nil {
				return err
			}

			zipEntry, err := zipWriter.Create(filepath.ToSlash(relPath))
			if err != nil {
				return err
			}

			_, err = io.Copy(zipEntry, file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Zip(sourceDir, targetZip string) error {
	zipFile, err := os.Create(targetZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	basePath := filepath.Dir(sourceDir)

	err = zipDir(sourceDir, zipWriter, basePath)
	if err != nil {
		return err
	}

	return nil
}

func Unzip(zipFile, destDir string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(destDir, file.Name)

		if file.FileInfo().IsDir() {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			return err
		}

		zipEntry, err := file.Open()
		if err != nil {
			return err
		}
		defer zipEntry.Close()

		targetFile, err := os.Create(path)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, zipEntry)
		if err != nil {
			return err
		}
	}

	return nil
}
