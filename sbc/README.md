The CLI client for SIPB is a script called "sbc", for "simple bin commands". 

# Installation:

For now, clone the repo, get the script and move it somewhere into your $PATH (Personally, I use ~/.local/bin) for all of my scripts. 

After installation, you have to create ~/.config/sbc/sbcrc, in which you have to specify the following in the following manner:

```
username:dolt
password:password123
domain:https://iamagiantdolt.com:67
```

# TODO:

* Clean up the script, some parts are hacky (OK, it is all very hacky, but it works, at least for me).
* Create a makefile that does a proper installation of the script
* Create a manual page

# Dependencies:
* Coreutils
* cURL
