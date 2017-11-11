# picol.go

### Motivation

Wanted to a minimal Tcl implementation, which could be used as an experimental
extention mechanism for golang projects.

#### TODO
+ Support interface{} as the command argument type of string
+ Reduce memory allocations, like instead of using string type use byte
  array as parser input and tokens would be byte slices on the source array.
+ More test-cases


#### Won't do
+ Comprehensiveness & compatibility with the standard Tcl implementation
+ Tcl standard library


### Old Readme

Original http://oldblog.antirez.com/post/picol.html

Sample use:
```golang
func CommandPuts(i *picol.Interp, argv []string, pd interface{}) (string, error) {
	if len(argv) != 2 {
		return "", fmt.Errorf("Wrong number of args for %s %s", argv[0], argv)
	}
	fmt.Println(argv[1])
	return "", nil
}
...
	interp := picol.InitInterp()
	// add core functions
	interp.RegisterCoreCommands()
	// add user function
	interp.RegisterCommand("puts", CommandPuts, nil)
	// eval
	result, err := interp.Eval(string(buf))
	if err != nil {
		fmt.Println("ERROR", err, result)
	} else {
		fmt.Println(result)
	}
```

