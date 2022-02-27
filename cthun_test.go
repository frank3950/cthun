package cthun

import (
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
