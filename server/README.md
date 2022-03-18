# Server

The backend for `sipb` is written in Golang. Its config file contains multiple elements, each of which are mandatory:

`WebpageDir` is the directory on the server's filesystem that contains files to be served.

`BinPath` is the absolute URL path to the bin. The server creates the directory if necessary. This field must always be expressed as an absolute path.
The server evaluates the path to the bin on the local filesystem by prepending `WebpageDir`.

As an example, let's say `WebpageDir` is `/var/www/html`. The contents of this directory correspond to `example.com/`.
So, if your bin is to be at `example.com/content/bin`, you would set `BinPath` to `/content/bin`.
