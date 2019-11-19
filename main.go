package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %v <charCount>\n", filepath.Base(os.Args[0]))
	os.Exit(2)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	charCount := 8

	if flag.NArg() > 0 {
		arg := flag.Arg(0)
		if count, err := strconv.ParseInt(arg, 10, 64); err == nil {
			charCount = int(count)
		}
	}

	searchDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	processFile := func(path string, fileInfo os.FileInfo, err error) error {
		if fileInfo.IsDir() {
			return nil
		}

		name := fileInfo.Name()

		groupName := ""
		if len(name) < charCount {
			groupName = name
		} else {
			groupName = name[:charCount]
		}

		fmt.Printf("%v -> %v\n", name, groupName)

		runProgram("mkdir", "-p", groupName)
		runProgram1("mv", name, groupName)

		return nil
	}

	err = filepath.Walk(searchDir, processFile)
}

func runProgram1(program string, args ...string) (string, error) {
	outStrings, err := runProgram(program, args...)
	if err != nil {
		return "", err
	}
	if len(outStrings) == 0 {
		return "", nil
	}

	return outStrings[0], nil
}

func runProgram(program string, args ...string) ([]string, error) {

	// dlog("Running %v", cmdString(program, args))

	cmd := exec.Command(program, args...)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	outString := string(out)
	if len(outString) == 0 {
		return nil, nil
	}

	outStrings := strings.Split(outString, "\n")
	return outStrings, nil
}

// Utils

func cmdString(program string, args []string) string {
	comps := make([]string, 0)

	comps = append(comps, program)

	for _, arg := range args {
		if strings.Contains(arg, " ") {
			comps = append(comps, fmt.Sprintf("'%v'", arg))
		} else {
			comps = append(comps, arg)
		}
	}

	result := strings.Join(comps, " ")
	return result
}
