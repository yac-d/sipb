package main

import (
	"fmt"
	"os"

	"github.com/yac-d/sbc/conf"
	"github.com/yac-d/sbc/pastebin"
)

const usage = `Simple Bin Commands
USAGE: sbc [options] [arguments]
Options:
    -l, --list [n]
        Displays details of n most recent files
    -f, -u, --upload <File> [Note]
        Uploads File to the pastebin with Note
    -o <Pattern>
        Downloads the most recent file whose name matches the given pattern
    -d, --download <n>
        Downloads the nth most recent file
    -c, --config
        Reconfigure sbc settings`

func main() {
	config, err := conf.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pb := pastebin.New(config.Server)
	fmt.Println(pb.Count())

	var path string
	fmt.Scanf("%s", &path)
	fmt.Println(pb.Upload(path, "This is a note"))
}
