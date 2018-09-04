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
			log.Printf("Height: %d Battery: %d Remaining Flytime: %d\n", fd.Height, fd.BatteryPercentage, fd.DroneFlyTimeLeft)
		})

		log.Print("droneTakeOff")
		err := drone.TakeOff()
		if err != nil {
			log.Print(errors.Wrap(err, "droneTakeOff"))
		}

		gobot.After(2*time.Second, func() {
			err := drone.Forward(2)
			if err != nil {
				log.Print(errors.Wrap(err, "drone.Forward"))
			}
		})

		gobot.After(time.Second, func() {
			err := drone.Backward(3)
			if err != nil {
				log.Print(errors.Wrap(err, "drone.Backward"))
			}
		})

		gobot.After(5*time.Second, func() {
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
