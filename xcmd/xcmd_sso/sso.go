package main

import (
	"github.com/pubgo/g/xservice/sso/api"
	"github.com/pubgo/g/xservice/sso/app"
	"github.com/pubgo/g/xservice/sso/model"
	"github.com/pubgo/g/xservice/sso/utils"
	"github.com/pubgo/g/xservice/sso/web"
	"os"

	l4g "github.com/alecthomas/log4go"
	_ "github.com/pubgo/g/xservice/sso/model/wechat"
	"github.com/spf13/cobra"
	"os/signal"
	"syscall"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "single-sign-on",
	Short: "sso",
	Long:  `single-sign-on`,
	RunE:  runServerCmd,
}
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "运行服务器",
	RunE:  runServerCmd,
}

func runServerCmd(cmd *cobra.Command, args []string) error {
	config, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	utils.CfgDisableConfigWatch, _ = cmd.Flags().GetBool("disableconfigwatch")

	runServer(config)
	return nil
}
func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "config.json", "Configuration file to use.")
	rootCmd.PersistentFlags().Bool("disableconfigwatch", false, "When set config.json will not be loaded from disk when the file is changed.")
	rootCmd.AddCommand(serverCmd)
}
func runServer(configFileLocation string) {
	if err := utils.InitAndLoadConfig(configFileLocation); err != nil {
		l4g.Exit("Unable to load single-sign-on configuration file: ", err)
		return
	}

	if err := utils.InitTranslations(utils.Cfg.LocalizationSettings); err != nil {
		l4g.Exit("Unable to load single-sign-on translation files: %v", err)
		return
	}

	pwd, _ := os.Getwd()
	l4g.Info(utils.T("single-sign-on.current_version"), model.CurrentVersion, model.BuildNumber, model.BuildDate, model.BuildHash)
	l4g.Info(utils.T("single-sign-on.working_dir"), pwd)
	l4g.Info(utils.T("single-sign-on.config_file"), utils.FindConfigFile(configFileLocation))

	if model.BuildNumber == "dev" {
		*utils.Cfg.ServiceSettings.EnableDeveloper = true
	}

	app.NewServer()
	app.InitStores()
	api.InitRouter()
	api.InitApi()
	web.InitWeb()

	app.ReloadConfig()

	app.StartServer()

	go runTokenCleanupJob()

	utils.RegenerateClientConfig()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	app.StopServer()
}

func runTokenCleanupJob() {
	doTokenCleanup()
	model.CreateRecurringTask("Token Cleanup", doTokenCleanup, time.Hour*24)
}

func doTokenCleanup() {
	app.Srv.SqlStore.Token().Cleanup()
}
