package file_service

import (
	"archive/zip"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func GetCRC32(filePath string) (uint32, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	hash := crc32.NewIEEE()
	if _, err := io.Copy(hash, file); err != nil {
		return 0, err
	}
	checksum := hash.Sum32()
	return checksum, nil
}

func CreateDirectoryIfNotExists(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
		fmt.Println("Directory created:", dirPath)
	}
	return err
}

func IsDirectoryExists(dirPath string) bool {
	_, err := os.Stat(dirPath)
	dirExists := os.IsExist(err)
	return dirExists

}

func Unzip(zipPath string, output string) {
	archive, err := zip.OpenReader(zipPath)
	if err != nil {
		panic(err)
	}
	defer archive.Close()
	for _, f := range archive.File {
		filePath := filepath.Join(output, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(output)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}
