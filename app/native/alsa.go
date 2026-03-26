package native

import (
	"fmt"
	"os/exec"
)

func SetVolume(volume int32) {
	cmd := exec.Command("amixer", "sset", "Master", fmt.Sprintf("%d%", volume))
	cmd.Run()
}
