### readenv

Tool for reading and executing command with .env

Make sure you have a SHELL environment variable.

Lines starting with `#` are ignoreds.


#### Install

```bash
go install github.com/alexisvisco/readenv@v0.0.7
```

Make sure ~/go/bin is in your path variable.


#### usage

```bash
readenv .env                                         make tests
         ^ your .env file can be any filename        ^ your command
```