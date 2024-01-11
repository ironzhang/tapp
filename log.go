package tapp

import (
	"github.com/ironzhang/tlog/iface"
	"github.com/ironzhang/tlog/zaplog"
)

var DefaultLogConfig = zaplog.Config{
	Level: iface.INFO,
	Cores: []zaplog.CoreConfig{
		{
			Name:     "Debug",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.DEBUG,
			MaxLevel: iface.DEBUG,
			URLs:     []string{"rfile://$workdir/log/debug.log?cut=day"},
		},
		{
			Name:     "Info",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.INFO,
			MaxLevel: iface.FATAL,
			URLs:     []string{"rfile://$workdir/log/info.log?cut=day"},
		},
		{
			Name:     "Warn",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.WARN,
			MaxLevel: iface.FATAL,
			URLs:     []string{"rfile://$workdir/log/warn.log?cut=day"},
		},
		{
			Name:     "Error",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.ERROR,
			MaxLevel: iface.FATAL,
			URLs:     []string{"rfile://$workdir/log/error.log?cut=day"},
		},
		{
			Name:     "Fatal",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.PANIC,
			MaxLevel: iface.FATAL,
			URLs:     []string{"rfile://$workdir/log/fatal.log?cut=day"},
		},
	},
	Loggers: []zaplog.LoggerConfig{
		{
			Name:            "",
			DisableCaller:   false,
			StacktraceLevel: zaplog.PanicStacktrace,
			Cores:           []string{"Debug", "Info", "Warn", "Error", "Fatal"},
		},
	},
}
