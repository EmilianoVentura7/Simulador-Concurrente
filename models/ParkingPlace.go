package models

import (
	"sync"
)

type ParkingPlace struct {
	spaces        []*ParkingSpace
	mu            sync.Mutex
	availableCond *sync.Cond
}

func NewParking(spots []*ParkingSpace) *ParkingPlace {
	p := &ParkingPlace{
		spaces: spots,
	}
	p.availableCond = sync.NewCond(&p.mu)

	return p
}

func (p *ParkingPlace) GetSpots() []*ParkingSpace {
	return p.spaces
}

func (p *ParkingPlace) GetParkingSpotAvailable() *ParkingSpace {
	p.mu.Lock()
	defer p.mu.Unlock()

	for {
		for _, spot := range p.spaces {
			if spot.GetIsAvailable() {
				spot.SetIsAvailable(false)
				return spot
			}
		}
		p.availableCond.Wait()
	}
}

func (p *ParkingPlace) ReleaseParkingSpace(spot *ParkingSpace) {
	p.mu.Lock()
	defer p.mu.Unlock()

	spot.SetIsAvailable(true)
	p.availableCond.Signal()
}
