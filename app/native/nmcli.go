package native

import "os/exec"

func BroadcastSSID(iface string, SSID string, password string) {
	cmd := exec.Command("nmcli", "con", "delete", "ap")
	cmd.Run()

	cmd = exec.Command("nmcli", "device", "wifi", "ap", "ifname", iface, "ssid", SSID, "password", password)
	cmd.Run()

	cmd = exec.Command("nmcli", "con", "modify", "ap", "ipv4.addresses", "172.16.0.0/24")
	cmd.Run()
}

func StopBroadcasting() {
	cmd := exec.Command("nmcli", "con", "down", "ap")
	cmd.Run()
}
