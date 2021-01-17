package cli

import (
	"bufio"
	"fmt"
	"klog/app"
	"os"
	"time"
)

var Start = Command{
	Name:        "start",
	Description: "Create a new entry",
	Main:        start,
}

func start(ctx app.Context, args []string) int {
	start := time.Now()
	//date, _ := record.NewDateFromTime(start)

	{
		//t, _ := record.NewTime(start.Hour(), start.Minute())
		//ctx.QuickStartAt(date, t)
	}

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

	{
		//end := time.Now()
		//t, _ := record.CreateTimeFromTime(end)
		//ctx.QuickStopAt(date, t)
	}
	return OK
}
