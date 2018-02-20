/*

=======================================================

gobatt - Lightweight battery tray icon for Linux.

Repository: https://github.com/solusipse/gobatt

=======================================================

The MIT License (MIT)

Copyright (c) 2013 solusipse

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

var acpiPaths = []string{}

const (
	ACPIROOT    = "/sys/class/power_supply/BAT"
	UPDATE_TIME = 1
)

var lastPercentage float64
var timeSlice [10]float64

func main() {
	if err := initAcpiPaths(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	icon := trayIconInit()

	glib.TimeoutAdd(UPDATE_TIME*1000, func() bool {
		batteryStatus, batteryPercentage := updateData()
		setTrayIcon(icon, batteryStatus, batteryPercentage)
		return true
	})

	glib.TimeoutAdd(10000, func() bool {
		batteryStatus, batteryPercentage := updateData()
		getRemainingTime(icon, batteryStatus, batteryPercentage)
		return true
	})

	gtk.Main()
}

func initAcpiPaths() error {
	items, err := filepath.Glob(ACPIROOT + "*")
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return errors.New("no batteries found")
	}
	for _, item := range items {
		stat, err := os.Stat(item)
		if err != nil {
			fmt.Println("Stat error with item:", item)
			continue
		}
		if stat.IsDir() {
			item += "/"
			fmt.Println("Found battery:", item)
			acpiPaths = append(acpiPaths, item)
		} else {
			fmt.Println("Skipping non-directory item:", item)
		}
	}
	return nil
}

func getFileContent(base, filename string) string {
	content, _ := ioutil.ReadFile(base + filename)
	return string(content)
}

func getBatteryState() string {
	return strings.TrimSuffix(getFileContent(acpiPaths[0], "status"), "\n")
}

func getBatteryPercentage() float64 {
	result := float64(0)
	for _, acpiPath := range acpiPaths {
		_fc := strings.TrimSuffix(getFileContent(acpiPath, "energy_full"), "\n")
		_nc := strings.TrimSuffix(getFileContent(acpiPath, "energy_now"), "\n")
		fullCap, _ := strconv.Atoi(_fc)
		nowCap, _ := strconv.Atoi(_nc)
		result += (float64(nowCap) / float64(fullCap))
		fmt.Println("Result for:", acpiPath, "is:", (float64(nowCap) / float64(fullCap)))
	}
	result /= float64(len(acpiPaths))
	fmt.Println("Result average:", result)
	return result
}

func updateData() (string, float64) {
	return getBatteryState(), getBatteryPercentage()
}

func trayIconInit() *gtk.StatusIcon {
	gtk.Init(nil)
	glib.SetApplicationName("gobatt")

	icon := gtk.NewStatusIcon()
	icon.SetTitle("gobatt")

	return icon
}

func getGtkIcon(percent float64, status string) string {
	percent = percent * 100
	if status == "Discharging" {
		if percent <= 10 {
			return "battery-caution-symbolic"
		} else if percent <= 20 {
			return "battery-empty-symbolic"
		} else if percent <= 45 {
			return "battery-low-symbolic"
		} else if percent <= 75 {
			return "battery-good-symbolic"
		} else if percent <= 100 {
			return "battery-full-symbolic"
		}
	}
	if status == "Charging" {
		if percent <= 10 {
			return "battery-caution-charging-symbolic"
		} else if percent <= 20 {
			return "battery-empty-charging-symbolic"
		} else if percent <= 45 {
			return "battery-low-charging-symbolic"
		} else if percent <= 75 {
			return "battery-good-charging-symbolic"
		} else if percent <= 99 {
			return "battery-full-charging-symbolic"
		} else if percent <= 100 {
			return "battery-full-charged-symbolic"
		}
	}
	if status == "Full" {
		return "battery-full-charged-symbolic"
	}

	return "battery-missing-symbolic"
}

func addTimeRecord(record float64) {
	if timeSlice[9] != 0 {
		var bufferSlice [10]float64
		for i := 0; i < 9; i++ {
			bufferSlice[i+1] = timeSlice[i]
		}
		timeSlice = bufferSlice
		timeSlice[0] = record
	} else {
		for i, j := range timeSlice {
			if j == 0 {
				timeSlice[i] = record
				break
			}
		}
	}
}

func getAverageTime() int {
	if timeSlice[9] != 0 {
		var buffer float64 = 0
		for _, j := range timeSlice {
			buffer += j
		}
		return int(buffer / 10)
	}
	return -1
}

func getRemainingTime(icon *gtk.StatusIcon, status string, percent float64) {
	if lastPercentage == 0 {
		lastPercentage = percent
	}

	if lastPercentage > percent {
		remaining := ((10 * percent) / (lastPercentage - percent)) / 60

		addTimeRecord(remaining)
		lastPercentage = percent
	}

	if lastPercentage < percent {
		remaining := ((10 * (1 - percent)) / (percent - lastPercentage)) / 60

		addTimeRecord(remaining)
		lastPercentage = percent
	}

}

func getTooltipString(percent float64, status string, time int) string {
	if percent*100 >= 99 {
		return "Battery is fully charged."
	}

	tooltipString := status
	tooltipString += ": " + strconv.Itoa(int(percent*100)) + "%\n"

	if time == -1 {
		tooltipString += "Remaining time: estimating."
	} else {
		hours := time / 60
		minutes := time - hours*60
		tooltipString += "Remaining time: " + strconv.Itoa(hours) + "h " +
			strconv.Itoa(minutes) + "m."
	}

	return tooltipString
}

func setToolTip(icon *gtk.StatusIcon, status string, percent float64, time int) {
	icon.SetTooltipMarkup(getTooltipString(percent, status, time))
}

func setTrayIcon(icon *gtk.StatusIcon, status string, percent float64) {
	iconName := getGtkIcon(percent, status)

	if icon.GetIconName() != iconName {
		icon.SetFromIconName(iconName)
	}
	setToolTip(icon, status, percent, getAverageTime())
}
