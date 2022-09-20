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

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage:\n  %v [options] <charCount> [prefixCharCountToRemove]\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func main() {
	verbose := flag.Bool("v", false, "Verbose")
	veryVerbose := flag.Bool("vv", false, "Very verbose")
	dryRun := flag.Bool("d", false, "Dry run")
	nonRecursive := flag.Bool("nr", false, "Non-recursive mode")
	processFilesOnly := flag.Bool("pf", false, "Process only files")
	processDirsOnly := flag.Bool("pd", false, "Process only directories")
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

		if *veryVerbose {
			fmt.Printf("path='%v', localPath='%v'\n", path, localPath)
		}

		if *nonRecursive && strings.Contains(localPath, "/") {
			return nil
		}

		if strings.HasPrefix(localPath, ".") || strings.HasPrefix(localPath, "#") {
			return nil
		}

		if (*processDirsOnly && !fileInfo.IsDir()) || (*processFilesOnly && fileInfo.IsDir()) {
			return nil
		}

		name := fileInfo.Name()

		groupName := ""
		if len(name) < charCount {
			groupName = name
		} else {
			groupName = name[:charCount]
		}

		runProgram(!*dryRun, *verbose || *veryVerbose, "mkdir", "-p", groupName)

		if prefixCharCountToRemove == 0 || len(name) <= prefixCharCountToRemove {
			fmt.Printf("%v -> %v\n", name, groupName)
			runProgram1(!*dryRun, *verbose || *veryVerbose, "mv", name, groupName)
		} else {
			outName := name[prefixCharCountToRemove:]
			outPath := fmt.Sprintf("%v/%v", groupName, outName)
			if len(groupName) == 0 {
				outPath = outName
			}

			slashOrEmpty := ""
			if fileInfo.IsDir() {
				slashOrEmpty = "/"
			}
			fmt.Printf("%v%v -> %v%v\n", name, slashOrEmpty, outPath, slashOrEmpty)
			runProgram1(!*dryRun, *verbose || *veryVerbose, "mv", name, outPath)
		}

		return nil
	}

	err = filepath.Walk(searchDir, processFile)
}

func runProgram1(run bool, verbose bool, program string, args ...string) (string, error) {
	outStrings, err := runProgram(run, verbose, program, args...)
	if err != nil {
		return "", err
	}
	if len(outStrings) == 0 {
		return "", nil
	}

	return outStrings[0], nil
}

func runProgram(run bool, verbose bool, program string, args ...string) ([]string, error) {
	if verbose {
		fmt.Printf("$ %v\n", cmdString(program, args))
	}

	if !run {
		return []string{}, nil
	}

	cmd := exec.Command(program, args...)

	out, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	outString := string(out)
	if len(outString) == 0 {
		return []string{}, nil
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
