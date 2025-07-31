## pcre


    NAME
	pcre - run a list of PCRE regexes across a data set

    SYNOPSIS
	pcre --regex REGEX_LIST --data FILE [ --alpha | --desc ] [ --dumpreg REGEX_NAME ] [ --unmatch FILE ]
	pcre --version | --help

    DESCRIPTION
       REQUIRED
	--regex REGEX_LIST
		Specify the location of the REGEX_LIST file
	--data  FILE
		Specify the location of the file containing the DATA set to test

       OPTIONAL

	--dumpreg REGEX_NAME
		Specify which REGEX_NAME results to dump to a file. The file will be
		named REGEX_NAME.txt

	--unmatch FILE
		Specify the name of the FILE where unmatched records will be placed.
		The default name is 'unmatched.txt'

	By default, results are printed in the order listed in REGEX_LIST. The result
        output can be modified with these next two flags

	--alpha
		Results will be printed with regex names from REGEX_LIST listed
                alphabetically
	--desc
		Results will be printed with hit counts listed in descending order

	--help
		Display this usage text and exit
	--version
		Display the program version and exit


### examples


    > ./pcre --regex ../regex/pcre.list --data ../data/data.txt
    SUMMARY

      run date :  2025-07-28T00:36:26
      0.000257 : elapsed total (seconds)
            29 : data lines read in
             9 : regexes loaded
            29 : matched lines
             0 : unmatched lines

    MATCHES

                             REGEX      ORD     HITS      DURATION
               sshd-sesson-close-1        1        4 0.000000954
               sshd-session-open-1        2        5 0.000000954
                 sshd-disconnect-1        3        1 0.000002146
              sshd-accept-passwd-1        4        3 0.000000954
              sshd-failed-passwd-1        5        2 0.000000954
                  sshd-auth-fail-1        6        1 0.000000954
                        sshd-other        7        2 0.000002146
                            cron-1        8        4 0.000002861
                      kad-notify-1        9        7 0.000001907


    > ./pcre --regex ../regex/pcre.list --data ../data/data.txt --alpha
    SUMMARY

      run date :  2025-07-28T00:36:33
      0.000404 : elapsed total (seconds)
            29 : data lines read in
             9 : regexes loaded
            29 : matched lines
             0 : unmatched lines

    MATCHES

                             REGEX      ORD     HITS      DURATION
                            cron-1        8        4 0.000001907
                      kad-notify-1        9        7 0.000002861
              sshd-accept-passwd-1        4        3 0.000000954
                  sshd-auth-fail-1        6        1 0.000002146
                 sshd-disconnect-1        3        1 0.000000954
              sshd-failed-passwd-1        5        2 0.000000954
                        sshd-other        7        2 0.000002146
               sshd-session-open-1        2        5 0.000000954
               sshd-sesson-close-1        1        4 0.000000954


    > ./pcre --regex ../regex/pcre.list --data ../data/data.txt --desc
    SUMMARY

      run date :  2025-07-28T00:37:14
      0.000234 : elapsed total (seconds)
            29 : data lines read in
             9 : regexes loaded
            29 : matched lines
             0 : unmatched lines

    MATCHES

                             REGEX      ORD     HITS      DURATION
                      kad-notify-1        9        7 0.000002861
               sshd-session-open-1        2        5 0.000000954
               sshd-sesson-close-1        1        4 0.000000000
                            cron-1        8        4 0.000003099
              sshd-accept-passwd-1        4        3 0.000001907
              sshd-failed-passwd-1        5        2 0.000000954
                        sshd-other        7        2 0.000001907
                 sshd-disconnect-1        3        1 0.000001192
                  sshd-auth-fail-1        6        1 0.000001907

