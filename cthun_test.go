package cthun

import (
	"reflect"
	"testing"
)

func TestExecCMD(t *testing.T) {
	var testCase = []struct {
		cmd string // cmd
	}{
		{"ls"},
		{"ls -l"},
		{"ls -lh"},
		{"ls -l -h"},
		{"test//" + "/ggsci<<EOF\ninfo * detail\nEOF\n"},
		{"test/" + "/ggsci<<EOF\ninfo * detail\nEOF\n"},
		{"test" + "/ggsci<<EOF\ninfo * detail\nEOF\n"},
	}
	for _, tc := range testCase {
		o, err := ExecCMD(tc.cmd)
		if err != nil {
			t.Errorf("\nexec %s with error:\n%s\n%s", tc.cmd, o, err)
		}
	}
	testCase = []struct {
		cmd string // cmd
	}{
		{"jhgjhgjg"},
		{"sss -l"},
		{"ls -l."},
		{"ls -l -."},
	}
	for _, tc := range testCase {
		_, err := ExecCMD(tc.cmd)
		if err == nil {
			t.Errorf("\nexec %s should return error", tc.cmd)
		}
	}
}

func TestSearchGG(t *testing.T) {
	i := ClassicGG{}
	i.exts = append(i.exts, ext{name: "E_TTT", tables: []string{"t1"}})
	i.pumps = append(i.pumps, pump{name: "P_TTT", tables: []string{"t1"}})
	i.reps = append(i.reps, rep{name: "R_TTT", maps: map[string]string{
		"a": "aa",
	}})
	s := SearchGG(i, "TTT")
	if len(s) != 3 {
		t.Errorf("search %v, expected %v", len(s), 3)
	}
}

func TestGetGGLag(t *testing.T) {
	gg := ClassicGG{}
	m1, m2 := getLagMap(testInfo)
	gg.setupLag(m1, m2)
	e1, e2 := GetGGLag(&gg)
	if reflect.DeepEqual(m1, e1) {
		t.Errorf("lagMap=%v, expected=%v", e1, m1)
	}
	if reflect.DeepEqual(m2, e2) {
		t.Errorf("ckpLagMap=%v, expected=%v", e2, m2)
	}
}
