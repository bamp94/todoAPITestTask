package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"cyberzilla_api_task/application"
	"cyberzilla_api_task/config"
	"cyberzilla_api_task/controller"
	"cyberzilla_api_task/docs"
	"cyberzilla_api_task/model"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// @title CYBERZILLA API task
// @description Документация для http сервера приложения
func main() {
	appCli := cli.NewApp()
	appCli.Version = controller.Branch
	if appCli.Version == "" {
		appCli.Version = "version not specified"
	}
	docs.SwaggerInfo.Version = appCli.Version

	appCli.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "/etc/cyberzilla_api_task/config.json",
			Usage: "optional config path",
		},
	}

	appCli.Action = func(cliContext *cli.Context) {
		cli.ShowVersion(cliContext)
		config := config.New(cliContext.String("config"))

		m := model.NewFromConfig(config.DB)
		if err := m.CheckMigrations(); err != nil {
			logrus.WithError(err).Fatal("invalid database condition")
		}

		shutdown := make(chan int, 16)
		wg := sync.WaitGroup{}
		appNew := application.New(m, config)
		c := controller.New(config, appNew)
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			for {
				select {
				case <-shutdown:
					cancel()
					return
				}
			}
		}()

		wg.Add(1)
		go func() {
			c.ServeHTTP(ctx, config.Port)
			wg.Done()
		}()

		gracefulClosing(&wg, shutdown)
	}

	appCli.Commands = []cli.Command{
		{
			Name:  "migrate",
			Usage: "update migrations to the latest stage",
			Action: func(cliContext *cli.Context) {
				config := config.New(cliContext.GlobalString("config"))
				m := model.NewFromConfig(config.DB)
				m.Migrate()
			},
		},
	}

	if err := appCli.Run(os.Args); err != nil {
		panic(err)
	}
}

func gracefulClosing(servicesWg *sync.WaitGroup, shutdown chan int) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	logrus.Info("stopping services... (Double enter Ctrl + C to force close)")

	close(shutdown)

	quit := make(chan struct{})
	go func() {
		<-sig
		<-sig
		logrus.Info("services unsafe stopped")
		<-quit
	}()

	go func() {
		servicesWg.Wait()
		logrus.Info("services gracefully stopped")
		<-quit
	}()

	quit <- struct{}{}
	close(quit)
}
