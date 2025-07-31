  ## re2

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


### examples


     > ./re2 -regex ../regex/pcre.list -data ../data/data.txt
    SUMMARY

      run date : 2025-07-28T00:21:37Z
      0.000562 : elapsed (seconds)
            29 : data lines read in
             9 : regexes loaded
            29 : matched lines
             0 : unmatched lines

    MATCHES

                             REGEX      ORD     HITS    DURATION
               sshd-sesson-close-1        1        4      0.000213
               sshd-session-open-1        2        5      0.000018
                 sshd-disconnect-1        3        1      0.000007
              sshd-accept-passwd-1        4        3      0.000033
              sshd-failed-passwd-1        5        2      0.000015
                  sshd-auth-fail-1        6        1      0.000010
                        sshd-other        7        2      0.000026
                            cron-1        8        4      0.000050
                      kad-notify-1        9        7      0.000037

    > ./re2 -regex ../regex/pcre.list -data ../data/data.txt -alpha
    SUMMARY

      run date : 2025-07-28T00:22:06Z
      0.000707 : elapsed (seconds)
            29 : data lines read in
             9 : regexes loaded
            29 : matched lines
             0 : unmatched lines

    MATCHES

                             REGEX      ORD     HITS    DURATION
                            cron-1        8        4      0.000053
                      kad-notify-1        9        7      0.000038
              sshd-accept-passwd-1        4        3      0.000035
                  sshd-auth-fail-1        6        1      0.000011
                 sshd-disconnect-1        3        1      0.000008
              sshd-failed-passwd-1        5        2      0.000017
                        sshd-other        7        2      0.000028
               sshd-session-open-1        2        5      0.000019
               sshd-sesson-close-1        1        4      0.000221


    > ./re2 -regex ../regex/pcre.list -data ../data/data.txt -desc
    SUMMARY

      run date : 2025-07-28T00:20:40Z
      0.000622 : elapsed (seconds)
            29 : data lines read in
             9 : regexes loaded
            29 : matched lines
             0 : unmatched lines

    MATCHES

                             REGEX      ORD     HITS    DURATION
                      kad-notify-1        9        7      0.000028
               sshd-session-open-1        2        5      0.000013
               sshd-sesson-close-1        1        4      0.000195
                            cron-1        8        4      0.000037
              sshd-accept-passwd-1        4        3      0.000026
              sshd-failed-passwd-1        5        2      0.000012
                        sshd-other        7        2      0.000019
                 sshd-disconnect-1        3        1      0.000005
                  sshd-auth-fail-1        6        1      0.000008

