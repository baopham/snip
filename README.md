snip
============
Command line interface to save snippets, shortcuts, commands, etc.


Usage
--------
```
NAME:
   snip - Save snippets: commands, texts, emoji, etc.

USAGE:
   snip [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     add, a      snip add -k="port" -c="lsof -i :{p}" -desc="List processes listening on a particular port"
     search, s   snip search port
     execute, x  get snippet
     list, l     list all saved snippets
     remove, r   remove a saved snippet
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Add

```bash
snip add -k="port" -c="lsof -i :{p}" -desc="List processes listening on a particular port"
```

Use `{placeholder}` for placeholders. See [Execute](#execute) for more on this

### List

```bash
snip list
```

you should see:

![list](screenshots/list.png)

### Search

```bash
snip search port
```

you should see:

![search](screenshots/search.png)

### Execute

```bash
snip execute port p=9000
```

This will replace the placeholder `{p}` with `9000`:

![execute](screenshots/execute.png)

### Remove

```bash
snip remove port
```

To remove a snippet by keyword


Install
------
```
go get github.com/baopham/snip
```

Make sure you have `$GOPATH` set and `$GOPATH/bin` is in the `$PATH`, e.g.:

```bash
export GOPATH=$HOME/Projects/Go
export PATH=$PATH:/usr/local/opt/go/libexec/bin:$GOPATH/bin
```

Requirements
-------------
* Go ^1.8

License
--------
MIT

Author
-------
Bao Pham
