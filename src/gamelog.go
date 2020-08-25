package main

import (
    "errors"
    "log"
    "regexp"
    "strings"
    "time"

    "github.com/hpcloud/tail"
)


func tailLog(filename string) ([]string, error) {
	result := []string{}

	t, err := tail.TailFile(config.FactorioLog, tail.Config{Follow: false})
	if err != nil {
		log.Printf("Error tailing log %s", err)
		return result, err
	}

	for line := range t.Lines {
		result = append(result, line.Text)
	}

  result = reformatTimestamps(result)

  return result, nil
}


func getOffset(line string) (string, error) {
  re, _ := regexp.Compile(`^\d+.\d+`)

  if !re.MatchString(line) {
    log.Printf("This line has no offset", line)
    return "error", errors.New(line)
  }

  offset := re.FindString(line)

  return offset, nil
}


func getStartTime(line string) (time.Time) {
  re, _ := regexp.Compile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)
  date := string(re.FindString(line))
  startTime, _ := time.Parse(time.RFC3339, strings.Replace(date, " ", "T", 1) + "Z")

  return startTime
}


func replaceTimestampInLine(line string, offset string, startTime time.Time) (string) {
    offset, err := getOffset(line)
    offsetDuration, _ := time.ParseDuration(offset + "s")
    timestamp := startTime.Add(offsetDuration)

    if err == nil  {
      return timestamp.Format("2006-01-02 03:04:05") +  ":" + strings.Replace(line, offset, "", 1)
    } else {
      return line
    }
}


func reformatTimestamps(log []string) ([]string) {
  firstLine := log[len(log)-1]
  startTime := getStartTime(firstLine)
  result := []string{}

  for i := range log {
    line := log[i]
    offset, _ := getOffset(line)
    result = append(result, replaceTimestampInLine(line, offset, startTime))
  }

  return result
}
