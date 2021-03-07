// +build mage

/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"

	"sigs.k8s.io/release-utils/mage"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Verify

const (
	binDir    = "bin"
	scriptDir = "scripts"
)

var boilerplateDir = filepath.Join(scriptDir, "boilerplate")

// Verify runs repository verification scripts
func Verify() error {
	fmt.Println("Running copyright header checks...")
	err := mage.VerifyBoilerplate("", binDir, boilerplateDir, false)
	if err != nil {
		return err
	}

	fmt.Println("Running external dependency checks...")
	err = mage.VerifyDeps("", "", "", true)
	if err != nil {
		return err
	}

	fmt.Println("Running go module linter...")
	err = mage.VerifyGoMod(scriptDir)
	if err != nil {
		return err
	}

	fmt.Println("Running golangci-lint...")
	err = mage.RunGolangCILint("", false)
	if err != nil {
		return err
	}

	fmt.Println("Running go build...")
	err = mage.VerifyBuild(scriptDir)
	if err != nil {
		return err
	}

	return nil
}

// Default targets

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "MyApp", ".")
	return cmd.Run()
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	return os.Rename("./MyApp", "/usr/bin/MyApp")
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "get", "github.com/stretchr/piglatin")
	return cmd.Run()
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("MyApp")
}
