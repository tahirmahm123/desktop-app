
package launcher

import (
	"fmt"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
)

// ServiceName -  name of the service
const _serviceName = "VPN"

var _evtlog *eventlog.Log
var _stopped chan struct{}

type vpnservice struct{}

// Prepare to start VPNe for Windows
func doPrepareToRun() error {

	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Error(fmt.Sprintf("failed to determine if we are running in an interactive session: %v", err))
	}

	log.Info("IsAnInteractiveSession: ", isIntSess)
	if isIntSess {
		log.Info("Starting as a console application (testing mode; InteractiveSession=true)")
		// It is interactive session. Continue as console application (testing mode)
		return nil
	}
	log.Info("Starting as a service (InteractiveSession=false)")

	// run a service handler (service is active until 'Execute(...)' method is running)
	go runWindowsService()

	// continue starting other stuff
	return nil
}

// inform OS-specific implementation about listener port
func doStartedOnPort(port int, secret uint64) {
}

// OS-specific service finalizer
func doStopped() {
	var stoppedChan = _stopped

	if stoppedChan != nil {
		// notify service handler that service is stopped
		stoppedChan <- struct{}{}
	}
}

// doCheckIsAdmin - check is application running with root privileges
func doCheckIsAdmin() bool {
	var sid *windows.SID

	// Although this looks scary, it is directly copied from the
	// official windows documentation. The Go API for this is a
	// direct wrap around the official C++ API.
	// See https://docs.microsoft.com/en-us/windows/desktop/api/securitybaseapi/nf-securitybaseapi-checktokenmembership
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		log.Error(fmt.Sprintf("SID Error: %s", err.Error()))
		return false
	}

	// This appears to cast a null pointer so I'm not sure why this
	// works, but this guy says it does and it Works for Me™:
	// https://github.com/golang/go/issues/28804#issuecomment-438838144
	token := windows.Token(0)

	// Also note that an admin is _not_ necessarily considered
	// elevated.
	// For elevation see https://github.com/mozey/run-as-admin
	//return token.IsElevated()

	member, err := token.IsMember(sid)
	if err != nil {
		log.Error(fmt.Sprintf("Token Membership Error: %s", err.Error()))
		return false
	}

	log.Info(fmt.Sprintf("IsAdmin=%v IsElevated=%v", member, token.IsElevated()))

	return member
}

func runWindowsService() {
	var err error
	_evtlog, err = eventlog.Open(_serviceName)
	if err != nil {
		log.Warning(fmt.Sprintf("Unable to initialize windows event log: %v", err))
		_evtlog = nil
	}
	defer func() {
		if _evtlog != nil {
			_evtlog.Close()
		}
	}()

	log.Info(fmt.Sprintf("starting %s service", _serviceName))
	if _evtlog != nil {
		_evtlog.Info(1, fmt.Sprintf("starting %s service", _serviceName))
	}

	// create stop-detection channel
	_stopped = make(chan struct{}, 1)

	// run windows-service-handler (func (m *vpnservice) Execute(...))
	err = svc.Run(_serviceName, &vpnservice{})
	if err != nil {
		log.Error(fmt.Sprintf("%s service failed: %v", _serviceName, err))
		if _evtlog != nil {
			_evtlog.Error(1, fmt.Sprintf("%s service failed: %v", _serviceName, err))
		}
		return
	}

	log.Info(fmt.Sprintf("%s service stopped", _serviceName))
	if _evtlog != nil {
		_evtlog.Info(1, fmt.Sprintf("%s service stopped", _serviceName))
	}
}

func (m *vpnservice) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	log.Info("Service handler started")
	defer func() {
		changes <- svc.Status{State: svc.StopPending}
		log.Info("Service handler: StopPending")

		// Stop the service (if not stopped yet)
		// This call should be performed at the end. Application will fully stop after that
		Stop()

		changes <- svc.Status{State: svc.Stopped}
		log.Info("Service handler: Stopped")
	}()

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		select {
		case <-_stopped:
			log.Info("Service stopped")
			break loop

		case c := <-r:
			switch c.Cmd {

			case svc.Interrogate:
				// SERVICE_CONTROL_INTERROGATE
				// 0x00000004
				// Notifies a service that it should report its current status information to the service control manager. The hService handle must have the SERVICE_INTERROGATE access right.
				// Note that this control is not generally useful as the SCM is aware of the current state of the service.

				log.Info("Service control request: ", "Interrogate", c.Cmd)

				changes <- c.CurrentStatus

				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				//time.Sleep(100 * time.Millisecond)
				//changes <- c.CurrentStatus

			case svc.Stop, svc.Shutdown:
				log.Info("Service control request: ", "Stop|Shutdown", c.Cmd)
				if _evtlog != nil {
					_evtlog.Info(1, fmt.Sprintf("Service control request: Stop|Shutdown %d", c.Cmd))
				}
				break loop

			default:
				log.Warning("Unexpected service control request: ", c.Cmd)
				if _evtlog != nil {
					_evtlog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
				}
			}
		}
	}

	return
}

func isNeedToSavePortInFile() bool {
	return true
}
