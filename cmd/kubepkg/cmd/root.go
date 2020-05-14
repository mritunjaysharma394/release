/*
Copyright 2020 The Kubernetes Authors.

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

package cmd

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"k8s.io/release/pkg/kubepkg"
	"k8s.io/release/pkg/kubepkg/options"
	"k8s.io/release/pkg/log"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "kubepkg",
	Short:             "kubepkg",
	PersistentPreRunE: initLogging,
}

var (
	opts     *options.Options = options.New()
	logLevel string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
	packages := []string{}
	rootCmd.PersistentFlags().StringArrayVar(
		&packages,
		"packages",
		opts.Packages(),
		"packages to build",
	)

	channels := []string{}
	rootCmd.PersistentFlags().StringArrayVar(
		&channels,
		"channels",
		opts.Channels(),
		"channels to build for",
	)

	architectures := []string{}
	rootCmd.PersistentFlags().StringArrayVar(
		&architectures,
		"arch",
		opts.Architectures(),
		"architectures to build for",
	)

	kubeVersion := ""
	rootCmd.PersistentFlags().StringVar(
		&kubeVersion,
		"kube-version",
		opts.KubeVersion(),
		"Kubernetes version to build",
	)

	revision := ""
	rootCmd.PersistentFlags().StringVar(
		&revision,
		"revision",
		opts.Revision(),
		"deb package revision.",
	)

	cniVersion := ""
	rootCmd.PersistentFlags().StringVar(
		&cniVersion,
		"cni-version",
		opts.CNIVersion(),
		"CNI version to build",
	)

	criToolsVersion := ""
	rootCmd.PersistentFlags().StringVar(
		&criToolsVersion,
		"cri-tools-version",
		opts.CRIToolsVersion(),
		"CRI tools version to build",
	)

	releaseDownloadLinkBase := ""
	rootCmd.PersistentFlags().StringVar(
		&releaseDownloadLinkBase,
		"release-download-link-base",
		opts.ReleaseDownloadLinkBase(),
		"release download link base",
	)

	templateDir := ""
	rootCmd.PersistentFlags().StringVar(
		&templateDir,
		"template-dir",
		opts.TemplateDir(),
		"template directory",
	)

	specOnly := false
	rootCmd.PersistentFlags().BoolVar(
		&specOnly,
		"spec-only",
		opts.SpecOnly(),
		"only create specs instead of building packages",
	)

	rootCmd.PersistentFlags().StringVar(
		&logLevel,
		"log-level",
		"info",
		"the logging verbosity, either 'panic', 'fatal', 'error', 'warn', 'warning', 'info', 'debug' or 'trace'",
	)

	opts = opts.WithPackages(packages...).
		WithChannels(channels...).
		WithArchitectures(architectures...).
		WithKubeVersion(kubeVersion).
		WithRevision(revision).
		WithCNIVersion(cniVersion).
		WithCRIToolsVersion(criToolsVersion).
		WithReleaseDownloadLinkBase(releaseDownloadLinkBase).
		WithTemplateDir(templateDir).
		WithSpecOnly(specOnly)
}

func initLogging(*cobra.Command, []string) error {
	return log.SetupGlobalLogger(logLevel)
}

func run(opts *options.Options) error {
	client := kubepkg.New(opts)
	builds, err := client.ConstructBuilds()
	if err != nil {
		return errors.Wrap(err, "running kubepkg")
	}
	return client.WalkBuilds(builds)
}
