package kde

import "github.com/godbus/dbus/v5"

func IsCompositorRunning() (compositorRunning bool, err error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return compositorRunning, err
	}
	defer conn.Close()

	obj := conn.Object("org.kde.KWin", "/Compositor")
	err = obj.Call(
		"org.freedesktop.DBus.Properties.Get",
		0,
		"org.kde.kwin.Compositing",
		"compositingActive").
		Store(&compositorRunning)
	return compositorRunning, err
}

func SuspendCompositor() error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	obj := conn.Object("org.kde.KWin", "/Compositor")
	call := obj.Call("org.kde.kwin.Compositing.suspend", 0)
	return call.Err
}

func ResumeCompositor() error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	obj := conn.Object("org.kde.KWin", "/Compositor")
	call := obj.Call("org.kde.kwin.Compositing.resume", 0)
	return call.Err
}
