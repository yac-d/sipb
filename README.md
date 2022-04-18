# sipb

SImple Image Paste Bin for the MysoreLUG

![Screenshot](/assets/Pastebin.png "Cool, right?")

---

For all server options and how to build natively, see [this](https://github.com/Eeshaan-rando/sipb/blob/main/server/README.md).

For running in a Docker container (recommended), do the following:
- Build the container with `docker build -t sipb .` in the repository's root directory
- Run the container forwarding port 80 and using `./bin` as the directory for bin contents \
`docker run -p 80:80 -v $PWD/bin:/var/www/bin sipb`
