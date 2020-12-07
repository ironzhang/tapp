package tapp

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"git.xiaojukeji.com/pearls/tlog"
	"git.xiaojukeji.com/pearls/tlog/zaplog"
)

var exit = os.Exit
var now = time.Now

// Flagger 命令行参数设置接口
//
// Application 如果实现了该接口，可设置自定义的命令行参数
type Flagger interface {
	SetFlags(fs *flag.FlagSet)
}

// Command 命令执行接口
//
// Application 如果实现了该接口，可执行自定义命令
type Command interface {
	DoCommand() (quit bool)
}

// Application 应用接口
type Application interface {
	Init() error
	Fini() error
}

// RunFunc 运行函数定义
type RunFunc func(ctx context.Context) error

// VersionInfo 版本信息
type VersionInfo struct {
	Version   string
	GitCommit string
	BuildTime string
}

// Options 命令行选项
type Options struct {
	Version          bool   // 输出版本信息
	ConfigFile       string // 应用配置文件
	ConfigExample    string // 生成应用配置文件示例
	LogConfigFile    string // 日志配置文件
	LogConfigExample string // 生成日志配置文件示例
}

// Framework 应用框架
type Framework struct {
	// 版本信息，如果为 nil，则不输出版本信息
	Version *VersionInfo

	// 命令行选项
	Options Options

	// 应用程序，不能为 nil，必须设置
	Application Application

	// 应用配置，可以为 nil
	Config interface{}

	// 运行函数列表，可以为空
	Runners []RunFunc

	// 日志上下文钩子，可以为 nil
	LoggerContextHook zaplog.ContextHook

	// 日志配置，如果为 nil，则使用 DefaultLogConfig
	LogConfig *zaplog.Config

	// 配置编解码器，如果为 nil，则使用 tapp.TOMLC
	Codec Codec

	// 命令行解析器，如果为 nil，则使用 flag.CommandLine
	CommandLine *flag.FlagSet
}

func (p *Framework) init() error {
	if p.Application == nil {
		return errors.New("Framework.Application is nil")
	}
	if p.LogConfig == nil {
		p.LogConfig = &DefaultLogConfig
	}
	if p.Codec == nil {
		p.Codec = TOMLC
	}
	if p.CommandLine == nil {
		p.CommandLine = flag.CommandLine
	}
	return nil
}

func (p *Framework) setupCommandLine() {
	// 版本命令行参数
	if p.Version != nil {
		p.CommandLine.BoolVar(&p.Options.Version, "version", p.Options.Version, "print version info")
	}

	// 应用配置命令行参数
	if p.Config != nil {
		p.CommandLine.StringVar(&p.Options.ConfigFile, "config", p.Options.ConfigFile,
			"app config file")
		p.CommandLine.StringVar(&p.Options.ConfigExample, "config-example", p.Options.ConfigExample,
			"app config example")
	}

	// 日志配置命令行参数
	p.CommandLine.StringVar(&p.Options.LogConfigFile, "log-config", p.Options.LogConfigFile,
		"log config file")
	p.CommandLine.StringVar(&p.Options.LogConfigExample, "log-config-example", p.Options.LogConfigExample,
		"log config example")

	// 应用自定义命令行参数
	if f, ok := p.Application.(Flagger); ok {
		f.SetFlags(p.CommandLine)
	}
}

func (p *Framework) parseCommandLine(args []string) error {
	p.setupCommandLine()
	return p.CommandLine.Parse(args[1:])
}

func (p *Framework) printVersion() error {
	if p.Version == nil {
		return errors.New("Framework.Version is nil")
	}
	info := p.Version
	fmt.Fprintf(os.Stdout, "version: %s\n", info.Version)
	fmt.Fprintf(os.Stdout, "git commit: %s\n", info.GitCommit)
	fmt.Fprintf(os.Stdout, "go version: %s\n", runtime.Version())
	fmt.Fprintf(os.Stdout, "go OS: %s\n", runtime.GOOS)
	fmt.Fprintf(os.Stdout, "go Arch: %s\n", runtime.GOARCH)
	fmt.Fprintf(os.Stdout, "build time: %s\n", info.BuildTime)
	return nil
}

func (p *Framework) generateAppConfig() error {
	if p.Config == nil {
		return errors.New("Framework.Config is nil")
	}
	if err := writeToFile(p.Codec, p.Options.ConfigExample, p.Config); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "generate app config file(%s) successful\n", p.Options.ConfigExample)
	return nil
}

func (p *Framework) generateLogConfig() error {
	if p.LogConfig == nil {
		return errors.New("Framework.LogConfig is nil")
	}
	if err := writeToFile(p.Codec, p.Options.LogConfigExample, p.LogConfig); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "generate log config file(%s) successful\n", p.Options.LogConfigExample)
	return nil
}

