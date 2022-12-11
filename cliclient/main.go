package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/yac-d/sbc/conf"
	"github.com/yac-d/sbc/pastebin"
)

const usage = `Simple Bin Commands
USAGE: sbc [options] [arguments]
Options:
    -l, --list [n]
        Displays details of files
    -f, -u, --upload <File> [Note]
        Uploads File to the pastebin with Note
    -o <Pattern>
        Downloads the most recent file whose name matches the given pattern
    -d, --download <n>
        Downloads the nth most recent file (n > 0)
    -c, --config
        Reconfigure sbc settings`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	config, err := conf.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pb := pastebin.New(config.Server)

	var (
		fileToUpload   string
		patternToDload string
		nToDload       int
		reconfigure    bool
		list           bool
	)

	flag.StringVar(&fileToUpload, "f", "", "Uploads file")
	flag.StringVar(&fileToUpload, "u", "", "Uploads file")
	flag.StringVar(&fileToUpload, "upload", "", "Uploads file")
	flag.StringVar(&patternToDload, "o", "", "Downloads most recent file matching pattern")
	flag.IntVar(&nToDload, "d", -1, "Downloads nth most recent file (n > 0)")
	flag.IntVar(&nToDload, "download", -1, "Downloads nth most recent file (n > 0)")
	flag.BoolVar(&reconfigure, "c", false, "Reconfigure sbc")
	flag.BoolVar(&reconfigure, "config", false, "Reconfigure sbc")
	flag.BoolVar(&list, "l", false, "Displays details of files")
	flag.BoolVar(&list, "list", false, "Displays details of files")

	flag.Usage = func() { fmt.Println(usage) }
	flag.Parse()

	if fileToUpload != "" {
		// flag.Arg(0) contains the note
		checkError(pb.Upload(fileToUpload, flag.Arg(0)))
	} else if patternToDload != "" {
		cnt, err := pb.Count()
		checkError(err)
		for i := 1; i < cnt+1; i++ {
			details, err := pb.DetailsOfNthNewest(i)
			checkError(err)
			if strings.Contains(details.Name, patternToDload) {
				checkError(pb.DownloadNth(i))
			}
		}
	} else if nToDload > 0 {
		checkError(pb.DownloadNth(nToDload))
	} else if list {
		n, err := pb.Count()
		checkError(err)
		if flag.Arg(0) != "" {
			cnt, err := strconv.Atoi(flag.Arg(0))
			checkError(err)
			n = int(math.Min(float64(cnt), float64(n)))
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"#", "Name", "Size", "Type", "Timestamp", "Note"})
		for i := 1; i < n+1; i++ {
			details, err := pb.DetailsOfNthNewest(i)
			checkError(err)
			t.AppendRow(table.Row{i, details.Name, prettySize(details.Size), details.Type, details.Timestamp, details.Note})
		}
		t.SetStyle(table.StyleRounded)
		t.Render()
	} else if reconfigure {
		checkError(config.PromptAndPersist())
	} else {
		fmt.Println(usage)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func prettySize(bytes int64) string {
	suffixes := [5]string{"B", "KiB", "MiB", "GiB", "TiB"}
	logB1024 := int(math.Floor(math.Log(float64(bytes)) / math.Log(1024)))
	suffix := suffixes[logB1024]
	num := float64(bytes) / math.Pow(1024, float64(logB1024))
	return fmt.Sprintf("%.2f %s", num, suffix)
}
