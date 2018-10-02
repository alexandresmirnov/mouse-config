package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/ghodss/yaml"
)

const CONSTANT_DECELERATION = "'Device Accel Constant Deceleration'"
const PROFILE = "'Device Accel Profile'"

const LIBINPUT_ACCEL_SPEED = "'libinput Accel Speed'"
const LIBINPUT_ACCEL_PROFILE_ENABLED = "'libinput Accel Profile Enabled'"

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

type ControlAccel struct {
	Button int
	Factor float32
	Type   string
}

type Mouse struct {
	Name           string
	Accel          float32
	AccelProfile   string
	ButtonMap      string
	ControlAccel   ControlAccel
	StabilizeClick ControlAccel
}

type Config struct {
	Mice []Mouse
}

func execCmd(cmd string) string {
	out, err := exec.Command("bash", "-c", cmd).Output()
	check(err)

	return string(out)
}

func execCmdAsync(cmd string) {
	execCmd := exec.Command("bash", "-c", cmd)
	execCmd.Start()
}

func setProp(id, prop, value string) {
	cmd := fmt.Sprintf("xinput set-prop %s %s %s", id, prop, value)
	fmt.Println(cmd)
	execCmd(cmd)
}

func mouseNameToID(name string) string {
	searchForID := "xinput --list | grep -i -m 1 \"" + name + "\" | grep -o \"id=[0-9]\\+\" | grep -o \"[0-9]\\+\""
	id := execCmd(searchForID)
	id = id[:len(id)-1] // trim newline char

	return id
}

func floatToString(n float32) string {
	return fmt.Sprintf("%.4f", n)
}

func intToString(n int) string {
	return fmt.Sprintf("%d", n)
}

func testMouseExists(name string) bool {
	out, err := exec.Command("bash", "-c", "xinput list | grep "+name).Output()

	if string(out) == "" {
		return false
	} else if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	return true
}

func main() {
	configFileName := "config.yaml"
	if len(os.Args) == 2 {
		configFileName = os.Args[1]
	} else if len(os.Args) >= 2 {
		fmt.Printf("Error: too many arguments")
	}

	raw, err := ioutil.ReadFile(configFileName)
	check(err)

	var config Config
	yamlErr := yaml.Unmarshal([]byte(raw), &config)
	check(yamlErr)

	// fmt.Printf("%v+\n", config)

	for _, mouse := range config.Mice {
		name := mouse.Name

		if !testMouseExists(name) {
			continue
		}

		accel := mouse.Accel
		accelProfile := mouse.AccelProfile
		buttonMap := mouse.ButtonMap
		controlAccel := mouse.ControlAccel
		stabilizeClick := mouse.StabilizeClick
		id := mouseNameToID(name)

		setProp(id, LIBINPUT_ACCEL_SPEED, floatToString(accel))

		if accelProfile != "" {
			if accelProfile == "linear" {
				setProp(id, LIBINPUT_ACCEL_PROFILE_ENABLED, "0 0")
			} else {
				setProp(id, LIBINPUT_ACCEL_PROFILE_ENABLED, accelProfile)
			}
		}

		if len(buttonMap) > 0 {
			setButtonMapCmd := fmt.Sprintf("xinput -set-button-map %s %s", id, buttonMap)
			execCmd(setButtonMapCmd)
		}

		if controlAccel.Button != 0 {

			if controlAccel.Type == "primary" {

				launchMouseControl := fmt.Sprintf("bash %s/dotfiles/scripts/primarybuttonmousecontrol.sh '%s' %s %s",
					os.Getenv("HOME"), name, intToString(controlAccel.Button), floatToString(controlAccel.Factor))

				execCmdAsync(launchMouseControl)

			} else if controlAccel.Type == "secondary" || controlAccel.Type == "" {
				launchMouseControl := fmt.Sprintf("bash %s/dotfiles/scripts/mousecontrol.sh '%s' %s %s",
					os.Getenv("HOME"), name, intToString(controlAccel.Button), floatToString(controlAccel.Factor))

				execCmdAsync(launchMouseControl)

			}
		}

		if stabilizeClick.Button != 0 {

			launchMouseControl := fmt.Sprintf("bash %s/dotfiles/scripts/stabilizeclick.sh '%s' %s %s",
				os.Getenv("HOME"), name, intToString(stabilizeClick.Button), floatToString(stabilizeClick.Factor))

			execCmdAsync(launchMouseControl)

		}

	}

}