func (p *Framework) doCommand() (err error) {
	quit := false
	if p.Options.Version {
		if err = p.printVersion(); err != nil {
			return fmt.Errorf("print version: %w", err)
		}
		quit = true
	} else if p.Options.ConfigExample != "" {
		if err = p.generateAppConfig(); err != nil {
			return fmt.Errorf("generate app config: %w", err)
		}
		quit = true
	} else if p.Options.LogConfigExample != "" {
		if err = p.generateLogConfig(); err != nil {
			return fmt.Errorf("generate log config: %w", err)
		}
		quit = true
	}
	if c, ok := p.Application.(Command); ok {
		if c.DoCommand() {
			quit = true
		}
	}
	if quit {
		exit(0)
	}
	return nil
}

func (p *Framework) openLogger() (*zaplog.Logger, error) {
	if p.Options.LogConfigFile != "" {
		if err := loadFromFile(p.Codec, p.Options.LogConfigFile, p.LogConfig); err != nil {
			return nil, err
		}
	}
	if p.LogConfig == nil {
		return nil, errors.New("Framework.LogConfig is nil")
	}
	return zaplog.New(*p.LogConfig, zaplog.SetContextHook(p.LoggerContextHook))
}

func (p *Framework) loadAppConfig() error {
	if p.Config == nil || p.Options.ConfigFile == "" {
		return nil
	}
	return loadFromFile(p.Codec, p.Options.ConfigFile, p.Config)
}

func (p *Framework) printStartInfo() {
	args := make([]interface{}, 0, 4)
	if p.Version != nil {
		args = append(args, "version", p.Version)
	}
	if p.Config != nil {
		args = append(args, "config", p.Config)
	}
	tlog.Infow("start", args...)
	fmt.Fprintf(os.Stdout, "[%v] start, version=%v, config=%v\n", now(), p.Version, p.Config)
}

func (p *Framework) run() (err error) {
	// init app
	if err = p.Application.Init(); err != nil {
		tlog.Errorw("init app", "error", err)
		return fmt.Errorf("init app: %w", err)
	}
	tlog.Debug("init app successful")

	// quit signal
	quit := make(chan error, len(p.Runners))
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// 由于无法捕捉 SIGKILL 信号，因此正常结束进程使用
		// kill pid
		// 轻易不要使用
		// kill -9 pid
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		sig := <-ch
		tlog.Infow("recv signal", "sig", sig)
		cancel()
		time.Sleep(10 * time.Second)
		quit <- errors.New("wait 10s, force exit")
	}()

	// run app
	var wg sync.WaitGroup
	for i, r := range p.Runners {
		wg.Add(1)
		go func(n int, f RunFunc) {
			if err := f(ctx); err != nil {
				tlog.Errorf("run %dth runner(%v): %v", n, f, err)
				quit <- fmt.Errorf("run %dth runner(%v): %w", n, f, err)
			}
			wg.Done()
		}(i, r)
	}
	go func() {
		wg.Wait()
		quit <- nil
	}()

	// wait quit
	qerr := <-quit
	if qerr != nil {
		tlog.Errorw("run app", "error", qerr)
	}

	// fini app
	if err = p.Application.Fini(); err != nil {
		tlog.Errorw("fini app", "error", err)
		if qerr != nil {
			return qerr
		}
		return fmt.Errorf("fini app: %w", err)
	}
	tlog.Debug("fini app successful")

	// return error
	return qerr
}

func (p *Framework) Main(args []string) {
	var err error

	// 初始化框架
	if err = p.init(); err != nil {
		fmt.Fprintf(os.Stderr, "[%v] init: %v\n", now(), err)
		exit(1)
	}

	// 解析命令行参数
	if err = p.parseCommandLine(args); err != nil {
		fmt.Fprintf(os.Stderr, "[%v] parse command line: %v\n", now(), err)
		exit(2)
	}

	// 执行命令
	if err = p.doCommand(); err != nil {
		fmt.Fprintf(os.Stderr, "[%v] do command: %v\n", now(), err)
		exit(3)
	}

	// 加载日志对象
	logger, err := p.openLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[%v] open logger: %v\n", now(), err)
		exit(4)
	}
	defer logger.Close()
	tlog.SetLogger(logger)

	// 关闭程序
	shutdown := func(code int) {
		logger.Close()                     // 关闭日志
		time.Sleep(200 * time.Millisecond) // 有些模块如监控上报是异步处理的，所以等待 200ms 再退出
		exit(code)
	}

	// 加载应用配置
	err = p.loadAppConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[%v] load app config: %v\n", now(), err)
		shutdown(5)
	}

	// 往日志中输出程序启动信息
	p.printStartInfo()

	// 运行程序
	if err = p.run(); err != nil {
		fmt.Fprintf(os.Stderr, "[%v] run: %v", now(), err)
		shutdown(6)
	}

	// 正常退出
	shutdown(0)
}
