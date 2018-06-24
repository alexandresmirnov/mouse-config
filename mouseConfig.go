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

func check(e error) {
  if e != nil {
    panic(e)
  }
}

type ControlAccel struct {
  Button int
  Factor float32
}

type Mouse struct {
  Name string
	Decel float32
  Linear bool
  ButtonMap string
  ControlAccel ControlAccel
}

type Config struct {
  Mice []Mouse
}

func execCmd(cmd string) string {
  out, err := exec.Command("bash", "-c", cmd).Output() ; check(err)

  return string(out)
}

func execCmdAsync(cmd string) {
  execCmd := exec.Command("bash", "-c", cmd)
  execCmd.Start()
}

func setProp(id, prop, value string){
  cmd := fmt.Sprintf("xinput set-prop %s %s %s", id, prop, value)
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

func main() {
  configFileName := "config.yaml"
  if len(os.Args) == 2 {
    configFileName = os.Args[1]
  } else if len(os.Args) >= 2 {
    fmt.Printf("Error: too many arguments")
  }

  raw, err := ioutil.ReadFile(configFileName) ; check(err)

	var config Config
  yamlErr := yaml.Unmarshal([]byte(raw), &config) ; check(yamlErr)

  //fmt.Printf("%v+\n", config)

  for _, mouse := range config.Mice {
    name := mouse.Name
    decel := mouse.Decel
    linear := mouse.Linear
    buttonMap := mouse.ButtonMap
    controlAccel := mouse.ControlAccel
    id := mouseNameToID(name)

    setProp(id, CONSTANT_DECELERATION, floatToString(decel))

    if linear {
      setProp(id, PROFILE, "-1")
    } else {
      setProp(id, PROFILE, "0")
    }

    if len(buttonMap) > 0 {
      setButtonMapCmd := fmt.Sprintf("xinput -set-button-map %s %s", id, buttonMap)
      execCmd(setButtonMapCmd)
    }

    if controlAccel.Button != 0 {
      fmt.Println("launching control script")
      launchMouseControl := fmt.Sprintf("bash %s/dotfiles/scripts/mousecontrol.sh '%s' %s %s", 
        os.Getenv("HOME"), name, intToString(controlAccel.Button), floatToString(controlAccel.Factor))

      execCmdAsync(launchMouseControl)
    }




  }

}
