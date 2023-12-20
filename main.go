package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/aws/aws-sdk-go/service/scheduler"
)

var (
	app = kingpin.New("otag", "One Time Auto Scaling Group")

	debug = app.Flag("debug", "Enable debug mode.").Bool()

	register   = app.Command("register", "Register one time auto scaling group")
	date       = register.Flag("date", "Specify the date and time in 2006-01-02 15:04 format (e.g. 2021-01-01 00:00)").Default(time.Now().Format("2006-01-02 15:04")).String()
	deregister = app.Command("deregister", "Deregister one time auto scaling group")
)

func main() {
	lv := new(slog.LevelVar)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lv}))
	slog.SetDefault(logger)
	slog.With(*debug)
	if *debug {
		lv.Set(slog.LevelDebug)
	}

	app.Version("0.0.1")
	// テキストファイルを読み取ります (例: input.txt)
	file, err := os.Open("input.txt")
	if err != nil {
		slog.Error("Error opening the file: %w", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isHeader := true
	for scanner.Scan() {
		line := scanner.Text()

		// ヘッダ行をスキップします
		if isHeader {
			isHeader = false
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 6 {
			slog.Warn("Invalid line: %s", line)
			continue
		}

		a := &autoScaling{}
		s := scheduler.CreateScheduleInput{}
		switch scalingType := parts[0]; scalingType {
		case "autoscaling":
			a = &autoScaling{
				asgName:         parts[1],
				actionType:      parts[2],
				capacityMin:     parts[4],
				capacityDesired: parts[5],
				capacityMax:     parts[6],
			}
		case "scheduler":
			t := scheduler.Target{}
			t.SetArn(parts[2])
			t.SetRoleArn(parts[4])
			s.SetTarget(&t)
		}
		switch kingpin.MustParse(app.Parse(os.Args[1:])) {

		case register.FullCommand():
			dateUTC, err := addMin(date, parts[3])
			if err != nil {
				slog.Error("Invalid date format: %w", err)
				continue
			}

			switch scalingType := parts[0]; scalingType {
			case "autoscaling":
				a.dateUTC = dateUTC
				slog.Debug("%v", a)
				fmt.Printf("%v\n", a.Register())
			case "scheduler":
				s.SetScheduleExpression(strings.TrimSuffix(dateUTC, "Z"))
				s.SetName(parts[1])
				f := scheduler.FlexibleTimeWindow{}
				f.SetMode(scheduler.FlexibleTimeWindowModeOff)
				s.SetActionAfterCompletion(scheduler.ActionAfterCompletionDelete)
				sc := Scheduler{
					Name:                  s.Name,
					ScheduleExpression:    s.ScheduleExpression,
					Target:                s.Target,
					FlexibleTimeWindow:    &f,
					ActionAfterCompletion: s.ActionAfterCompletion,
				}
				fmt.Printf("%v\n", sc.Register())
			}

		case deregister.FullCommand():
			switch scalingType := parts[0]; scalingType {
			case "autoscaling":
				fmt.Printf("%v\n", a.Deregister())
			case "scheduler":
				sc := Scheduler{
					Name:               s.Name,
					ScheduleExpression: s.ScheduleExpression,
					Target:             s.Target,
					FlexibleTimeWindow: s.FlexibleTimeWindow,
				}
				fmt.Printf("%v\n", sc.Deregister())
			}
		}
	}
}

func addMin(date *string, addMin string) (string, error) {
	a, err := strconv.ParseFloat(addMin, 64)
	if err != nil {
		slog.Error("Invalid a: %w", err)
		return "", err
	}

	loc := time.FixedZone("Asia/Tokyo", 9*60*60)
	d, err := time.ParseInLocation("2006-01-02 15:04", *date, loc)
	if err != nil {
		slog.Error("Invalid date format: %w", err)
		return "", err
	}
	slog.Debug("date: %t", date)

	dateJST := d.Add(time.Minute * time.Duration(a)).Local()
	slog.Debug("Date Asia/Tokyo: %v\n", dateJST)
	return dateJST.UTC().Format("2006-01-02T15:04:05Z"), nil

}
