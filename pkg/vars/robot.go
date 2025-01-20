package vars

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var robotsMu sync.Mutex

var aRobots []Robot
var sRobots []SavedRobot

type Robot struct {
	Active      bool   `json:"active"`
	IP          string `json:"ip"`
	ESN         string `json:"esn"`
	CurrentGUID string `json:"guid"`
	Name        string `json:"name"`
}

type SavedRobot struct {
	IP          string `json:"ip"`
	ESN         string `json:"esn"`
	CurrentGUID string `json:"guid"`
	Name        string `json:"name"`
}

type RobotToReturn struct {
	Active      bool   `json:"active"`
	IP          string `json:"ip"`
	ESN         string `json:"esn"`
	CurrentGUID string `json:"guid"`
	Name        string `json:"name"`
}

func LoadRobots() error {
	file, err := os.Open(SavedRobotsFilePath)
	if err != nil {
		return err
	}
	err = json.NewDecoder(file).Decode(&sRobots)
	if err != nil {
		return err
	}
	for _, r := range sRobots {
		aRobots = append(aRobots, Robot{
			Active:      false,
			IP:          r.IP,
			ESN:         r.ESN,
			CurrentGUID: r.CurrentGUID,
			Name:        r.Name,
		})
	}
	return nil
}

// ip OR esn OR name
func GetRobot(ip string, esn string, name string) (Robot, error) {
	robotsMu.Lock()
	defer robotsMu.Unlock()
	for _, r := range aRobots {
		if r.ESN == esn || r.IP == ip || r.Name == name {
			return r, nil
		}
	}
	return Robot{}, errors.New("robot not found")
}

func GetAllInactiveNames() []string {
	robotsMu.Lock()
	defer robotsMu.Unlock()
	var bots []string
	for _, r := range aRobots {
		if !r.Active {
			bots = append(bots, r.Name)
		}
	}
	return bots
}

func SetActive(esn string) error {
	robotsMu.Lock()
	defer robotsMu.Unlock()
	for i, r := range aRobots {
		if r.ESN == esn {
			aRobots[i].Active = true
			return nil
		}
	}
	return errors.New("setactive: bot " + esn + " did not exist...")
}

func SetInactive(esn string) error {
	robotsMu.Lock()
	defer robotsMu.Unlock()
	for i, r := range aRobots {
		if r.ESN == esn {
			aRobots[i].Active = false
			return nil
		}
	}
	return errors.New("setactive: bot " + esn + " did not exist...")
}

func IsBotInList(esn string) bool {
	robotsMu.Lock()
	defer robotsMu.Unlock()
	for _, r := range aRobots {
		if r.ESN == esn {
			return true
		}
	}
	return false
}

func SaveRobot(rIn Robot) error {
	robotsMu.Lock()
	defer robotsMu.Unlock()
	matched := false
	sMatched := false
	for i, r := range aRobots {
		if r.ESN == rIn.ESN {
			aRobots[i].Active = rIn.Active
			aRobots[i].CurrentGUID = rIn.CurrentGUID
			aRobots[i].ESN = rIn.ESN
			aRobots[i].IP = rIn.IP
			aRobots[i].Name = rIn.Name
			matched = true
		}
	}
	for i, r := range sRobots {
		if r.ESN == rIn.ESN {
			sRobots[i].CurrentGUID = rIn.CurrentGUID
			sRobots[i].ESN = rIn.ESN
			sRobots[i].IP = rIn.IP
			sRobots[i].Name = rIn.Name
			sMatched = true
		}
	}
	if !matched {
		aRobots = append(aRobots, rIn)
	}
	if !sMatched {
		sRobots = append(sRobots, SavedRobot{
			IP:          rIn.IP,
			ESN:         rIn.ESN,
			Name:        rIn.Name,
			CurrentGUID: rIn.CurrentGUID,
		})
	}
	f, err := os.Open(SavedRobotsFilePath)
	if err != nil {
		return err
	}
	err = json.NewEncoder(f).Encode(sRobots)
	if err != nil {
		return err
	}
	return nil
}
