package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	BlockHeader = "### hosts-go ###"
	BlockFooter = "### hosts-go end ###"
)

var (
	version         = "dev"
	systemHostsPath = "/etc/hosts"
)

func main() {
	i, _ := debug.ReadBuildInfo()
	if i != nil && i.Main.Version != "(devel)" {
		version = i.Main.Version
	}
	rootCmd := newCmd()
	rootCmd.Version = version

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
}

func newCmd() *cobra.Command {
	var (
		testOnly      bool
		contentOnly   bool
		serviceAction string
		urls          []string
		duration      time.Duration
		reloadCommand string
	)

	cmd := &cobra.Command{
		Use:   "hosts-go",
		Short: "Fetch and merge hosts files from the internet",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(urls) == 0 {
				return errors.New("please provide at least one URL using --url flag")
			}

			switch serviceAction {
			case "install":
				return installService(urls)
			case "uninstall":
				return uninstallService(urls)
			case "":
			default:
				return errors.New("invalid service action")
			}

			if !testOnly {
				runLoop(urls, duration, reloadCommand)
				return nil
			}

			hostsBytes, err := FetchAndMergeHosts(urls)
			if err != nil {
				return err
			}

			if contentOnly {
				fmt.Println(string(hostsBytes))
				return nil
			}

			hostsContent, err := RenderHostsFile(systemHostsPath, hostsBytes)
			if err != nil {
				return err
			}

			fmt.Println(hostsContent)
			return nil

		},
	}

	cmd.Flags().StringSliceVarP(&urls, "url", "u", nil, "URLs to fetch hosts files from")
	cmd.Flags().BoolVarP(&testOnly, "test", "t", false, "Only output the rendered hosts content")
	cmd.Flags().BoolVar(&contentOnly, "content-only", false, "Output only the fetched hosts content")
	cmd.Flags().StringVarP(&serviceAction, "service", "s", "", "Install or uninstall service")
	cmd.Flags().DurationVarP(&duration, "duration", "d", 1*time.Hour, "Duration between each fetch")
	cmd.Flags().StringVar(&reloadCommand, "reload-command", "", "Command to execute after hosts file updated")
	return cmd
}

func FetchAndMergeHosts(urls []string) ([]byte, error) {
	var merged bytes.Buffer

	for _, url := range urls {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("http status %s", resp.Status)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to fetch hosts file from %s: %v", url, err)
		}

		_, err = io.Copy(&merged, resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read hosts file from %s: %v", url, err)
		}
	}

	return merged.Bytes(), nil
}

func RenderHostsFile(hostsFile string, content []byte) (string, error) {
	file, err := os.OpenFile(hostsFile, os.O_RDONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open hosts file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var result, after strings.Builder
	var step int
	for scanner.Scan() {
		line := scanner.Text()
		if line == BlockHeader && step == 0 {
			step = 1
		} else if line == BlockFooter && step == 1 {
			step = 2
			continue
		}
		if step == 0 {
			result.WriteString(line + "\n")
		} else if step == 2 {
			after.WriteString(line + "\n")
		}
	}

	if scanner.Err() != nil {
		return "", fmt.Errorf("failed to read hosts file: %v", scanner.Err())
	}
	result.WriteString(BlockHeader + "\n")
	result.Write(content)
	result.WriteString("\n\n" + "# hosts-go updated at " + time.Now().Format(time.RFC3339) + "\n")
	result.WriteString(BlockFooter + "\n")
	result.WriteString(after.String())
	return result.String(), nil
}

func WriteHostsFile(hostsFile, content string) error {
	if err := os.WriteFile(hostsFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write hosts file: %v", err)
	}
	return nil
}

func runLoop(urls []string, duration time.Duration, reloadCommand string) {
	tick := time.Tick(duration)
	for ; true; <-tick {
		logrus.Info("start service")
		if err := update(urls); err != nil {
			logrus.Error(err)
			continue
		}
		logrus.Info("update hosts file successfully")
		if reloadCommand != "" {
			logrus.Info("run reload command")
			cmd := exec.Command("sh", "-c", reloadCommand)
			if err := cmd.Run(); err != nil {
				logrus.Error(err)
			}
		}
	}
}

func update(urls []string) error {
	hostsBytes, err := FetchAndMergeHosts(urls)
	if err != nil {
		return err
	}
	hostsContent, err := RenderHostsFile(systemHostsPath, hostsBytes)
	if err != nil {
		return err
	}
	return WriteHostsFile(systemHostsPath, hostsContent)
}
