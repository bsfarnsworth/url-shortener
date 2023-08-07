# url-shortener

`url-shortener` is a project to build an application that creates short URLs from long ones.

The repository is reached at `https://github.com/bsfarnsworth/url-shortener`.

## Overview

This [document](Wee.md) describes the concept and design of **wee**, a nickname for this application.

## Install

- Create a new directory and clone the project from git: e.g.

```
git clone ssh://git@github.com/bsfarnsworth/url-shortener .
cd url-shortener
```

- Build the executable

```
bin/build.sh
```

- Run it

```
./wee &
```

- Open a browser window to the developer deployment

```
localhost:3000
```

- Or, to the production cloud deployment

```
https://wee.fly.dev
```


## Usage

- Open a terminal to the directory where you installed `url-shortener`.

```
cd tests
```

- Choose either the developer (testconfig.cfg) or production (prodconfig.cfg) configuration

```
source testconfig.cfg
```

or

```
source prodconfig.cfg
```

- Shorten any URL of your choice.  For instance,

```
./shorten.sh 'https://apod.nasa.gov/apod/'
```

The JSON `weeUrl` value is your shortened URL.  On my development deploy it came up as `xqltu46`.

```
{"token":"79f99aed-2779-4bea-905c-5def4750340e","weeUrl":"xqltu46"
```

- Confirm the URL of a weeURL:

```
./lengthen.sh xqltu46
```

- Retire the URL.  You use the `token` value:

```
./retire.sh 79f99aed-2779-4bea-905c-5def4750340e
```

- Make use of the wee URL.  If we are using the development server,

```
localhost:3000/xqltu46
```

or the cloud server:

```
https://wee.fly.dev/xqltu46
```

*Important*, the weeURLs are used on the same server where they were shortened.  This is because the table of URLs lives on the server where it was created and they are not shared between the servers.

