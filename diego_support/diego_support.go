package diego_support

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

type DiegoSupport interface {
	SetDiegoFlag(string, bool) error
	HasDiegoEnabled(string) (bool, error)
	GetAppGuid(string) (string, error)
}

type diegoSupport struct {
	cli plugin.CliConnection
}

type AppSummary struct {
	Diego bool `json:"diego"`
}

func NewDiegoSupport(cli plugin.CliConnection) DiegoSupport {
	return &diegoSupport{
		cli: cli,
	}
}

func (d *diegoSupport) SetDiegoFlag(appGuid string, enable bool) error {
	_, err := d.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps/"+appGuid, "-X", "PUT", "-d", `{"diego":`+strconv.FormatBool(enable)+`}`)
	return err
}

func (d *diegoSupport) HasDiegoEnabled(appGuid string) (bool, error) {

	result, err := d.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps/"+appGuid+"/summary")
	if err != nil {
		return false, err
	}

	if !strings.Contains(result[0], `"diego": `) {
		return false, errors.New(fmt.Sprintf("%s\nJSON:\n%v\n\n", "'diego' flag is not found in json response", result))
	}

	b := []byte(result[0])
	summary := AppSummary{}
	err = json.Unmarshal(b, &summary)
	if err != nil {
		return false, err
	}

	return summary.Diego, nil
}

func (d *diegoSupport) GetAppGuid(appName string) (string, error) {
	result, err := d.cli.CliCommandWithoutTerminalOutput("app", appName, "--guid")
	if err != nil {
		if strings.Contains(result[0], "FAILED") {
			return "", errors.New("App " + appName + " not found.")
		}

		return "", err
	}

	return strings.TrimSpace(result[0]), nil
}
