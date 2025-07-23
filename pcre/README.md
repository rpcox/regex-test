## pcre

    > ./pcre --regex ../regex/pcre.list --data ../data/data.txt
    SUMMARY

        29 : data lines read in
         9 : regexes loaded
        29 : matched lines
         0 : unmatched lines

    MATCHES

                             REGEX  SEQ     HITS     SECONDS
               sshd-sesson-close-1    1        4 0.000000954
               sshd-session-open-1    2        5 0.000000954
                 sshd-disconnect-1    3        1 0.000001907
              sshd-accept-passwd-1    4        3 0.000002146
              sshd-failed-passwd-1    5        2 0.000001907
                  sshd-auth-fail-1    6        1 0.000003099
                        sshd-other    7        2 0.000002861
                            cron-1    8        4 0.000005007
                      kad-notify-1    9        7 0.000004053

    0.000226 seconds elasped
