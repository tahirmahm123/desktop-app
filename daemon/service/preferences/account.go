
package preferences

// AccountStatus contains information about current account
type AccountStatus struct {
	Active      bool
	ActiveUntil int64
	IsFreeTrial bool
	CurrentPlan string
}
