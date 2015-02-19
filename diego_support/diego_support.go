package diego_support

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

type DiegoSupport interface {
	EnableDiego(string) error
	DisableDiego(string) error
}

type diegoSupport struct {
	cli plugin.CliConnection
}

func NewDiegoSupport(cli plugin.CliConnection) DiegoSupport {
	return &diegoSupport{
		cli: cli,
	}
}

func (d *diegoSupport) EnableDiego(app string) error {
	appGuid, err := d.getAppGuid(app)
	if err != nil {

		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	_, err = d.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps/"+appGuid, "-X", "PUT", "-d", `{"diego":true}`)
	return err
}

func (d *diegoSupport) DisableDiego(app string) error {
	appGuid, err := d.getAppGuid(app)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	_, err = d.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps/"+appGuid, "-X", "PUT", "-d", `{"diego":false}`)
	return err
}

func (d *diegoSupport) getAppGuid(appName string) (string, error) {
	result, err := d.cli.CliCommandWithoutTerminalOutput("app", appName, "--guid")
	if err != nil {
		if strings.Contains(result[0], "FAILED") {
			return "", errors.New("App " + appName + " not found.")
		}

		return "", err
	}

	return result[0], nil
}
