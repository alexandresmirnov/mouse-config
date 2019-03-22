package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/ghodss/yaml"
)

const libinputAccelSpeed = "libinput Accel Speed"
const libinputAccelProfileEnabled = "libinput Accel Profile Enabled"

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}

type config struct {
	Mice []mouse
}

type mouse struct {
	Name           string
	Accel          float32
	AccelProfile   string
	ButtonMap      string
	ControlAccel   controlAccel
	StabilizeClick controlAccel
	CustomProps    []customProp
}

type controlAccel struct {
	Button int
	Factor float32
	Type   string
}

type customProp struct {
	Name  string
	Value string
}

func execCmd(cmd string) string {
	out, err := exec.Command("bash", "-c", cmd).Output()

	if err != nil {
		fmt.Printf("FAILURE in \"%s\": %s\n", cmd, err)
	}

	return string(out)
}

func execCmdAsync(cmd string) {
	execCmd := exec.Command("bash", "-c", cmd)
	execCmd.Start()
}

func setProp(id, prop, value string) {
	cmd := fmt.Sprintf("xinput set-prop %s '%s' %s", id, prop, value)
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
	out, err := exec.Command("bash", "-c", "xinput list | grep "+"\""+name+"\"").Output()

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

	var config config
	yamlErr := yaml.Unmarshal([]byte(raw), &config)
	check(yamlErr)

	for _, mouse := range config.Mice {
		name := mouse.Name

		if !testMouseExists(name) {
			continue
		}

		fmt.Printf("\n%s:\n", name)

		// obtain device ID
		id := mouseNameToID(name)

		// easier reference to props
		accel := mouse.Accel
		accelProfile := mouse.AccelProfile
		buttonMap := mouse.ButtonMap
		controlAccel := mouse.ControlAccel
		stabilizeClick := mouse.StabilizeClick
		customProps := mouse.CustomProps

		// set accel
		setProp(id, libinputAccelSpeed, floatToString(accel))

		// set accelProfile
		if accelProfile != "" {
			if accelProfile == "linear" {
				setProp(id, libinputAccelProfileEnabled, "0 0")
			} else {
				setProp(id, libinputAccelProfileEnabled, accelProfile)
			}
		}

		// set buttonMap
		if len(buttonMap) > 0 {
			setButtonMapCmd := fmt.Sprintf("xinput -set-button-map %s %s", id, buttonMap)
			fmt.Println(setButtonMapCmd)
			execCmd(setButtonMapCmd)
		}

		// set up control accel script
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

		// set up stabilize click script
		if stabilizeClick.Button != 0 {

			launchMouseControl := fmt.Sprintf("bash %s/dotfiles/scripts/stabilizeclick.sh '%s' %s %s",
				os.Getenv("HOME"), name, intToString(stabilizeClick.Button), floatToString(stabilizeClick.Factor))

			execCmdAsync(launchMouseControl)

		}

		// check presence of customProps
		if len(customProps) > 0 {

			for _, prop := range customProps {

				name := prop.Name
				value := prop.Value

				setProp(id, name, value)

			}

		}

	}

}
