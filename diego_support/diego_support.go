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
	SetDiegoFlag(string, bool) ([]string, error)
	HasDiegoEnabled(string) (bool, error, []string)
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

func (d *diegoSupport) SetDiegoFlag(appGuid string, enable bool) ([]string, error) {
	return d.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps/"+appGuid, "-X", "PUT", "-d", `{"diego":`+strconv.FormatBool(enable)+`}`)
}

func (d *diegoSupport) HasDiegoEnabled(appGuid string) (bool, error, []string) {

	result, err := d.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps/"+appGuid+"/summary")
	if err != nil {
		return false, err, result
	}

	if !strings.Contains(result[0], `"diego": `) {
		return false, errors.New(fmt.Sprintf("%s\nJSON:\n", "'diego' flag is not found in json response")), result
	}

	b := []byte(result[0])
	summary := AppSummary{}
	err = json.Unmarshal(b, &summary)
	if err != nil {
		return false, err, result
	}

	return summary.Diego, nil, []string{}
}
