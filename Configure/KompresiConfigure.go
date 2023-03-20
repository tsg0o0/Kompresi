package main

import (
	"bufio"
	"fmt"
	"runtime"
	"encoding/json"
	"io/ioutil"
	"strings"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	Version      int    `json:"version"`
	InputDir     string `json:"inputDir"`
	OutputDir    string `json:"outputDir"`
	DeleteOrigin bool   `json:"deleteOrigin"`
	OptimLv      int    `json:"optimLv"`
}

func main() {
	exedir, err := os.Executable()
	if err != nil {
		fmt.Println("Fatal error: Could not obtain the location of the executable file.")
		fmt.Println("", err, "")
	}
	exedir = filepath.Dir(exedir)
	
	//load cnfig.json
	config, err := readConfig(exedir + "/config.json")
	if err != nil {
		log.Fatal(err)
	}

	//make app
	a := app.New()
	w := a.NewWindow("Kompresi Configure")
	
	//widgets
	inputDir := widget.NewEntry()
	outputDir := widget.NewEntry()
	deleteOrigin := widget.NewCheck("Delete original file", nil)
	optimLv := widget.NewRadioGroup([]string{"Fast but low compression", "Auto", "Slow but high compression"}, nil)
	optimLvInt := 1

	//set widget default value
	inputDir.SetText(config.InputDir)
	outputDir.SetText(config.OutputDir)
	deleteOrigin.SetChecked(config.DeleteOrigin)
	if config.OptimLv == 0 {
		optimLv.SetSelected("Fast but low compression")
	}else if config.OptimLv == 1 {
		optimLv.SetSelected("Auto")
	}else if config.OptimLv == 2 {
		optimLv.SetSelected("Slow but high compression")
	}
	
	//save button
	saveButton := widget.NewButton("Save", func() {
		config.InputDir = inputDir.Text
		config.OutputDir = outputDir.Text
		config.DeleteOrigin = deleteOrigin.Checked
		if optimLv.Selected == "Fast but low compression" {
			optimLvInt = 0
		}else if optimLv.Selected == "Auto" {
			optimLvInt = 1
		}else if optimLv.Selected == "Slow but high compression" {
			optimLvInt = 2
		}else{
			optimLvInt = 1
		}
		config.OptimLv = optimLvInt
		
		if inputDir.Text == "" {
			dialog.ShowError(fmt.Errorf("Input Directory is not entered."), w)
		}else if outputDir.Text == "" {
			dialog.ShowError(fmt.Errorf("Output Directory is not entered."), w)
		}else if outputDir.Text == inputDir.Text {
			dialog.ShowError(fmt.Errorf("Input Directory and Output Directory cannot be the same."), w)
		}else{
			//output config.json
			err := writeConfig(exedir + "/config.json", config)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Failed to save."), w)
			}else{
				dialog.ShowInformation("Saved!", "Setup is complete.", w)
			}
		}
	})
	
	inputSelectFolder := widget.NewButton("Select Folder", func() {
        dialog := dialog.NewFolderOpen(func(selected fyne.ListableURI, err error) {
            if err == nil {
            	if selected != nil {
            		if runtime.GOOS == "windows" {
            			inputDir.SetText(strings.Replace(selected.Path(), "/", "\\", -1))
            		}else{
                		inputDir.SetText(selected.Path())
                	}
                }
            } else {
                dialog.ShowError(fmt.Errorf("Could not retrieve the folder."), w)
            }
        }, w)
        dialog.Show()
        dialog.Resize(fyne.NewSize(1400, 1400))
    })
    outputSelectFolder := widget.NewButton("Select Folder", func() {
        dialog := dialog.NewFolderOpen(func(selected fyne.ListableURI, err error) {
            if err == nil {
            	if selected != nil {
            		if runtime.GOOS == "windows" {
            			outputDir.SetText(strings.Replace(selected.Path(), "/", "\\", -1))
            		}else{
                		outputDir.SetText(selected.Path())
                	}
                }
            } else {
                dialog.ShowError(fmt.Errorf("Could not retrieve the folder."), w)
            }
        }, w)
        dialog.Show()
        dialog.Resize(fyne.NewSize(1400, 1400))
    })
	
	//RUN
	logWin := a.NewWindow("Kompresi Daemon")
    text := widget.NewTextGrid()
    text.ShowWhitespace = true
    txScroll := container.NewScroll(text)
	runButton := widget.NewButton("Run", func() {
		config.InputDir = inputDir.Text
		config.OutputDir = outputDir.Text
		config.DeleteOrigin = deleteOrigin.Checked
		if optimLv.Selected == "Fast but low compression" {
			optimLvInt = 0
		}else if optimLv.Selected == "Auto" {
			optimLvInt = 1
		}else if optimLv.Selected == "Slow but high compression" {
			optimLvInt = 2
		}else{
			optimLvInt = 1
		}
		config.OptimLv = optimLvInt
		
		if inputDir.Text == "" {
			dialog.ShowError(fmt.Errorf("Input Directory is not entered."), w)
		}else if outputDir.Text == "" {
			dialog.ShowError(fmt.Errorf("Output Directory is not entered."), w)
		}else if outputDir.Text == inputDir.Text {
			dialog.ShowError(fmt.Errorf("Input Directory and Output Directory cannot be the same."), w)
		}else{
			//output config.json
			err := writeConfig(exedir + "/config.json", config)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Failed to save."), w)
			}else{
				
				
		
		exedir, _ := os.Executable()
		exedir = filepath.Dir(exedir)
		
		cmd := exec.Command("kompresi")
		if runtime.GOOS == "darwin" {
			cmd = exec.Command(exedir + "/kompresi")
		}else if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd.exe", "/c", exedir + "/kompresi.exe", "daemon")
		}else{
			fmt.Println("Fatal error: This operating system could not be recognized.")
		}
		
		outPipe, err := cmd.StdoutPipe()
        if err != nil {
            fmt.Println(err)
            return
        }

        errPipe, err := cmd.StderrPipe()
        if err != nil {
            fmt.Println(err)
            return
        }

        go func() {
            outScanner := bufio.NewScanner(outPipe)
            for outScanner.Scan() {
                text.SetText(text.Text() + "\n " + outScanner.Text())
                txScroll.ScrollToBottom()
            }

            errScanner := bufio.NewScanner(errPipe)
            for errScanner.Scan() {
                text.SetText(text.Text() + "\n " + errScanner.Text())
                txScroll.ScrollToBottom()
            }
        }()

        exeerr := cmd.Start()
        if exeerr != nil {
        	dialog.ShowError(fmt.Errorf("Could not start the application."), w)
            fmt.Println(exeerr)
            return
        }

        go func() {
            cmd.Wait()
            execErr := cmd.Run()
            if execErr != nil {
                fmt.Println(execErr)
                return
            }
        }()
    
        //exit
        stopButton := widget.NewButton("Stop", func() {
        	killErr := cmd.Process.Kill()
            if killErr != nil {
                fmt.Println(killErr)
            }
            os.Exit(0)
    	})
        
        header := container.NewGridWithColumns(2, widget.NewLabel("Kompresi"), stopButton)
        header.Resize(fyne.NewSize(700, 30))
        txScroll.Resize(fyne.NewSize(700, 452))
        txScroll.Move(fyne.NewPos(0, 36))
    	view := container.NewWithoutLayout(header, txScroll)
        
        logWin.SetCloseIntercept(func() {
            killErr := cmd.Process.Kill()
            if killErr != nil {
                fmt.Println(killErr)
            }
            os.Exit(0)
        })
		
        logWin.SetContent(view)
        logWin.Resize(fyne.NewSize(712, 500))
        logWin.Show()
        w.Close()
        
      }}
	})

	//window
	w.Resize(fyne.NewSize(712, -1))
	w.SetContent(container.NewVBox(
		widget.NewLabel("Input Directory:"),
		inputDir,
		inputSelectFolder,
		widget.NewLabel("Output Directory:"),
		outputDir,
		outputSelectFolder,
		deleteOrigin,
		widget.NewLabel("Optimize Level:"),
		optimLv,
		saveButton,
		runButton,
	))
	w.ShowAndRun()
}

func writeConfig(confpath string, config *Config) error {
	//write config
	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(confpath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func readConfig(confpath string) (*Config, error) {
	//load config
	bytes, err := ioutil.ReadFile(confpath)
	if err != nil {
		return nil, err
	}
	
	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
