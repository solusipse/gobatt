package main

import (
        "fmt"
        "strings"
        "strconv"
        "io/ioutil"
        "github.com/mattn/go-gtk/glib"
        "github.com/mattn/go-gtk/gtk"
)

const (
        ACPIPATH = "/sys/class/power_supply/BAT1/"
        UPDATE_TIME = 1
)

var lastPercentage float64

var timeSlice [10] float64

func main() {
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

func getFileContent(filename string) string {
    content, _ := ioutil.ReadFile(ACPIPATH + filename)
    return string(content)
}

func getBatteryState() string {
    return strings.TrimSuffix(getFileContent("status"), "\n")
}

func getBatteryPercentage() float64 {
    _fc := strings.TrimSuffix(getFileContent("energy_full"), "\n")
    _nc := strings.TrimSuffix(getFileContent("energy_now"), "\n")
    fullCap, _ := strconv.Atoi(_fc)
    nowCap,  _ := strconv.Atoi(_nc)
    return (float64(nowCap)/float64(fullCap))
}

func updateData() (string, float64) {
    return getBatteryState(), getBatteryPercentage()
}

func trayIconInit() *gtk.StatusIcon {
    gtk.Init(nil)
    glib.SetApplicationName("wm-batt-tray")

    //popupInfo := gtk.NewMenuItemWithLabel("")

    icon := gtk.NewStatusIcon()
    icon.SetTitle("wm-batt-tray")

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
        return int(buffer/10)
    }
    return -1
}

func getRemainingTime(icon *gtk.StatusIcon, status string, percent float64) {
    /* TODO: this method */
    if lastPercentage == 0 {
        lastPercentage = percent
    }

    if lastPercentage > percent {
        remainingFloat := ((10 * percent)/(lastPercentage - percent))/60

        addTimeRecord(remainingFloat)

        if getAverageTime() == -1 {
            fmt.Println("Estimating")
        } else {
            fmt.Println(getAverageTime())
        }

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
        hours := time/60
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
    // TODO: don't update when icon hasn't changed
    iconName := getGtkIcon(percent, status)
    icon.SetFromIconName(iconName)
    setToolTip(icon, status, percent, getAverageTime())
}