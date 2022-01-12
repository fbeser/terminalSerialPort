package main

import (
	"errors"
	"io"
	"time"

	"github.com/tarm/serial"
)

var serialOpen bool

var serialPortItems []string

var baudRateItems []int = []int{
	600,
	1200,
	2400,
	4800,
	9600,
	14400,
	19200,
	38400,
	57600,
	115200,
	230400,
	460800,
	500000,
	576000,
	921600,
	1000000,
	1152000,
	1500000,
	2000000,
	2500000,
	3000000,
	3500000,
	4000000,
}

var dataSizeItems []int = []int{
	5,
	6,
	7,
	8,
}

var parityItems []string = []string{
	"None",
	"Odd",
	"Even",
	"Mark",
	"Space",
}

var stopBitItems []int = []int{
	1,
	2,
}

var serialSelected int32 = 0
var baudSelected int32 = 4
var dataSizeSelected int32 = 3
var paritySelected int32 = 0
var stopbitSelected int32 = 0

var sPort *serial.Port

func openPortSw() error {

	if serialOpen {
		return serialClose(sPort)
	}

	if int32(len(serialPortItems)-1) < serialSelected || len(serialPortItems[serialSelected]) == 0 {
		return errors.New("port not found")
	}

	c := &serial.Config{Name: serialPortItems[serialSelected], Baud: baudRateItems[baudSelected], Size: byte(dataSizeItems[dataSizeSelected]), Parity: serial.Parity(paritySelected), StopBits: serial.StopBits(stopBitItems[stopbitSelected]), ReadTimeout: time.Millisecond}
	var err error
	if sPort, err = serial.OpenPort(c); err != nil {
		return err
	}

	serialOpen = true

	go rawDataTimers()

	infoCreate("Serial Port Opened.")

	go serialRead(sPort)

	return nil
}

func serialWrite(message string) error {
	if !serialOpen {
		return errors.New("Port is not open")
	}
	if len(message) == 0 {
		return errors.New("Empty Send Message")
	}

	b := []byte(message)
	_, err := sPort.Write(b)
	if err != nil {
		return err
	}

	//.........
	rawDataAdd(rawDataStruct{data: message, typ: "TX", color: redText, time: time.Now()})
	return nil
}

func serialRead(p *serial.Port) {
	//timer := true

	//readBytes := make([]byte, 0)
	readByte := make([]byte, 1)

	//	dataNow := false
	for {
		_, err := p.Read(readByte)
		if err != nil && err != io.EOF {
			err = serialClose(p)
			if err != nil {
				errorCreate(err)
			}
			break
		} else if err == io.EOF {
			continue
		}
		rawDataAdd(rawDataStruct{data: string(readByte), typ: "RX", color: greenText, time: time.Now()})
		// dataNow = true
		// readBytes = append(readBytes, readByte[0])
		// if timer {
		// 	timer = false
		// 	go func() {
		// 		for dataNow {
		// 			dataNow = false
		// 			t := time.Now()
		// 			time.Sleep(time.Second * 1)
		// 			if !dataNow {
		// 				rawDataAdd(rawDataStruct { data : string(readBytes), color : greenText, time : t })
		// 				readBytes = make([]byte, 0)
		// 			}
		// 		}
		// 		timer = true
		// 	}()
		// }
	}
}

func serialClose(p *serial.Port) (err error) {
	err = p.Close()
	if err != nil {
		return err
	}
	serialOpen = false
	infoCreate("Serial Port Closed.")
	return nil
}
