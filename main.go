package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/simonleung8/diego-beta/diego_support"
	"github.com/simonleung8/diego-beta/docker"
	"github.com/simonleung8/diego-beta/utils"
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
			{
				Name:     "docker-push",
				HelpText: "push a docker image from docker hub as an app",
				UsageDetails: plugin.Usage{
					Usage: "cf docker-push APP_NAME DOCKER_IMAGE",
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
	} else if args[0] == "docker-push" && len(args) == 3 {
		c.dockerPush(cliConnection, args)
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
	u := utils.NewUtils(cliConnection)

	appGuid, err, output := u.GetAppGuid(appName)
	if err != nil {
		exitWithError(err, output)
	}

	output, err = d.SetDiegoFlag(appGuid, on)

	if err != nil {
		exitWithError(err, output)
	}

	fmt.Printf("Diego support for %s is set to %t\n\n", appName, on)
}

func (c *DiegoBeta) checkDiegoSupport(checkEnable bool, cliConnection plugin.CliConnection, appName string) {
	var err error
	var output []string

	d := diego_support.NewDiegoSupport(cliConnection)
	u := utils.NewUtils(cliConnection)

	appGuid, err, output := u.GetAppGuid(appName)
	if err != nil {
		exitWithError(err, output)
	}

	var result bool
	result, err, output = d.HasDiegoEnabled(appGuid)

	if err != nil {
		exitWithError(err, output)
	}

	fmt.Println(checkEnable == result)
}

func (c *DiegoBeta) dockerPush(cliConnection plugin.CliConnection, args []string) {
	u := utils.NewUtils(cliConnection)
	space, err, output := u.GetTargetSpace()
	if err != nil {
		exitWithError(err, output)
	}

	spaceGuid, err, output := u.GetSpaceGuid(space)
	if err != nil {
		exitWithError(err, output)
	}

	appName := args[1]
	dockerImg := args[2]

	d := docker.NewDocker(cliConnection)

	fmt.Println("Creating app", appName, "...")
	output, err = d.CreateApp(appName, dockerImg, spaceGuid)
	if err != nil {
		exitWithError(err, output)
	}
	sayOk()

	fmt.Println("Creating route for", appName, "...")
	domain, err, output := u.FindDomain()
	if err != nil {
		exitWithError(err, output)
	}
	sayOk()

	output, err = u.CreateRoute(space, domain, appName)
	if err != nil {
		exitWithError(err, output)
	}
	fmt.Println("Route " + appName + "." + domain + " created")
	sayOk()

	fmt.Println("Mapping route to", appName, "...")
	output, err = u.MapRoute(appName, domain, appName)
	if err != nil {
		exitWithError(err, output)
	}
	fmt.Println("Mapped " + appName + "." + domain + " route to " + appName)
	sayOk()

	fmt.Println("Start app", appName, "...")
	output, err = u.StartApp(appName)
	if err != nil {
		exitWithError(err, output)
	}
}

func exitWithError(err error, output []string) {
	sayFailed()
	fmt.Println("Error: ", err)
	for _, str := range output {
		fmt.Println(str)
	}
	os.Exit(1)
}

func say(message string, color uint, bold int) string {
	return fmt.Sprintf("\033[%d;%dm%s\033[0m", bold, color, message)
}

func sayOk() {
	fmt.Println(say("Ok\n", 32, 1))
}

func sayFailed() {
	fmt.Println(say("FAILED", 31, 1))
}
