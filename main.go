package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func main() {
	info, err := host.Info()
	temps, err := host.SensorsTemperatures()
	cpuInfo, err := cpu.Info()
	cpuUsage, err := cpu.Percent(0, false)
	memoryInfo, err := mem.VirtualMemory()
	diskUsage, err := disk.Usage("/")
	netConns, err := net.Connections("")
	throw(err)

	var highestTemp float64 = 0
	for _, v := range temps {
		if v.Temperature > highestTemp {
			highestTemp = v.Temperature
		}
	}

	// Literal spaghetti:
	res := fmt.Sprintf(`
%s

%s
Platform:   %s
Version:    %s
Kernel Ver: %s
Arch:       %s
CPU:        %s
Cores:      %s
RAM:        %s
Disk Space: %s

%s
Boot Time:    %s
Uptime:       %s
CPU Usage:    %s
RAM Usage:    %s
Used Storage: %s
Free Storage: %s
Connections:  %s

%s
Temperature: %s
`,
		bb(info.Hostname),
		rb("System"),
		gb(info.OS),
		gb(info.PlatformVersion),
		gb(info.KernelVersion),
		gb(info.KernelArch),
		gb(cpuInfo[0].ModelName),
		gb(strconv.Itoa(int(cpuInfo[0].Cores))),
		gb(strconv.Itoa(int(memoryInfo.Total))+" bytes"),
		gb(strconv.Itoa(int(diskUsage.Total))+" bytes"),
		rb("Statistics"),
		gb(fmt.Sprint(time.Duration(info.BootTime))),
		gb(fmt.Sprint(time.Duration(info.Uptime)*time.Second)),
		gb(strconv.Itoa(int(cpuUsage[0]))+"%"),
		gb(strconv.Itoa(int(memoryInfo.UsedPercent))+"%"),
		gb(strconv.Itoa(int(diskUsage.Used))+" bytes (in / path)"),
		gb(strconv.Itoa(int(diskUsage.Free))+" bytes (in / path)"),
		gb(strconv.Itoa(len(netConns))),
		rb("Sensors"),
		gb(strconv.Itoa(int(highestTemp))+"°"),
	)

	fmt.Println(res)
}

func rb(s string) (r aurora.Value) {
	r = aurora.Bold(aurora.Red(s))
	return
}

func gb(s string) (r aurora.Value) {
	r = aurora.Bold(aurora.Green(s))
	return
}

func bb(s string) (r aurora.Value) {
	r = aurora.Bold(aurora.Blue(s))
	return
}

func throw(e error) {
	if e != nil {
		panic(e)
	}
}
