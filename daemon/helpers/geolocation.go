
package helpers

import "math"

// GetDistanceFromLatLonInKm calc distance betwee two points
func GetDistanceFromLatLonInKm(lat1, lon1, lat2, lon2 float64) float64 {
	var R = float64(6371)           // Radius of the earth in km
	var dLat = deg2rad(lat2 - lat1) // deg2rad below
	var dLon = deg2rad(lon2 - lon1)
	var a = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(deg2rad(lat1))*
			math.Cos(deg2rad(lat2))*
			math.Sin(dLon/2)*
			math.Sin(dLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	var d = R * c // Distance in km
	return d
}

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
