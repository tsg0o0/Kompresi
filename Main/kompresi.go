package main

import (
	"fmt"
	"runtime"
	"os"
	"os/exec"
	"encoding/json"
	"io/ioutil"
	"strings"
	"image"
	_ "image/png"
	_ "image/jpeg"
	"path/filepath"
	"github.com/fsnotify/fsnotify"
)

type Config struct {
    Version      int    `json:"version"`
    InputDir     string `json:"inputDir"`
    OutputDir    string `json:"outputDir"`
    DeleteOrigin bool   `json:"deleteOrigin"`
    OptimLv      int   `json:"optimLv"`
    //0-Low 2-High
}
var config Config

func main() {
	exedir, err := os.Executable()
	if err != nil {
		fmt.Println("\x1b[33mFatal error: Could not obtain the location of the executable file.\x1b[0m")
		fmt.Println("\x1b[33m", err, "\x1b[0m")
		os.Exit(1)
	}
	exedir = filepath.Dir(exedir)
	
	//bin check
	if runtime.GOOS == "darwin" {
		_, err = os.Stat(exedir + "/resources/mac/guetzli")
		_, err = os.Stat(exedir + "/resources/mac/zopflipng")
	}else if runtime.GOOS == "linux" {
		_, err = os.Stat(exedir + "/resources/linux/guetzli")
		_, err = os.Stat(exedir + "/resources/linux/zopflipng")
	}else if runtime.GOOS == "windows" {
		_, err = os.Stat(exedir + "/resources/win/guetzli.exe")
		_, err = os.Stat(exedir + "/resources/win/zopflipng.exe")
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
	if len(arg) == 1 {
		//Run daemon
		fmt.Println("\x1b[32mBooting Daemon...\x1b[0m")
		loadConfig(false)
		watcherDaemon()
	}else if len(arg) > 2 {
		loadConfig(true)
		//Edit config
		if arg[1] == "inputDir" {
			if arg[2] == config.OutputDir {
				fmt.Println("\x1b[31mThe input and output directories cannot be the same.\x1b[0m")
			}else{
				config.InputDir = arg[2]
			}
		}else if arg[1] == "outputDir" {
			if arg[2] == config.InputDir {
				fmt.Println("\x1b[31mThe input and output directories cannot be the same.\x1b[0m")
			}else{
				config.OutputDir = arg[2]
			}
		}else if arg[1] == "deleteOrigin" {
			if arg[2] == "Yes" || arg[2] == "y" || arg[2] == "true" {
				config.DeleteOrigin = true
			}else if arg[2] == "No" || arg[2] == "n" || arg[2] == "false" {
				config.DeleteOrigin = false
			}else{
				fmt.Println("\x1b[31mThe only values that disappear with the input are Yes or No.\x1b[0m")
			}
		}else if arg[1] == "optimLv" {
			if arg[2] == "0" {
				config.OptimLv = 0
			}else if arg[2] == "1" {
				config.OptimLv = 1
			}else if arg[2] == "2" {
				config.OptimLv = 2
			}else{
				fmt.Println("\x1b[31mThe only values that disappear with the input are Yes or No.\x1b[0m")
			}
		}
		configJSON, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			fmt.Println("\x1b[31mError encoding config file\x1b[0m")
			fmt.Println("\x1b[31m", err, "\x1b[0m")
			os.Exit(1)
		}
		if err := ioutil.WriteFile("config.json", configJSON, 0644); err != nil {
			fmt.Println("\x1b[31mError writing config file\x1b[0m")
			fmt.Println("\x1b[31m", err, "\x1b[0m")
		}
		fmt.Println("\n\x1b[32mConfig file updated.\x1b[0m")
		os.Exit(1)
	}else{
		if arg[1] == "help" {
			fmt.Println("\n\x1b[35m==Change Settings==\x1b[0m")
			fmt.Println("Use the following command or rewrite the JSON file directly to complete the setup.")
			fmt.Println("\nCommand (argument):")
			fmt.Println("\n	inputDir      'YOUR INPUT DIRECTRY PATH'")
			fmt.Println("	- Select the directory to load the images.")
			fmt.Println("\n	outputDir     'YOUR OUTPUT DIRECTRY PATH'")
			fmt.Println("	- Select a directory to output compressed images.")
			fmt.Println("\n	deleteOrigin  'Yes or No'")
			fmt.Println("	- Delete original files after compression.")
			fmt.Println("\n	optimLv  '0 - 2'")
			fmt.Println("	- Change the compression level.")
			fmt.Println("	  0: Fast but low compression")
			fmt.Println("	  1: Auto")
			fmt.Println("	  2: Slow but high compression")
			fmt.Println("\n\x1b[35m==Compress images by themselves==\x1b[0m")
			fmt.Println("\n	Argument: 'YOUR INPUT IMAGE PATH'")
			fmt.Println("\n\x1b[32m==Starts the daemon with no arguments!==\x1b[0m")
			os.Exit(1)
		}else if arg[1] == "license" {
			fmt.Println("\n\x1b[32m==Kompresi by tsg0o0==\x1b[0m")
			fmt.Println("\nGo application for lossless compression of PNG and JPEG images.")
			fmt.Println("\nThis software is licensed under the terms of the Mozilla Public License 2.0.")
			fmt.Println("(https://www.mozilla.org/en-US/MPL/2.0/)")
			fmt.Println("")
			os.Exit(1)
		}else{
			//Compress Image
			imgCatch(arg[1])
		}
	}
}

