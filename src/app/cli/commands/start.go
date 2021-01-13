package commands

import (
	"bufio"
	"fmt"
	"klog/app"
	"klog/app/cli"
	"klog/datetime"
	"klog/project"
	"os"
	"time"
)

var Start cli.Command

func init() {
	Start = cli.Command{
		Name:        "start",
		Alias:       []string{},
		Description: "Create a new entry",
		Main:        start,
	}
}

func start(env app.Environment, project project.Project, args []string) int {
	start := time.Now()
	wd, _ := app.Start(project, start)
	ticker := time.NewTicker(1 * time.Second)
	fmt.Print("\n")
	go func() {
		for {
			select {
			case <-ticker.C:
				diff := time.Now().Sub(start)
				out := time.Time{}.Add(diff)
				fmt.Printf("\033[F\b\b\b\b\b\b\b\b\b%02d:%02d:%02d\n", out.Hour(), out.Minute(), out.Second())
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Print(line)
		if line == "^Q" {
			break
		}
	}

	later := time.Now()
	laterTime, _ := datetime.CreateTimeFromTime(later)
	wd.EndOpenRange(laterTime)
	project.Save(wd)
	return cli.OK
}
