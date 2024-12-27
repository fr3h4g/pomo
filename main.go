package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

type Settings struct {
	DefaultWorkTime int
}

var settings = Settings{DefaultWorkTime: 1}

var rootCmd = &cobra.Command{
	Use:   "pomo",
	Short: "Pomo is a simple pomodoro timer",
	Long:  `Pomo is a simple pomodoro timer.`,
	Run: func(cmd *cobra.Command, args []string) {
		status()
	},
}

func loadSettings() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	pomodir := path.Join(userHomeDir, ".pomo")
	err = os.MkdirAll(pomodir, 0755)
	if err != nil {
		panic(err)
	}
	os.MkdirAll(pomodir, 0755)
	settingsFile := path.Join(pomodir, "settings.json")
	if _, err := os.Stat(settingsFile); os.IsNotExist(err) {
		os.WriteFile(settingsFile, []byte(`{"DefaultWorkTime": 25}`), 0644)
	}
	file, err := os.Open(settingsFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&settings)
	if err != nil {
		panic(err)
	}
}

func setStartTime(args []string) {
	workTimeVal := settings.DefaultWorkTime
	if len(args) > 0 {
		workTime, err := strconv.Atoi(args[0])
		if err != nil {
			panic(err)
		}
		workTimeVal = workTime
	}
	unixtime := time.Now().Unix() + int64(workTimeVal*60)
	saveStartTime(unixtime)
}

func saveStartTime(unixtime int64) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	pomodir := path.Join(userHomeDir, ".pomo")
	err = os.MkdirAll(pomodir, 0755)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(path.Join(pomodir, "start_time"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("%d", unixtime))
	if err != nil {
		panic(err)
	}

}

func getTimeRemaining() int64 {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	pomodir := path.Join(userHomeDir, ".pomo")
	startTimePath := path.Join(pomodir, "start_time")
	startTime, err := os.ReadFile(startTimePath)
	if err != nil {
		panic(err)
	}
	i, err := strconv.ParseInt(string(startTime), 10, 64)
	if err != nil {
		panic(err)
	}
	return i - time.Now().Unix()
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(statusCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the pomodoro timer",
	RunE: func(cmd *cobra.Command, args []string) error {
		setStartTime(args)
		return nil
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the pomodoro timer",
	RunE: func(cmd *cobra.Command, args []string) error {
		unixtime := time.Now().Unix() - 1
		saveStartTime(unixtime)
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the pomodoro timer",
	RunE: func(cmd *cobra.Command, args []string) error {
		return status()
	},
}

func status() error {
	reminingTime := getTimeRemaining()

	if reminingTime <= 0 {
		return nil
	}

	seconds := reminingTime % 60
	minutes := (reminingTime - seconds) / 60
	fmt.Printf("🍅%d:%02d\n", minutes, seconds)

	return nil
}

func main() {
	loadSettings()
	rootCmd.Execute()
}
