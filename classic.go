package cthun

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
)

type mgr struct {
	state string
}

type ext struct {
	name   string
	state  string
	tables []string
}

type pump struct {
	name   string
	state  string
	rHost  string
	tables []string
}

type rep struct {
	name  string
	state string
	maps  map[string]string
}

type Inst struct {
	Home string
	mgr
	exts  []ext
	pumps []pump
	reps  []rep
}

func (i *Inst) addExt(eChan <-chan ext) {
	for e := range eChan {
		i.exts = append(i.exts, e)
	}
}

func (i *Inst) addPump(pChan <-chan pump) {
	for p := range pChan {
		i.pumps = append(i.pumps, p)
	}
}

func (i *Inst) addRep(rChan <-chan rep) {
	for r := range rChan {
		i.reps = append(i.reps, r)
	}
}

func (i Inst) TakeInfoDetailString() (string, error) {
	out, err := ExecCMD(i.Home + "/ggsci<<EOF\ninfo * detail\nEOF\n")
	if err != nil {
		LogError.Println(err)
		return "", err
	}
	return out, nil
}

func cutInfoDetailString(infoDetail string) <-chan string {
	c := make(chan string)
	go func() {
		var builder strings.Builder
		match := false
		buf := bufio.NewScanner(strings.NewReader(infoDetail))
		for buf.Scan() {
			line := buf.Text()
			if match {
				builder.WriteString(line + "\n")
			}
			if strings.HasPrefix(line, "EXTRACT") || strings.HasPrefix(line, "REPLICAT") {
				match = true
				builder.WriteString(line + "\n")
			}
			if strings.HasPrefix(line, "Current directory") {
				c <- builder.String()
				match = false
				builder.Reset()
			}
		}
		close(c)
	}()
	return c
}

func (i Inst) parseParamFile(e <-chan ext, p <-chan pump, r <-chan rep) (<-chan ext, <-chan pump, <-chan rep) {
	echan := make(chan ext)
	pchan := make(chan pump)
	rchan := make(chan rep)
	go func() {
		// parse extrace
		for e1 := range e {
			var builder strings.Builder
			builder.WriteString(i.Home + "/dirprm/")
			builder.WriteString(strings.ToLower(e1.name))
			builder.WriteString(".prm")
			fileName := builder.String()
			f, err := os.Open(fileName)
			defer f.Close()
			if err != nil {
				LogError.Printf("open param file error: %s", err)
			}
			bParam, err := io.ReadAll(f)
			sParam := string(bParam)
			if err != nil {
				LogError.Printf("read param file error: %s", err)
			}
			// replace comments
			rComment := regexp.MustCompile(`--.*`)
			param := rComment.ReplaceAllString(sParam, "")
			rTable := regexp.MustCompile(`(?i)table[^;]*;`)
			tList := rTable.FindAllString(param, -1)
			for _, line := range tList {
				t := strings.TrimSpace(strings.Replace(strings.ToUpper(strings.ReplaceAll(line, ";", "")), "TABLE", "", 1))
				e1.tables = append(e1.tables, t)
			}
			echan <- e1
		}
		close(echan)
	}()
	go func() {
		for p1 := range p {
			var builder strings.Builder
			builder.WriteString(i.Home + "/dirprm/")
			builder.WriteString(strings.ToLower(p1.name))
			builder.WriteString(".prm")
			fileName := builder.String()
			f, err := os.Open(fileName)
			defer f.Close()
			if err != nil {
				LogError.Printf("open param file error: %s", err)
			}
			bParam, err := io.ReadAll(f)
			sParam := string(bParam)
			if err != nil {
				LogError.Printf("read param file error: %s", err)
			}
			// replace comments
			rComment := regexp.MustCompile(`--.*`)
			param := rComment.ReplaceAllString(sParam, "")

			rHost := regexp.MustCompile(`(?i)rmthost[^,]*,`)
			host := rHost.FindString(param)
			hStr := strings.ReplaceAll(host, ",", "")
			p1.rHost = strings.Fields(hStr)[1]

			rTable := regexp.MustCompile(`(?i)table[^;]*;`)
			tList := rTable.FindAllString(param, -1)
			for _, line := range tList {
				t := strings.TrimSpace(strings.Replace(strings.ToUpper(strings.ReplaceAll(line, ";", "")), "TABLE", "", 1))
				p1.tables = append(p1.tables, t)
			}
			pchan <- p1
		}
		close(pchan)
	}()
	go func() {
		for r1 := range r {
			r1.maps = make(map[string]string)
			var builder strings.Builder
			builder.WriteString(i.Home + "/dirprm/")
			builder.WriteString(strings.ToLower(r1.name))
			builder.WriteString(".prm")
			fileName := builder.String()
			f, err := os.Open(fileName)
			defer f.Close()
			if err != nil {
				LogError.Printf("open param file error: %s", err)
			}
			bParam, err := io.ReadAll(f)
			sParam := string(bParam)
			if err != nil {
				LogError.Printf("read param file error: %s", err)
			}
			// replace comments
			rComment := regexp.MustCompile(`--.*`)
			param := rComment.ReplaceAllString(sParam, "")
			rTable := regexp.MustCompile(`(?i)map[^;]*;`)
			tList := rTable.FindAllString(param, -1)
			for _, line := range tList {
				// upper string
				upperStr := strings.ToUpper(line)
				// cut ;
				nStr := strings.ReplaceAll(upperStr, ";", "")
				// cur map prefix
				noMapStr := strings.ReplaceAll(nStr, "MAP", "")
				// trim space
				tStr := strings.TrimSpace(noMapStr)
				// cut ,
				cStr := strings.ReplaceAll(tStr, ",", " ")
				// field
				fStr := strings.Fields(cStr)
				r1.maps[fStr[0]] = fStr[2]
			}
			rchan <- r1
		}
		close(rchan)
	}()
	return echan, pchan, rchan
}

