package jago

import (
	"strings"
	"os"
	"path"
	"strconv"
	"github.com/urfave/cli"
	"github.com/pkg/profile"
	"github.com/fatih/color"
	"fmt"
)

func Go() {
	app := cli.NewApp()
	app.Name = "Jago"
	app.Usage = "A simplified Java Virtual Machine for the educational purpose"
	app.Description = ""
	app.UsageText = "jago [-options] class [args...]"
	app.Version = "1.0.0"
	app.Author = "Chao Yang"
	app.Email = "chaoyangnz@gmail.com"
	app.Description = "A Java Virtual Machine demonstrating the basic features of execution engine, class loading, type/value system, exception handling, native methods etc."

	var classpath string
	var noLogo = false
	var logThread string
	var logClassloader string
	var debug = false
	var trace = false
	var profiling = false
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "classpath, cp",
			Value:       "",
			Usage:       "application classpath separated by colon",
			Destination: &classpath,
		},
		cli.BoolFlag{
			Name:        "noLogo, nl",
			Usage:       "don't show logo",
			Destination: &noLogo,
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "debug mode",
			Destination: &debug,
		},
		cli.BoolFlag{
			Name:        "trace, t",
			Usage:       "trace mode",
			Destination: &trace,
		},
		cli.StringFlag{
			Name:        "log:thread",
			Value:       "",
			Usage:       "log level of instruction execution in a thread context, options: info, debug, trace",
			Destination: &logThread,
		},
		cli.StringFlag{
			Name:        "log:classloader",
			Value:       "",
			Usage:       "log level of class loading, options: info, debug, trace",
			Destination: &logClassloader,
		},
		cli.BoolFlag{
			Name:        "profile, p",
			Usage:       "profile jago",
			Destination: &profiling,
		},
	}

	app.Action = func(c *cli.Context) error {
		if !noLogo {
			color.HiMagenta("    _                   \r\n   (_) __ _  __ _  ___  \r\n   | |/ _` |/ _` |/ _ \\ \r\n   | | (_| | (_| | (_) |   version %s\r\n  _/ |\\__,_|\\__, |\\___/    \r\n |__/       |___/     \n", app.Version)
		}
		commandEcho := fmt.Sprintf("Command: %s", strings.Join(os.Args, " "))
		color.Cyan("\n\n%s\n", commandEcho)

		args := c.Args()
		if c.NArg() < 1 {
			return nil
		}

		jagoHome := os.Getenv("JAGO_HOME")
		if jagoHome == "" {
			color.Yellow("Please set JAGO_HOME environment variable first.")
			return nil
		} else if !path.IsAbs(jagoHome) {
			color.Yellow("JAGO_HOME needs to be absolute path.")
			return nil
		}

		javaClassName := args.Get(0)
		javaArgs := make([]string, c.NArg()-1)
		for i := 0; i < len(args)-1; i++ {
			javaArgs[i] = args[i+1]
		}

		if classpath != "" {
			color.Blue("Add a new classpath: %s\n", classpath)
			VM.SetSystemSetting("classpath.application", classpath)
		}

		switch logThread {
		case "error":
			VM.SetSystemSetting("log.level.thread", strconv.Itoa(ERROR))
		case "warn":
			VM.SetSystemSetting("log.level.thread", strconv.Itoa(WARN))
		case "info":
			VM.SetSystemSetting("log.level.thread", strconv.Itoa(INFO))
		case "debug":
			VM.SetSystemSetting("log.level.thread", strconv.Itoa(DEBUG))
		case "trace":
			VM.SetSystemSetting("log.level.thread", strconv.Itoa(TRACE))
		default:
			if logThread != "" {
				color.Yellow("Invalid log:thread option: %s\n", logThread)
			}
		}

		switch logClassloader {
		case "error":
			VM.SetSystemSetting("log.level.classloader", strconv.Itoa(ERROR))
		case "warn":
			VM.SetSystemSetting("log.level.classloader", strconv.Itoa(WARN))
		case "info":
			VM.SetSystemSetting("log.level.classloader", strconv.Itoa(INFO))
		case "debug":
			VM.SetSystemSetting("log.level.classloader", strconv.Itoa(DEBUG))
		case "trace":
			VM.SetSystemSetting("log.level.classloader", strconv.Itoa(TRACE))
		default:
			if logClassloader != "" {
				color.Yellow("Invalid log:classloader option: %s\n", logClassloader)
			}
		}

		if debug {
			VM.SetSystemSetting("log.level.classloader", strconv.Itoa(DEBUG))
			VM.SetSystemSetting("log.level.thread", strconv.Itoa(DEBUG))
			VM.SetSystemSetting("log.level.io", strconv.Itoa(DEBUG))
			VM.SetSystemSetting("log.level.misc", strconv.Itoa(DEBUG))
		}

		if trace {
			VM.SetSystemSetting("log.level.classloader", strconv.Itoa(TRACE))
			VM.SetSystemSetting("log.level.thread", strconv.Itoa(TRACE))
			VM.SetSystemSetting("log.level.io", strconv.Itoa(TRACE))
			VM.SetSystemSetting("log.level.misc", strconv.Itoa(TRACE))
		}

		if profiling {
			// CPU profiling
			defer profile.Start(profile.CPUProfile).Stop()
			// Memory profiling
			//defer profile.Start(profile.MemProfile).Stop()
		}

		color.Green(strings.Repeat("·", len(commandEcho)) + "\n")

		VM.Init()
		VM.Startup(javaName2BinaryName0(javaClassName), javaArgs...)
		return nil
	}

	app.Run(os.Args)
}