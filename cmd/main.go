package main

import (
	"fmt"
	"os"

	"github.com/frank3950/cthun"
	flag "github.com/spf13/pflag"
)

var help = flag.BoolP("help", "h", false, "Print help menu")
var table = flag.StringSliceP("table", "t", nil, "Check table is in use,for example: --check-table-in-use t1,t2")
var home = flag.String("home", os.Getenv("OGG_HOME"), "Specify the ggs home directory, default $OGG_HOME")

func main() {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
	} else if *home == "" {
		fmt.Println("ERROR: can not find gg home. use --home or set $OGG_HOME")
	} else if *table != nil {
		i := cthun.Inst{Home: *home}
		s := cthun.CheckTableInUse(&i, *table)
		fmt.Println(s)
	}
}
