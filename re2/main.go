package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rpcox/pkg/exit"
)

const _tool = "re2"
const _version = "0.1.0"

func Version() {
	fmt.Printf("%s v%s\n", _tool, _version)
	os.Exit(0)
}

func Usage(code int, msg string) {

	if msg != "" {
		fmt.Println(msg)
	}

	text := `
 NAME
	re2 - run a list of RE2 regexes across a data set

SYNOPSIS
	re2 -regex REGEX_LIST -data FILE [ -alpha | -desc ]
	re2 -version | -help

DESCRIPTION
`
	fmt.Println(text)
	flag.PrintDefaults()
	os.Exit(code)
}

func HelpVerCheck(v, h bool) {
	if v {
		Version()
	}

	if h {
		Usage(0, "")
	}
}

func LoadRegexList(fileName string) (map[string]*regexp.Regexp, map[string]int, *[]string) {
	regexFh, err := os.OpenFile(fileName, os.O_RDONLY, 0640)
	exit.IfErr(err != nil, err, 1)
	defer regexFh.Close()

	lineCount := 0
	r := make(map[string]*regexp.Regexp)
	seq := make(map[string]int)
	var sequencedRegexes []string

	scanner := bufio.NewScanner(regexFh)
	for scanner.Scan() {
		line := scanner.Text()
		if match, _ := regexp.MatchString(`^#`, line); match {
			continue
		}
		lineCount++
		s := strings.Split(line, "\t")
		exit.If(len(s) < 2, "err regex #"+strconv.Itoa(lineCount)+": "+line, 1)
		r[s[0]] = regexp.MustCompile(s[1])
		seq[s[0]] = lineCount
		sequencedRegexes = append(sequencedRegexes, s[0])
	}

	return r, seq, &sequencedRegexes
}

func main() {
	/*	_alpha := flag.Bool("alpha", false, "Results will be printed with regex names from REGEX_LIST listed alphabetically")
		_desc := flag.Bool("desc", false, "Results will be printed with hit counts listed in descending order")
	*/
	_data := flag.String("data", "", "Specify the location of the data set")
	_regex := flag.String("regex", "", "Specify the location of the regex list")
	_nomatch := flag.String("nomatch", "unmatched.txt", "Specify the location to place unmatched data")
	_help := flag.Bool("help", false, "Display usage and exit")
	_ver := flag.Bool("version", false, "Display version and exit")
	flag.Parse()
	HelpVerCheck(*_ver, *_help)

	umFh, err := os.OpenFile(*_nomatch, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
	exit.IfErr(err != nil, err, 1)
	defer umFh.Close()
	dataFh, err := os.OpenFile(*_data, os.O_RDONLY, 0640)
	exit.IfErr(err != nil, err, 1)
	defer dataFh.Close()
	scanner := bufio.NewScanner(dataFh)

	testReg, seq, sequencedRegexes := LoadRegexList(*_regex)
	stats := make(map[string]int)
	lineCount, matchCount, unmatchedCount := 0, 0, 0
	start := time.Now()

NEXT_LINE:
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		for _, regName := range *sequencedRegexes {
			if s := testReg[regName].FindString(line); s != "" {
				stats[regName]++
				matchCount++
				goto NEXT_LINE
			}
		}

		unmatchedCount++
		umFh.WriteString(line + "\n")
	}

	elapsed := time.Since(start)
	fmt.Printf("%8d : data lines read in\n", lineCount)
	fmt.Printf("%8d : matched lines\n", matchCount)
	fmt.Printf("%8d : unmatched lines\n\n", unmatchedCount)
	fmt.Println(" SEQ    COUNT  REGEX")

	//for i, regName := range *sequencedRegexes {
	for _, regName := range *sequencedRegexes {
		//	fmt.Printf("%4d %8d  %s\n", i, stats[regName], regName)
		fmt.Printf("%4d %8d  %s\n", seq[regName], stats[regName], regName)
	}

	fmt.Println("elapsed:", elapsed)
}
