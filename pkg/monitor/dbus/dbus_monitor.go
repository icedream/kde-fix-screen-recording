package dbus

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/godbus/dbus/v5"

	"github.com/icedream/kde-fix-screen-recording/pkg/monitor"
)

type dbusMonitor struct {
	dbus               *dbus.Conn
	lastInhibitMessage *dbus.Message
}

func NewDBusMonitor() (monitor.Monitor, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}

	return &dbusMonitor{
		dbus: conn,
	}, nil
}

func (m *dbusMonitor) convertEvent(dbusMessageC <-chan *dbus.Message, eventC chan<- monitor.Event) {
	log.Println("Going to convert events now")
	for dbusMessage := range dbusMessageC {
		// pathName := dbusMessage.Headers[dbus.FieldPath].String()
		interfaceName := dbusMessage.Headers[dbus.FieldInterface].Value()
		memberName := dbusMessage.Headers[dbus.FieldMember].Value()

		if interfaceName == "org.gnome.SessionManager" && memberName == "Inhibit" && len(dbusMessage.Body) >= 3 {
			spew.Dump(dbusMessage)
			// inhibitBinary, ok := dbusMessage.Body[0].(string)
			// if !ok {
			// 	continue
			// }
			inhibitReason, ok := dbusMessage.Body[2].(string)
			if !ok {
				continue
			}
			if inhibitReason == "Native desktop capture" {
				m.lastInhibitMessage = dbusMessage
				eventC <- monitor.Event{
					RecordingStatus: monitor.Recording,
				}
			}
			continue
		}

		if interfaceName == "org.gnome.SessionManager" && memberName == "Uninhibit" &&
			m.lastInhibitMessage != nil &&
			m.lastInhibitMessage.Headers[dbus.FieldSender].Value() == dbusMessage.Headers[dbus.FieldSender].Value() {
			spew.Dump(dbusMessage)

			m.lastInhibitMessage = nil
			eventC <- monitor.Event{
				RecordingStatus: monitor.NotRecording,
			}
		}
	}
}

func (m *dbusMonitor) Start() (<-chan monitor.Event, error) {
	conn := m.dbus

	for _, v := range []string{
		"method_call",
		// "method_return",
		// "error",
		// "signal",
	} {
		call := conn.BusObject().Call(
			"org.freedesktop.DBus.AddMatch", 0,
			"eavesdrop='true',type='"+v+"'")
		if call.Err != nil {
			return nil, call.Err
		}
	}

	c := make(chan *dbus.Message, 10)
	e := make(chan monitor.Event, 1)

	log.Println("Going to listen for DBus events now")
	conn.Eavesdrop(c)

	go m.convertEvent(c, e)

	return e, nil
}

func (m *dbusMonitor) Stop() error {
	return m.dbus.Close()
}
