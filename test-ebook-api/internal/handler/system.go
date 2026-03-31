package handler

import (
	"fmt"
	"runtime"
	"test-ebook-api/internal/pkg"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemHandler struct {
	startTime time.Time
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{
		startTime: time.Now(),
	}
}

func (h *SystemHandler) GetSystemStatus(c *gin.Context) {
	// CPU usage (all cores combined)
	percent, _ := cpu.Percent(time.Second, false)
	cpuUsage := 0.0
	if len(percent) > 0 {
		cpuUsage = percent[0]
	}

	// Memory usage
	v, _ := mem.VirtualMemory()
	memUsage := 0.0
	if v != nil {
		memUsage = v.UsedPercent
	}

	// Disk usage (cross-platform)
	diskPath := "/"
	if runtime.GOOS == "windows" {
		diskPath = "C:\\"
	}
	d, _ := disk.Usage(diskPath)
	diskUsage := 0.0
	if d != nil {
		diskUsage = d.UsedPercent
	}

	// Uptime
	hostInfo, _ := host.Info()
	uptimeStr := "N/A"
	if hostInfo != nil {
		uptimeDuration := time.Duration(hostInfo.Uptime) * time.Second
		days := int(uptimeDuration.Hours() / 24)
		hours := int(uptimeDuration.Hours()) % 24
		minutes := int(uptimeDuration.Minutes()) % 60
		uptimeStr = fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}

	// Runtime stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	pkg.Success(c, gin.H{
		"cpu":       cpuUsage,
		"memory":    memUsage,
		"disk":      diskUsage,
		"uptime":    uptimeStr,
		"version":   "v0.7.0",
		"db_status": "healthy",
		"app_mem":   float64(m.Alloc) / 1024 / 1024, // MB
	})
}
