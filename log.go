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
			URLs:     []string{"rfile://$workdir/log/debug.log?cut=hour"},
		},
		{
			Name:     "Info",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.INFO,
			MaxLevel: iface.FATAL,
			URLs:     []string{"rfile://$workdir/log/info.log?cut=hour"},
		},
		{
			Name:     "Warn",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.WARN,
			MaxLevel: iface.FATAL,
			URLs:     []string{"rfile://$workdir/log/warn.log?cut=hour"},
		},
		{
			Name:     "Error",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.ERROR,
			MaxLevel: iface.FATAL,
			URLs:     []string{"rfile://$workdir/log/error.log?cut=hour"},
		},
		{
			Name:     "Fatal",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.PANIC,
			MaxLevel: iface.FATAL,
			URLs:     []string{"rfile://$workdir/log/fatal.log?cut=hour"},
		},
		{
			Name:     "Access",
			Encoding: "",
			Encoder:  zaplog.NewConsoleEncoderConfig(),
			MinLevel: iface.DEBUG,
			MaxLevel: iface.FATAL,
			URLs:     []string{"rfile://$workdir/log/access.log?cut=hour"},
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
