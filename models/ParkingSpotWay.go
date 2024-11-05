package models

type ParkingSpotWay struct {
	Orientation string
	Point       float64
}

func newParkingSpotWay(orientation string, point float64) *ParkingSpotWay {
	return &ParkingSpotWay{
		Orientation: orientation,
		Point:       point,
	}
}
