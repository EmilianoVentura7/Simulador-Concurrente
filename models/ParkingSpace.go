package models

import (
	"github.com/oakmound/oak/v4/alg/floatgeom"
)

type ParkingSpace struct {
	hitbox                 *floatgeom.Rect2
	orientationsForParking *[]ParkingSpotWay
	orientationsForLeaving *[]ParkingSpotWay
	number                 int
	isAvailable            bool
}

func NewParkingSpace(x, y, x2, y2 float64, row, number int) *ParkingSpace {
	orientationsForParking := getOrientationsForParking(x, y, row)
	orientationsForLeaving := getOrientationsForLeaving()
	hitbox := floatgeom.NewRect2(x, y, x2, y2)

	return &ParkingSpace{
		hitbox:                 &hitbox,
		orientationsForParking: orientationsForParking,
		orientationsForLeaving: orientationsForLeaving,
		number:                 number,
		isAvailable:            true,
	}
}

func getOrientationsForParking(x, y float64, row int) *[]ParkingSpotWay {
	var orientations []ParkingSpotWay

	if row == 1 {
		orientations = append(orientations, *newParkingSpotWay("right", 115))
	} else if row == 2 {
		orientations = append(orientations, *newParkingSpotWay("right", 205))
	} else if row == 3 {
		orientations = append(orientations, *newParkingSpotWay("right", 295))
	} else if row == 4 {
		orientations = append(orientations, *newParkingSpotWay("right", 385))
	} else if row == 5 {
		orientations = append(orientations, *newParkingSpotWay("right", 475))
	}

	orientations = append(orientations, *newParkingSpotWay("up", y+5))
	orientations = append(orientations, *newParkingSpotWay("right", x+5))

	return &orientations
}

func getOrientationsForLeaving() *[]ParkingSpotWay {
	var orientations []ParkingSpotWay

	orientations = append(orientations, *newParkingSpotWay("up", 10))
	orientations = append(orientations, *newParkingSpotWay("left", 85))
	orientations = append(orientations, *newParkingSpotWay("down", 200))

	return &orientations
}

func (p *ParkingSpace) GetHitbox() *floatgeom.Rect2 {
	return p.hitbox
}

func (p *ParkingSpace) GetNumber() int {
	return p.number
}

func (p *ParkingSpace) GetOrientationsForParking() *[]ParkingSpotWay {
	return p.orientationsForParking
}

func (p *ParkingSpace) GetOrientationsForLeaving() *[]ParkingSpotWay {
	return p.orientationsForLeaving
}

func (p *ParkingSpace) GetIsAvailable() bool {

	return p.isAvailable
}

func (p *ParkingSpace) SetIsAvailable(isAvailable bool) {
	p.isAvailable = isAvailable
}
