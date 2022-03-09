package cthun

import (
	"bufio"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	i := ClassicGG{}
	i.exts = append(i.exts, ext{name: "E_TTT", tables: []string{"t1"}})
	i.pumps = append(i.pumps, pump{name: "P_TTT", tables: []string{"t1"}})
	i.reps = append(i.reps, rep{name: "R_TTT", maps: map[string]string{
		"a": "aa",
	}})
	s := i.Search("TTT")
	LogInfo.Println(i)
	LogInfo.Println(s)
	if len(s) != 3 {
		t.Errorf("search %v, expected %v", len(s), 3)
	}
}

func TestGetDatSize(t *testing.T) {
	i := ClassicGG{Home: "test"}
	bSize, err := i.GetDatSize()
	if err != nil {
		t.Errorf("TestGetDirSize err:%s", err)
	}
	if bSize != 8 {
		t.Errorf("GetDirSize()=%v expected=%v", bSize, 8)
	}
}

func TestAddExt(t *testing.T) {
	i := ClassicGG{}
	eChan := make(chan ext)
	go func() {
		eChan <- ext{name: "E_CXA"}
		eChan <- ext{name: "E_CXB"}
		eChan <- ext{name: "E_CXC"}
		close(eChan)
	}()
	i.addExt(eChan)
	eCount := len(i.exts)
	if eCount != 3 {
		t.Errorf("get %d extract, expected %d", eCount, 3)
	}
}

func TestAddPump(t *testing.T) {
	i := ClassicGG{}
	pChan := make(chan pump)
	go func() {
		pChan <- pump{name: "P_CXA"}
		pChan <- pump{name: "P_CXB"}
		close(pChan)
	}()
	i.addPump(pChan)
	pCount := len(i.pumps)
	if pCount != 2 {
		t.Errorf("get %d pump, expected %d", pCount, 2)
	}
}

func TestAddRep(t *testing.T) {
	i := ClassicGG{}
	rChan := make(chan rep)
	go func() {
		rChan <- rep{name: "R_CXA"}
		rChan <- rep{name: "R_CXB"}
		rChan <- rep{name: "R_CXC"}
		rChan <- rep{name: "R_CXD"}
		close(rChan)
	}()
	i.addRep(rChan)
	rCount := len(i.reps)
	if rCount != 4 {
		t.Errorf("get %d replicat, expected %d", rCount, 4)
	}
}

func TestParseParamFile(t *testing.T) {
	i := ClassicGG{Home: "test/"}
	echan := make(chan ext)
	pchan := make(chan pump)
	rchan := make(chan rep)
	go func() {
		echan <- ext{name: "E_TEST"}
		pchan <- pump{name: "P_TEST"}
		rchan <- rep{name: "R_TEST"}
		close(echan)
		close(pchan)
		close(rchan)
	}()
	e, p, r := i.parseParamFile(echan, pchan, rchan)

	re := <-e
	rp := <-p
	rr := <-r
	LogInfo.Println(rr)
	if len(re.tables) != 3 {
		t.Errorf("get %d extract table, excepted %d", len(re.state), 3)
	}
	if len(rp.tables) != 2 {
		t.Errorf("get %d pump table, excepted %d", len(re.state), 2)
	}
	if len(rr.maps) != 3 {
		t.Errorf("get %d replicat table, excepted %d", len(rr.maps), 3)
	}
}

func TestParseInfoDetailString(t *testing.T) {
	c := make(chan string)
	go func() {
		c <- testExtStr
		c <- testPumpStr
		c <- testRepStr
		close(c)
	}()
	e, p, r := parseInfoDetailString(c)
	e1 := <-e
	if e1.name != "E_CXH" {
		t.Errorf("extract name %s != %s", e1.name, "E_CXH")
	}
	p1 := <-p
	if p1.name != "P_CXE" {
		t.Errorf("pump name %s != %s", p1.name, "P_CXE")
	}
	r1 := <-r
	if r1.name != "R_JKDCN" {
		t.Errorf("replicat name %s != %s", r1.name, "R_JKDCN")
	}
}

