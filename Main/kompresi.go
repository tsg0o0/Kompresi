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

var exedir = ""

func main() {
	arg := os.Args[1]
	exedir, _ := os.Executable()
	exedir = filepath.Dir(exedir)
	fmt.Println("EXE Directory: ", exedir)
	imgCatch(arg)
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
		fmt.Println("Failed. (by zopfli): ", exeerr)
	}else{
		//Success
		resultInfo, _ := os.Stat(inputFile)
		fmt.Println("Success. (by zopfli)")
		fmt.Println("Original file size:", originalInfo.Size())
		fmt.Println("Result file size:", resultInfo.Size())
		fmt.Println("Percentage of original", ( ( 100 * resultInfo.Size() ) / originalInfo.Size() ), "%")
	}
	
}

func jpegCompress(inputFile string) {
	originalInfo, _ := os.Stat(inputFile)
	fmt.Println("Compressing... (by guetzli)")
	
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
		fmt.Println("Failed. (by guetzli): ", exeerr)
	}else{
		//Success
		resultInfo, _ := os.Stat(inputFile)
		fmt.Println("Success. (by guetzli)")
		fmt.Println("Original file size:", originalInfo.Size())
		fmt.Println("Result file size:", resultInfo.Size())
		fmt.Println("Percentage of original", ( ( 100 * resultInfo.Size() ) / originalInfo.Size() ), "%")
	}
	
}

