goconfig
========

[![Travis Build Status](https://travis-ci.org/muja/goconfig.svg?branch=master)](https://travis-ci.org/muja/goconfig)

# Table of contents

1. Introduction
2. Usage
3. Contributing
4. Reporting bugs

-------------------

# 1. Introduction

This project parses config files that have the same syntax as gitconfig files. It is a
minimal parser, and maybe a writer sometime in the future. It has no knowledge of git-specific
keys and as such, does not provide any convenience methods like  `config.GetUserName()`.
For these, look into [go-gitconfig](https://github.com/tcnksm/go-gitconfig)

Most of the code was copied and translated to Go from [git/config.c](https://github.com/git/git/blob/95ec6b1b3393eb6e26da40c565520a8db9796e9f/config.c)

# 2. Usage

Currently, there is only one function: `Parse`.

```go
import "os/user"
import "path/filepath"
import "io/ioutil"
import "github.com/muja/goconfig"

user, _ := user.Current()
// don't forget to handle error!
gitconfig := filepath.Join(user.HomeDir, ".gitconfig")
bytes, _ := ioutil.ReadFile(gitconfig)

config, lineno, err := goconfig.Parse(bytes)
if err != nil {
  // Note: config is non-nil and contains successfully parsed values
  log.Fatalf("Error on line %d: %v.\n", err)
}
fmt.Println(config["user.name"])
fmt.Println(config["user.email"])
```

# 3. Contributing

Contributions are welcome! Fork -> Push -> Pull request.

# 4. Bug report / suggestions

Just create an issue! I will try to reply as soon as possible.
