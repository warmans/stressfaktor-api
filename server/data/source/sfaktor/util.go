package sfaktor

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

var localTime *time.Location

func init() {
	var err error
	localTime, err = time.LoadLocation("Europe/Berlin")
	if err != nil {
		log.Fatal("Cannot load localtime")
	}
}

//ParseTime parses times in the format e.g. "Montag, 21.12.2015", "21:00 Uhr"
func ParseTime(dateString, timeString string) (time.Time, error) {

	if strings.Contains(dateString, ",") {
		cleanDate := fmt.Sprintf(
			"%s %s",
			strings.Trim(timeString, " Uhr"),
			strings.TrimSpace(strings.Split(dateString, ",")[1]),
		)
		return time.ParseInLocation("15.04 02.01.2006", cleanDate, localTime)
	}

	return time.Time{}, errors.New("header date format looks wrong")
}

//HTML Parsing partially stolen from: https://github.com/kennygrant/sanitize/blob/master/sanitize.go
func StripHTML(s string) string {
	output := ""
	if !strings.ContainsAny(s, "<>") {
		output = s
	} else {

		s = strings.Replace(s, "\n", " ", -1)

		// Walk through the string removing all tags
		b := bytes.NewBufferString("")
		inTag := false
		for _, r := range s {
			switch r {
			case '<':
				inTag = true
			case '>':
				inTag = false
			default:
				if !inTag {
					b.WriteRune(r)
				}
			}
		}
		output = b.String()
	}
	return output
}