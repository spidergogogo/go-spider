package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/sevlyar/go-daemon"
)

var (
	cmdFlag = flag.String("s", "", "start - start process\nstop - fast shutdown")
)

func main() {
	flag.Parse()
	if *cmdFlag != "start" && *cmdFlag != "stop" {
		fmt.Println("# Command params: -s start/stop")
		flag.Usage()
		return
	}
	initLog()
	cfgDaemon, err := daemonConfig()
	if err != nil {
		log.Error(fmt.Sprintf("获取daemon配置失败:%s", err))
		return
	}
	daemon.AddCommand(daemon.StringFlag(cmdFlag, "stop"), syscall.SIGTERM, termHandler)
	fmt.Println("logFile", cfgDaemon.logFile)
	cntxt := &daemon.Context{
		PidFileName: cfgDaemon.PidFile,
		PidFilePerm: 0644,
		LogFileName: cfgDaemon.logFile,
		LogFilePerm: 0644,
		WorkDir:     cfgDaemon.workDir,
		Umask:       027,
	}
	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			fmt.Println("Unable send signal to the daemon: ", err)
			return
		}
		err = daemon.SendCommands(d)
		if err != nil {
			fmt.Println("Send Command Error", err)
		}
		return
	}
	d, err := cntxt.Reborn()
	if err != nil {
		fmt.Println("Reborn Error: ", err)
		return
	}
	if d != nil {
		return
	}
	defer func(cntxt *daemon.Context) {
		err := cntxt.Release()
		if err != nil {
			log.Error(fmt.Sprintf("Context Realse Error: %s", err))
		}
	}(cntxt)
	log.Info("==================")
	log.Info("Daemon Started")

	go terminateHelper()
	go process()

	err = daemon.ServeSignals()
	if err != nil {
		log.Error(fmt.Sprintf("Daemon Terminat Error: %s", err))
	}
	log.Info("Daemon Terminated.")
}

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func terminateHelper() {
	func() {
		for {
			time.Sleep(time.Second)
			select {
			case <-stop:
				return
			default:
			}
		}
	}()
	done <- struct{}{}
}

func termHandler(sig os.Signal) error {
	// log.Info("terminating...")
	log.Info("terminating...")
	stop <- struct{}{}
	<-done
	return daemon.ErrStop
}
