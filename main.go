package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return nil
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	var dir, output string
	if runtime.GOOS == "windows" {
		dir = filepath.Join(wd, "temp")
		output = filepath.Join(wd, "magic-output")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error reading home directory:", err)
			return
		}
		dir = filepath.Join(homeDir, "Desktop", "workspace", "magic", "temp")
		output = filepath.Join(homeDir, "Desktop", "workspace", "magic", "magic-output")
	}

	// Error handling for reading directory
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Ask for user input for the date
	var inputDate string
	fmt.Println("Enter the date (YYYY-MM-DD):")
	if _, err := fmt.Scanln(&inputDate); err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Parse the input date
	date, err := time.Parse("2006-01-02", inputDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	formattedDate := date.Format("060102")

	f := make(map[string]map[string]string)
	for _, file := range files {
		filename := file.Name()
		size := filename[0:1]
		extension := filepath.Ext(filename)
		name := strings.TrimSuffix(filename[1:], extension)

		if _, ok := f[name]; !ok {
			f[name] = make(map[string]string)
		}
		f[name][size] = filepath.Join(dir, filename)
	}

	if err := os.RemoveAll(output); err != nil {
		fmt.Println("Error removing output directory:", err)
		return
	}
	if err := os.Mkdir(output, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	for name, sizes := range f {
		for s, path := range sizes {
			if len(name) < 2 {
				name = "0" + name
			}
			fileName := filepath.Join(output, fmt.Sprintf("%s自制-vid%s-%s.mp4", formattedDate, name, s))
			if err := copyFile(path, fileName); err != nil {
				fmt.Println("Error copying file:", err)
				return
			}
		}
	}
	fmt.Println("Done!")
}
