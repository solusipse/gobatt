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
var lastTime float64

func main() {
    icon := trayIconInit()

    lastPercentage = -1
    lastTime = -1

    glib.TimeoutAdd(UPDATE_TIME*1000, func() bool {
        batteryStatus, batteryPercentage := updateData()
        setTrayIcon(icon, batteryStatus, batteryPercentage)
        return true
    })

    glib.TimeoutAdd(5000, func() bool {
        batteryStatus, batteryPercentage := updateData()
        getRemainingTime(batteryStatus, batteryPercentage)
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
        } else if percent <= 25 {
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

func getRemainingTime(status string, percent float64) int {
    /* TODO: this method */
    if lastPercentage == -1 {
        lastPercentage = percent
        lastTime = -1
        return -1
    }

    if lastPercentage > percent {
        remainingFloat := ((lastPercentage - percent)*5000)
        fmt.Println((remainingFloat+lastTime)/2)

        lastPercentage = percent
        lastTime = remainingFloat
    }

    return 5
}

func getTooltipString(percent float64, status string) string {
    if percent*100 >= 99 {
        return "Battery is fully charged."
    }

    tooltipString := status
    tooltipString += ": " + strconv.FormatFloat(percent*100, 'g', 2, 64) + "%\n"
    tooltipString += "Remaining time: "
    return tooltipString
}

func setTrayIcon(icon *gtk.StatusIcon, status string, percent float64) {
    iconName := getGtkIcon(percent, status)
    icon.SetFromIconName(iconName)
    icon.SetTooltipMarkup(getTooltipString(percent, status))
}