func watcherDaemon() {
	//watch dir
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error creating watcher:", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add(config.InputDir)
	if err != nil {
		fmt.Println("Error adding watch:", err)
		return
	}

	//event
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Println("File Detect:", event.Name)
				imgCatch(event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Error:", err)
		}
	}
}

func loadConfig(ignore bool) {
	exedir, _ := os.Executable()
	exedir = filepath.Dir(exedir)
	
	fmt.Println("Loading config...")
	configFile, err := os.Open(exedir + "/config.json")
		if err != nil {
			fmt.Println("\x1b[31mFatal error: Config file Cannot found.\x1b[0m")
			fmt.Println("\x1b[31m", err, "\x1b[0m")
			os.Exit(1)
		}
		defer configFile.Close()
		jsonParser := json.NewDecoder(configFile)
		if err := jsonParser.Decode(&config); err != nil {
			fmt.Println("\x1b[31mFatal error: Error decoding config file.\x1b[0m")
			fmt.Println("\x1b[31m", err, "\x1b[0m")
			os.Exit(1)
		}
		
		if ignore == false {
		fmt.Println("Input directory:", config.InputDir)
		fmt.Println("Output directory:", config.OutputDir)
		fmt.Println("Delete original files:", config.DeleteOrigin)
		fmt.Println("Optimize Level:", config.OptimLv)
		}
		
		if config.InputDir == "" && config.OutputDir == "" && ignore == false {
			fmt.Println("\n\x1b[35m==Complete the setup!==\x1b[0m")
			fmt.Println("Use the following command or rewrite the JSON file directly to complete the setup.")
			fmt.Println("\nCommand:")
			fmt.Println("\n	inputDir      'YOUR INPUT DIRECTRY PATH'")
			fmt.Println("	- Select the directory to load the images.")
			fmt.Println("\n	outputDir     'YOUR OUTPUT DIRECTRY PATH'")
			fmt.Println("	- Select a directory to output compressed images.")
			fmt.Println("\n	deleteOrigin  'Yes or No'")
			fmt.Println("	- Delete original files after compression.")
			fmt.Println("\n	optimLv  '0 - 2'")
			fmt.Println("	- Change the compression level.")
			fmt.Println("	  0: Fast but low compression")
			fmt.Println("	  1: Auto")
			fmt.Println("	  2: Slow but high compression")
			fmt.Println("\n\x1b[35mPlease change the settings and try again.\x1b[0m")
			os.Exit(1)
		}
		if config.InputDir == config.OutputDir && ignore == false {
			fmt.Println("\x1b[31mThe input and output directories cannot be the same.\x1b[0m")
			os.Exit(1)
		}
}

func imgCatch(inputFile string) {
	file, _ := os.Open(inputFile)
	defer file.Close()
	_, format, err := image.DecodeConfig(file)
	if err != nil {
    	fmt.Println("\x1b[31mFile Error: ", err, "\x1b[0m")
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
	
	outputFile := strings.Replace(inputFile, config.InputDir, config.OutputDir, 1)
	
	optimArg := "--iterations=1"
	if config.OptimLv == 0 {
		optimArg = "--iterations=1"
	}else if config.OptimLv == 1 {
		if originalInfo.Size() >= 524288 {
			optimArg = "--iterations=1"
		}else if originalInfo.Size() >= 65536 {
			optimArg = "--iterations=5"
		}else if originalInfo.Size() >= 8192 {
			optimArg = "--iterations=10"
		}else{
			optimArg = "--iterations=15"
		}
	}else if config.OptimLv == 2 {
		optimArg = "--iterations=15"
	}
	fmt.Println("OptimArg:", optimArg)
	
	cmd := exec.Command("zopflipng", "-m", "-y", inputFile, inputFile)
	// run zopflipng
	if runtime.GOOS == "darwin" {
		cmd = exec.Command(exedir + "/resources/mac/zopflipng", optimArg, "-y", inputFile, outputFile)
	}else if runtime.GOOS == "linux" {
		cmd = exec.Command(exedir + "/resources/linux/zopflipng", optimArg, "-y", inputFile, outputFile)
	}else if runtime.GOOS == "windows" {
		cmd = exec.Command(exedir + "/resources/win/zopflipng.exe", optimArg, "-y", inputFile, outputFile)
	}
	
	//RUN
	exeerr := cmd.Run()
	if exeerr != nil {
		//Failed
		fmt.Println("\x1b[33mFailed. (by zopfli): ", exeerr, "\x1b[0m")
	}else{
		//Success
		if config.DeleteOrigin == true {
			os.Remove(inputFile)
		}
		resultInfo, _ := os.Stat(outputFile)
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
		cmd = exec.Command(exedir + "/resources/win/guetzli.exe", inputFile, inputFile)
	}
	
	//RUN
	exeerr := cmd.Run()
	if exeerr != nil {
		//Failed
		fmt.Println("\x1b[33mFailed. (by guetzli): ", exeerr, "\x1b[0m")
	}else{
		//Success
		resultInfo, _ := os.Stat(inputFile)
		fmt.Println("\x1b[32mSuccess. (by guetzli)\x1b[0m")
		fmt.Println("Original file size:", originalInfo.Size())
		fmt.Println("Result file size:", resultInfo.Size())
		fmt.Println("Percentage of original", ( ( 100 * resultInfo.Size() ) / originalInfo.Size() ), "%")
	}
	
}

