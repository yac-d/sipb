# Server

The backend for `sipb` is written in Golang.

## Configuration

The config file `config.yaml` contains multiple elements, each of which are mandatory:

- `WebpageDir` is the directory on the server's filesystem that contains files to be served.
- `BinPath` is the absolute URL path to the bin. The server creates the directory if necessary. This field must always be expressed as an absolute path.
The server evaluates the path to the bin on the local filesystem by prepending `WebpageDir`. \
As an example, let's say `WebpageDir` is `/var/www/html`. The contents of this directory correspond to `example.com/`.
So, if your bin is to be at `example.com/content/bin`, you would set `BinPath` to `/content/bin`.
- `BindAddr` is the address the server should bind to. For the service to be available on your local network, this should be the server's IP address on the local network.
- `Port` is the network port for the server to bind to.
- `MaxFileCnt` is the maximum number of files the bin is allowed to store. When the limit is reached, the oldest file is removed after each upload. This can be set to `-1` for no limit.
- `MaxFileSize` is the maximum allowed file size in bytes. Files will be truncated if exceeding this limit. This can be set to `-1` for no limit.

`config.yaml` must be in the folder where the server binary is run.
Environment variables, if defined, can override these values. They are:

- `SIPB_WEBPAGE_DIR`
- `SIPB_BIN_PATH`
- `SIPB_BIND_ADDR`
- `SIPB_PORT`
- `SIPB_MAX_FILE_CNT`
- `SIPB_MAX_FILE_SIZE`

## Build and run

- Install `go`.
- Navigate to this directory in the repository.
- Build the binary with `go build`.
- If you have configured the server to bind to a privileged port (<1024), run `sudo setcap 'CAP_NET_BIND_SERVICE=+ep' ./sipb` to allow the server to run as a regular user. `setcap` can be installed from the appropriate package for your distribution (`libcap-progs` on openSUSE).
- Run the server with `./sipb`. Ensure your `config.yaml` is in the folder from which you are running the server.
- Launch a browser and go to `address-you-bound-to.com/pastebin`.
- Paste away!

## HTTP Requests

A `GET` request to `/retrieve/fileCount` returns the file count as a decimal-formatted string.

Fetching details of the last nth uploaded file is done by sending a `POST` request, whose content is n as decimal-formatted string, to `/retrieve`.
The response is in JSON format, with the following fields:

- `Path`: URL to the file itself
- `Size`: File size in bytes
- `Type`: MIME type of the file

Uploading a file is done by `POST`ing a multipart form data request to `/upload`. The key/name for the file in the form should be `file`.
