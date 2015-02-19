package main

import (
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/simonleung8/diego-beta/diego_support"
)

type DiegoBeta struct{}

func (c *DiegoBeta) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "Diego-Beta",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "enable-diego",
				HelpText: "enable Diego support for an app",
				UsageDetails: plugin.Usage{
					Usage: "cf enable-diego APP_NAME",
				},
			},
			{
				Name:     "disable-diego",
				HelpText: "disable Diego support for an app",
				UsageDetails: plugin.Usage{
					Usage: "cf disable-diego APP_NAME",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(DiegoBeta))
}

func (c *DiegoBeta) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "enable-diego" && len(args) == 2 {
		c.switchDiegoSupport(true, cliConnection, args)
	} else if args[0] == "disable-diego" && len(args) == 2 {
		c.switchDiegoSupport(false, cliConnection, args)
	} else {
		c.showUsage(args)
	}
}

func (c *DiegoBeta) showUsage(args []string) {
	for _, cmd := range c.GetMetadata().Commands {
		if cmd.Name == args[0] {
			fmt.Println("Invalid Usage: \n", cmd.UsageDetails.Usage)
		}
	}
}

func (c *DiegoBeta) switchDiegoSupport(on bool, cliConnection plugin.CliConnection, args []string) {
	if on {
		d := diego_support.NewDiegoSupport(cliConnection)
		if err := d.EnableDiego(args[1]); err == nil {
			fmt.Println("Diego support enabled for " + args[1])
		} else {
			fmt.Println("Error: ", err)
		}
	} else {
		d := diego_support.NewDiegoSupport(cliConnection)
		if err := d.DisableDiego(args[1]); err == nil {
			fmt.Println("Diego support disabled for " + args[1])
		} else {
			fmt.Println("Error: ", err)
		}
	}
}
