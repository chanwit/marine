package main

import (
	"os"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/chanwit/marine"
	"github.com/codegangsta/cli"
)

func prepare(c *cli.Context) {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	filename := c.String("image")
	basename := path.Base(filename)
	name := strings.SplitN(basename, "-", 2)[0]

	outfile := c.String("out")
	if outfile == "" {
		outfile = name + ".ova"
	}

	force := c.Bool("force")

	exist, err := marine.Exist(name)
	if exist && force {
		marine.Remove(name)
	}

	if (err == nil && exist == false) || force {
		_, err = marine.Import(filename,
			512,
			"docker",
			"golang")
		err = marine.Export(name, outfile)
		if err != nil {
			log.Errorf("Export error: %s", err)
		}
	} else {
		log.Infof("Image existed: %s", name)
	}
}

var flImage = cli.StringFlag{
	Name:  "image, i",
	Usage: "image file name <.ova>",
}

var flForce = cli.BoolFlag{
	Name:  "force, f",
	Usage: "force action",
}

var flOutfile = cli.StringFlag {
	Name: "out, o",
	Usage: "output image file",
	Value: "",
}

func main() {

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "functional testing for Swarm"
	app.Version = "0.1.0"

	app.Author = ""
	app.Email = ""

	app.Commands = []cli.Command{
		{
			Name:      "prepare",
			ShortName: "p",
			Usage:     "prepare a base image",
			Flags: []cli.Flag{
				flImage, flForce, flOutfile,
			},
			Action: prepare,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
