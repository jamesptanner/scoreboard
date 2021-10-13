package scoreboard

type Config struct {
	Background     string `json:"background"`
	HomeTeam       string `json:"hometeam"`
	AwayTeam       string `json:"awayteam"`
	HomeColour     string
	AwayColour     string
	BarWidth       int
	Framerate      int `json:"framerate"`
	HomeStartGoals int `json:"homeStartGoals"`
	AwayStartGoals int `json:"homeAwayGoals"`
	Width          int
	Height         int
	Duration       int
	Goals          []Goal `json:"goals"`
	Margin         int
	FontSize       int
}

type Goal struct {
	HomeGoal  bool   `json:"hometeam"`
	TimeStamp string `json:"timestamp"`
	Frame     int    `json:"frame"`
}