func TestTakeInfoDetailString(t *testing.T) {
	var rightTestCase = []struct {
		in ClassicGG
	}{
		{in: ClassicGG{Home: "test//"}},
		{in: ClassicGG{Home: "test/"}},
		{in: ClassicGG{Home: "test"}},
	}
	for _, tc := range rightTestCase {
		_, err := tc.in.getInfoDetail()
		if err != nil {
			t.Errorf("TakeInfoDetailString(%v) err=%v", tc.in, err)
		}
	}

	var wrongTestCase = []struct {
		in ClassicGG
	}{
		{in: ClassicGG{Home: "/tmp/"}},
		{in: ClassicGG{Home: "test/gg_home1/"}},
	}
	for _, tc := range wrongTestCase {
		_, err := tc.in.getInfoDetail()
		if err == nil {
			t.Errorf("TakeInfoDetailString(%v) should return err", tc.in)
		}
	}
}

func TestCutInfoDetailString(t *testing.T) {
	c := cutInfoDetail(testInfoDetailStr1)
	for str := range c {
		buf := bufio.NewScanner(strings.NewReader(str))
		for buf.Scan() {
			if buf.Text() == "this line should not match, match fail" {
				t.Error("match fail")
			}
		}
	}
}

func TestSetup(t *testing.T) {
	var testCase = []string{"test/", "test//", "test"}
	for _, tc := range testCase {
		i := ClassicGG{Home: tc}
		err := i.Setup()
		if err != nil {
			t.Errorf("setup error: %s", err)
		}
	}
}

