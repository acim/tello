package main

import (
	"log"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type flightData struct {
	batteryPercentage int8
	Height            int16
}

var flightDataState flightData

func main() {
	drone := tello.NewDriver("8888")

	flightDataState = flightData{}

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
			logFlightData(data.(*tello.FlightData))

		})

		err := drone.TakeOff()
		if err != nil {
			log.Fatal(errors.Wrap(err, "droneTakeOff"))
		}
		time.Sleep(5 * time.Second)

		err = drone.Forward(10)
		if err != nil {
			log.Print(errors.Wrap(err, "drone.Forward"))
		}
		time.Sleep(2 * time.Second)

		err = drone.Backward(10)
		if err != nil {
			log.Print(errors.Wrap(err, "drone.Backward"))
		}
		time.Sleep(2 * time.Second)

		err = drone.Left(10)
		if err != nil {
			log.Print(errors.Wrap(err, "drone.Left"))
		}
		time.Sleep(2 * time.Second)

		err = drone.Right(10)
		if err != nil {
			log.Print(errors.Wrap(err, "drone.Right"))
		}
		time.Sleep(2 * time.Second)

		err = drone.Up(10)
		if err != nil {
			log.Print(errors.Wrap(err, "drone.Up"))
		}
		time.Sleep(2 * time.Second)

		err = drone.Down(10)
		if err != nil {
			log.Print(errors.Wrap(err, "drone.Down"))
		}
		time.Sleep(2 * time.Second)

		err = drone.FrontFlip()
		if err != nil {
			log.Print(errors.Wrap(err, "drone.FrontFlip"))
		}
		time.Sleep(2 * time.Second)

		err = drone.BackFlip()
		if err != nil {
			log.Print(errors.Wrap(err, "drone.BackFlip"))
		}
		time.Sleep(2 * time.Second)

		err = drone.LeftFlip()
		if err != nil {
			log.Print(errors.Wrap(err, "drone.LeftFlip"))
		}
		time.Sleep(2 * time.Second)

		err = drone.RightFlip()
		if err != nil {
			log.Print(errors.Wrap(err, "drone.RightFlip"))
		}
		time.Sleep(2 * time.Second)

		err = drone.Land()
		if err != nil {
			log.Print(errors.Wrap(err, "mplayer.Start"))
		}
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

func logFlightData(fd *tello.FlightData) {
	newFD := flightData{
		batteryPercentage: fd.BatteryPercentage,
		Height:            fd.Height,
	}
	if newFD != flightDataState {
		flightDataState = newFD
		log.Printf("Height: %d Battery: %d\n", fd.Height, fd.BatteryPercentage)
	}
}
