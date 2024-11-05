package scenes

import (
	"SimuladorConcurrente/models"
	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/scene"
	"image/color"
	"math/rand"
	"sync"
	"time"
)

var (
	spaces = []*models.ParkingSpace{
		models.NewParkingSpace(140, 35, 170, 65, 1, 1),
		models.NewParkingSpace(230, 35, 260, 65, 2, 2),
		models.NewParkingSpace(320, 35, 350, 65, 3, 3),
		models.NewParkingSpace(410, 35, 440, 65, 4, 4),
		models.NewParkingSpace(500, 35, 530, 65, 5, 5),
		models.NewParkingSpace(140, 80, 170, 110, 1, 1),
		models.NewParkingSpace(230, 80, 260, 110, 2, 2),
		models.NewParkingSpace(320, 80, 350, 110, 3, 3),
		models.NewParkingSpace(410, 80, 440, 110, 4, 4),
		models.NewParkingSpace(500, 80, 530, 110, 5, 5),
		models.NewParkingSpace(140, 125, 170, 155, 1, 1),
		models.NewParkingSpace(230, 125, 260, 155, 2, 2),
		models.NewParkingSpace(320, 125, 350, 155, 3, 3),
		models.NewParkingSpace(410, 125, 440, 155, 4, 4),
		models.NewParkingSpace(500, 125, 530, 155, 5, 5),
		models.NewParkingSpace(140, 170, 170, 200, 1, 1),
		models.NewParkingSpace(230, 170, 260, 200, 2, 2),
		models.NewParkingSpace(320, 170, 350, 200, 3, 3),
		models.NewParkingSpace(410, 170, 440, 200, 4, 4),
		models.NewParkingSpace(500, 170, 530, 200, 5, 5),
	}
	parking       = models.NewParking(spaces)
	doorMutex     sync.Mutex
	carController = models.NewCarController()
)

type ParkingScene struct {
}

func NewParkingScene() *ParkingScene {
	return &ParkingScene{}
}

func (ps *ParkingScene) Start() {
	isFirstTime := true

	_ = oak.AddScene("parkingScene", scene.Scene{
		Start: func(ctx *scene.Context) {
			setUpScene(ctx)

			event.GlobalBind(ctx, event.Enter, func(enterPayload event.EnterPayload) event.Response {
				if !isFirstTime {
					return 0
				}

				isFirstTime = false

				for {
					go carCycle(ctx)

					time.Sleep(time.Millisecond * time.Duration(getRandomNumber(1000, 2000)))
				}

				return 0
			})
		},
	})
}

func setUpScene(ctx *scene.Context) {
	parkingHitbox := floatgeom.NewRect2(0, 0, 640, 480)
	entities.New(ctx, entities.WithRect(parkingHitbox), entities.WithColor(color.RGBA{86, 101, 115, 255}), entities.WithDrawLayers([]int{0}))

	entities.New(ctx, entities.WithRect(floatgeom.NewRect2(80, 5, 560, 10)), entities.WithColor(color.RGBA{255, 255, 255, 255}), entities.WithDrawLayers([]int{0}))
	entities.New(ctx, entities.WithRect(floatgeom.NewRect2(140, 230, 560, 235)), entities.WithColor(color.RGBA{255, 255, 255, 255}), entities.WithDrawLayers([]int{0}))
	entities.New(ctx, entities.WithRect(floatgeom.NewRect2(75, 5, 80, 235)), entities.WithColor(color.RGBA{255, 255, 255, 255}), entities.WithDrawLayers([]int{0}))
	entities.New(ctx, entities.WithRect(floatgeom.NewRect2(560, 5, 565, 235)), entities.WithColor(color.RGBA{255, 255, 255, 255}), entities.WithDrawLayers([]int{0}))

	for _, space := range spaces {
		entities.New(ctx, entities.WithRect(floatgeom.NewRect2(space.GetHitbox().Min.X(), space.GetHitbox().Min.Y(), space.GetHitbox().Max.X(), space.GetHitbox().Min.Y()+2.5)), entities.WithColor(color.RGBA{255, 255, 255, 255}))
		entities.New(ctx, entities.WithRect(floatgeom.NewRect2(space.GetHitbox().Min.X(), space.GetHitbox().Max.Y()-2.5, space.GetHitbox().Max.X(), space.GetHitbox().Max.Y())), entities.WithColor(color.RGBA{255, 255, 255, 255}))
	}
}

func carCycle(ctx *scene.Context) {
	car := models.NewCar(ctx)

	carController.AddCar(car)

	car.EnterQueue(carController)

	spaceAvailable := parking.GetParkingSpotAvailable()

	doorMutex.Lock()

	car.EnterDoor(carController)

	doorMutex.Unlock()

	car.GetParked(spaceAvailable, carController)

	time.Sleep(time.Millisecond * time.Duration(getRandomNumber(40000, 50000)))

	car.LeaveSpace(carController)

	parking.ReleaseParkingSpace(spaceAvailable)

	car.Leave(spaceAvailable, carController)

	doorMutex.Lock()

	car.LeaveDoor(carController)

	doorMutex.Unlock()

	car.Vanish(carController)

	car.Disappear()

	carController.DeleteCar(car)
}

func getRandomNumber(min, max int) float64 {
	source := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(source)
	return float64(generator.Intn(max-min+1) + min)
}
