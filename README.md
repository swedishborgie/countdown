# countdown
A cross platform countdown ticker in Golang for OBS Studio.

To install:

    go get github.com/swedishborgie/countdown

Usage:

    $countdown -help
    Usage of countdown:
      -complete string
            after the target time is hit this will be written to the file (if set, no prefix/postfix)
      -format string
            the duration format to output (see https://github.com/davidscholberg/go-durationfmt) (default "%00h:%00m:%00s")
      -output string
            the output file to update (default "countdown.txt")
      -postfix string
            a postfix to add to end of the output
      -prefix string
            a prefix to add beginning of the output
      -target string
            the time to count down to in hh:mm:ss format (24-hr) (default "10:00:00")
      -update string
            the amount of time to wait until checking to see if the file needs updating (default "100ms")