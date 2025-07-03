package options

import (
	httpservertype "frisboo-bank/pkg/http/http_server/options/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"time"
)

type HttpServerOption = func(options *HttpServerOptions)

func UseServerType(t httpservertype.HttpServerType) HttpServerOption {
	return func(o *HttpServerOptions) { o.Type = t }
}

func UseBasePath(base string) HttpServerOption {
	return func(o *HttpServerOptions) { o.BasePath = base }
}

func UseDevelopment(dev bool) HttpServerOption {
	return func(o *HttpServerOptions) { o.Development = dev }
}

func UseHost(host string) HttpServerOption {
	return func(o *HttpServerOptions) { o.Host = host }
}

func UsePort(port string) HttpServerOption {
	return func(o *HttpServerOptions) { o.Port = port }
}

func UseIgnoreLogUrls(urls []string) HttpServerOption {
	return func(o *HttpServerOptions) { o.IgnoreLogUrls = urls }
}

func UseBodyLimit(limit string) HttpServerOption {
	return func(o *HttpServerOptions) { o.BodyLimit = limit }
}

func UseIdleTimeout(timeout time.Duration) HttpServerOption {
	return func(o *HttpServerOptions) { o.IdleTimeout = timeout }
}

func UseMaxHeaderBytes(max int) HttpServerOption {
	return func(o *HttpServerOptions) { o.MaxHeaderBytes = max }
}

func UseReadHeaderTimeout(timeout time.Duration) HttpServerOption {
	return func(o *HttpServerOptions) { o.ReadHeaderTimeout = timeout }
}

func UseReadTimeout(timeout time.Duration) HttpServerOption {
	return func(o *HttpServerOptions) { o.ReadTimeout = timeout }
}

func UseServerShutdownTimeout(timeout time.Duration) HttpServerOption {
	return func(o *HttpServerOptions) { o.ServerShutdownTimeout = timeout }
}

func UseWriteTimeout(timeout time.Duration) HttpServerOption {
	return func(o *HttpServerOptions) { o.WriteTimeout = timeout }
}

func UseLogger(logger loggerContracts.Logger) HttpServerOption {
	return func(o *HttpServerOptions) { o.Logger = logger }
}
