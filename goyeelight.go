// Package goyeelight - Control the Yeelight LED Bulb with Go
package goyeelight

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
)

const timeout = time.Duration(10 * time.Second)

// Yeelight instance.
// Create an instance of Yeelight, by using New()
type Yeelight struct {
	host, port string
}

type (
	// Result struct is used on the standard response message
	Result struct {
		Status bool            `json:"status"`
		Data   json.RawMessage `json:"data"`
	}

	// ResponseOk struct is used on the success responses
	ResponseOk struct {
		ID     int             `json:"id"`
		Result json.RawMessage `json:"result"`
	}

	// ResponseError struct is used on the error responses
	ResponseError struct {
		ID    int   `json:"id"`
		Error Error `json:"error"`
	}

	// Error struct is used on the ResponseError payload
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

// Makes the request
func (y *Yeelight) request(cmd string) (string, error) {
	conn, err := net.DialTimeout("tcp", y.host+":"+y.port, timeout)
	if err != nil {
		return "", err
	}

	conn.SetReadDeadline(time.Now().Add(timeout))
	fmt.Fprintf(conn, cmd+"\r\n")

	data, err := bufio.NewReader(conn).ReadString('\n')
	conn.Close()
	if err != nil {
		return "", err
	}
	return response(data)
}

// Handles the response
func response(data string) (string, error) {
	res := ResponseOk{}
	json.Unmarshal([]byte(data), &res)

	if res.Result == nil {
		// error
		res := ResponseError{}
		json.Unmarshal([]byte(data), &res)

		err := errors.New(res.Error.Message)
		return "", err
	}

	// okay
	return string(res.Result), nil
}

// New returns a new Yeelight instance.
func New(host, port string) *Yeelight {
	y := &Yeelight{host: host, port: port}
	return y
}

// GetProp method is used to retrieve current property of smart LED.
func (y *Yeelight) GetProp(values ...string) (map[string]string, error) {
	cmd := `{"id":1,"method":"get_prop","params":[`
	for _, value := range values {
		cmd += `"` + string(value) + `",`
	}
	cmd += `]}`

	res, err := y.request(cmd)
	if err != nil {
		return nil, err
	}

	props := make([]string, 0)
	json.Unmarshal([]byte(res), &props)

	if len(props) != len(values) {
		err := errors.New("Wrong response")
		return nil, err
	}

	m := make(map[string]string, 0)
	for i, prop := range props {
		m[values[i]] = prop
	}

	return m, nil
}

// SetCtAbx method is used to change the color temperature of a smart LED.
func (y *Yeelight) SetCtAbx(value, effect, duration string) (string, error) {
	cmd := `{"id":2,"method":"set_ct_abx","params":[` + value + `,"` + effect + `",` + duration + `]}`
	return y.request(cmd)
}

// SetRGB method is used to change the color RGB of a smart LED.
func (y *Yeelight) SetRGB(value, effect, duration string) (string, error) {
	cmd := `{"id":3,"method":"set_rgb","params":[` + value + `,"` + effect + `",` + duration + `]}`
	return y.request(cmd)
}

// SetHSV method is used to change the color of a smart LED.
func (y *Yeelight) SetHSV(hue, sat, effect, duration string) (string, error) {
	cmd := `{"id":4,"method":"set_hsv","params":[` + hue + `,` + sat + `,"` + effect + `",` + duration + `]}`
	return y.request(cmd)
}

// SetBright method is used to change the brightness of a smart LED.
func (y *Yeelight) SetBright(brightness, effect, duration string) (string, error) {
	cmd := `{"id":5,"method":"set_bright","params":[` + brightness + `,"` + effect + `",` + duration + `]}`
	return y.request(cmd)
}

// SetPower method is used to switch on or off the smart LED (software managed on/off).
func (y *Yeelight) SetPower(power, effect, duration string) (string, error) {
	cmd := `{"id":6,"method":"set_power","params":["` + power + `","` + effect + `",` + duration + `]}`
	return y.request(cmd)
}

// Toogle method is used to toggle the smart LED.
// Note: This method is defined because sometimes user may just want
// to flip the state without knowing the current state.
func (y *Yeelight) Toogle() (string, error) {
	cmd := `{"id":7,"method":"toggle","params":[]}`
	return y.request(cmd)
}

// SetDefault method is used to save current state of smart LED in persistent
// memory. So if user powers off and then powers on the smart LED again (hard power reset),
// the smart LED will show last saved state.
func (y *Yeelight) SetDefault() (string, error) {
	cmd := `{"id":8,"method":"set_default","params":[]}`
	return y.request(cmd)
}

// StartCf method is used to start a color flow. Color flow is a series of smart
// LED visible state changing. It can be brightness changing, color changing or color
// temperature changing.This is the most powerful command. All our recommended scenes,
// e.g. Sunrise/Sunset effect is implemented using this method. With the flow expression, user
// can actually “program” the light effect.
func (y *Yeelight) StartCf(count, action, flowExpression string) (string, error) {
	cmd := `{"id":9,"method":"start_cf","params":[` + count + `,` + action + `,"` + flowExpression + `"]}`
	return y.request(cmd)
}

// StopCf method is used to stop a running color flow.
func (y *Yeelight) StopCf() (string, error) {
	cmd := `{"id":10,"method":"stop_cf","params":[]}`
	return y.request(cmd)
}

// SetScene method is used to set the smart LED directly to specified state.
// If the smart LED is off, then it will turn on the smart LED firstly and then
// apply the specified command.
func (y *Yeelight) SetScene(class, values string) (string, error) {
	cmd := `{"id":11,"method":"set_scene","params":["` + class + `",` + values + `]}`
	fmt.Println(cmd)
	return y.request(cmd)
}

// CronAdd method is used to start a timer job on the smart LED.
func (y *Yeelight) CronAdd(t, value string) (string, error) {
	cmd := `{"id":12,"method":"cron_add","params":[` + t + `,` + value + `]}`
	return y.request(cmd)
}

// CronGet method is used to retrieve the setting of the current cron job of the specified type.
func (y *Yeelight) CronGet(t string) (string, error) {
	cmd := `{"id":13,"method":"cron_get","params":[` + t + `]}`
	return y.request(cmd)
}

// CronDel method is used to stop the specified cron job.
func (y *Yeelight) CronDel(t string) (string, error) {
	cmd := `{"id":14,"method":"cron_del","params":[` + t + `]}`
	return y.request(cmd)
}

// SetAdjust method is used to change brightness, CT or color of a smart LED
// without knowing the current value, it's main used by controllers.
func (y *Yeelight) SetAdjust(action, prop string) (string, error) {
	cmd := `{"id":15,"method":"set_adjust","params":["` + action + `","` + prop + `"]}`
	return y.request(cmd)
}

// SetName method is used to name the device. The name will be stored on the
// device and reported in discovering response. User can also read the name
// through “get_prop” method
func (y *Yeelight) SetName(name string) (string, error) {
	cmd := `{"id":16,"method":"set_name","params":["` + name + `"]}`
	return y.request(cmd)
}

// On method is used to switch on the smart LED
func (y *Yeelight) On() (string, error) {
	return y.SetPower("on", "smooth", "1000")
}

// Off method is used to switch off the smart LED
func (y *Yeelight) Off() (string, error) {
	return y.SetPower("off", "smooth", "1000")
}
