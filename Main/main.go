package main

import (
	"fmt"
	"runtime"
	"os"
	"os/exec"
	"image"
	_ "image/png"
	_ "image/jpeg"
)

func main() {
	arg := os.Args[1]
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
		}
	}
}

func pngCompress(inputFile string) {
	originalInfo, _ := os.Stat(inputFile)
	fmt.Println("Compressing... (by zopfli)")
	
	cmd := exec.Command("zopflipng", "-m", "-y", inputFile, inputFile)
	
	// run zopflipng
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("./resources/mac/zopflipng", "-m", "-y", inputFile, inputFile)
	}else if runtime.GOOS == "linux" {
		cmd = exec.Command("./resources/linux/zopflipng", "-m", "-y", inputFile, inputFile)
	}else if runtime.GOOS == "windows" {
		cmd = exec.Command("./resources/win/zopflipng", "-m", "-y", inputFile, inputFile)
	}
	err := cmd.Run()
	if err != nil {
		//Failed
		fmt.Println("Failed. (by zopfli): ", err)
	}else{
		//Success
		resultInfo, _ := os.Stat(inputFile)
		fmt.Println("Success. (by zopfli)")
		fmt.Println("Original file size:", originalInfo.Size())
		fmt.Println("Result file size:", resultInfo.Size())
		fmt.Println("Percentage of original", ( ( 100 * resultInfo.Size() ) / originalInfo.Size() ), "%")
	}
	
}

