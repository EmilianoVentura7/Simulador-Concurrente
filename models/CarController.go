package models

import "sync"

type CarController struct {
	Cars  []*Car
	Mutex sync.Mutex
}

func NewCarController() *CarController {
	return &CarController{
		Cars: make([]*Car, 0),
	}
}

func (cm *CarController) AddCar(car *Car) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	cm.Cars = append(cm.Cars, car)
}

func (cm *CarController) DeleteCar(car *Car) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	for i, c := range cm.Cars {
		if c == car {
			cm.Cars = append(cm.Cars[:i], cm.Cars[i+1:]...)
			break
		}
	}
}

func (cm *CarController) GetCars() []*Car {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	return cm.Cars
}
