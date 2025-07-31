package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rpcox/pkg/exit"
)

const _tool = "re2"
const _version = "0.3.0"

var (
	_commit string
	_branch string
)

func Version() {
	if _commit != "" {
		fmt.Fprintf(os.Stdout, "%s v%s (commit:%s branch:%s)\n", _tool, _version, _commit, _branch)
	} else {
		fmt.Fprintf(os.Stdout, "%s v%s)\n", _tool, _version)
	}

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
	re2 --regex REGEX_LIST --data FILE [ --alpha | --desc ] [ --dumpreg REGEX_NAME ] [ --unmatch FILE ]
	re2 --version | --help

 DESCRIPTION
    REQUIRED
	-regex REGEX_LIST
		Specify the location of the REGEX_LIST file
	-data  FILE
		Specify the location of the file containing the DATA set to test

    OPTIONAL

	-dumpreg REGEX_NAME
		Specify which REGEX_NAME results to dump to a file. The file will be
		named REGEX_NAME.txt

	-unmatch FILE
		Specify the name of the FILE where unmatched records will be placed.
		The default name is 'unmatched.txt'

	By default, results are printed in the order listed in REGEX_LIST. The result
        output can be modified with these next two flags

	-alpha
		Results will be printed with regex names from REGEX_LIST listed
                alphabetically
	-desc
		Results will be printed with hit counts listed in descending order

	-help
		Display this usage text and exit
	-version
		Display the program version and exit

SEE ALSO
	https://github.com/google/re2/wiki/syntax

`
	fmt.Println(text)
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

type RegexData struct {
	Elapsed  time.Duration  // duration of time to match on regex Name
	HitCount int64          // number of hits for a given regex
	Name     string         // name of the regex
	Re       *regexp.Regexp // RE2 regex defined for Name
	Sequence int            // the ordinal number (position) of the regex in rule list
}

func LoadRegexList(fileName string) (map[string]*RegexData, *[]string) {
	regexFh, err := os.OpenFile(fileName, os.O_RDONLY, 0640)
	exit.IfErr(err, 1)
	defer regexFh.Close()

	lineCount := 0
	m := make(map[string]*RegexData)
	var sequencedRegexes []string

	comment := regexp.MustCompile(`^#`)
	scanner := bufio.NewScanner(regexFh)
	for scanner.Scan() {
		line := scanner.Text()
		if match := comment.FindString(line); match != "" {
			continue
		}
		lineCount++
		var rd RegexData
		s := strings.Split(line, "\t")
		exit.If(len(s) < 2, "err regex #"+strconv.Itoa(lineCount)+": "+line, 1)
		rd.Re = regexp.MustCompile(s[1])
		rd.Sequence = lineCount
		m[s[0]] = &rd
		sequencedRegexes = append(sequencedRegexes, s[0])
	}

	return m, &sequencedRegexes
}

func SortByKey(m map[string]*RegexData) *[]string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return &keys
}

func SortByHitCount(m map[string]*RegexData) *[]string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]].HitCount > m[keys[j]].HitCount
	})

	return &keys
}

func main() {
	_alpha := flag.Bool("alpha", false, "Results will be printed by alphabetical REGEX_NAME")
	_desc := flag.Bool("desc", false, "Results will be printed with regex hit counts listed in descending order")
	_dumpreg := flag.String("dumpreg", "", "Dump the records that match a particular regex")
	_data := flag.String("data", "", "Specify the location of the data set")
	_regex := flag.String("regex", "", "Specify the location of the regex list")
	_unmatch := flag.String("unmatch", "unmatched.txt", "Specify the location to place unmatched data")
	_help := flag.Bool("help", false, "Display usage and exit")
	_ver := flag.Bool("version", false, "Display version and exit")
	flag.Parse()
	HelpVerCheck(*_ver, *_help)
	exit.If(*_alpha && *_desc, "-alpha and -desc are mutually exclusive", 1)

	umFh, err := os.OpenFile(*_unmatch, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
	exit.IfErr(err, 1)
	defer umFh.Close()
	dataFh, err := os.OpenFile(*_data, os.O_RDONLY, 0640)
	exit.IfErr(err, 1)
	defer dataFh.Close()

	regexMap, regexSequence := LoadRegexList(*_regex)
	var dumpFh *os.File
	if *_dumpreg != "" {
		fileName := *_dumpreg + ".txt"
		if _, ok := regexMap[*_dumpreg]; ok {
			dumpFh, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0640)
			exit.IfErr(err, 1)
			fmt.Printf("dumping '%s' to %s\n", *_dumpreg, fileName)
		} else {
			fmt.Fprintf(os.Stderr, "no matching regex name found for '$opt{dumpreg}'. nothing to dump\n")
			_dumpreg = nil
		}
	}

	lineCount, matchCount, unmatchedCount := 0, 0, 0
	start := time.Now()

	scanner := bufio.NewScanner(dataFh)
NEXT_LINE:
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		lineStart := time.Now()
		for _, regName := range *regexSequence {
			if s := regexMap[regName].Re.FindString(line); s != "" {
				regexMap[regName].Elapsed += time.Since(lineStart)
				regexMap[regName].HitCount++
				matchCount++
				if _dumpreg != nil && regName == *_dumpreg {
					dumpFh.WriteString(line + "\n")
				}
				goto NEXT_LINE
			}
		}

		unmatchedCount++
		umFh.WriteString(line + "\n")
	}

	elapsed := time.Since(start)

	var matchSum int64
	for _, v := range regexMap {
		if v.HitCount > 0 {
			matchSum += v.HitCount
		}
	}

	fmt.Printf("SUMMARY\n\n")
	fmt.Printf("regex list : %s\n", *_regex)
	fmt.Printf("  data set : %s\n", *_data)
	fmt.Printf("  run date : %20s\n", time.Now().UTC().Format(time.RFC3339))
	fmt.Printf("%10.6f : elapsed (seconds)\n", elapsed.Seconds())
	fmt.Printf("%10d : data lines read in\n", lineCount)
	fmt.Printf("%10d : regexes loaded\n", len(*regexSequence))
	fmt.Printf("%10d : matched lines\n", matchSum)
	if matchSum != int64(matchCount) {
		fmt.Fprintf(os.Stderr, "\t!!matchSum (%d) != matchCount (%d)\n", matchSum, matchCount)
	}
	fmt.Printf("%10d : unmatched lines\n\n", unmatchedCount)
	fmt.Printf("MATCHES\n\n")
	fmt.Printf("%30s %8s %8s   %11s\n", "REGEX", "ORD", "HITS", "DURATION")

	var sp *[]string
	if *_alpha {
		sp = SortByKey(regexMap)
	} else if *_desc {
		sp = SortByHitCount(regexMap)
	} else {
		sp = regexSequence
	}

	for _, regName := range *sp {
		fmt.Printf("%30s %8d %8d    %10.6f\n", regName, regexMap[regName].Sequence, regexMap[regName].HitCount, regexMap[regName].Elapsed.Seconds())
	}
	fmt.Println()
}
