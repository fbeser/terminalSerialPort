package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
	"go.bug.st/serial.v1"
)

var senderPreviousView string = "Sender 1"

var previousWindowsWidth int
var previousWindowsHeight int
var startSize bool = true

var resizeLock bool = false

var optionsViewAccessRow []int = []int{
	1,
	4,
	7,
	10,
	13,
	15,
	17,
	21,
	23,
	25,
}
var optionsViewAccessRowSelected int
var optionsViewAccessColSelected int
var optionsViewAccessColLimit int = 2

type rawDataStruct struct {
	data          string
	typ           string
	color         string
	time          time.Time
	consolMessage bool
}

var rawData []rawDataStruct
var rawDataHex []rawDataStruct
var rawDataMutex sync.Mutex

var rawDataProccesed int

var bufferSizeItems []string = []string{
	"100",
	"200",
	"300",
	"400",
	"500",
	"500+",
}
var bufferSizeSelected int

var cursorLock bool = true
var cursorKey []bool = []bool{false, false}

var senderBackSpaceMutex sync.Mutex

var hexEnabled bool
var timeEnabled bool

var sender1HexEnabled bool
var sender2HexEnabled bool

var sender1HexChanged bool
var sender2HexChanged bool

var sender1LineFeed bool
var sender2LineFeed bool

var sender1CarriageReturn bool
var sender2CarriageReturn bool

var hexFilter *regexp.Regexp

const (
	whiteText  = "\x1b[0;29m"
	redText    = "\x1b[0;31m"
	greenText  = "\x1b[0;32m"
	yellowText = "\x1b[0;33m"
	blueText   = "\x1b[0;34m"
	purpleText = "\x1b[0;35m"

	greenRedText = "\x1b[31;42m"
)

func main() {

	// rawDataAdd(rawDataStruct { data : "Fatih", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "Beser", color : greenText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "Cemal", color : yellowText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "mehmet\r\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "Fatih", color : greenText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "Beser", color : yellowText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "1IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "2IndexAny returns the index of the first instance of any Unicode code point from chars in s, or -1 if no Unicode code point from chars is present in s.", color : greenText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "3IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "4IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "5IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "6IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "7IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "8IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "9IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "10IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "11IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "12IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "13IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "14IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "15IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })
	// rawDataAdd(rawDataStruct { data : "16IndexAny returns the index of the first\n instance of any Unicode code point from chars in s,\r\n or -1 if no Unicode code point from chars is present in s.\n", color : redText, time : time.Now() })

	scanPorts()

	var err error
	hexFilter, err = regexp.Compile("[^a-fA-F0-9 ]+")
	if err != nil {
		errorCreate(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.InputEsc = true

	g.SetManagerFunc(layout)

	if err = initKeybindings(g); err != nil {
		log.Panicln(err)
	}

	if err = g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func scanPorts() {
	var err error
	serialPortItems, err = serial.GetPortsList()
	if err != nil {
		errorCreate(err)
	}

	// serialPortItems = append(serialPortItems, "/dev/ttySC0")
	// serialPortItems = append(serialPortItems, "/dev/ttySC1")
}

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, cursorLeft); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, cursorRight); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, selected); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, goBack); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, clearData); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, rawDataView); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlO, gocui.ModNone, optionsView); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlN, gocui.ModNone, senderView); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, openPortSc); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlB, gocui.ModNone, bufferSizeView); err != nil {
		return err
	}

	return nil
}

