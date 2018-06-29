[![Go Report Card](https://goreportcard.com/badge/github.com/ThomasK81/SCTAExtract)](https://goreportcard.com/report/github.com/ThomasK81/SCTAExtract)
[![DOI](https://zenodo.org/badge/138982644.svg)](https://zenodo.org/badge/latestdoi/138982644)

# SCTAExtract
Extracting SCTA texts via the SCTA CSV API

# USAGE OSX (and some Linux distros)

1. Clone this repository.
2. Open a terminal in the folder of the repository and type: `./SCTAExtract [output-filename] [optionally: endpoint]` (you might have to `chmod +x ./SCTAExtraxt` before you can use the executable)
3. If you do not set `[output-filename]` then `scta.csv` will be used as default. If you want to change the endpoint (away from the default: `http://scta.info/csv/scta`) you have to set both parameters!
4. Enjoy your new SCTA collection file!

# USAGE Windows

The same but do `WinSCTAExtract.exe [output-filename] [optionally: endpoint]`

Sample output:

```
./SCTAExtract mytext.csv
No parameters given. Using default:
Endpoint: http://scta.info/csv/scta
Output-File: mytext.csv
main references fetched: 100%
sub references fetched: 100%201322 passages extracted
writing to file...

lines written: 100%
wrote 13027646 Latin words from 201322 passages to mytext.csv
also 477 Greek words
wrote 0 Arabic words
```

# Other Operating Systtems

`SCTAExtract.go` is written in Go and can be easily compiled for your system. Flick me a message if you are interested or install the programming language Go and easily build it for your OS with `go build SCTAExtract.go`.

