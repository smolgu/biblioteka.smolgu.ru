// +build go1.3

// Copyright 2014 Kirill Zhuharev. All rights reserved.

package main

import (
	"os"

	"github.com/codegangsta/cli"

	"github.com/smolgu/lib/cmd"
	"github.com/smolgu/lib/modules/setting"
)

const APP_VER = "0.0.1.0001 Pre Alpha"

func init() {
	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "SmolSU library"
	app.Usage = "SmoSU library site"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.CmdWeb,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
