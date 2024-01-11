# tapp

team application

## Example

```
package main

import (
	"context"
	"os"

	"github.com/ironzhang/tapp"
	"github.com/ironzhang/tlog"
)

type Config struct {
	Enviroment string
	GRPCAddr   string
	HTTPAddr   string
}

type Application struct {
	config Config
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
	return nil
}

func (p *Application) RunHTTPServer(ctx context.Context) error {
	tlog.Infow("run http server", "addr", p.config.HTTPAddr)
	select {
	case <-ctx.Done():
		break
	}
	tlog.Infow("stop http server", "addr", p.config.HTTPAddr)
	return nil
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
```
