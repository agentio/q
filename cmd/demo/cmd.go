package demo

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/agentio/q/pkg/gcloud"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Set up a sample managed service",
		RunE:  action,
	}
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	info, err := gcloud.GetInfo(false)
	if err != nil {
		return err
	}
	type DemoContext struct {
		Project string
	}
	c := &DemoContext{}
	if c.Project, err = info.Project(); err != nil {
		return err
	}

	dir := "stores-demo"
	if err = os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	files := []struct {
		name     string
		template string
	}{
		{"SETUP.sh", SETUP_SH},
		{"api_config.yaml", API_CONFIG_YAML},
		{"service.yaml", SERVICE_YAML},
		{"iam.yaml", IAM_YAML},
		{"CHECK.sh", CHECK_SH},
	}
	for _, f := range files {
		tmpl, err := template.New(f.name).Parse(f.template)
		if err != nil {
			return err
		}
		file, err := os.Create(filepath.Join(dir, f.name))
		if err != nil {
			return err
		}
		defer file.Close()
		err = tmpl.Execute(file, c)
		if err != nil {
			return err
		}
	}

	response, err := http.Get("https://github.com/bobadojo/descriptor/raw/main/descriptor.pb")
	if err != nil {
		return err
	}
	output, err := os.Create(filepath.Join(dir, "descriptor.pb"))
	if err != nil {
		return err
	}
	_, err = output.ReadFrom(response.Body)
	if err != nil {
		return err
	}

	fmt.Printf("To run the demo, see the SETUP.sh script in %s\n", dir)
	return nil
}

const SETUP_SH = `#!/bin/sh

echo 'Enabling the Service Control API.'
gcloud services enable servicecontrol.googleapis.com

echo 'Enabling the Service Management API.'
gcloud services enable servicemanagement.googleapis.com

echo 'Enabling the Cloud Run Admin API.'
gcloud services enable run.googleapis.com

echo 'Creating the service from the service descriptor and API config file.'
gcloud endpoints services deploy descriptor.pb api_config.yaml

echo 'Creating a service account to run the server container.'
gcloud iam service-accounts create stores

echo 'Giving the service account roles to call the Service Control APIs.'
gcloud projects add-iam-policy-binding {{ .Project }} \
	--member serviceAccount:stores@{{ .Project }}.iam.gserviceaccount.com \
	--role roles/servicemanagement.serviceController
gcloud projects add-iam-policy-binding {{ .Project }} \
	--member serviceAccount:stores@{{ .Project }}.iam.gserviceaccount.com \
  --role roles/cloudtrace.agent

echo 'Creating the Cloud Run container.'
gcloud run services replace service.yaml

echo 'Configuring IAM to allow outside access to the container.'
gcloud run services set-iam-policy --quiet stores iam.yaml
`

const API_CONFIG_YAML = `type: google.api.Service
config_version: 3

#
# Name of the service configuration.
#
name: stores.endpoints.{{ .Project }}.cloud.goog

#
# API title to appear in the user interface (Google Cloud Console).
#
title: Boba Dojo Stores API
apis:
- name: bobadojo.stores.v1.Stores

#
# API usage restrictions.
#
usage:
  rules:
  # ListStores methods can be called without an API Key.
  - selector: bobadojo.stores.v1.Stores.ListStores
    allow_unregistered_calls: true
  - selector: bobadojo.stores.v1.Stores.GetStore
    allow_unregistered_calls: true
  - selector: bobadojo.stores.v1.Stores.FindStores
    allow_unregistered_calls: true
`

const SERVICE_YAML = `apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: stores
spec:
  template:
    spec:
      serviceAccountName: stores@{{ .Project }}.iam.gserviceaccount.com
      containers:
      - image: gcr.io/endpoints-release/endpoints-runtime:2
        name: espv2
        args:
        - --listener_port=8081
        - --backend=grpc://localhost:8080
        - --service=stores.endpoints.{{ .Project }}.cloud.goog
        - --rollout_strategy=managed
        ports:
        - name: http1
          containerPort: 8081
      - image: us-west1-docker.pkg.dev/bobadojo/stores/stores:latest
        name: stores
`

const IAM_YAML = `bindings:
- members:
  - allUsers
  role: roles/run.invoker
version: 1
`

const CHECK_SH = `#/bin/sh

gcloud projects get-iam-policy {{ .Project }} \
--flatten='bindings[].members' \
--format='table(bindings.role)' \
--filter='bindings.members:stores@{{ .Project }}.iam.gserviceaccount.com'
`