var val int

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if maxX < 30 {
		maxX = 30
	}
	if maxY < 30 {
		maxY = 30
	}

	if startSize {
		previousWindowsWidth = maxX
		previousWindowsHeight = maxY
		startSize = false
	}

	if previousWindowsWidth > maxX {
		v, err := g.View("Raw Data")
		if err != nil {
			v.Clear()
		}
		previousWindowsWidth -= ((previousWindowsWidth - maxX) / 10) + 1
		maxX = previousWindowsWidth
	}
	if previousWindowsHeight > maxY {
		v, err := g.View("Raw Data")
		if err != nil {
			v.Clear()
		}
		previousWindowsHeight -= ((previousWindowsHeight - maxY) / 10) + 1
		maxY = previousWindowsHeight
	}

	if previousWindowsWidth < maxX {
		v, err := g.View("Raw Data")
		if err != nil {
			v.Clear()
		}
		previousWindowsWidth += ((maxX - previousWindowsWidth) / 10) + 1
		maxX = previousWindowsWidth
	}
	if previousWindowsHeight < maxY {
		v, err := g.View("Raw Data")
		if err != nil {
			v.Clear()
		}
		previousWindowsHeight += ((maxY - previousWindowsHeight) / 10) + 1
		maxY = previousWindowsHeight
	}

	if v, err := g.SetView("Raw Data", 0, 0, int(float32(maxX)*0.8), int(float32(maxY)*0.8)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = cursorLock
		v.Title = "Raw Data"
		fmt.Fprint(v)
	}

	if v, err := g.SetView("Options", int(float32(maxX)*0.8)+1, 0, maxX-1, int(float32(maxY)*0.8)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//v.Highlight = false
		//v.SelBgColor = gocui.ColorGreen
		//v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 1)
		v.Title = "Options"
		fmt.Fprint(v)
		if _, err := g.SetCurrentView("Options"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("Sender 1", 0, int(float32(maxY)*0.8)+1, int(float32(maxX)*0.8), int(float32(maxY)*0.8)+((maxY-1-int(float32(maxY)*0.8))/2)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Title = "Sender 1 [ ]"
		fmt.Fprint(v)
	}

	if v, err := g.SetView("Sender 2", 0, int(float32(maxY)*0.8)+1+((maxY-1-int(float32(maxY)*0.8))/2), int(float32(maxX)*0.8), maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Title = "Sender 2 [ ]"
		fmt.Fprint(v)
	}

	if v, err := g.SetView("Shortcut", int(float32(maxX)*0.8)+1, int(float32(maxY)*0.8)+1, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Shortcut"
		fmt.Fprintln(v, "Raw Data       = Ctrl + R")
		fmt.Fprintln(v, "Options View   = Ctrl + O")
		fmt.Fprintln(v, "Sender 1 & 2   = Ctrl + N")
		fmt.Fprintln(v, "Next View      = Ctrl + Space")
		fmt.Fprintln(v, "Port Switch    = Ctrl + P")
		fmt.Fprintln(v, "Buffer Size    = Ctrl + B")
		fmt.Fprintln(v, "Clear          = Ctrl + C")
		fmt.Fprintln(v, "Quit           = Ctrl + Q")
	}

	time.Sleep(time.Second / 30) // kare sayısı
	update(g)

	return nil
}

var test string

func update(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {

		currentView := g.CurrentView()
		switch currentView.Name() {
		case "Raw Data":
			g.Cursor = true
		case "Options":
			currentView.Highlight = false
			g.Cursor = false
		case "Sender 1":
			g.Cursor = true
			currentView.Title = "Sender 1 [*]"
		case "Sender 2":
			g.Cursor = true
			currentView.Title = "Sender 2 [*]"
		}

		v, err := g.View("Raw Data")
		if err != nil {
			// handle error
		}
		v.Clear()

		sizeX, sizeY := v.Size()

		///// GELEN-GİDEN-ERROR EKRANA YAZDIRMA /////

		var tempRawData []rawDataStruct

		if hexEnabled {
			tempRawData = rawDataHex
		} else {
			tempRawData = rawData
		}

		cy := 0
		var tempStr rawDataStruct
		for i := 0; i < len(tempRawData); i++ {
			tempStr = (tempRawData)[i]

			tempStr.data = strings.Replace(tempStr.data, "\r", "", -1)

			var timeStringOffset int

			if timeEnabled {
				timeStringOffset = 8
			}

			_, offsetRowOrigin := v.Origin()
			sLine, _ := v.Line(cy - offsetRowOrigin)

			var haveEnter int

			if sizeX-(len(sLine)+len(tempStr.data)+timeStringOffset) < 0 {
				for { // dizinin tek bir yapısı için döngü
					if timeEnabled {
						if sizeX-(len(sLine)+len(tempStr.data))*-1 >= 0 && sizeX-(len(sLine)+len(tempStr.data))*-1 < 8 {
							timeStringOffset = 0
						} else {
							timeStringOffset = 8
						}
					}
					offsetVal := (sizeX - (len(sLine) + len(tempStr.data))) * -1
					if offsetVal < 0 {
						offsetVal = 0
					}
					if len(tempStr.data)-offsetVal > 0 {
						var tempS string
						if timeEnabled && strings.Index(tempStr.data[:len(tempStr.data)-offsetVal], "\n") != -1 {
							tempS = strings.Replace(tempStr.data[:len(tempStr.data)-offsetVal], "\n", " ", -1)
							haveEnter = 1
						} else {
							tempS = tempStr.data[:len(tempStr.data)-offsetVal]
						}
						fmt.Fprintln(v, tempStr.color+tempS)
						tempStr.data = tempStr.data[len(tempStr.data)-offsetVal:]
						cy++
						sLine, _ = v.Line(cy - offsetRowOrigin)

					}
					if sizeX-len(tempStr.data)-timeStringOffset >= 0 {

						if timeEnabled && !tempStr.time.IsZero() {
							if last := len(tempStr.data) - 1; last >= 0 && tempStr.data[last] == '\n' || haveEnter == 1 {
								if haveEnter == 0 {
									tempStr.data = tempStr.data[:last]
								}
								tempStr.data += yellowText + fmt.Sprintf("%02d", tempStr.time.Hour()) + ":" + fmt.Sprintf("%02d", tempStr.time.Minute()) + ":" + fmt.Sprintf("%02d", tempStr.time.Second()) + "\n"
							} else {
								tempStr.data += yellowText + fmt.Sprintf("%02d", tempStr.time.Hour()) + ":" + fmt.Sprintf("%02d", tempStr.time.Minute()) + ":" + fmt.Sprintf("%02d", tempStr.time.Second())
							}
						}
						fmt.Fprint(v, tempStr.color+tempStr.data)
						if strings.Index(tempStr.data, "\n") != -1 {
							cy++
						}
						sLine, _ = v.Line(cy - offsetRowOrigin)

						break
					}
				}
			} else {
				if timeEnabled && !tempStr.time.IsZero() {
					if last := len(tempStr.data) - 1; last >= 0 && tempStr.data[last] == '\n' {
						tempStr.data = tempStr.data[:last]
						tempStr.data += yellowText + fmt.Sprintf("%02d", tempStr.time.Hour()) + ":" + fmt.Sprintf("%02d", tempStr.time.Minute()) + ":" + fmt.Sprintf("%02d", tempStr.time.Second()) + "\n"
					} else {
						tempStr.data += yellowText + fmt.Sprintf("%02d", tempStr.time.Hour()) + ":" + fmt.Sprintf("%02d", tempStr.time.Minute()) + ":" + fmt.Sprintf("%02d", tempStr.time.Second())
					}
				}
				fmt.Fprint(v, tempStr.color+tempStr.data)
				if strings.Index(tempStr.data, "\n") != -1 {
					cy++
				}
				sLine, _ = v.Line(cy - offsetRowOrigin)
				if sizeX-len(sLine) <= 0 {
					fmt.Fprintln(v, " ")
					cy++
				}

			}
		}
		if cursorLock {
			_, offsetRowOrigin := v.Origin()
			sLine, _ := v.Line(cy - offsetRowOrigin)
			if cy > sizeY {
				v.SetCursor(len(sLine), sizeY-1)
			} else {
				v.SetCursor(len(sLine), cy)
			}
		}

		v, err = g.View("Options")
		if err != nil {
			// handle error
		}
		v.Clear()

		sizeX, _ = v.Size()

		_, cy = v.Cursor()

		fmt.Fprintln(v, centerAligmentText("PORT", sizeX, false))
		if len(serialPortItems) > 0 {
			if cy == optionsViewAccessRow[0] && currentView.Name() == "Options" {
				fmt.Fprintln(v, centerAligmentText(serialPortItems[serialSelected], sizeX, true))
			} else {
				fmt.Fprintln(v, centerAligmentText(serialPortItems[serialSelected], sizeX, false))
			}
		} else {
			fmt.Fprintln(v, " ")
		}

		fmt.Fprintln(v, " ")
		fmt.Fprintln(v, centerAligmentText("BAUD RATE", sizeX, false))
		if cy == optionsViewAccessRow[1] && currentView.Name() == "Options" {
			fmt.Fprintln(v, centerAligmentText(fmt.Sprintf("%d", baudRateItems[baudSelected]), sizeX, true))
		} else {
			fmt.Fprintln(v, centerAligmentText(fmt.Sprintf("%d", baudRateItems[baudSelected]), sizeX, false))
		}

		fmt.Fprintln(v, " ")
		fmt.Fprintln(v, centerAligmentText("DATA SIZE", sizeX, false))
		if cy == optionsViewAccessRow[2] && currentView.Name() == "Options" {
			fmt.Fprintln(v, centerAligmentText(fmt.Sprintf("%d", dataSizeItems[dataSizeSelected]), sizeX, true))
		} else {
			fmt.Fprintln(v, centerAligmentText(fmt.Sprintf("%d", dataSizeItems[dataSizeSelected]), sizeX, false))
		}

		fmt.Fprintln(v, " ")
		fmt.Fprintln(v, centerAligmentText("PARITY", sizeX, false))
		if cy == optionsViewAccessRow[3] && currentView.Name() == "Options" {
			fmt.Fprintln(v, centerAligmentText(parityItems[paritySelected], sizeX, true))
		} else {
			fmt.Fprintln(v, centerAligmentText(parityItems[paritySelected], sizeX, false))
		}

		fmt.Fprintln(v, " ")
		fmt.Fprintln(v, centerAligmentText("STOPBIT", sizeX, false))
		if cy == optionsViewAccessRow[4] && currentView.Name() == "Options" {
			fmt.Fprintln(v, centerAligmentText(fmt.Sprintf("%d", stopBitItems[stopbitSelected]), sizeX, true))
		} else {
			fmt.Fprintln(v, centerAligmentText(fmt.Sprintf("%d", stopBitItems[stopbitSelected]), sizeX, false))
		}

		fmt.Fprintln(v, " ")
		serialOpenText := "OPEN PORT [ ]"
		if serialOpen {
			serialOpenText = "OPEN PORT [*]"
		}
		if cy == optionsViewAccessRow[5] && currentView.Name() == "Options" {
			fmt.Fprintln(v, centerAligmentText(serialOpenText, sizeX, true))
		} else {
			fmt.Fprintln(v, centerAligmentText(serialOpenText, sizeX, false))
		}

		fmt.Fprintln(v, " ")
		hexEnabledText := "Hex[ ]"
		if hexEnabled {
			hexEnabledText = "Hex[*]"
		}
		if cy == optionsViewAccessRow[6] && optionsViewAccessColSelected == 0 && currentView.Name() == "Options" {
			fmt.Fprint(v, centerAligmentText(hexEnabledText, sizeX/2, true))
		} else {
			fmt.Fprint(v, centerAligmentText(hexEnabledText, sizeX/2, false))
		}

		timeEnabledText := "Time[ ]"
		if timeEnabled {
			timeEnabledText = "Time[*]"
		}

		if cy == optionsViewAccessRow[6] && optionsViewAccessColSelected == 1 && currentView.Name() == "Options" {
			fmt.Fprintln(v, centerAligmentText(timeEnabledText, sizeX/2, true))
		} else {
			fmt.Fprintln(v, centerAligmentText(timeEnabledText, sizeX/2, false))
		}

		fmt.Fprintln(v, " ")
		fmt.Fprint(v, centerAligmentText("Sender 1", sizeX/2, false))
		fmt.Fprintln(v, centerAligmentText("Sender 2", sizeX/2, false))
		fmt.Fprintln(v, " ")

		sender1HexEnabledText := "Hex[ ]"
		if sender1HexEnabled {
			sender1HexEnabledText = "Hex[*]"
		}
		if cy == optionsViewAccessRow[7] && optionsViewAccessColSelected == 0 && currentView.Name() == "Options" {
			fmt.Fprint(v, centerAligmentText(sender1HexEnabledText, sizeX/2, true))
		} else {
			fmt.Fprint(v, centerAligmentText(sender1HexEnabledText, sizeX/2, false))
		}

		sender2HexEnabledText := "Hex[ ]"
		if sender2HexEnabled {
			sender2HexEnabledText = "Hex[*]"
		}
		if cy == optionsViewAccessRow[7] && optionsViewAccessColSelected == 1 && currentView.Name() == "Options" {
			fmt.Fprintln(v, centerAligmentText(sender2HexEnabledText, sizeX/2, true))
		} else {
			fmt.Fprintln(v, centerAligmentText(sender2HexEnabledText, sizeX/2, false))
		}

		fmt.Fprintln(v, " ")
		sender1CarriageReturnText := "\\r[ ]"
		if sender1CarriageReturn {
			sender1CarriageReturnText = "\\r[*]"
		}
		if cy == optionsViewAccessRow[8] && optionsViewAccessColSelected == 0 && currentView.Name() == "Options" {
			fmt.Fprint(v, centerAligmentText(sender1CarriageReturnText, sizeX/2, true))
		} else {
			fmt.Fprint(v, centerAligmentText(sender1CarriageReturnText, sizeX/2, false))
		}

		sender2CarriageReturnText := "\\r[ ]"
		if sender2CarriageReturn {
			sender2CarriageReturnText = "\\r[*]"
		}
		if cy == optionsViewAccessRow[8] && optionsViewAccessColSelected == 1 && currentView.Name() == "Options" {
			fmt.Fprintln(v, centerAligmentText(sender2CarriageReturnText, sizeX/2, true))
		} else {
			fmt.Fprintln(v, centerAligmentText(sender2CarriageReturnText, sizeX/2, false))
		}

		fmt.Fprintln(v, " ")
		sender1LineFeedText := "\\n[ ]"
		if sender1LineFeed {
			sender1LineFeedText = "\\n[*]"
		}
		if cy == optionsViewAccessRow[9] && optionsViewAccessColSelected == 0 && currentView.Name() == "Options" {
			fmt.Fprint(v, centerAligmentText(sender1LineFeedText, sizeX/2, true))
		} else {
			fmt.Fprint(v, centerAligmentText(sender1LineFeedText, sizeX/2, false))
		}

		sender2LineFeedText := "\\n[ ]"
		if sender2LineFeed {
			sender2LineFeedText = "\\n[*]"
		}
		if cy == optionsViewAccessRow[9] && optionsViewAccessColSelected == 1 && currentView.Name() == "Options" {
			fmt.Fprintln(v, centerAligmentText(sender2LineFeedText, sizeX/2, true))
		} else {
			fmt.Fprintln(v, centerAligmentText(sender2LineFeedText, sizeX/2, false))
		}

		fmt.Fprintln(v, test)

		v, err = g.View("Sender 1")
		if err != nil {
			errorCreate(err)
		}
		sLine, _ := v.Line(0)
		if sender1HexChanged {
			sender1HexChanged = false
			v.Clear()
			if sender1HexEnabled {
				stringToHexString(&sLine)
				fmt.Fprint(v, sLine)
				v.SetCursor(len(sLine), 0)
			} else {
				hexStringToString(&sLine)
				fmt.Fprint(v, sLine)
				v.SetCursor(len(sLine), 0)
			}
		}

		if sender1HexEnabled && hexFilter.MatchString(sLine) {
			v.Clear()
			sLine := hexFilter.ReplaceAllString(sLine, "")
			fmt.Fprint(v, sLine)
			v.SetCursor(len(sLine), 0)
		}

		v, err = g.View("Sender 2")
		if err != nil {
			errorCreate(err)
		}
		sLine, _ = v.Line(0)
		if sender2HexChanged {
			sender2HexChanged = false
			v.Clear()
			if sender2HexEnabled {
				stringToHexString(&sLine)
				fmt.Fprint(v, sLine)
				v.SetCursor(len(sLine), 0)
			} else {
				hexStringToString(&sLine)
				fmt.Fprint(v, sLine)
				v.SetCursor(len(sLine), 0)
			}
		}

		if sender2HexEnabled && hexFilter.MatchString(sLine) {
			v.Clear()
			sLine := hexFilter.ReplaceAllString(sLine, "")
			fmt.Fprint(v, sLine)
			v.SetCursor(len(sLine), 0)
		}

		return nil

	})
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {

	currentView := v.Name()

	if currentView == "Raw Data" {
		return rawDataCursorDown(g, v)
	}

	if currentView == "Options" {
		return optionsCursorDown(g, v)
	}

	if v != nil {
		cx, cy := v.Cursor()
		_, sizeY := v.Size()
		for i := 1; i < sizeY-cy+1; i++ {
			if str, err := v.Line(cy + i); err == nil {
				if str != " " && str != "" {
					if err := v.SetCursor(cx, cy+i); err == nil {
						break
					}
				}
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {

	currentView := v.Name()

	if currentView == "Raw Data" {
		return rawDataCursorUp(g, v)
	}

	if currentView == "Options" {
		return optionsCursorUp(g, v)
	}

	if v != nil {
		cx, cy := v.Cursor()
		for cy > 0 {
			cy--
			if str, err := v.Line(cy); err == nil {
				if str != " " && str != "" {
					if err := v.SetCursor(cx, cy); err == nil {
						break
					}
				}
			}
		}
	}
	return nil
}

func cursorLeft(g *gocui.Gui, v *gocui.View) error {

	currentView := v.Name()

	if currentView == "Raw Data" {
		return rawDataCursorLeft(g, v)
	}

	if currentView == "Options" {
		return optionsCursorLeft(g, v)
	}

	if strings.Index(currentView, "Sender") != -1 {
		return senderCursorLeft(g, v)
	}

	return nil

}

func cursorRight(g *gocui.Gui, v *gocui.View) error {

	currentView := v.Name()

	if currentView == "Raw Data" {
		return rawDataCursorRight(g, v)
	}

	if currentView == "Options" {
		return optionsCursorRight(g, v)
	}

	if strings.Index(currentView, "Sender") != -1 {
		return senderCursorRight(g, v)
	}

	return nil

}

func rawDataCursorDown(g *gocui.Gui, v *gocui.View) error {

	cursorKey[0] = true
	if !cursorLock {
		go cursorLocker(g, v)
	}

	if v != nil {
		cx, cy := v.Cursor()
		sLine, _ := v.Line(cy + 1)
		if cx > len(sLine)-1 {
			cx = len(sLine) - 1
		}
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func rawDataCursorUp(g *gocui.Gui, v *gocui.View) error {
	cursorLock = false
	v.Autoscroll = false
	cursorKey[1] = true
	go func() {
		time.Sleep(time.Millisecond * 200)
		cursorKey[1] = false
	}()
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		sLine, _ := v.Line(cy - 1)
		if cx > len(sLine)-1 {
			cx = len(sLine) - 1
		}
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func rawDataCursorLeft(g *gocui.Gui, v *gocui.View) error {
	cursorLock = false
	cx, cy := v.Cursor()
	cx--
	if cx < 0 {
		cy--
		sLine, _ := v.Line(cy)
		cx = len(sLine) - 1
	}
	v.SetCursor(cx, cy)
	return nil
}

func rawDataCursorRight(g *gocui.Gui, v *gocui.View) error {
	cursorLock = false
	cx, cy := v.Cursor()
	sLine, _ := v.Line(cy)
	cx++
	if cx > len(sLine)-1 {
		cy++
		cx = 0
	}
	v.SetCursor(cx, cy)
	return nil
}

func optionsCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		_, sizeY := v.Size()
		if optionsViewAccessRowSelected < len(optionsViewAccessRow)-1 && optionsViewAccessRow[optionsViewAccessRowSelected+1] < sizeY {
			optionsViewAccessRowSelected++
		}
		cy = optionsViewAccessRow[optionsViewAccessRowSelected]
		v.SetCursor(cx, cy)
	}
	return nil
}

func optionsCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		_, sizeY := v.Size()
		if optionsViewAccessRowSelected > 0 {
			optionsViewAccessRowSelected--
			for optionsViewAccessRowSelected > 0 && sizeY <= optionsViewAccessRow[optionsViewAccessRowSelected] {
				optionsViewAccessRowSelected--
			}
		}
		cy = optionsViewAccessRow[optionsViewAccessRowSelected]
		v.SetCursor(cx, cy)
	}
	return nil
}

func optionsCursorLeft(g *gocui.Gui, v *gocui.View) error {
	if optionsViewAccessColSelected > 0 {
		optionsViewAccessColSelected--
	}
	return nil
}

func optionsCursorRight(g *gocui.Gui, v *gocui.View) error {
	if optionsViewAccessColSelected < optionsViewAccessColLimit-1 {
		optionsViewAccessColSelected++
	}
	return nil
}

func senderCursorLeft(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	if cx > 0 {
		cx--
		v.SetCursor(cx, cy)
	}
	return nil
}

func senderCursorRight(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	sLine, _ := v.Line(cy)
	if cx < len(sLine) {
		cx++
		v.SetCursor(cx, cy)
	}
	return nil
}

func cursorLocker(g *gocui.Gui, v *gocui.View) {
	time.Sleep(time.Millisecond * 100)
	if cursorKey[0] && cursorKey[1] {
		cursorLock = true
		v.Autoscroll = true
	}
	cursorKey[0] = false
	cursorKey[1] = false
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	var err error
	switch v.Name() {
	case "Sender 1":
		v.Title = "Sender 1 [ ]"
		_, err = g.SetCurrentView("Sender 2")
	case "Sender 2":
		v.Title = "Sender 2 [ ]"
		_, err = g.SetCurrentView("Raw Data")
	case "Raw Data":
		_, err = g.SetCurrentView("Options")
		return err
	case "Options":
		_, err = g.SetCurrentView("Sender 1")
		return err
	case "Serial Ports":
		if err = g.DeleteView("Serial Ports"); err != nil {
			return err
		}
		baudRateSelect(g, v)
	case "Baud Rates":
		if err = g.DeleteView("Baud Rates"); err != nil {
			return err
		}
		dataSizeSelect(g, v)
	case "Data Sizes":
		if err = g.DeleteView("Data Sizes"); err != nil {
			return err
		}
		paritySelect(g, v)
	case "Parity":
		if err = g.DeleteView("Parity"); err != nil {
			return err
		}
		stopbitSelect(g, v)
	case "Stopbit":
		if err = g.DeleteView("Stopbit"); err != nil {
			return err
		}
		serialPortSelect(g, v)
	}
	return err
}

func selected(g *gocui.Gui, v *gocui.View) error {
	switch v.Name() {
	case "Sender 1":
		senderEnter(g, v)
	case "Sender 2":
		senderEnter(g, v)
	case "Options":
		_, cy := v.Cursor()
		switch cy {
		case 1:
			if !serialOpen {
				serialPortSelect(g, v)
			}
		case 4:
			if !serialOpen {
				baudRateSelect(g, v)
			}
		case 7:
			if !serialOpen {
				dataSizeSelect(g, v)
			}
		case 10:
			if !serialOpen {
				paritySelect(g, v)
			}
		case 13:
			if !serialOpen {
				stopbitSelect(g, v)
			}
		case 15:
			openPortSc(g, v)
		case 17:
			if optionsViewAccessColSelected == 0 {
				hexEnabled = !hexEnabled
				// if hexEnabled {
				// 	rawDatatoHex()
				// }
			}
			if optionsViewAccessColSelected == 1 {
				timeEnabled = !timeEnabled
			}
		case 21:
			if optionsViewAccessColSelected == 0 {
				sender1HexEnabled = !sender1HexEnabled
				sender1HexChanged = true
			}
			if optionsViewAccessColSelected == 1 {
				sender2HexEnabled = !sender2HexEnabled
				sender2HexChanged = true
			}
		case 23:
			if optionsViewAccessColSelected == 0 {
				sender1CarriageReturn = !sender1CarriageReturn
			}
			if optionsViewAccessColSelected == 1 {
				sender2CarriageReturn = !sender2CarriageReturn
			}
		case 25:
			if optionsViewAccessColSelected == 0 {
				sender1LineFeed = !sender1LineFeed
			}
			if optionsViewAccessColSelected == 1 {
				sender2LineFeed = !sender2LineFeed
			}
		}
	case "Serial Ports":
		_, cy := v.Cursor()
		if cy < len(serialPortItems) {
			serialSelected = int32(cy)
		}
		if err := g.DeleteView("Serial Ports"); err != nil {
			return err
		}
		if _, err := g.SetCurrentView("Options"); err != nil {
			return err
		}
	case "Baud Rates":
		_, cy := v.Cursor()
		if cy < len(baudRateItems) {
			baudSelected = int32(cy)
		}
		if err := g.DeleteView("Baud Rates"); err != nil {
			return err
		}
		if _, err := g.SetCurrentView("Options"); err != nil {
			return err
		}
	case "Data Sizes":
		_, cy := v.Cursor()
		if cy < len(dataSizeItems) {
			dataSizeSelected = int32(cy)
		}
		if err := g.DeleteView("Data Sizes"); err != nil {
			return err
		}
		if _, err := g.SetCurrentView("Options"); err != nil {
			return err
		}
	case "Parity":
		_, cy := v.Cursor()
		if cy < len(parityItems) {
			paritySelected = int32(cy)
		}
		if err := g.DeleteView("Parity"); err != nil {
			return err
		}
		if _, err := g.SetCurrentView("Options"); err != nil {
			return err
		}
	case "Stopbit":
		_, cy := v.Cursor()
		if cy < len(stopBitItems) {
			stopbitSelected = int32(cy)
		}
		if err := g.DeleteView("Stopbit"); err != nil {
			return err
		}
		if _, err := g.SetCurrentView("Options"); err != nil {
			return err
		}
	case "Buffer Size":
		_, cy := v.Cursor()
		if cy < len(bufferSizeItems) {
			bufferSizeSelected = cy
		}
		if err := g.DeleteView("Buffer Size"); err != nil {
			return err
		}
		if _, err := g.SetCurrentView("Options"); err != nil {
			return err
		}
	}
	update(g)
	return nil
}

func senderEnter(g *gocui.Gui, v *gocui.View) error {

	if !serialOpen {
		errorCreate(errors.New("Port is not open"))
		return nil
	}

	senderName := v.Name()

	_, cy := v.Cursor()
	sLine, _ := v.Line(cy)

	if senderName == "Sender 1" && sender1HexEnabled {
		sLine = strings.Replace(sLine, " ", "", -1)
		if len(sLine)%2 != 0 {
			sLine += "0"
		}
		decoded, err := hex.DecodeString(sLine)
		if err != nil {
			errorCreate(err)
			return nil
		}

		sLine = string(decoded)
	}
	if senderName == "Sender 1" && sender1CarriageReturn {
		sLine += "\r"
	}
	if senderName == "Sender 1" && sender1LineFeed {
		sLine += "\n"
	}

	if senderName == "Sender 2" && sender2HexEnabled {
		sLine = strings.Replace(sLine, " ", "", -1)
		if len(sLine)%2 != 0 {
			sLine += "0"
		}
		decoded, err := hex.DecodeString(sLine)
		if err != nil {
			errorCreate(err)
			return nil
		}

		sLine = string(decoded)
	}
	if senderName == "Sender 2" && sender2CarriageReturn {
		sLine += "\r"
	}
	if senderName == "Sender 2" && sender2LineFeed {
		sLine += "\n"
	}

	if len(sLine) > 0 {
		if err := serialWrite(sLine); err != nil {
			errorCreate(err)
		}
	} else {
		errorCreate(errors.New("Empty " + v.Name()))
	}
	return nil
}

func goBack(g *gocui.Gui, v *gocui.View) error {

	currentView := v.Name()

	if currentView == "Serial Ports" || currentView == "Baud Rates" || currentView == "Data Sizes" || currentView == "Parity" || currentView == "Stopbit" || currentView == "Buffer Size" {
		if err := g.DeleteView(currentView); err != nil {
			return err
		}
		if _, err := g.SetCurrentView("Options"); err != nil {
			return err
		}
	}

	return nil
}

func rawDataView(g *gocui.Gui, v *gocui.View) error {

	currentView := v.Name()

	if currentView == "Serial Ports" || currentView == "Baud Rates" || currentView == "Data Sizes" || currentView == "Parity" || currentView == "Stopbit" {
		if err := g.DeleteView(currentView); err != nil {
			return err
		}
	}

	if strings.Index(currentView, "Sender") != -1 {
		v.Title = currentView + " [ ]"
		senderPreviousView = currentView
	}

	if _, err := g.SetCurrentView("Raw Data"); err != nil {
		return err
	}
	return nil
}

func optionsView(g *gocui.Gui, v *gocui.View) error {

	currentView := v.Name()

	if currentView == "Serial Ports" || currentView == "Baud Rates" || currentView == "Data Sizes" || currentView == "Parity" || currentView == "Stopbit" {
		if err := g.DeleteView(currentView); err != nil {
			return err
		}
	}

	if strings.Index(currentView, "Sender") != -1 {
		v.Title = currentView + " [ ]"
		senderPreviousView = currentView
	}

	if _, err := g.SetCurrentView("Options"); err != nil {
		return err
	}

	return nil
}

func senderView(g *gocui.Gui, v *gocui.View) error {

	currentView := v.Name()

	if currentView == "Serial Ports" || currentView == "Baud Rates" || currentView == "Data Sizes" || currentView == "Parity" || currentView == "Stopbit" {
		if err := g.DeleteView(currentView); err != nil {
			return err
		}
	}

	if currentView == "Sender 1" {
		v.Title = "Sender 1 [ ]"
		if _, err := g.SetCurrentView("Sender 2"); err != nil {
			return err
		}
	} else if currentView == "Sender 2" {
		v.Title = "Sender 2 [ ]"
		if _, err := g.SetCurrentView("Sender 1"); err != nil {
			return err
		}
	} else {
		if _, err := g.SetCurrentView(senderPreviousView); err != nil {
			return err
		}
	}

	return nil
}

func serialPortSelect(g *gocui.Gui, v *gocui.View) error {
	scanPorts()
	g.Cursor = false
	maxX, maxY := g.Size()
	if v, err := g.SetView("Serial Ports", maxX/2-30, maxY/4, maxX/2+30, maxY/4+1+len(serialPortItems)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Serial Ports"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		sizeX, _ := v.Size()
		for i := 0; i < len(serialPortItems); i++ {
			fmt.Fprintln(v, centerAligmentText(serialPortItems[i], sizeX, false))
		}

		if _, err := g.SetCurrentView("Serial Ports"); err != nil {
			return err
		}

		v.SetCursor(0, int(serialSelected))

	}
	return nil
}

func baudRateSelect(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	maxX, maxY := g.Size()
	if v, err := g.SetView("Baud Rates", maxX/2-10, maxY/4, maxX/2+11, maxY/4+1+len(baudRateItems)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Baud Rates"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		sizeX, _ := v.Size()
		for i := 0; i < len(baudRateItems); i++ {
			fmt.Fprintln(v, centerAligmentText(fmt.Sprintf("%d", baudRateItems[i]), sizeX, false))
		}

		if _, err := g.SetCurrentView("Baud Rates"); err != nil {
			return err
		}

		v.SetCursor(0, int(baudSelected))

	}
	return nil
}

func dataSizeSelect(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	maxX, maxY := g.Size()
	if v, err := g.SetView("Data Sizes", maxX/2-6, maxY/4, maxX/2+6, maxY/4+1+len(dataSizeItems)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Data Sizes"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		sizeX, _ := v.Size()

		for i := 0; i < len(dataSizeItems); i++ {
			fmt.Fprintln(v, centerAligmentText(fmt.Sprintf("%d", dataSizeItems[i]), sizeX, false))
		}

		if _, err := g.SetCurrentView("Data Sizes"); err != nil {
			return err
		}

		v.SetCursor(0, int(dataSizeSelected))

	}
	return nil
}

func paritySelect(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	maxX, maxY := g.Size()
	if v, err := g.SetView("Parity", maxX/2-6, maxY/4, maxX/2+6, maxY/4+1+len(parityItems)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Parity"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		sizeX, _ := v.Size()

		for i := 0; i < len(parityItems); i++ {
			fmt.Fprintln(v, centerAligmentText(parityItems[i], sizeX, false))
		}

		if _, err := g.SetCurrentView("Parity"); err != nil {
			return err
		}

		v.SetCursor(0, int(paritySelected))

	}
	return nil
}

func stopbitSelect(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	maxX, maxY := g.Size()
	if v, err := g.SetView("Stopbit", maxX/2-6, maxY/4, maxX/2+6, maxY/4+1+len(stopBitItems)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Stopbit"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		sizeX, _ := v.Size()

		for i := 0; i < len(stopBitItems); i++ {
			fmt.Fprintln(v, centerAligmentText(fmt.Sprintf("%d", stopBitItems[i]), sizeX, false))
		}

		if _, err := g.SetCurrentView("Stopbit"); err != nil {
			return err
		}

		v.SetCursor(0, int(stopbitSelected))

	}
	return nil
}

func bufferSizeView(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	maxX, maxY := g.Size()
	if v, err := g.SetView("Buffer Size", maxX/2-8, maxY/4, maxX/2+6, maxY/4+1+len(bufferSizeItems)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Buffer Size"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		sizeX, _ := v.Size()

		for i := 0; i < len(bufferSizeItems); i++ {
			fmt.Fprintln(v, centerAligmentText(bufferSizeItems[i], sizeX, false))
		}

		if _, err := g.SetCurrentView("Buffer Size"); err != nil {
			return err
		}

		v.SetCursor(0, int(bufferSizeSelected))
	}

	return nil
}

func openPortSc(g *gocui.Gui, v *gocui.View) error {
	var err error
	if err = openPortSw(); err != nil {
		errorCreate(err)
	}
	return nil
}

func clearData(g *gocui.Gui, v *gocui.View) error {
	rawDataMutex.Lock()
	defer rawDataMutex.Unlock()

	rawDataProccesed = 0
	rawData = []rawDataStruct{}
	rawDataHex = []rawDataStruct{}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func centerAligmentText(text string, maxLength int, selected bool) string {
	if len(text) >= maxLength {
		if selected {
			return greenRedText + text + whiteText
		}
		return text

	}
	str := whiteText + strings.Repeat(" ", (maxLength-len(text))/2)
	if selected {
		str += greenRedText + text + whiteText
	} else {
		str += text
	}
	str += strings.Repeat(" ", (maxLength-len(text))/2+1)

	return str
}

func rawDataAdd(str rawDataStruct) {

	rawDataMutex.Lock()
	defer rawDataMutex.Unlock()

	if len(strings.Split(str.data, "\n")) > 2 {
		for i := 0; i < len(strings.Split(str.data, "\n")); i++ {
			if strings.Split(str.data, "\n")[i] != "" {
				rawData = append(rawData, rawDataStruct{data: strings.Split(str.data, "\n")[i], color: str.color, time: str.time})
				if i != len(strings.Split(str.data, "\n"))-1 {
					rawData[len(rawData)-1].data += "\n"
				}
				if hexEnabled {
					rawDataHexAdd(rawData[len(rawData)-1])
				}
			}
		}
	} else {
		rawData = append(rawData, str)
		if hexEnabled {
			rawDataHexAdd(str)
		}
	}
}

func rawDataArray() {
	rawDataMutex.Lock()
	defer rawDataMutex.Unlock()

	if len(rawData) == rawDataProccesed {
		return
	}

	tempRawData := rawData[:rawDataProccesed]

	//fmt.Println(len(rawData), rawDataProccesed)

	for i := rawDataProccesed; i < len(rawData); i++ {
		if i == 0 {
			tempRawData = append(tempRawData, rawData[0])
			continue
			//tempTime = tempRawData[0].time
		}
		if int(rawData[i].time.Sub(tempRawData[len(tempRawData)-1].time).Seconds()) == 0 && rawData[i].color == tempRawData[len(tempRawData)-1].color {
			tempRawData[len(tempRawData)-1].data += rawData[i].data
		} else {
			tempRawData = append(tempRawData, rawData[i])
		}
	}
	rawData = tempRawData
	for i := rawDataProccesed; i < len(rawData); i++ {
		rawDataHexAdd(rawData[i])
	}

	rawDataProccesed = len(tempRawData)
	if bufferSizeSelected != 5 {
		if val, _ := strconv.Atoi(bufferSizeItems[bufferSizeSelected]); len(rawData) > val {
			rawData = rawData[len(rawData)-val:]
			rawDataHex = rawDataHex[len(rawDataHex)-val:]
			rawDataProccesed = len(rawData)
		}
	}

}

func rawDataHexAdd(str rawDataStruct) {
	var tempStr string
	var specialChar string
	if str.consolMessage {
		rawDataHex = append(rawDataHex, str)
		return
	}
	for i := 0; i < len(str.data); i++ {
		if str.data[i] == '\n' || str.data[i] == '\r' {
			specialChar += "{" + fmt.Sprintf("%02X", str.data[i]) + "}" + string(str.data[i])
		} else {
			tempStr += fmt.Sprintf("%02X ", str.data[i])
		}
	}
	if specialChar != "" {
		var t time.Time
		rawDataHex = append(rawDataHex, rawDataStruct{data: tempStr, color: str.color, time: t})
		rawDataHex = append(rawDataHex, rawDataStruct{data: specialChar, color: whiteText, time: str.time})
	} else {
		rawDataHex = append(rawDataHex, rawDataStruct{data: tempStr, color: str.color, time: str.time})
	}
}

func rawDatatoHex() {
	rawDataMutex.Lock()
	defer rawDataMutex.Unlock()

	rawDataHex = make([]rawDataStruct, 0)
	for i := 0; i < len(rawData); i++ {
		rawDataHexAdd(rawData[i])
	}
}

func stringToHexString(s *string) {
	*s = fmt.Sprintf("% X", *s)
}

func hexStringToString(s *string) {
	src := []byte(*s)
	dst := make([]byte, hex.DecodedLen(len(src)))

	_, err := hex.Decode(dst, src)
	if err != nil {
		errorCreate(err)
	}
	*s = string(dst)
}

func infoCreate(message string) {
	rawDataAdd(rawDataStruct{data: message + "\n", color: greenRedText, time: time.Now(), consolMessage: true})
}

func errorCreate(err error) {
	rawDataAdd(rawDataStruct{data: err.Error() + "\n", color: purpleText, time: time.Now(), consolMessage: true})
}
