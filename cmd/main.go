package main

import (
	"fmt"
	"os"

	"github.com/frank3950/cthun"
	flag "github.com/spf13/pflag"
)

var help = flag.BoolP("help", "h", false, "Print help menu")
var s = flag.StringP("search", "s", "", "Search parameter, support name/rmthost/table/map")
var home = flag.String("home", os.Getenv("OGG_HOME"), "Specify the ggs home directory, default $OGG_HOME")
var size = flag.Bool("size", false, "Get dirdat size")

func main() {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
	} else if *home == "" {
		fmt.Println("ERROR: can not find gg home. use --home or set $OGG_HOME")
	} else if *s != "" {
		i := cthun.ClassicGG{Home: *home}
		cthun.SetupGG(&i)
		result := cthun.SearchGG(i, *s)
		for _, r := range result {
			fmt.Println(r)
		}
	} else if *size {
		i := cthun.ClassicGG{Home: *home}
		cthun.SetupGG(&i)
		s, err := cthun.GetGGDatSize(i)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(s)
	}
}
