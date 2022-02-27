package cthun

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	LogInfo  *log.Logger = log.New(os.Stdout, "[INFO]  ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	LogWarn  *log.Logger = log.New(os.Stdout, "[WARN]  ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	LogError *log.Logger = log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
)

func ExecCMD(c string) (string, error) {
	cmd := exec.Command("sh", "-c", c)
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogError.Println(string(out) + err.Error())
	}
	return string(out), err
}

type GGHandler interface {
	Setup() error
	GetAllTab() []string
}

func CheckTableInUse(gg GGHandler, tList []string) string {
	gg.Setup()
	var builder strings.Builder
	allTab := gg.GetAllTab()
	for _, t := range tList {
		for _, str := range allTab {
			if strings.Contains(str, strings.ToUpper(t)) {
				builder.WriteString(str + "\n")
			}
		}
	}
	return builder.String()
}
