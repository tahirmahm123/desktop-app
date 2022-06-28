//go:build darwin && !nowifi
// +build darwin,!nowifi
package wifiNotifier

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework SystemConfiguration -framework CoreWLAN -framework Foundation

#include <stdlib.h>
#import <Foundation/Foundation.h>
#import <CoreWLAN/CoreWLAN.h>
#import <SystemConfiguration/SystemConfiguration.h>

static inline  char* nsstring2cstring(NSString *s){
    if (s == NULL) { return NULL; }

	char *cstr = strdup([s UTF8String]);
    return cstr;
}

#define NOT_CONNECTED @""

static inline CWInterface * getCWInterface() {
	CWWiFiClient *swc = [CWWiFiClient sharedWiFiClient];
	if (swc == nil) return nil;
	return [swc interface];
}

static inline void wifi_network_changed(SCDynamicStoreRef store, CFArrayRef changedKeys, void *ctx)
{
	CWInterface * WiFiInterface = getCWInterface();
	if (WiFiInterface == nil) return;

	NSString *currentSSID = [WiFiInterface ssid] ? [WiFiInterface ssid] : NOT_CONNECTED;
	extern void __onWifiChanged(char *);
	__onWifiChanged(nsstring2cstring(currentSSID));
}

static inline char * getCurrentSSID(void) {
	CWInterface * WiFiInterface = getCWInterface();
	if (WiFiInterface == nil) return nsstring2cstring(NOT_CONNECTED);

	NSString *ssid = [WiFiInterface ssid] ? [WiFiInterface ssid] : NOT_CONNECTED;
	return nsstring2cstring(ssid);
}

static inline int getCurrentNetworkSecurity() {
	CWInterface * WiFiInterface = getCWInterface();
	if (WiFiInterface == nil) return 0xFFFFFFFF;

	return [WiFiInterface security];
}

static inline char* getAvailableSSIDs(void) {
	CWInterface * WiFiInterface = getCWInterface();
	if (WiFiInterface == nil) return nil;

	NSError *err = nil;
	NSSet *scanset = [WiFiInterface scanForNetworksWithSSID:Nil error:&err];
    if (err!=nil || scanset == nil || scanset.count == 0) return nil;

	NSString *retString = nil;
	int i=0;
	for (CWNetwork * nw in scanset)
    {
		if (nw == nil || [nw ssid] == nil) continue;
		NSString * ssid = [[[nw ssid] componentsSeparatedByCharactersInSet:[NSCharacterSet newlineCharacterSet]] componentsJoinedByString:@" "];
		if (i++ == 0)
			retString = ssid;
		else
			retString = [NSString stringWithFormat:@"%@\n%@", retString , ssid];
	}

	return nsstring2cstring(retString);
}

static inline void setWifiNotifier(void) {
    CWInterface * WiFiInterface = getCWInterface();
	if (WiFiInterface == nil) return;

	NSArray* arr = [CWWiFiClient interfaceNames];
	NSSet *wifiInterfaces = [NSSet setWithArray:arr];

    NSMutableArray *scKeys = [[NSMutableArray alloc] init];
    [wifiInterfaces enumerateObjectsUsingBlock:^(NSString *ifName, BOOL *stop)
     {
         [scKeys addObject: [NSString stringWithFormat:@"State:/Network/Interface/%@/AirPort", ifName]];
     }];

    SCDynamicStoreContext ctx = { 0, NULL, NULL, NULL, NULL };
    SCDynamicStoreRef store = SCDynamicStoreCreate(kCFAllocatorDefault, CFSTR("VPN"), wifi_network_changed, &ctx);

    SCDynamicStoreSetNotificationKeys(store, (__bridge CFArrayRef)scKeys, NULL);

    CFRunLoopSourceRef src = SCDynamicStoreCreateRunLoopSource(kCFAllocatorDefault, store, 0);
	CFRunLoopAddSource([[NSRunLoop currentRunLoop] getCFRunLoop], src, kCFRunLoopCommonModes);

	CFRunLoopRun();
}
*/
import "C"
import (
	"strings"
	"time"
	"unsafe"

	"github.com/tahirmahm123/vpn-desktop-app/daemon/logger"
)

var log *logger.Logger

func init() {
	log = logger.NewLogger("wifi")
}

var internalOnWifiChangedCb func(string)

//export __onWifiChanged
func __onWifiChanged(ssid *C.char) {
	goSsid := C.GoString(ssid)
	C.free(unsafe.Pointer(ssid))

	if internalOnWifiChangedCb != nil {
		internalOnWifiChangedCb(goSsid)
	}
}

// GetAvailableSSIDs returns the list of the names of available Wi-Fi networks
func GetAvailableSSIDs() []string {
	ssidList := C.getAvailableSSIDs()
	goSsidList := C.GoString(ssidList)
	C.free(unsafe.Pointer(ssidList))
	return strings.Split(goSsidList, "\n")
}

// GetCurrentSSID returns current WiFi SSID
func GetCurrentSSID() string {
	ssid := C.getCurrentSSID()
	goSsid := C.GoString(ssid)
	C.free(unsafe.Pointer(ssid))
	return goSsid
}

// GetCurrentNetworkIsInsecure returns current security mode
func GetCurrentNetworkIsInsecure() bool {
	const (
		CWSecurityNone               = 0
		CWSecurityWEP                = 1
		CWSecurityWPAPersonal        = 2
		CWSecurityWPAPersonalMixed   = 3
		CWSecurityWPA2Personal       = 4
		CWSecurityPersonal           = 5
		CWSecurityDynamicWEP         = 6
		CWSecurityWPAEnterprise      = 7
		CWSecurityWPAEnterpriseMixed = 8
		CWSecurityWPA2Enterprise     = 9
		CWSecurityEnterprise         = 10
		CWSecurityWPA3Personal       = 11
		CWSecurityWPA3Enterprise     = 12
		CWSecurityWPA3Transition     = 13
	)

	security := C.getCurrentNetworkSecurity()
	switch security {
	case CWSecurityNone,
		CWSecurityWEP,
		CWSecurityPersonal,
		CWSecurityDynamicWEP:
		return true
	}
	return false
}

// SetWifiNotifier initializes a handler method 'OnWifiChanged'
func SetWifiNotifier(cb func(string)) error {
	internalOnWifiChangedCb = cb
	go func() {
		log.Info("WiFi notifier enter")
		defer log.Error("WiFi notifier exit")

		for {
			// Detection WiFi status change in infinite loop.
			// C.setWifiNotifier() should never return.
			//
			// BUT! It can return in some corner cases (e.g. we call it on system boot when WiFi interface still not initialized)
			// In this case - we waiting some delay and trying to call this function again
			C.setWifiNotifier()
			log.Info("Unexpected WiFi notifier exit")
			time.Sleep(time.Second)
			log.Info("WiFi notifier enter. Retry...")
		}
	}()
	return nil
}
