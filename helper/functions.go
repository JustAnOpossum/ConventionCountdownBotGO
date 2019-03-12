package helper

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"path"
	"strconv"
	"time"
)

//ConfigStruct configuration info
type ConfigStruct struct {
	Con              string
	WelcomeMsg       string
	SubMsg           string
	AlreadySubMsg    string
	UnsubMsg         string
	NotSubMsg        string
	CmdMsg           string
	InfoMsg          string
	GroupAddMsg      string
	GroupNotAdminMsg string
	Date             time.Time
	DBName           string
}

//GetDays gets days until con
func GetDays(day time.Time) string {
	timeUntil := day.Sub(time.Now())
	daysUntil := timeUntil.Hours() / 24
	daysRounded := math.Round(daysUntil)
	return strconv.Itoa(int(daysRounded))
}

//LoadConfig Takes datadir and loads config
func LoadConfig(dataDir string, configVar *ConfigStruct) {
	configFile, err := ioutil.ReadFile(path.Join(dataDir, "config.json"))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(configFile, configVar)
	if err != nil {
		panic(err)
	}
}
