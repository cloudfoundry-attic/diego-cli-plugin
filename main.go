package main

import (
	"fmt"
	"os"

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
			{
				Name:     "has-diego-enabled",
				HelpText: "Check if Diego support is enabled for an app",
				UsageDetails: plugin.Usage{
					Usage: "cf has-diego-enabled APP_NAME",
				},
			},
			{
				Name:     "has-diego-disabled",
				HelpText: "Check if Diego support is disabled for an app",
				UsageDetails: plugin.Usage{
					Usage: "cf has-diego-disabled APP_NAME",
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
		c.toggleDiegoSupport(true, cliConnection, args[1])
	} else if args[0] == "disable-diego" && len(args) == 2 {
		c.toggleDiegoSupport(false, cliConnection, args[1])
	} else if args[0] == "has-diego-enabled" && len(args) == 2 {
		c.checkDiegoSupport(true, cliConnection, args[1])
	} else if args[0] == "has-diego-disabled" && len(args) == 2 {
		c.checkDiegoSupport(false, cliConnection, args[1])
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

func (c *DiegoBeta) toggleDiegoSupport(on bool, cliConnection plugin.CliConnection, appName string) {
	d := diego_support.NewDiegoSupport(cliConnection)
	appGuid, err := d.GetAppGuid(appName)
	if err != nil {
		exitWithError(err)
	}

	err = d.SetDiegoFlag(appGuid, on)

	if err != nil {
		exitWithError(err)
	}

	fmt.Printf("Diego support for %s is set to %t\n\n", appName, on)
}

func (c *DiegoBeta) checkDiegoSupport(checkEnable bool, cliConnection plugin.CliConnection, appName string) {
	var err error

	d := diego_support.NewDiegoSupport(cliConnection)
	appGuid, err := d.GetAppGuid(appName)
	if err != nil {
		exitWithError(err)
	}

	var result bool
	result, err = d.HasDiegoEnabled(appGuid)

	if err != nil {
		exitWithError(err)
	}

	fmt.Println(checkEnable == result)
}

func exitWithError(err error) {
	fmt.Println("Error: ", err)
	os.Exit(1)
}
