package main

import (
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"
)

var (
	logger *GoDNSLogger
)

func main() {

	initLogger()

	initHostFiles()

	server := &Server{
		host:     settings.Server.Host,
		port:     settings.Server.Port,
		rTimeout: 5 * time.Second,
		wTimeout: 5 * time.Second,
	}

	server.Run()

	logger.Info("godns %s start", settings.Version)

	if settings.Debug {
		go profileCPU()
		go profileMEM()
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

forever:
	for {
		select {
		case <-sig:
			logger.Info("signal received, stopping")
			break forever
		}
	}

}

func profileCPU() {
	f, err := os.Create("godns.cprof")
	if err != nil {
		logger.Error("%s", err)
		return
	}

	pprof.StartCPUProfile(f)
	time.AfterFunc(6*time.Minute, func() {
		pprof.StopCPUProfile()
		f.Close()

	})
}

func profileMEM() {
	f, err := os.Create("godns.mprof")
	if err != nil {
		logger.Error("%s", err)
		return
	}

	time.AfterFunc(5*time.Minute, func() {
		pprof.WriteHeapProfile(f)
		f.Close()
	})

}

func initLogger() {
	logger = NewLogger()

	if settings.Log.Stdout {
		logger.SetLogger("console", nil)
	}

	if settings.Log.File != "" {
		config := map[string]interface{}{"file": settings.Log.File}
		logger.SetLogger("file", config)
	}

	logger.SetLevel(settings.Log.LogLevel())
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func exits(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
func initHostFiles() {
	if err := os.Mkdir("./conf", 0775); err != nil {
		logger.Info("create conf dir error:", err.Error())
	}
	if !exits(settings.ResolvConfig.ResolvFile) {
		if _, err := os.Create(settings.ResolvConfig.ResolvFile); err != nil {
			logger.Error("create ResolvFile error", err.Error())
		}
	}
	if !exits(settings.Hosts.HostsFile) {
		if _, err := os.Create(settings.Hosts.HostsFile); err != nil {
			logger.Error("create HostsFile error", err.Error())
		}
	}
}
