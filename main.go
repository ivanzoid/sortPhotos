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
	fmt.Fprintf(os.Stderr, "usage: %v <charCount> [prefixCharCountToRemove]\n", filepath.Base(os.Args[0]))
	os.Exit(2)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	charCount := 8
	prefixCharCountToRemove := 0

	if flag.NArg() > 0 {
		arg := flag.Arg(0)
		if count, err := strconv.ParseInt(arg, 10, 64); err == nil {
			charCount = int(count)
		}
		if flag.NArg() > 1 {
			arg2 := flag.Arg(1)
			if prefixCount, err := strconv.ParseInt(arg2, 10, 64); err == nil {
				prefixCharCountToRemove = int(prefixCount)
			}
		}
	}

	searchDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	pwd, _ := os.Getwd()

	processFile := func(path string, fileInfo os.FileInfo, err error) error {
		localPath := strings.TrimPrefix(path, pwd)
		if len(localPath) == 0 {
			return nil
		}
		if len(localPath) > 0 {
			localPath = strings.TrimPrefix(localPath, "/")
		}
		if strings.Contains(localPath, "/") {
			return nil
		}
		if strings.HasPrefix(localPath, ".") || strings.HasPrefix(localPath, "#") {
			return nil
		}

		// if fileInfo.IsDir() {
		// 	return nil
		// }

		// fmt.Printf("path='%v'\n", path)
		// fmt.Printf("localPath='%v'\n", localPath)

		name := fileInfo.Name()

		groupName := ""
		if len(name) < charCount {
			groupName = name
		} else {
			groupName = name[:charCount]
		}

		runProgram("mkdir", "-p", groupName)

		if prefixCharCountToRemove == 0 || len(name) <= prefixCharCountToRemove {
			fmt.Printf("%v -> %v\n", name, groupName)
			runProgram1("mv", name, groupName)
		} else {
			outName := name[prefixCharCountToRemove:]
			outPath := fmt.Sprintf("%v/%v", groupName, outName)

			fmt.Printf("%v -> %v\n", name, outPath)
			runProgram1("mv", name, outPath)
		}

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
	// fmt.Printf("Running %v\n", cmdString(program, args))
	// return nil, nil

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
