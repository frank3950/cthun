package cthun

import (
	"log"
	"os"
	"os/exec"
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

type Setuper interface {
	Setup() error
}

type Searcher interface {
	Search(keyword string) []string
}

type Lager interface {
	GetAllLag() map[string]int
	GetAllCkpLag() map[string]int
}

type Sizer interface {
	GetDatSize() (int, error)
}

func GetGGDatSize(s Sizer) (int, error) {
	return s.GetDatSize()
}

func SearchGG(s Searcher, k string) []string {
	result := s.Search(k)
	return result
}

func SetupGG(s Setuper) {
	s.Setup()
}
