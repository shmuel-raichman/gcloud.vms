package utils

import (
	"fmt"
	"html"
	"strconv"
)

const (
	Thunderstorm = "0001F4A8" //# Code: 200's, 900, 901, 902, 905
	Drizzle      = "0001F4A7" //# Code: 300's
	Rain         = "00002614" //# Code: 500's
	Snowflake    = "00002744" //# Code: 600's snowflake
	Snowman      = "000026C4" //# Code: 600's snowman, 903, 906
	Atmosphere   = "0001F301" //# Code: 700's foogy
	ClearSky     = "00002600" //# Code: 800 clear sky
	FewClouds    = "000026C5" //# Code: 801 sun behind clouds
	Clouds       = "00002601" //# Code: 802-803-804 clouds general
	Hot          = "0001F525" //# Code: 904
	DefaultEmoji = "0001F300" //# default emojis
)

func GetEmoji(emojiName string) {

	// Hex to Int
	i, _ := strconv.ParseInt(emojiName, 0, 64)

	// Unescape the string (HTML Entity -> String).
	str := html.UnescapeString(string(i))

	// Display the emoji.
	fmt.Println(str)
}
