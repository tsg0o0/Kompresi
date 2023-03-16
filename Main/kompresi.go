package main

import (
	"fmt"
	"runtime"
	"os"
	"os/exec"
	"image"
	_ "image/png"
	_ "image/jpeg"
	"path/filepath"
)

func main() {
	exedir, err := os.Executable()
	if err != nil {
		fmt.Println("\x1b[33mFatal error: Could not obtain the location of the executable file.\x1b[0m")
		fmt.Println("\x1b[33m", err, "\x1b[0m")
		os.Exit(1)
	}
	exedir = filepath.Dir(exedir)
	fmt.Println("EXE Directory: ", exedir)
	
	//bin check
	if runtime.GOOS == "darwin" {
		_, err = os.Stat(exedir + "/resources/mac/guetzli")
		_, err = os.Stat(exedir + "/resources/mac/zopflipng")
	}else if runtime.GOOS == "linux" {
		_, err = os.Stat(exedir + "/resources/linux/guetzli")
		_, err = os.Stat(exedir + "/resources/linux/zopflipng")
	}else if runtime.GOOS == "windows" {
		_, err = os.Stat(exedir + "/resources/win/guetzli")
		_, err = os.Stat(exedir + "/resources/win/zopflipng")
	}else{
		fmt.Println("\x1b[31mFatal error: This operating system could not be recognized.\x1b[0m")
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("\x1b[31mFatal error: Binary required for execution not found.\x1b[0m")
		fmt.Println("\x1b[31m", err, "\x1b[0m")
		os.Exit(1)
	}
	
	arg := os.Args
	fmt.Println(arg)
	if len(arg) == 1 {
		fmt.Println("\x1b[32mBooting Daemon...\x1b[0m")
	}else{
		imgCatch(arg[1])
	}
}

func imgCatch(inputFile string) {
	file, _ := os.Open(inputFile)
	defer file.Close()
	_, format, err := image.DecodeConfig(file)
	if err != nil {
    	fmt.Println("Error: ", err)
	}else{
		if format == "png" {
			//PNG
			pngCompress(inputFile)
		}else if format == "jpeg" {
			//JPEG
			jpegCompress(inputFile)
		}
	}
}

func pngCompress(inputFile string) {
	originalInfo, _ := os.Stat(inputFile)
	fmt.Println("Compressing... (by zopfli)")
	exedir, _ := os.Executable()
	exedir = filepath.Dir(exedir)
	
	cmd := exec.Command("zopflipng", "-m", "-y", inputFile, inputFile)
	// run zopflipng
	if runtime.GOOS == "darwin" {
		cmd = exec.Command(exedir + "/resources/mac/zopflipng", "-m", "-y", inputFile, inputFile)
	}else if runtime.GOOS == "linux" {
		cmd = exec.Command(exedir + "/resources/linux/zopflipng", "-m", "-y", inputFile, inputFile)
	}else if runtime.GOOS == "windows" {
		cmd = exec.Command(exedir + "/resources/win/zopflipng", "-m", "-y", inputFile, inputFile)
	}
	
	//RUN
	exeerr := cmd.Run()
	if exeerr != nil {
		//Failed
		fmt.Println("\x1b[33mFailed. (by zopfli): ", exeerr, "\x1b[0m")
	}else{
		//Success
		resultInfo, _ := os.Stat(inputFile)
		fmt.Println("\x1b[32mSuccess. (by zopfli)\x1b[0m")
		fmt.Println("Original file size:", originalInfo.Size())
		fmt.Println("Result file size:", resultInfo.Size())
		fmt.Println("Percentage of original", ( ( 100 * resultInfo.Size() ) / originalInfo.Size() ), "%")
	}
	
}

func jpegCompress(inputFile string) {
	originalInfo, _ := os.Stat(inputFile)
	fmt.Println("Compressing... (by guetzli)")
	exedir, _ := os.Executable()
	exedir = filepath.Dir(exedir)
	
	cmd := exec.Command("guetzli", inputFile, inputFile)
	// run guetzli
	if runtime.GOOS == "darwin" {
		cmd = exec.Command(exedir + "/resources/mac/guetzli", inputFile, inputFile)
	}else if runtime.GOOS == "linux" {
		cmd = exec.Command(exedir + "/resources/linux/guetzli", inputFile, inputFile)
	}else if runtime.GOOS == "windows" {
		cmd = exec.Command(exedir + "/resources/win/guetzli", inputFile, inputFile)
	}
	
	//RUN
	exeerr := cmd.Run()
	if exeerr != nil {
		//Failed
		fmt.Println("\x1b[33mFailed. (by guetzli): ", exeerr, "\x1b[0m")
	}else{
		//Success
		resultInfo, _ := os.Stat(inputFile)
		fmt.Println("Success. (by guetzli)")
		fmt.Println("Original file size:", originalInfo.Size())
		fmt.Println("Result file size:", resultInfo.Size())
		fmt.Println("Percentage of original", ( ( 100 * resultInfo.Size() ) / originalInfo.Size() ), "%")
	}
	
}