func parseInfoDetailString(c <-chan string) (<-chan ext, <-chan pump, <-chan rep) {
	echan := make(chan ext)
	pchan := make(chan pump)
	rchan := make(chan rep)
	go func() {
		for infoDetail := range c {
			if strings.HasPrefix(infoDetail, "REPLICAT") {
				buf := bufio.NewScanner(strings.NewReader(infoDetail))
				r := rep{}
				for buf.Scan() {
					line := buf.Text()
					if strings.HasPrefix(line, "REPLICAT") {
						s := strings.Fields(strings.TrimSpace(line))
						r.name = s[1]
						r.state = s[7]
						rchan <- r
					}
				}
			} else if strings.HasPrefix(infoDetail, "EXTRACT") && strings.Contains(infoDetail, "Log Read Checkpoint  Oracle Redo Logs") {
				buf := bufio.NewScanner(strings.NewReader(infoDetail))
				e := ext{}
				for buf.Scan() {
					line := buf.Text()
					if strings.HasPrefix(line, "EXTRACT") {
						s := strings.Fields(strings.TrimSpace(line))
						e.name = s[1]
						e.state = s[7]
						echan <- e
					}
				}
			} else {
				buf := bufio.NewScanner(strings.NewReader(infoDetail))
				p := pump{}
				for buf.Scan() {
					line := buf.Text()
					if strings.HasPrefix(line, "EXTRACT") {
						s := strings.Fields(strings.TrimSpace(line))
						p.name = s[1]
						p.state = s[7]
						pchan <- p
					}
				}
			}
		}
		close(echan)
		close(pchan)
		close(rchan)
	}()
	return echan, pchan, rchan
}

func (i *Inst) Setup() error {
	s, err := i.TakeInfoDetailString()
	if err != nil {
		LogError.Println(err)
		return err
	}
	c := cutInfoDetailString(s)
	eChan, pChan, rChan := i.parseParamFile(parseInfoDetailString(c))
	var wg = sync.WaitGroup{}
	wg.Add(3)
	go func() {
		i.addExt(eChan)
		wg.Done()
	}()
	go func() {
		i.addPump(pChan)
		wg.Done()
	}()
	go func() {
		i.addRep(rChan)
		wg.Done()
	}()
	wg.Wait()
	return nil
}

func (i Inst) GetAllTab() []string {
	var tList []string
	for _, e := range i.exts {
		for _, et := range e.tables {
			str := e.name + ": " + et
			tList = append(tList, str)
		}
	}
	for _, p := range i.pumps {
		for _, pt := range p.tables {
			str := p.name + ": RMTHOST=" + p.rHost + " " + pt
			tList = append(tList, str)
		}
	}
	for _, r := range i.reps {
		for rk, rv := range r.maps {
			str := r.name + ": " + "MAP " + rk + ",TABLE " + rv
			tList = append(tList, str)
		}
	}
	return tList
}