var testInfoDetailStr1 = `EXTRACT    E_CXH     Last Started 2021-12-06 11:42   Status RUNNING
Checkpoint Lag       00:00:01 (updated 00:00:01 ago)
Process ID           166878
Log Read Checkpoint  Oracle Redo Logs
                     2022-01-21 10:10:54  Seqno 374830, RBA 1116801040
                     SCN 74.3145158248 (320972738152)

  Target Extract Trails:

  Trail Name                                       Seqno        RBA     Max MB Trail Type

  /u01/app/ogg/dirdat/xh                            3325  459193485       1024 EXTTRAIL  

  Extract Source                          Begin             End             

  +DATA/orcl/onlinelog/group_14.276.1073659605  2021-12-06 11:25  2022-01-21 10:10
  +DATA/orcl/onlinelog/group_7.269.1073659567  2021-12-06 11:04  2021-12-06 11:42
  +DATA/orcl/onlinelog/group_7.269.1073659567  2021-12-06 11:04  2021-12-06 11:42
  +DATA/orcl/onlinelog/group_17.279.1073659665  2021-12-03 11:07  2021-12-06 11:22
  +DATA/orcl/onlinelog/group_17.279.1073659665  2021-12-03 11:07  2021-12-06 11:22
  +DATA/orcl/onlinelog/group_14.276.1073659605  2021-12-03 09:35  2021-12-03 11:16
  +DATA/orcl/onlinelog/group_14.276.1073659605  2021-12-03 09:35  2021-12-03 11:16
  +DATA/orcl/onlinelog/group_3.2407.1073934761  2021-11-18 17:47  2021-12-03 09:45
  +DATA/orcl/onlinelog/group_3.2407.1073934761  2021-11-18 17:47  2021-12-03 09:45
  +DATA/orcl/onlinelog/group_17.279.1073659665  2021-10-26 16:03  2021-11-18 18:06
  +DATA/orcl/onlinelog/group_17.279.1073659665  2021-10-26 16:03  2021-11-18 18:06
  +DATA/orcl/onlinelog/group_3.2407.1073934761  2021-10-13 11:34  2021-10-26 16:13
  +DATA/orcl/onlinelog/group_3.2407.1073934761  2021-10-13 11:34  2021-10-26 16:13
  +DATA/orcl/onlinelog/group_13.275.1073659599  2021-09-15 11:17  2021-10-13 11:44
  +DATA/orcl/onlinelog/group_13.275.1073659599  2021-09-15 11:17  2021-10-13 11:44
  +DATA/orcl/onlinelog/group_16.278.1073659661  2021-09-13 16:50  2021-09-15 11:37
  +DATA/orcl/onlinelog/group_16.278.1073659661  2021-09-13 16:50  2021-09-15 11:37
  +DATA/orcl/onlinelog/group_14.276.1073659605  2021-09-10 10:21  2021-09-13 17:06
  +DATA/orcl/onlinelog/group_14.276.1073659605  2021-09-10 10:21  2021-09-13 17:06
  +DATA/orcl/onlinelog/group_25.287.1073659709  2021-09-09 15:29  2021-09-10 10:49
  +DATA/orcl/onlinelog/group_25.287.1073659709  2021-09-09 15:29  2021-09-10 10:49
  Not Available                           * Initialized *   2021-09-09 15:29
  Not Available                           * Initialized *   2021-09-09 15:29
  Not Available                           * Initialized *   2021-09-09 15:29


Current directory    /u01/app/ogg
this line should not match, match fail
Report file          /u01/app/ogg/dirrpt/E_CXH.rpt
Parameter file       /u01/app/ogg/dirprm/e_cxh.prm
Checkpoint file      /u01/app/ogg/dirchk/E_CXH.cpe
Process file         
Error log            /u01/app/ogg/ggserr.log
this line should not match, match fail
EXTRACT    P_CXE     Last Started 2021-12-22 09:51   Status RUNNING
Checkpoint Lag       00:00:00 (updated 00:00:09 ago)
Process ID           247756
Log Read Checkpoint  File /u01/app/ogg/dirdat/xe000000432
                     2022-01-21 10:10:44.000000  RBA 330140419

  Target Extract Trails:

  Trail Name                                       Seqno        RBA     Max MB Trail Type

  /u01/app/ogg/dirdat/xe                             431  330145242       1024 RMTTRAIL  

  Extract Source                          Begin             End             

  /u01/app/ogg/dirdat/xe000000432         2021-12-22 09:49  2022-01-21 10:10
  /u01/app/ogg/dirdat/xe000000363         2021-06-11 19:34  2021-12-22 09:49
  /u01/app/ogg/dirdat/xe000000001         2021-06-11 19:34  2021-06-11 19:34
  /u01/app/ogg/dirdat/xe000000001         * Initialized *   2021-06-11 19:34
  /u01/app/ogg/dirdat/xe000000000         * Initialized *   First Record    


Current directory    /u01/app/ogg

Report file          /u01/app/ogg/dirrpt/P_CXE.rpt
Parameter file       /u01/app/ogg/dirprm/p_cxe.prm
Checkpoint file      /u01/app/ogg/dirchk/P_CXE.cpe
Process file         
Error log            /u01/app/ogg/ggserr.log

REPLICAT   R_JKDCN   Last Started 2021-11-05 22:08   Status RUNNING
Checkpoint Lag       00:00:03 (updated 00:00:04 ago)
Process ID           201445
Log Read Checkpoint  File /u01/app/ogg/dirdat/cn000016977
                     2022-01-21 10:10:48.999055  RBA 277031097

Current Log BSN value: (requires database login)

Last Committed Transaction CSN value: (requires database login)

  Extract Source                          Begin             End             

  /u01/app/ogg/dirdat/cn000016977         * Initialized *   2022-01-21 10:10
  /u01/app/ogg/dirdat/cn000005904         * Initialized *   First Record    
  /u01/app/ogg/dirdat/cn000005904         * Initialized *   First Record    
  /u01/app/ogg/dirdat/cn000000000         * Initialized *   First Record    


Current directory    /u01/app/ogg
this line should not match, match fail
`
var testExtStr = `EXTRACT    E_CXH     Last Started 2021-12-06 11:42   Status RUNNING
Checkpoint Lag       00:00:01 (updated 00:00:01 ago)
Process ID           166878
Log Read Checkpoint  Oracle Redo Logs
                     2022-01-21 10:10:54  Seqno 374830, RBA 1116801040
                     SCN 74.3145158248 (320972738152)

  Target Extract Trails:

  Trail Name                                       Seqno        RBA     Max MB Trail Type

  /u01/app/ogg/dirdat/xh                            3325  459193485       1024 EXTTRAIL  

  Extract Source                          Begin             End             

  +DATA/orcl/onlinelog/group_14.276.1073659605  2021-12-06 11:25  2022-01-21 10:10
  +DATA/orcl/onlinelog/group_7.269.1073659567  2021-12-06 11:04  2021-12-06 11:42
  +DATA/orcl/onlinelog/group_7.269.1073659567  2021-12-06 11:04  2021-12-06 11:42
  +DATA/orcl/onlinelog/group_17.279.1073659665  2021-12-03 11:07  2021-12-06 11:22
  +DATA/orcl/onlinelog/group_17.279.1073659665  2021-12-03 11:07  2021-12-06 11:22
  +DATA/orcl/onlinelog/group_14.276.1073659605  2021-12-03 09:35  2021-12-03 11:16
  +DATA/orcl/onlinelog/group_14.276.1073659605  2021-12-03 09:35  2021-12-03 11:16
  +DATA/orcl/onlinelog/group_3.2407.1073934761  2021-11-18 17:47  2021-12-03 09:45
  +DATA/orcl/onlinelog/group_3.2407.1073934761  2021-11-18 17:47  2021-12-03 09:45
  +DATA/orcl/onlinelog/group_17.279.1073659665  2021-10-26 16:03  2021-11-18 18:06
  +DATA/orcl/onlinelog/group_17.279.1073659665  2021-10-26 16:03  2021-11-18 18:06
  +DATA/orcl/onlinelog/group_3.2407.1073934761  2021-10-13 11:34  2021-10-26 16:13
  +DATA/orcl/onlinelog/group_3.2407.1073934761  2021-10-13 11:34  2021-10-26 16:13
  +DATA/orcl/onlinelog/group_13.275.1073659599  2021-09-15 11:17  2021-10-13 11:44
  +DATA/orcl/onlinelog/group_13.275.1073659599  2021-09-15 11:17  2021-10-13 11:44
  +DATA/orcl/onlinelog/group_16.278.1073659661  2021-09-13 16:50  2021-09-15 11:37
  +DATA/orcl/onlinelog/group_16.278.1073659661  2021-09-13 16:50  2021-09-15 11:37
  +DATA/orcl/onlinelog/group_14.276.1073659605  2021-09-10 10:21  2021-09-13 17:06
  +DATA/orcl/onlinelog/group_14.276.1073659605  2021-09-10 10:21  2021-09-13 17:06
  +DATA/orcl/onlinelog/group_25.287.1073659709  2021-09-09 15:29  2021-09-10 10:49
  +DATA/orcl/onlinelog/group_25.287.1073659709  2021-09-09 15:29  2021-09-10 10:49
  Not Available                           * Initialized *   2021-09-09 15:29
  Not Available                           * Initialized *   2021-09-09 15:29
  Not Available                           * Initialized *   2021-09-09 15:29


Current directory    /u01/app/ogg
this line should not match, match fail
Report file          /u01/app/ogg/dirrpt/E_CXH.rpt
Parameter file       /u01/app/ogg/dirprm/e_cxh.prm
Checkpoint file      /u01/app/ogg/dirchk/E_CXH.cpe
Process file         
Error log            /u01/app/ogg/ggserr.log
this line should not match, match fail`

