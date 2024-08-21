package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	ktx := kong.Parse(&Time{})
	err := ktx.Run()
	if err != nil {
		log.Fatalf("%++v", err)
	}
}

type Time struct {
	Topic string
	Args  []string `arg:""`
}

func (t Time) Run() error {
	fmt.Println(t.Args)
	var cmd = exec.Command(t.Args[0], t.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	var failed bool

	err := cmd.Run()
	if err != nil {
		failed = true
		fmt.Printf("%++v\n", err)
	}

	duration := time.Since(start)

	defer func() {
		usage, hasUsage, err := GetChildrenUsage()
		if err != nil {
			failed = false
			fmt.Printf("%++v\n", err)

		}
		fmt.Println("duration", duration)
		if hasUsage {
			fmt.Println("user", usage.UserTime)
			fmt.Println("system", usage.SystemTime)
			fmt.Println("max-rss", usage.MaxRSS)
		}
		if t.Topic != "" {
			msg := ""
			if failed {
				msg = "FAILED\n"
			}

			msg += fmt.Sprintf("duration: %s\n", duration)
			msg += fmt.Sprintf("user: %s\n", usage.UserTime)
			msg += fmt.Sprintf("system: %s\n", usage.SystemTime)
			msg += fmt.Sprintf("mem: %d\n", usage.MaxRSS)

			req, err := http.NewRequest("POST", "https://ntfy.sh/"+t.Topic, strings.NewReader(msg))
			if err != nil {
				fmt.Println(err)
				return
			}
			req.Header.Set("Title", fmt.Sprintf("Command is finished: %s", strings.Join(t.Args, " ")))
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}
			if resp.StatusCode != http.StatusOK {
				panic(resp.Status)
			}
			_ = resp.Body.Close()
		}
	}()

	return nil
}
