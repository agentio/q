package gcloud

import (
	"encoding/json"
	"errors"
	"log"
	"os/exec"
	"strings"
)

func GetInfo(verbose bool) (*Info, error) {
	output, err := exec.Command("gcloud", "info", "--format=json").Output()
	if err != nil {
		return nil, err
	}
	if verbose {
		log.Printf("%s", string(output))
	}
	info := &Info{}
	if err = json.Unmarshal(output, info); err != nil {
		return nil, err
	}
	return info, nil
}

type Info struct {
	Basic  Basic  `json:"basic"`
	Config Config `json:"config"`
}

func (info *Info) Account() (string, error) {
	if info.Config.Account != "" {
		return info.Config.Account, nil
	}
	return "", errors.New("user required")
}

func (info *Info) Project() (string, error) {
	if info.Config.Project != "" {
		return info.Config.Project, nil
	}
	return "", errors.New("project required")
}

func (info *Info) RunRegion() (string, error) {
	if info.Config.Properties.Run != nil &&
		info.Config.Properties.Run["region"] != nil &&
		info.Config.Properties.Run["region"].Value != "" {
		return info.Config.Properties.Run["region"].Value, nil
	}
	return "", errors.New("run/region required")
}

type Basic struct {
	Version string `json:"version"`
}

type Config struct {
	Account    string     `json:"account"`
	Project    string     `json:"project"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	Core map[string]*Property `json:"core"`
	Run  map[string]*Property `json:"run"`
}

type Property struct {
	Value string `json:"value"`
}

func GetADCToken(verbose bool) (string, error) {
	output, err := exec.Command("gcloud", "auth", "application-default", "print-access-token").Output()
	if err != nil {
		return "", err
	}
	if verbose {
		log.Printf("%s", string(output))
	}
	token := string(output)
	token = strings.TrimSpace(token)
	return token, nil
}
