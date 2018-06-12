/*
 for each line in parseFile(text) {
   t := newText(line)
   result := t.Download()
   print(result)
   result.WriteToJson()
   result.WriteToMongoDB()
 }
*/

package main

import (
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "time"
  "strings"
  "regexp"
  "encoding/json"
  "github.com/gutengo/tagen/io/ioutil2"
  "github.com/gutengo/shell"
  "github.com/GeertJohan/go.rice"
  "github.com/mxk/go-flowrate/flowrate"
  "gopkg.in/mgo.v2"
)

var TestDuration = 1*time.Minute
//var TestDuration = 1*time.Second

func SpeedTest(inputPath, outputPath string) {
  text := readInput(inputPath)

  for _, line := range parseFile(text) {
    name, url := line[0], line[1]
    t := NewTest(name, url)
    result := t.Download()
    shell.Say(result)
    switch outputPath {
    case "":
      // skip
    case "mongodb":
      result.WriteToMongoDB()
    default:
      result.WriteToJson(outputPath)
    }
  }
}

func readInput(path string) (text string) {
  if path == "" {
    text = rice.MustFindBox("assets").MustString("Speedtestfile")
  } else {
    bytes, err := ioutil.ReadFile(path)
    if err != nil {
      shell.ErrorExit(err)
    }
    text = string(bytes)
  }
  return text
}

// [ [name, url], ... ]
func parseFile(text string) (ret [][]string) {
  for _, line := range strings.Split(text, "\n") {
    if regexp.MustCompile(`^\s*$`).MatchString(line) {
      continue
    }
    parts := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(line), -1)
    ret = append(ret, parts)
  }
  return ret
}

type Test struct {
  Name string
  Url string
  result *Result
}

func NewTest(name, url string) *Test {
  return &Test{name, url, nil}
}

func (t *Test) Download() *Result {
  resp, err := http.Get(t.Url)
  if err != nil {
    shell.Error("Skip "+t.Name+": "+err.Error())
    t.result = nil
    return NewEmptyResult(t.Name)
  }
	r := flowrate.NewReader(resp.Body, -1)
  defer r.Close()
  r.SetTransferSize(resp.ContentLength)

  timeout := false
  time.AfterFunc(TestDuration, func() {
    timeout = true
    r.Close()
  })
  _, err = io.Copy(ioutil.Discard, r)
  if err != nil && !timeout {
    shell.Error("Get "+t.Name+": "+err.Error())
    t.result = nil
    return NewEmptyResult(t.Name)
  }

  status := r.Status()
  return &Result{
    Date: newDate(status.Start),
    Name: t.Name,
    Avg: status.AvgRate / 1024,
    Peak: status.PeakRate / 1024,
    Duration: status.Duration.Seconds(),
  }
}

type Result struct {
  Date time.Time
  Name string
  Avg int64
  Peak int64
  Duration float64
}

func NewEmptyResult(name string) *Result {
  return &Result{
    Date: newDate(time.Now()),
    Name: name,
    Avg: 0,
    Peak: 0,
    Duration: 0,
  }
}

func (r *Result) String() string {
  return fmt.Sprintf("%v %20v    %5v %5v %v\n",  r.Date.Format("2006-01-02T15:04:05"), r.Name, r.Avg, r.Peak, r.Duration)
}

func (r *Result) WriteToJson(path string) {
  bytes, err := json.Marshal(r)
  if err != nil {
    shell.ErrorExit(err)
  }
  bytes = append(bytes, '\n')
  err = ioutil2.AppendFile(path, bytes, 0644)
  if err != nil {
    shell.ErrorExit(err)
  }
}

func (r *Result) WriteToMongoDB() {
  session, err := mgo.Dial("localhost")
  if err != nil {
    shell.ErrorExit(err)
  }
  defer session.Close()
  collection := session.DB("freedom_speedtest").C("result")
  err = collection.Insert(r)
  if err != nil {
    shell.ErrorExit(err)
  }
}

// set minutes and seconds to 00.
func newDate(date time.Time) time.Time {
  t, _ := time.Parse(time.RFC3339, time.Now().Format("2006-01-02T15:00:00+08:00"))
  return t
}
