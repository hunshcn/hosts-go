package main

import (
	"fmt"

	"github.com/kardianos/service"
	"github.com/sirupsen/logrus"
)

// copy from https://github.com/jeessy2/ddns-go/blob/1576b7bfd8272bb1dda1352919d367bab2f4ce94/main.go#L327
const sysvScript = `#!/bin/sh /etc/rc.common
DESCRIPTION="{{.Description}}"
cmd="{{.Path}}{{range .Arguments}} {{.|cmd}}{{end}}"
name="ddns-go"
pid_file="/var/run/$name.pid"
stdout_log="/var/log/$name.log"
stderr_log="/var/log/$name.err"
START=99
get_pid() {
    cat "$pid_file"
}
is_running() {
    [ -f "$pid_file" ] && cat /proc/$(get_pid)/stat > /dev/null 2>&1
}
start() {
	if is_running; then
		echo "Already started"
	else
		echo "Starting $name"
		{{if .WorkingDirectory}}cd '{{.WorkingDirectory}}'{{end}}
		$cmd >> "$stdout_log" 2>> "$stderr_log" &
		echo $! > "$pid_file"
		if ! is_running; then
			echo "Unable to start, see $stdout_log and $stderr_log"
			exit 1
		fi
	fi
}
stop() {
	if is_running; then
		echo -n "Stopping $name.."
		kill $(get_pid)
		for i in $(seq 1 10)
		do
			if ! is_running; then
				break
			fi
			echo -n "."
			sleep 1
		done
		echo
		if is_running; then
			echo "Not stopped; may still be shutting down or shutdown may have failed"
			exit 1
		else
			echo "Stopped"
			if [ -f "$pid_file" ]; then
				rm "$pid_file"
			fi
		fi
	else
		echo "Not running"
	fi
}
restart() {
	stop
	if is_running; then
		echo "Unable to stop, will not attempt to start"
		exit 1
	fi
	start
}
`

func getService(urls []string) service.Service {
	var depends []string
	options := make(service.KeyValue)
	switch service.ChosenSystem().String() {
	case "unix-systemv":
		options["SysvScript"] = sysvScript
	case "windows-service":
		options["DelayedAutoStart"] = true
	default:
		depends = append(depends, "Requires=network.target",
			"After=network-online.target")
	}
	svcConfig := &service.Config{
		Name:         "hosts-go",
		DisplayName:  "hosts-go",
		Description:  "Fetch and merge hosts files from the internet",
		Dependencies: depends,
		Option:       options,
	}
	for _, u := range urls {
		svcConfig.Arguments = append(svcConfig.Arguments, "-u", u)
	}

	s, err := service.New(nil, svcConfig)
	if err != nil {
		logrus.Fatal(err)
	}
	return s
}

func installService(urls []string) error {
	s := getService(urls)
	err := s.Install()
	if err != nil {
		return fmt.Errorf("failed to install service: %v", err)

	}
	fmt.Println("Service installed successfully!")
	return nil
}

func uninstallService(urls []string) error {
	s := getService(urls)
	err := s.Uninstall()
	if err != nil {
		return fmt.Errorf("failed to install service: %v", err)
	}

	logrus.Info("Service uninstalled successfully!")
	return nil
}
