# sjtuxiangxiang
[![label](https://img.shields.io/badge/SJTU-SEIEE-blue.svg)]()
[![label](https://img.shields.io/badge/SJTU-IEEE-brightgreen.svg)]()

This is the project for object-oriented software engineering course.

## How to run the server?

Mac OS / Linux user can open your terminal and run commands as below:

1. Download the source code

```
$ git clone https://github.com/Faldict/sjtuxiangxiang.git
$ cd sjtuxiangxiang
```

2. run the server

If you have installed [golang build tools](http://golangtc.com/download), you can directly run the command:

```
$ ./run.sh
```

Otherwise, I have pre-built and uploaded an execute bin. You can download a pre-build
binary file and may run this command:

```
$ ./main
```

The binary file's name is not so standard that I will rename it some time later.

If you are windows user, you should follow above steps and run the similar commands.

## Code Structure

I have only written 3 lines backend, and a hello-world front-end. The backend is
organized by `main.go` while the front-end
is in the static file folder.
