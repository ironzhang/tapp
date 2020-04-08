package tapp

import (
	"context"
	"os"
	"time"

	"git.xiaojukeji.com/pearls/tlog"
	"git.xiaojukeji.com/pearls/tlog/zaplog"
)

type TConfig struct {
	Enviroment string
	GRPCAddr   string
	HTTPAddr   string
}

type TApplication struct {
	config TConfig
}

func (p *TApplication) Init() error {
	return nil
}

func (p *TApplication) Fini() error {
	return nil
}

func (p *TApplication) RunGRPCServer(ctx context.Context) error {
	tlog.Infow("run grpc server", "addr", p.config.GRPCAddr)
	select {
	case <-ctx.Done():
		break
	case <-time.After(2 * time.Second):
		break
	}
	tlog.Infow("stop grpc server", "addr", p.config.GRPCAddr)
	return nil
}

func (p *TApplication) RunHTTPServer(ctx context.Context) error {
	tlog.Infow("run http server", "addr", p.config.HTTPAddr)
	select {
	case <-ctx.Done():
		break
	case <-time.After(2 * time.Second):
		break
	}
	tlog.Infow("stop http server", "addr", p.config.HTTPAddr)
	return nil
}

func Example_tapp() {
	now = func() time.Time {
		return time.Time{}
	}

	app := &TApplication{
		config: TConfig{
			Enviroment: "development",
			GRPCAddr:   ":7000",
			HTTPAddr:   ":8000",
		},
	}
	logcfg := zaplog.NewDevelopmentConfig()
	f := Framework{
		Version: &VersionInfo{
			Version:   "0.0.1",
			GitCommit: "git commit",
			BuildTime: "build time",
		},
		Application: app,
		Config:      app.config,
		Runners:     []RunFunc{app.RunGRPCServer, app.RunHTTPServer},
		LogConfig:   &logcfg,
	}
	f.Main(os.Args)

	// output:
	// 0001-01-01 00:00:00 +0000 UTC start, version=&{0.0.1 git commit build time}, config={development :7000 :8000}
}
