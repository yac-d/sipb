# CLI Client

```
USAGE: sbc [options] [arguments]
Options:
    -l, --list [n]
        Displays details of n most recent files
    -f, -u, --upload <File> [Note]
        Uploads File with Note
    -o <Pattern>
        Downloads all files whose name matches the given pattern
    -d, --download <n>
        Downloads the nth most recent file
    -c, --config
        Reconfigure sbc settings"""
```

To use, place the program in a location on `$PATH`.

A configuration file is stored at `~/.config/sbcrc`.
