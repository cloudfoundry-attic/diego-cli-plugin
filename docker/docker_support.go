// curl (make app)
// set diego flag
// create and map route
// start app

package docker

import "github.com/cloudfoundry/cli/plugin"

type Docker interface {
	CreateApp(string, string, string) ([]string, error)
}

type docker struct {
	cli plugin.CliConnection
}

func NewDocker(cli plugin.CliConnection) Docker {
	return &docker{
		cli: cli,
	}
}

func (d *docker) CreateApp(appName, dockerImg, spaceGuid string) ([]string, error) {
	return d.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps", "-X", "POST", "-d", `{"name":"`+appName+`","space_guid":"`+spaceGuid+`","docker_image":"`+dockerImg+`", "diego": true}`)
}
