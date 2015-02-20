package main

import (
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/diego-cli-plugin/diego_support"
	"github.com/cloudfoundry-incubator/diego-cli-plugin/docker"
	"github.com/cloudfoundry-incubator/diego-cli-plugin/utils"
	"github.com/cloudfoundry/cli/plugin"
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
				Name:     "set-health-check",
				HelpText: "set health_check_type flag to either port or none",
				UsageDetails: plugin.Usage{
					Usage: `cf set-health-check APP_NAME port /
 cf set-health-check APP_NAME none`,
				},
			},
			{
				Name:     "get-health-check",
				HelpText: "get health_check_type flag of an app",
				UsageDetails: plugin.Usage{
					Usage: "cf get-health-check APP_NAME",
				},
			},
			{
				Name:     "docker-push",
				HelpText: "push a docker image from docker hub as an app",
				UsageDetails: plugin.Usage{
					Usage: `cf docker-push APP_NAME DOCKER_IMAGE [OPTIONS]

Options
-c         : Startup command, set to 'null' to reset to default start command
--no-start : Do not start an app after pushing
--no-route : Do not map a route to this app and remove routes from previous pushes of this app
`,
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
		c.isDiegoEnabled(cliConnection, args[1])
	} else if args[0] == "docker-push" && len(args) >= 3 {
		c.dockerPush(cliConnection, args)
	} else if args[0] == "set-health-check" && len(args) == 3 && (args[2] == "port" || args[2] == "none") {
		c.setHealthCheck(cliConnection, args[1], args[2])
	} else if args[0] == "get-health-check" && len(args) == 2 {
		c.getHealthCheck(cliConnection, args[1])
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

	if output, err = d.SetDiegoFlag(appGuid, on); err != nil {
		exitWithError(err, output)
	}

	fmt.Printf("Diego support for %s is set to %t\n\n", appName, on)
}

func (c *DiegoBeta) isDiegoEnabled(cliConnection plugin.CliConnection, appName string) {
	var err error
	var output []string
	var appGuid string

	d := diego_support.NewDiegoSupport(cliConnection)
	u := utils.NewUtils(cliConnection)

	if appGuid, err, output = u.GetAppGuid(appName); err != nil {
		exitWithError(err, output)
	}

	var result bool
	if result, err, output = d.HasDiegoEnabled(appGuid); err != nil {
		exitWithError(err, output)
	}

	fmt.Println(result)
}

func (c *DiegoBeta) dockerPush(cliConnection plugin.CliConnection, args []string) {
	var err error
	var space, spaceGuid string
	var output []string

	u := utils.NewUtils(cliConnection)
	if space, err, output = u.GetTargetSpace(); err != nil {
		exitWithError(err, output)
	}

	if spaceGuid, err, output = u.GetSpaceGuid(space); err != nil {
		exitWithError(err, output)
	}

	appName := args[1]
	dockerImg := args[2]

	d := docker.NewDocker(cliConnection)

	//creating app
	fmt.Println("Creating app", appName, "...")
	if output, err = d.CreateApp(appName, dockerImg, spaceGuid); err != nil {
		exitWithError(err, output)
	}

	appGuid, err, output := u.GetAppGuid(appName)
	if err != nil {
		exitWithError(err, output)
	}
	sayOk()

	if isFlagExist(args[3:], "-c") {
		command := getFlagValue(args, "-c")
		fmt.Println("Updating start command: " + command)
		if command == "null" {
			command = ""
		}
		if output, err = u.UpdateApp(appGuid, "command", command); err != nil {
			exitWithError(err, output)
		}
		sayOk()
	}

	//creating route
	var domain string
	if isFlagExist(args[3:], "--no-route") {
		fmt.Println("Removing app routes if any ...")
		if output, err = u.DetachAppRoutes(appGuid); err != nil {
			exitWithError(err, output)
		}
		sayOk()
		return
	} else {
		fmt.Println("Creating route for", appName, "...")
		if domain, err, output = u.FindDomain(); err != nil {
			exitWithError(err, output)
		}

		if output, err = u.CreateRoute(space, domain, appName); err != nil {
			exitWithError(err, output)
		}
		fmt.Println("Route " + appName + "." + domain + " created")
		sayOk()
	}

	//mapping route
	fmt.Println("Mapping route to", appName, "...")
	if output, err = u.MapRoute(appName, domain, appName); err != nil {
		exitWithError(err, output)
	}
	fmt.Println("Mapped " + appName + "." + domain + " route to " + appName)
	sayOk()

	//starting app
	if isFlagExist(args[3:], "--no-start") {
		fmt.Println("Stop operation before starting '" + appName + "'")
		sayOk()
	} else {
		fmt.Println("Start app", appName, "...")
		if output, err = u.StartApp(appName); err != nil {
			exitWithError(err, output)
		}
	}
}

func (c *DiegoBeta) setHealthCheck(cliConnection plugin.CliConnection, appName string, value string) {
	u := utils.NewUtils(cliConnection)

	appGuid, err, output := u.GetAppGuid(appName)
	if err != nil {
		exitWithError(err, output)
	}

	fmt.Println("Setting health_check_type for " + appName + " to '" + value + "'")
	if output, err = u.UpdateApp(appGuid, "health_check_type", value); err != nil {
		exitWithError(err, output)
	}
	sayOk()
}

func (c *DiegoBeta) getHealthCheck(cliConnection plugin.CliConnection, appName string) {
	u := utils.NewUtils(cliConnection)

	appGuid, err, output := u.GetAppGuid(appName)
	if err != nil {
		exitWithError(err, output)
	}

	fmt.Println("Getting health_check_type for " + appName)
	healthCheckType, output, err := u.GetHealthCheck(appGuid)
	if err != nil {
		exitWithError(err, output)
	}

	sayOk()
	fmt.Println("health_check_type for "+appName+":", healthCheckType)
}

func exitWithError(err error, output []string) {
	sayFailed()
	fmt.Println("Error: ", err)
	for _, str := range output {
		fmt.Println(str)
	}
	os.Exit(1)
}

func isFlagExist(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}

func getFlagValue(args []string, flag string) string {
	for i, arg := range args {
		if arg == flag {
			if len(args) >= i+1 {
				return args[i+1]
			}
			break
		}
	}
	return ""
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
