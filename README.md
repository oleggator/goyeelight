# goyeelight
Control the Yeelight LED Bulb with Go

[![GoDoc](https://godoc.org/github.com/nunows/goyeelight?status.svg)](https://godoc.org/github.com/nunows/goyeelight)

## Usage

The "Developer Mode" need to be enabled to discover and operate the device.

### Quick Start

#### Install

``` bash

go get github.com/nunows/goyeelight

```

#### Example to control the Yeelight WiFi LED



``` go
package main

import "fmt"
import "github.com/oleggator/goyeelight"

func main() {
	// new Yeelight instance
	lamp := goyeelight.New("192.168.0.27", "55443")

	var (
		r   string
		err error
		m   map[string]string
	)

	// turn on the smart LED
	r, err = lamp.On()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)

	// get the "power" and "bright" propertys
	m, err = lamp.GetProp("power", "bright")
	if err != nil {
		fmt.Println(err)
	}
	for prop, value := range m {
		fmt.Println(prop, "=", value)
	}

	// set the bright to 50
	lamp.SetBright("50", "smooth", "500")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)

	// turn off the smart LED
	r, err = lamp.Off()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)
}
```

### Methods available in the package:

* GetProp
* SetCtAbx
* SetRGB
* SetHSV
* SetBright
* SetPower
* Toogle
* SetDefault
* StartCf
* StopCf
* SetScene
* CronAdd
* CronDel
* CronGet
* SetAdjust
* SetName
* On
* Off

Note: Only tested with Yeelight LED Bulb (Color) with the firmware 1.4.1_45
