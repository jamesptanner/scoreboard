package scoreboard

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/wargarblgarbl/libgosubs/srt"
)

func ProcessGoals(config *Config, hometeam bool) {
	goals := []Goal{}

	reverse := func(s []string) []string {
		a := make([]string, len(s))
		copy(a, s)

		for i := len(a)/2 - 1; i >= 0; i-- {
			opp := len(a) - 1 - i
			a[i], a[opp] = a[opp], a[i]
		}

		return a
	}

	safeIndex := func(s []string, i int, d string) string {
		if len(s) > i {
			return s[i]
		}
		return d
	}

	for _, v := range config.Goals {
		if v.HomeGoal == hometeam {
			if v.Frame == nil {
				if v.TimeStamp == nil {
					continue
				}
				//convert timestamp to frame.
				//xx:xx:xx.xxx
				parts := regexp.MustCompile(`[\.:]`).Split(*v.TimeStamp, -1)
				fixedOrder := reverse(parts)

				ms, _ := strconv.Atoi(safeIndex(fixedOrder, 0, "0"))
				s, _ := strconv.Atoi(safeIndex(fixedOrder, 1, "0"))
				m, _ := strconv.Atoi(safeIndex(fixedOrder, 2, "0"))
				h, _ := strconv.Atoi(safeIndex(fixedOrder, 3, "0"))

				totalSeconds := s + (60 * m) + (3600 * h)

				v.Frame = new(int)
				*(v.Frame) = (totalSeconds * config.Framerate) + (int)((ms*config.Framerate)/1000)

			}
			goals = append(goals, v)

		}
	}

	sort.Slice(goals, func(i, j int) bool {
		return *(goals[i].Frame) < *(goals[j].Frame)
	})

	toTimestamp := func(frame *int, rate int) string {
		ms := (((*frame) % rate) * 1000) / rate
		s := *frame / rate
		m, s := s/60, s%60
		h, m := m/60, m%60
		return fmt.Sprintf("%02d:%02d:%02d,%03d", h, m, s, ms)
	}

	subs := srt.SubRip{}
	last := "00:00:00,000"
	for i := 0; i < len(goals); i++ {
		next := toTimestamp(goals[i].Frame, config.Framerate)
		sub := srt.CreateSubtitle(i+1, last, next, []string{fmt.Sprintf("%d", i)})
		subs.Subtitle.Content = append(subs.Subtitle.Content, *sub)
		last = next
		log.Println(last)
	}
	sub := srt.CreateSubtitle(len(goals), last, toTimestamp(&config.Duration, config.Framerate), []string{fmt.Sprintf("%d", len(goals))})
	subs.Subtitle.Content = append(subs.Subtitle.Content, *sub)
}

func RenderBoard(config *Config, outFileName *string) {

	//homeScore, awayScore := generateScores(config)
	ProcessGoals(config, true)
	err := ffmpeg.Input("/dev/zero",
		ffmpeg.KwArgs{
			"t":       config.Duration,
			"s":       fmt.Sprintf("%dx%d", config.Width, config.Height),
			"f":       "rawvideo",
			"pix_fmt": "rgb24",
			"r":       config.Framerate,
		}).
		DrawBox(0, 0, config.Width, config.Height, config.Background, config.Height).
		DrawBox(0, 0, config.BarWidth, config.Height, config.HomeColour, config.BarWidth).
		DrawBox(config.Width-config.BarWidth, 0, config.BarWidth, config.Height, config.AwayColour, config.BarWidth).
		Drawtext(config.HomeTeam, config.BarWidth+config.Margin, 0, false, ffmpeg.KwArgs{
			"y":        fmt.Sprintf("main_h/2-text_h/2"),
			"fontsize": config.FontSize,
		}).
		Drawtext(config.AwayTeam, 0, 0, false, ffmpeg.KwArgs{
			"x":        fmt.Sprintf("main_w-%d-text_w", config.BarWidth+config.Margin),
			"y":        fmt.Sprintf("main_h/2-text_h/2"),
			"fontsize": config.FontSize,
		}).
		Output(*outFileName, ffmpeg.KwArgs{"frames": config.Duration}).
		GlobalArgs("-report").
		OverWriteOutput().
		Run()

	if err != nil {
		log.Printf("Failed to create file: %v\n", err)
	}
}
