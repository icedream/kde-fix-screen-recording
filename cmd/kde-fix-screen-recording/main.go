package main

import (
	"log"

	"github.com/icedream/kde-fix-screen-recording/pkg/kde"
	"github.com/icedream/kde-fix-screen-recording/pkg/monitor"
	"github.com/icedream/kde-fix-screen-recording/pkg/monitor/dbus"
)

func main() {
	m, err := dbus.NewDBusMonitor()
	if err != nil {
		log.Fatal(err)
	}
	e, err := m.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer m.Stop()

	defer kde.ResumeCompositor()

	log.Println("Now monitoring")
	for ev := range e {
		log.Println(ev)
		switch ev.RecordingStatus {
		case monitor.Recording:
			if err := kde.SuspendCompositor(); err != nil {
				log.Println("Failed to suspend compositor:", err)
			}
		case monitor.NotRecording:
			if err := kde.ResumeCompositor(); err != nil {
				log.Println("Failed to resume compositor:", err)
			}
		}
	}
}
