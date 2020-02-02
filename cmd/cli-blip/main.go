package main

import (
	"fmt"
	"git.ronaksoftware.com/blip/server/internal/tools"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io/ioutil"
	"strings"
)

func main() {
	RootCmd.AddCommand(ExitCmd)
	p := prompt.New(executor, completer)
	p.Run()
}

func executor(s string) {
	RootCmd.SetArgs(strings.Fields(s))
	_ = RootCmd.Execute()
}

func completer(d prompt.Document) []prompt.Suggest {
	suggests := make([]prompt.Suggest, 0, 10)
	cols := d.TextBeforeCursor()
	currCmd := RootCmd
	for _, col := range strings.Fields(cols) {
		for _, cmd := range currCmd.Commands() {
			if cmd.Name() == col {
				currCmd = cmd
				break
			}
		}
	}

	currWord := d.GetWordBeforeCursor()
	if strings.HasPrefix(currWord, "--") {
		// Search in Flags
		RootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
			if strings.HasPrefix(flag.Name, currWord[2:]) {
				suggests = append(suggests, prompt.Suggest{
					Text:        fmt.Sprintf("--%s", flag.Name),
					Description: flag.Usage,
				})
			}
		})
		currCmd.Flags().VisitAll(func(flag *pflag.Flag) {
			if strings.HasPrefix(flag.Name, currWord[2:]) {
				suggests = append(suggests, prompt.Suggest{
					Text:        fmt.Sprintf("--%s", flag.Name),
					Description: flag.Usage,
				})
			}
		})

	} else {
		for _, cmd := range currCmd.Commands() {
			if strings.HasPrefix(cmd.Name(), currWord) {
				suggests = append(suggests, prompt.Suggest{
					Text:        cmd.Name(),
					Description: cmd.Short,
				})
			}
		}
	}

	return suggests
}

var RootCmd = &cobra.Command{
	Use: "Root",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flag(FlagServerUrl).Changed {
			_ = ioutil.WriteFile(".blip-url", tools.StrToByte(cmd.Flag(FlagServerUrl).Value.String()), os.ModePerm)
		} else {
			baseUrlBytes, err := ioutil.ReadFile(".blip-url")
			if err == nil {
				baseUrl = tools.ByteToStr(baseUrlBytes)
			}
		}
		accessTokenBytes, err := ioutil.ReadFile(".blip-accessToken")
		if err == nil {
			accessToken = tools.ByteToStr(accessTokenBytes)
		}
		sessionIDBytes, err := ioutil.ReadFile(".blip-session")
		if err == nil {
			sessionID = tools.ByteToStr(sessionIDBytes)
		}
	},
}

var ExitCmd = &cobra.Command{
	Use: "exit",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func init() {
	fs := RootCmd.PersistentFlags()
	fs.String(FlagServerUrl, "https://v2.blipapi.xyz", "")
	fs.String(FlagPhone, "", "")
	fs.String(FlagPhoneCode, "", "")
	fs.String(FlagPhoneCodeHash, "", "")
	fs.String(FlagOtpID, "", "")
}
