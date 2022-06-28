
package obfsproxy_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/logger"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/obfsproxy"
	"github.com/tahirmahm123/vpn-desktop-app/daemon/service/platform"
)

func TestStart(t *testing.T) {
	platform.Init()
	logger.Enable(true)
	obfsp := obfsproxy.CreateObfsproxy(platform.ObfsproxyStartScript())

	port, err := obfsp.Start()
	if err != nil {
		fmt.Println("ERROR:", err)
	} else {
		fmt.Println("Started on:", port)
	}

	go func() {
		time.Sleep(time.Second * 5)
		obfsp.Stop()
	}()

	if err := obfsp.Wait(); err != nil {
		fmt.Println("STOP ERROR:", err)
	}
	fmt.Println("STOPED")
}
