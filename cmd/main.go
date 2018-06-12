package main

import (
  "os"
  "github.com/codegangsta/cli"
)

func main() {
  cli.AppHelpTemplate = `{{.Name}} v{{.Version}} - {{.Usage}}

USAGE:
   {{.Name}} {{if .Flags}}[options] {{end}}<template ..>

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

  app := cli.NewApp()
  app.Name = "freedom-speedtest"
  app.Usage = "test internet bandwidth speed"
  app.Version = "VERSION"

  app.Flags = []cli.Flag {
    cli.StringFlag{
      Name: "output, o",
      Value: "",
      Usage: "append to the file as json format",
    },
    cli.StringFlag{
      Name: "file, f",
      Value: "",
      Usage: "read urls from file",
    },
  }

  app.Action = func(c *cli.Context) {
    SpeedTest(c.String("file"), c.String("output"))
  }

  app.Run(os.Args)
}
