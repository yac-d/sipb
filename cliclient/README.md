# CLI Client

So far, only a small python API has been written. It depends on the `requests` library. Refer to the following to use it.

---

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

Returns the file count on the server. Raises an exception for HTTP response codes <400.

```
Pastebin.detailsOfNthNewest(n)
```

Returns a dictionary containing the URL `Path`, `Size` in bytes, and MIME `Type` of the last nth uploaded file.
0<`n`<=`Pastebin.count()`. Raises an exception for HTTP response codes <400.

```
Pastebin.upload(filepath)
```

Uploads the file at `filepath` and returns the number of bytes truncated. Raises an exception for HTTP response codes <400 and !=413.

---
