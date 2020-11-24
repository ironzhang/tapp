package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"git.xiaojukeji.com/pearls/tapp"
	"git.xiaojukeji.com/pearls/tlog"
)

type Config struct {
	Enviroment string
	GRPCAddr   string
	HTTPAddr   string
}

type Options struct {
	OnlyVersion bool
}

type Application struct {
	config  Config
	options Options
}

func (p *Application) SetFlags(fs *flag.FlagSet) {
	fs.BoolVar(&p.options.OnlyVersion, "only-version", false, "only print version")
}

func (p *Application) DoCommand() (exit bool) {
	if p.options.OnlyVersion {
		fmt.Printf("%s\n", Version)
		return true
	}
	return false
}

func (p *Application) Init() error {
	return nil
}

func (p *Application) Fini() error {
	return nil
}

func (p *Application) RunGRPCServer(ctx context.Context) error {
	tlog.Infow("run grpc server", "addr", p.config.GRPCAddr)
	select {
	case <-ctx.Done():
		break
	}
	tlog.Infow("stop grpc server", "addr", p.config.GRPCAddr)
	return errors.New("failed to stop grpc server")
}

func (p *Application) RunHTTPServer(ctx context.Context) error {
	tlog.Infow("run http server", "addr", p.config.HTTPAddr)
	select {
	case <-ctx.Done():
		break
	}
	tlog.Infow("stop http server", "addr", p.config.HTTPAddr)
	return errors.New("failed to stop http server")
}

var (
	Version   = "Unknown"
	GitCommit = "Unknown"
	BuildTime = "Unknown"
)

func main() {
	app := &Application{
		config: Config{
			Enviroment: "development",
			GRPCAddr:   ":7000",
			HTTPAddr:   ":8000",
		},
	}
	tapp.DefaultLogConfig.Level = tlog.DEBUG
	f := tapp.Framework{
		Version: &tapp.VersionInfo{
			Version:   Version,
			GitCommit: GitCommit,
			BuildTime: BuildTime,
		},
		Application: app,
		Config:      app.config,
		Runners:     []tapp.RunFunc{app.RunGRPCServer, app.RunHTTPServer},
	}
	f.Main(os.Args)
}
