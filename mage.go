//+build mage

package main

import (
	"errors"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"

	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
)

// Start the grafana container.
func Start() error {
	_, err := os.Stat("started")
	if !os.IsNotExist(err) {
		return nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := sh.RunV("podman", "run", "-d", "-p", "3000:3000", "-e", "GF_DATAPROXY_LOGGING=true", "-e", "GF_LOG_LEVEL=debug", "-v", cwd+"/custom.ini:/etc/grafana/grafana.ini", "-v", cwd+"/dist:/var/lib/grafana/plugins/departureboard.io-datasource", "--name=grafana", "grafana/grafana:7.0.0-beta3"); err != nil {
		return err
	}
	return sh.RunV("touch", "started")
}

// Stop and remove the grafana container.
func Stop() error {
	_ = sh.RunV("podman", "kill", "grafana")
	if err := sh.RunV("podman", "rm", "grafana"); err != nil {
		return err
	}
	return sh.RunV("rm", "started")
}

// Restart the grafana container if there is a new build.
func Restart() error {
	mg.Deps(build.BuildAll, BuildFrontend)
	changed, err := target.Dir("started", "dist")
	if !(os.IsNotExist(err) || (err == nil && changed)) {
		return nil
	}
	if err := sh.RunV("podman", "restart", "grafana"); err != nil {
		return err
	}
	return sh.RunV("touch", "started")
}

func TestE2E() error {
	if _, ok := os.LookupEnv("NATIONALRAIL_API_KEY"); !ok {
		return errors.New("e2e tests require the NATIONALRAIL_API_KEY environment variable to be set")
	}
	return sh.RunV("go", "test", "-tags=e2e", "./...")
}
func Build() error {
	changed, err := target.Dir("dist", "src", "pkg")
	if !(os.IsNotExist(err) || (err == nil && changed)) {
		return nil
	}
	mg.Deps(build.BuildAll, BuildFrontend)
	return nil
}

func BuildFrontend() error {
	changed, err := target.Dir("dist/module.js", "src")
	if !(os.IsNotExist(err) || (err == nil && changed)) {
		return nil
	}
	return sh.RunV("yarn", "build")
}

var Default = Build