var testPumpStr = `EXTRACT    P_CXE     Last Started 2021-12-22 09:51   Status RUNNING
Checkpoint Lag       00:00:00 (updated 00:00:09 ago)
Process ID           247756
Log Read Checkpoint  File /u01/app/ogg/dirdat/xe000000432
                     2022-01-21 10:10:44.000000  RBA 330140419

  Target Extract Trails:

  Trail Name                                       Seqno        RBA     Max MB Trail Type

  /u01/app/ogg/dirdat/xe                             431  330145242       1024 RMTTRAIL  

  Extract Source                          Begin             End             

  /u01/app/ogg/dirdat/xe000000432         2021-12-22 09:49  2022-01-21 10:10
  /u01/app/ogg/dirdat/xe000000363         2021-06-11 19:34  2021-12-22 09:49
  /u01/app/ogg/dirdat/xe000000001         2021-06-11 19:34  2021-06-11 19:34
  /u01/app/ogg/dirdat/xe000000001         * Initialized *   2021-06-11 19:34
  /u01/app/ogg/dirdat/xe000000000         * Initialized *   First Record    


Current directory    /u01/app/ogg

Report file          /u01/app/ogg/dirrpt/P_CXE.rpt
Parameter file       /u01/app/ogg/dirprm/p_cxe.prm
Checkpoint file      /u01/app/ogg/dirchk/P_CXE.cpe
Process file         
Error log            /u01/app/ogg/ggserr.log
`

var testRepStr = `REPLICAT   R_JKDCN   Last Started 2021-11-05 22:08   Status RUNNING
Checkpoint Lag       00:00:03 (updated 00:00:04 ago)
Process ID           201445
Log Read Checkpoint  File /u01/app/ogg/dirdat/cn000016977
                     2022-01-21 10:10:48.999055  RBA 277031097

Current Log BSN value: (requires database login)

Last Committed Transaction CSN value: (requires database login)

  Extract Source                          Begin             End             

  /u01/app/ogg/dirdat/cn000016977         * Initialized *   2022-01-21 10:10
  /u01/app/ogg/dirdat/cn000005904         * Initialized *   First Record    
  /u01/app/ogg/dirdat/cn000005904         * Initialized *   First Record    
  /u01/app/ogg/dirdat/cn000000000         * Initialized *   First Record    


Current directory    /u01/app/ogg
this line should not match, match fail`
