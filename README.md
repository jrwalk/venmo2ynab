# venmo2ynab
quick script for converting Venmo's CSV download to a format readable by YNAB's 
web client.

## building the tool

With Go installed, simply run

```
go build venmo2ynab.go
```

to build a binary executable.

## running the tool

The compiled binary can be run as `./venmo2ynab` with the following flags:

```
Usage of ./venmo2ynab:
  -dir string
        working directory (default "./")
  -inFile string
        name of input file
  -outFile string
        name of output file
```