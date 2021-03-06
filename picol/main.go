package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/shanmuganandh/picol"
	"io/ioutil"
	"os"
)

var fname = flag.String("f", "", "file name")

func CommandPuts(i *picol.Interp, argv []interface{}, pd interface{}) (interface{}, error) {
	if len(argv) != 2 {
		return "", fmt.Errorf("Wrong number of args for %s %s", argv[0].(string), argv)
	}
	fmt.Println(argv[1])
	return "", nil
}

func LoadedInterp() *picol.Interp {
	i := picol.InitInterp()
	i.RegisterCoreCommands()
	i.RegisterCommand("puts", CommandPuts, nil)

	return i
}

func main() {
	flag.Parse()

	interp := LoadedInterp()

	buf, err := ioutil.ReadFile(*fname)
	if err == nil {
		result, err := interp.Eval(string(buf))
		if err != nil {
			fmt.Println("ERRROR", result, err)
		}
	} else {
		for {
			fmt.Print("picol> ")
			scanner := bufio.NewReader(os.Stdin)
			clibuf, _ := scanner.ReadString('\n')
			result, err := interp.Eval(clibuf[:len(clibuf)-1])
			if len(result.(string)) != 0 {
				fmt.Println("ERRROR", result, err)
			}
		}
	}
}
