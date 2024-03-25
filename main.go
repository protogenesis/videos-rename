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

func stop() {
	fmt.Println("Press 'Enter' to exit...")
	fmt.Scanln()
	os.Exit(1)
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		stop()
	}

	var dir, output string
	if runtime.GOOS == "windows" {
		dir = filepath.Join(wd, "temp")
		output = filepath.Join(wd, "magic-output")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error reading home directory:", err)
			stop()
		}
		dir = filepath.Join(homeDir, "Desktop", "workspace", "magic", "temp")
		output = filepath.Join(homeDir, "Desktop", "workspace", "magic", "magic-output")
	}

	// Error handling for reading directory
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		stop()
	}

	// Ask for user input for the date
	var inputDate string
	var date time.Time
	fmt.Println("Enter the date (YYYY-MM-DD):")
	if _, err := fmt.Scanln(&inputDate); err != nil {
		date = time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1)
		if time.Now().Weekday() == time.Sunday {
			date = time.Now().AddDate(0, 0, -6)
		}
	} else {
		// Parse the input date
		date, err = time.Parse("2006-01-02", inputDate)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			stop()
		}
	}

	formattedDate := date.Format("060102")

	f := make(map[string]map[string]string)
	for _, file := range files {
		filename := file.Name()
		size := filename[0:1]
		extension := filepath.Ext(filename)
		name := strings.TrimSuffix(filename[2:], extension)

		if _, ok := f[name]; !ok {
			f[name] = make(map[string]string)
		}
		f[name][size] = filepath.Join(dir, filename)
	}

	if err := os.RemoveAll(output); err != nil {
		fmt.Println("Error removing output directory:", err)
		stop()
	}
	if err := os.Mkdir(output, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory:", err)
		stop()
	}

	for name, sizes := range f {
		for s, path := range sizes {
			if len(name) < 2 {
				name = "0" + name
			}
			fileName := filepath.Join(output, fmt.Sprintf("%s自制-vid%s-%s.mp4", formattedDate, name, s))
			if err := copyFile(path, fileName); err != nil {
				fmt.Println("Error copying file:", err)
				stop()
			}
		}
	}
	fmt.Println("Done!")
}
