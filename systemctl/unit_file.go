package systemctl

import "os"

func servicePath(name string) string {
	return `/etc/systemd/system/` + name + `.service`
}
func binaryPath(name string) string {
	return `/home/` + name + `/main`
}
func serviceFileContents(name string) (string, error) {
	// make sure that the binary exists
	_, err := os.Stat(binaryPath(name))
	if err != nil {
		return "", err
	}
	return `[Unit]
Description=` + name + `
After=network.target

[Service]
Type=simple
ExecStart=` + binaryPath(name) + `
Restart=on-failure
User=` + name + `
Group=` + name + `
StandardOutput=journal
StandardError=journal
AmbientCapabilities=CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target`, nil
}
