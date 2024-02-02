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

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	/*
		创建一个对象 f，{fileName: {z: path, s: path}}
		遍历文件列表
		获取当前文件的文件名
		对文件名截取第一个字符和后面的字符，第一个字符是文件的尺寸，后面的字符是文件的名称
		在 f 中添加一个对象，对象的 key 是文件的名称，value 是一个对象，对象的 key 是文件的尺寸，value 是文件的路径
		创建一个新的文件夹
		遍历 f
		获取当前文件的文件名
		遍历 f 中的对象
		获取当前文件的文件尺寸
		获取当前文件的文件路径
		复制这个文件到新的文件夹中并且重命名
	*/

	wd, _ := os.Getwd()

	f := make(map[string]map[string]string)

	var dir string
	var output string
	if runtime.GOOS == "windows" {
		dir = filepath.Join(wd, "/temp")
		output = filepath.Join(wd, "/output")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return
		}

		dir = filepath.Join(homeDir, "/Desktop/workspace/file-rename/temp")
		output = filepath.Join(homeDir, "/Desktop/workspace/file-rename/output")
	}

	files, err := os.ReadDir(dir)

	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		filename := file.Name()
		size := filename[0:1]

		extension := filepath.Ext(filename)
		name := filename[1:]
		name = strings.Replace(name, extension, "", -1)

		if _, ok := f[name]; !ok {
			f[name] = make(map[string]string)
		}
		// f[name][size] = fmt.Sprintf("%s/%s", dir, filename)
		f[name][size] = filepath.Join(dir, filename)
	}

	os.RemoveAll(output)
	os.Mkdir(output, os.ModePerm)

	for name, size := range f {
		for s, path := range size {
			date := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1)
			if time.Now().Weekday() == time.Sunday {
				date = date.AddDate(0, 0, -6)
			}

			formattedDate := date.Format("2006-01-02")
			if len(name) < 2 {
				name = fmt.Sprintf("0%s", name)
			}

			fileName := filepath.Join(output, fmt.Sprintf("%s自制-vid%s-%s.mp4", formattedDate, name, s))

			err = copyFile(path, fileName)
			if err != nil {
				fmt.Println("Error copying file:", err)
				return
			}
		}
	}
	fmt.Println("Done!")

	// dir, _ := os.ReadDir("./temp")
	// for _, file := range dir {
	// 	fmt.Println(file.Name())
	// }

	// // Replace with the path to your video file
	// videoPath := "./assets/before/1月22日 (2).mp4"

	// // FFmpeg command to get video resolution
	// cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", videoPath)

	// // Run the command
	// output, err := cmd.Output()
	// if err != nil {
	// 	fmt.Println("Error executing command:", err)
	// 	return
	// }

	// fmt.Println(string(output))
	// // Print the resolution
	// resolution := strings.TrimSpace(string(output))
	// fmt.Println("Resolution:", resolution)

	// create z1-z10 and s1-s10 files in the temp folder
	// for i := 1; i <= 10; i++ {
	// 	// create z1-z10
	// 	z := fmt.Sprintf("z%d", i)
	// 	zFile, err := os.Create(fmt.Sprintf("./temp/%s.mp4", z))
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	defer zFile.Close()

	// 	// create s1-s10
	// 	s := fmt.Sprintf("s%d", i)
	// 	sFile, err := os.Create(fmt.Sprintf("./temp/%s.mp4", s))
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	defer sFile.Close()
	// }
}
