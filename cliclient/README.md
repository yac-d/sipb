# CLI Client

```
USAGE: sbc [options] [arguments]
Options:
    -l, --list [n]
        Displays details of n most recent files
    -f, -u, --upload <File>...
        Uploads the given files to the pastebin
    -o, -d, --download <Pattern>
        Downloads all files whose name matches the given pattern
```

To use, clone the repository and make the script available to run through an alias or with `$PATH`.

Depends on `prettytable` and `sipb_api`.

A configuration file is stored at `~/.config/sbcrc`.

---

## Python API

Depends on the `requests` library.

```
class Pastebin(host, port=80, basicauth=None, https=False, pastebinroot="/")
```

Creates a pastebin instance.

`host`, `port`, and `https` are fairly self-explanatory. `basicauth` is a tuple in the form `(username, password)` for HTTP basic authentication.

`pastebinroot` is used only when a proxy forwards a particular path to the SIPB server's root.
An SIPB server serves the webpage at `/pastebin`. So, for example, if the pastebin's webpage through your proxy is at `my-domain.com/irc/pastebin/`, `pastebinroot` should be set to `/irc/`. 
`pastebinroot` will be prepended (opposite of "append"; why isn't this a definition of the word) to the URL paths of requests made.

```
Pastebin.count()
```

Returns the file count on the server. Raises an exception for HTTP response codes >=400.

```
Pastebin.detailsOfNthNewest(n)
```

Returns a dictionary containing the URL `Path`, `Size` in bytes, and MIME `Type` of the last nth uploaded file.
0<`n`<=`Pastebin.count()`. Raises an exception for HTTP response codes >=400.

```
Pastebin.upload(filepath)
```

Uploads the file at `filepath` and returns the number of bytes truncated. Raises an exception for HTTP response codes >=400 and !=413.
