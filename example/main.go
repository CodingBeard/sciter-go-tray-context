package main

import (
	"fmt"
	"github.com/codingbeard/tray"
	"time"
	"github.com/codingbeard/sciter-go-tray-context"
)

func main() {
	menu := context.Menu{
		Items: []context.MenuItem{
			{
				Text: "Settings",
				ClickCallback: func() {
					fmt.Println("Settings clicked")
				},
			},
			{
				Text: "About",
				ClickCallback: func() {
					fmt.Println("About clicked")
				},
			},
			{
				Text: "Check For Updates..",
				ClickCallback: func() {
					fmt.Println("Updates clicked")
				},
			},
			{
				Text: "Exit",
				ClickCallback: func() {
					fmt.Println("Exit clicked")
				},
			},
		},
	}

	trayIcon := tray.ClickableIcon{
		IconData: iconData,
		ClickHandler: func(x, y int, rightClick bool) {
			fmt.Println("x", x, "y", y, "isRightClick", rightClick)
			if rightClick {
				menu.DisplayContextMenu(x, y, 100)
			}
		},
	}
	trayIcon.Initialise()

	fmt.Println("Clickable Icon initialised")

	for {
		time.Sleep(time.Second)
	}
}
