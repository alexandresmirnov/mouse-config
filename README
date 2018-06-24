### Mouse config

This is a script that reads in a YAML config file and sets various xinput options.  Currently dependent on my scripts/mousecontrol.sh, I'll probably move that into here at some point in the near future.

### Requirements

- `github.com/ghodss/yaml`

### Usage

`mouseconfig [control.yaml]`

### YAML structure

Structure:

```yaml
mice:
    - name: string
      decel: int
      linear: bool
      buttonMap: string
      controlAccel:
        button: int
        factor: float
```

Example:

```yaml
mice:
    - name: "Contour"
      decel: 4 

    - name: "Kingsis Peripherals Evoluent VerticalMouse 4"
      decel: 3.5
      linear: true
      buttonMap: "1 9 3 4 5 6 7 6 9 7"
      controlAccel:
        button: 2 
        factor: 0.5 

```
