package models

import (
	"fmt"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/render/mod"
	"math/rand"
	"sync"
	"time"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/scene"
)

const (
	entranceSpotY = 205.00
	speed         = 10
)

type Car struct {
	hitbox floatgeom.Rect2
	entity *entities.Entity
	mu     sync.Mutex
}

func NewCar(ctx *scene.Context) *Car {
	hitbox := floatgeom.NewRect2(115, 490, 135, 510)

	mobile, _ := render.LoadSprite("assets/images/auto.png")

	newSwitch := render.NewSwitch("up", map[string]render.Modifiable{
		"up":    mobile,
		"down":  mobile.Copy().Modify(mod.FlipY),
		"left":  mobile.Copy().Modify(mod.Rotate(90)),
		"right": mobile.Copy().Modify(mod.Rotate(-90)),
	})

	entity := entities.New(ctx, entities.WithRect(hitbox), entities.WithRenderable(newSwitch), entities.WithDrawLayers([]int{1, 2}))

	return &Car{
		hitbox: hitbox,
		entity: entity,
	}
}

func (c *Car) moveToPosition(orientation string, point float64, controller *CarController) {
	if orientation == "left" {
		for c.X() > point {
			if !c.isCollision("left", controller.GetCars()) {
				_ = c.entity.Renderable.(*render.Switch).Set("left")
				c.ShiftX(-1)
				time.Sleep(speed * time.Millisecond)
			}
		}
	} else if orientation == "right" {
		for c.X() < point {
			if !c.isCollision("right", controller.GetCars()) {
				_ = c.entity.Renderable.(*render.Switch).Set("right")
				c.ShiftX(1)
				time.Sleep(speed * time.Millisecond)
			}
		}
	} else if orientation == "up" {
		for c.Y() > point {
			if !c.isCollision("up", controller.GetCars()) {
				_ = c.entity.Renderable.(*render.Switch).Set("up")
				c.ShiftY(-1)
				time.Sleep(speed * time.Millisecond)
			}
		}
	} else if orientation == "down" {
		for c.Y() < point {
			if !c.isCollision("down", controller.GetCars()) {
				_ = c.entity.Renderable.(*render.Switch).Set("down")
				c.ShiftY(1)
				time.Sleep(speed * time.Millisecond)
			}
		}
	}
}

func (c *Car) EnterQueue(controller *CarController) {
	c.moveToPosition("up", 275, controller)
}

func (c *Car) EnterDoor(controller *CarController) {
	c.moveToPosition("up", entranceSpotY, controller)
}

func (c *Car) LeaveDoor(controller *CarController) {
	c.moveToPosition("down", 235, controller)
}

func (c *Car) GetParked(space *ParkingSpace, controller *CarController) {
	for index := 0; index < len(*space.GetOrientationsForParking()); index++ {
		orientations := *space.GetOrientationsForParking()
		fmt.Println("Orientation: ", orientations[index].Orientation, "Point: ", orientations[index].Point)
		if orientations[index].Orientation == "left" {
			c.moveToPosition("left", orientations[index].Point, controller)
		} else if orientations[index].Orientation == "right" {
			c.moveToPosition("right", orientations[index].Point, controller)
		} else if orientations[index].Orientation == "up" {
			c.moveToPosition("up", orientations[index].Point, controller)
		} else if orientations[index].Orientation == "down" {
			c.moveToPosition("down", orientations[index].Point, controller)
		}
	}
}

func (c *Car) Leave(space *ParkingSpace, controller *CarController) {
	for index := 0; index < len(*space.GetOrientationsForLeaving()); index++ {
		orientations := *space.GetOrientationsForLeaving()
		if orientations[index].Orientation == "left" {
			c.moveToPosition("left", orientations[index].Point, controller)
		} else if orientations[index].Orientation == "right" {
			c.moveToPosition("right", orientations[index].Point, controller)
		} else if orientations[index].Orientation == "up" {
			c.moveToPosition("up", orientations[index].Point, controller)
		} else if orientations[index].Orientation == "down" {
			c.moveToPosition("down", orientations[index].Point, controller)
		}
	}
}

func (c *Car) LeaveSpace(controller *CarController) {
	spotX := c.X()
	c.moveToPosition("right", spotX+35, controller)
}

func (c *Car) Vanish(controller *CarController) {
	c.moveToPosition("down", 645, controller)
}

func (c *Car) ShiftY(dy float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.ShiftY(dy)
}

func (c *Car) ShiftX(dx float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.ShiftX(dx)
}

func (c *Car) X() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entity.X()
}

func (c *Car) Y() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entity.Y()
}

func (c *Car) Disappear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.Destroy()
}

func (c *Car) isCollision(orientation string, cars []*Car) bool {
	minDistance := 30.0
	for _, car := range cars {
		if orientation == "left" {
			if c.X() > car.X() && c.X()-car.X() < minDistance && c.Y() == car.Y() {
				return true
			}
		} else if orientation == "right" {
			if c.X() < car.X() && car.X()-c.X() < minDistance && c.Y() == car.Y() {
				return true
			}
		} else if orientation == "up" {
			if c.Y() > car.Y() && c.Y()-car.Y() < minDistance && c.X() == car.X() {
				return true
			}
		} else if orientation == "down" {
			if c.Y() < car.Y() && car.Y()-c.Y() < minDistance && c.X() == car.X() {
				return true
			}
		}
	}
	return false
}

func getRandomNumber(min, max int) float64 {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	return float64(generator.Intn(max-min+1) + min)
}
