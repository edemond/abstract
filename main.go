// Abstract implements an interpreter for a higher-level music notation and a player
// that emits MIDI.
package main

import (
	"github.com/edemond/abstract/ast"
	"github.com/edemond/abstract/drivers"
	"github.com/edemond/abstract/drivers/alsa"
	"github.com/edemond/abstract/drivers/jack"
	"github.com/edemond/abstract/parser"
	"github.com/edemond/abstract/types"
	"github.com/edemond/midi"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

const VERSION = "0.0.2"

// TODO: We need a way of listing potential drivers (and whether or not they'd work on the target system?)
var driverFlag = flag.String("d", "rawmidi", "\tMIDI driver (rawmidi).")
var modeListDevices = flag.Bool("a", false, "\tList available ALSA MIDI device names.")
var modeVersion = flag.Bool("v", false, "\tPrint version information.")
var loopFlag = flag.Bool("l", false, "\tLoop (Ctrl+C to stop).")

func listDevices() error {
	// TODO: Let the user choose which driver to list devices for at the command line.
	// If they don't pass the -d option, list all devices on all drivers that we can use.
	fmt.Println("Available ALSA MIDI devices:")
	devices, err := midi.GetDevices("alsa")
	if err != nil {
		return err
	}
	for _, device := range devices {
		// TODO: Show if device is input or output.
		fmt.Println(device.Name())
	}
	return nil
}

// Open a file and parse it into an AST.
func parse(filename string) (*ast.PlayStatement, error) {
	fmt.Printf("Compiling %v...\n", filename)
	p, err := parser.FromFile(filename)
	if err != nil {
		return nil, err
	}

	return p.Parse()
}

// Perform semantic analysis on an AST.
// Returns root part, total polyphony, BPM, PPQ, error.
func analyze(stmt *ast.PlayStatement, driver drivers.Driver) (types.Part, int, int, int, error) {
	a := NewAnalyzer()
	part, err := a.Analyze(stmt)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	// Open instruments.
	insts, err := a.OpenInstruments(driver)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return part, types.TotalVoices(insts), a.bpm, a.ppq, nil
}

func printUsage() {
	fmt.Println("usage: abstract <file.abs>\n")
	fmt.Println("commands:\n")
	flag.PrintDefaults()
	fmt.Println("")
}

type supportedDriver struct {
	description string
	constructor func() (drivers.Driver, error)
}

var _drivers = map[string]supportedDriver{
	"rawmidi": {
		"ALSA 'rawmidi' driver.",
		alsa.NewRawMidiDriver,
	},
	"jack": {
		"JACK 1.x driver.",
		jack.NewJACKDriver,
	},
}

func printSupportedDrivers() {
	fmt.Println("Supported drivers:\n")
	for name, driver := range _drivers {
		fmt.Printf("\t%v\t\t%v\n", name, driver.description)
	}
}

// Instantiate the driver of the given name.
// Returns nil, error if there's no such driver.
func getDriver(name string) (drivers.Driver, error) {
	sd, ok := _drivers[name]
	if !ok {
		return nil, fmt.Errorf("No such driver '%v'.", name)
	}
	driver, err := sd.constructor()
	if err != nil {
		return nil, err
	}
	return driver, nil
}

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	if *modeListDevices {
		err := listDevices()
		if err != nil {
			fmt.Println(err)
		}
		return
	} else if *modeVersion {
		fmt.Printf("Abstract v%v\n", VERSION)
		return
	}

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		return
	}
	filename := args[0]

	stmt, err := parse(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: To fail fastest, we should validate the driver before parsing, then open it afterwards.
	driver, err := getDriver(*driverFlag)
	if err != nil {
		fmt.Println(err)
		printSupportedDrivers()
		return
	}
	defer driver.Close()

	// TODO: Really, five return values? Can we put this into a struct
	// that analyze and driver.Play use to communicate?
	part, polyphony, bpm, ppq, err := analyze(stmt, driver)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = driver.Play(part, bpm, ppq, *loopFlag, polyphony)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Done.")
}
