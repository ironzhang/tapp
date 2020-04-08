package tapp

import (
	"git.xiaojukeji.com/pearls/tlog/iface"
	"git.xiaojukeji.com/pearls/tlog/zaplog"
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
			URLs:     []string{"rfile://$workdir/log/debug.log?suffix=hour&period=1h"},
		},
		{
			Name:     "Info",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.INFO,
			MaxLevel: iface.FATAL,
			URLs:     []string{"rfile://$workdir/log/info.log?suffix=hour&period=1h"},
		},
		{
			Name:     "Warn",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.WARN,
			MaxLevel: iface.FATAL,
			URLs:     []string{"rfile://$workdir/log/warn.log?suffix=hour&period=1h"},
		},
		{
			Name:     "Error",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.ERROR,
			MaxLevel: iface.FATAL,
			URLs:     []string{"rfile://$workdir/log/error.log?suffix=hour&period=1h"},
		},
		{
			Name:     "Fatal",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.PANIC,
			MaxLevel: iface.FATAL,
			URLs:     []string{"rfile://$workdir/log/fatal.log?suffix=hour&period=1h"},
		},
		{
			Name:     "Access",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.DEBUG,
			MaxLevel: iface.FATAL,
			URLs:     []string{"rfile://$workdir/log/access.log?suffix=hour&period=1h"},
		},
	},
	Loggers: []zaplog.LoggerConfig{
		{
			Name:            "",
			DisableCaller:   false,
			StacktraceLevel: zaplog.PanicStacktrace,
			Cores:           []string{"Debug", "Info", "Warn", "Error", "Fatal"},
		},
		{
			Name:            "access",
			DisableCaller:   false,
			StacktraceLevel: zaplog.PanicStacktrace,
			Cores:           []string{"Access"},
		},
	},
}
