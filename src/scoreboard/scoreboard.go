package scoreboard

import (
	"fmt"
	"log"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func RenderBoard(config *Config, outFileName *string) {

	//homeScore, awayScore := generateScores(config)

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
