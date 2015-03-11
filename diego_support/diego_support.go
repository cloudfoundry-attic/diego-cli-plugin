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
	CheckDiegoError(string) error
}

type diegoSupport struct {
	cli plugin.CliConnection
}

type AppSummary struct {
	Diego bool `json:"diego,omitempty"`
}

type diegoError struct {
	Code        int64  `json:"code;omitempty"`
	Description string `json:"description;omitempty"`
	ErrorCode   string `json:"error_code"`
}

func NewDiegoSupport(cli plugin.CliConnection) DiegoSupport {
	return &diegoSupport{
		cli: cli,
	}
}

func (d *diegoSupport) CheckDiegoError(jsonRsp string) error {
	b := []byte(jsonRsp)
	diegoErr := diegoError{}
	err := json.Unmarshal(b, &diegoErr)
	if err != nil {
		return err
	}

	if diegoErr.ErrorCode != "" || diegoErr.Code != 0 {
		return errors.New(diegoErr.ErrorCode + " - " + diegoErr.Description)
	}

	return nil
}

func (d *diegoSupport) SetDiegoFlag(appGuid string, enable bool) ([]string, error) {
	output, err := d.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps/"+appGuid, "-X", "PUT", "-d", `{"diego":`+strconv.FormatBool(enable)+`}`)
	if err != nil {
		return output, err
	}

	if err = d.CheckDiegoError(output[0]); err != nil {
		return output, err
	}

	return output, nil
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
