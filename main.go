package main

import (
	"log"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver("8888")

	work := func() {
		mplayer := exec.Command("mplayer", "-fps", "25", "-")
		mplayerIn, _ := mplayer.StdinPipe()
		if err := mplayer.Start(); err != nil {
			log.Print(errors.Wrap(err, "mplayer.Start"))
		}

		drone.On(tello.ConnectedEvent, func(data interface{}) {
			log.Print("Connected")
			err := drone.StartVideo()
			if err != nil {
				log.Print(errors.Wrap(err, "drone.StartVideo"))
			}
			err = drone.SetVideoEncoderRate(4)
			if err != nil {
				log.Print(errors.Wrap(err, "drone.SetVideoEncoderRate"))
			}
			gobot.Every(100*time.Millisecond, func() {
				err := drone.StartVideo()
				if err != nil {
					log.Print(errors.Wrap(err, "drone.StartVideo"))
				}
			})
		})

		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := mplayerIn.Write(pkt); err != nil {
				log.Print(errors.Wrap(err, "mplayerIn.Write"))
			}
		})

		drone.On(tello.FlightDataEvent, func(data interface{}) {
			fd := data.(*tello.FlightData)
			log.Printf("Height: %d Battery: %d\n", fd.Height, fd.BatteryPercentage)
		})

		drone.On(tello.TakeoffEvent, func(data interface{}) {
			log.Printf("Takeoff: %#v", data)
		})

		drone.On(tello.LandingEvent, func(data interface{}) {
			log.Printf("Landing: %#v", data)
		})

		drone.On(tello.FlipEvent, func(data interface{}) {
			log.Printf("Flip: %#v", data)
		})

		drone.On(tello.BounceEvent, func(data interface{}) {
			log.Printf("Bounce: %#v", data)
		})

		log.Print("droneTakeOff")
		err := drone.TakeOff()
		if err != nil {
			log.Print(errors.Wrap(err, "droneTakeOff"))
		}

		gobot.After(5*time.Second, func() {
			err := drone.Forward(10)
			if err != nil {
				log.Print(errors.Wrap(err, "drone.Forward"))
			}
		})

		gobot.After(7*time.Second, func() {
			err := drone.Backward(10)
			if err != nil {
				log.Print(errors.Wrap(err, "drone.Backward"))
			}
		})

		gobot.After(9*time.Second, func() {
			err := drone.Left(10)
			if err != nil {
				log.Print(errors.Wrap(err, "drone.Left"))
			}
		})

		gobot.After(12*time.Second, func() {
			err := drone.Right(10)
			if err != nil {
				log.Print(errors.Wrap(err, "drone.Right"))
			}
		})

		gobot.After(15*time.Second, func() {
			err := drone.FrontFlip()
			if err != nil {
				log.Print(errors.Wrap(err, "drone.FrontFlip"))
			}
		})

		gobot.After(16*time.Second, func() {
			err := drone.BackFlip()
			if err != nil {
				log.Print(errors.Wrap(err, "drone.BackFlip"))
			}
		})

		gobot.After(17*time.Second, func() {
			err := drone.Land()
			if err != nil {
				log.Print(errors.Wrap(err, "mplayer.Start"))
			}
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	err := robot.Start()
	if err != nil {
		log.Print(errors.Wrap(err, "robot.Start"))
	}
}
