/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package app

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"tkestack.io/tke/pkg/app/version"
)

var (
	progressMessage = color.GreenString("==>")
	usageTemplate   = fmt.Sprintf(`%s{{if .Runnable}}
  %s{{end}}{{if .HasAvailableSubCommands}}
  %s{{end}}{{if gt (len .Aliases) 0}}

%s
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

%s{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  %s {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

%s{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "%s --help" for more information about a command.{{end}}
`,
		color.CyanString("Usage:"),
		color.GreenString("{{.UseLine}}"),
		color.GreenString("{{.CommandPath}} [command]"),
		color.CyanString("Aliases:"),
		color.CyanString("Examples:"),
		color.CyanString("Available Commands:"),
		color.GreenString("{{rpad .Name .NamePadding }}"),
		color.CyanString("Flags:"),
		color.CyanString("Global Flags:"),
		color.CyanString("Additional help topics:"),
		color.GreenString("{{.CommandPath}} [command]"),
	)
)

// App is the main structure of a cli application.
// It is recommended that an app be created with the app.NewApp() function.
type App struct {
	basename    string
	name        string
	description string
	options     CliOptions
	runFunc     RunFunc
	silence     bool
	noVersion   bool
	commands    []*Command
}

// Option defines optional parameters for initializing the application
// structure.
type Option func(*App)

// WithOptions to open the application's function to read from the command line
// or read parameters from the configuration file.
func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

// RunFunc defines the application's startup callback function.
type RunFunc func(basename string) error

// WithRunFunc is used to set the application startup callback function option.
func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

// WithDescription is used to set the description of the application.
func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

// WithSilence sets the application to silent mode, in which the program startup
// information, configuration information, and version information are not
// printed in the console.
func WithSilence() Option {
	return func(a *App) {
		a.silence = true
	}
}

// WithNoVersion set the application does not provide version flag.
func WithNoVersion() Option {
	return func(a *App) {
		a.noVersion = true
	}
}

// NewApp creates a new application instance based on the given application name,
// binary name, and other options.
func NewApp(name string, basename string, opts ...Option) *App {
	a := &App{
		name:     name,
		basename: basename,
	}

	for _, o := range opts {
		o(a)
	}

	return a
}

// Run is used to launch the application.
func (a *App) Run() {
	initFlag()

	cmd := cobra.Command{
		Use:           FormatBaseName(a.basename),
		Long:          a.description,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.SetUsageTemplate(usageTemplate)
	cmd.SetOutput(os.Stdout)
	cmd.Flags().SortFlags = false
	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(a.name))
	}
	if a.runFunc != nil {
		cmd.Run = a.runCommand
	}

	if a.options != nil {
		if _, ok := a.options.(ConfigurableOptions); ok {
			addConfigFlag(a.basename, cmd.Flags())
		}
		a.options.AddFlags(cmd.Flags())
	}

	if !a.noVersion {
		version.AddFlags(cmd.Flags())
	}
	addHelpFlag(a.name, cmd.Flags())

	if err := cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

func (a *App) runCommand(cmd *cobra.Command, args []string) {
	printWorkingDir()
	if !a.noVersion {
		// display application version information
		version.PrintAndExitIfRequested(a.name)
	}
	if !a.silence {
		fmt.Printf("%v Starting %s...\n", progressMessage, a.name)
	}
	// merge configuration and print it
	if a.options != nil {
		if configurableOptions, ok := a.options.(ConfigurableOptions); ok {
			if errs := configurableOptions.ApplyFlags(); len(errs) > 0 {
				for _, err := range errs {
					fmt.Printf("%v %v\n", color.RedString("Error:"), err)
				}
				os.Exit(1)
			}
			if !a.silence {
				printConfig()
			}
		}
	}
	if !a.silence && !a.noVersion {
		printVersion()
	}
	// run application
	if a.runFunc != nil {
		if !a.silence {
			fmt.Printf("%v Log data will now stream in as it occurs:\n", progressMessage)
		}
		if err := a.runFunc(a.basename); err != nil {
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
			os.Exit(1)
		}
	}
}

func printVersion() {
	fmt.Printf("%v Version:\n", progressMessage)
	fmt.Printf("%s\n", version.Get())
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	fmt.Printf("%v WorkingDir: %s\n", progressMessage, wd)
}
