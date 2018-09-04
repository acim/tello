package main

import (
	"fmt"
	"os/exec"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver("8888")

	work := func() {
		mplayer := exec.Command("mplayer", "-fps", "25", "-")
		mplayerIn, _ := mplayer.StdinPipe()
		if err := mplayer.Start(); err != nil {
			fmt.Println(err)
			fmt.Println(err)
		}

		drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")
			err := drone.StartVideo()
			if err != nil {
				fmt.Println(err)
			}
			err = drone.SetVideoEncoderRate(4)
			if err != nil {
				fmt.Println(err)
			}
			gobot.Every(100*time.Millisecond, func() {
				err := drone.StartVideo()
				if err != nil {
					fmt.Println(err)
				}
			})
		})

		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := mplayerIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		})

		// err := drone.TakeOff()
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// gobot.After(100*time.Millisecond, func() {
		// 	err := drone.Forward(1)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// })

		// gobot.After(100*time.Millisecond, func() {
		// 	err := drone.Backward(1)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// })

		// gobot.After(1*time.Second, func() {
		// 	err := drone.Land()
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// })
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
}
