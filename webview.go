/*

Spectral analysis of a wav file. This time outputting via websockets to a web page.

Open the page draw.html it will connect via websocket to the localhost that should be runing this binary.

You get a graph of applitude/frequency

The intention is to expand to read streaming audio content.

*/

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mjibson/go-dsp/spectral"
	"github.com/mjibson/go-dsp/wav"
	webskt "inf-stash-01.gammatelecom.com/scm/df/pkg_websocket"
)

// Trap SIGINT to trigger a shutdown
var signals = make(chan os.Signal, 1)

func main() {

	// Start Websocket Listener\
	go webskt.Init()

	filePtr, err := os.Open("service-bell.wav")
	//filePtr, err := os.Open("audiocheck.net_dtmf_123456789_#.wav")
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		os.Exit(0)
	}

	thisWav, err := wav.New(filePtr)
	if err != nil {
		fmt.Printf("Error interprorating WAV: %v", err)
	}

	sampRate := thisWav.Header.SampleRate
	numSamples := thisWav.Samples
	fmt.Println(thisWav)

	a32, err := thisWav.ReadFloats(numSamples)
	if err != nil {
		fmt.Printf("Error reading in values from WAV: %v", err)
	}

	a := make([]float64, sampRate/10)

	timer := time.Tick(time.Millisecond * 100)
	fmt.Printf("Sample rate %v, and samples contained %v\n", sampRate, numSamples)
	for {
		if ws, isset := webskt.ListenFor[""]; isset {
			for t := 0; t < numSamples; t = t + int(sampRate/10) {
				//mark := time.Now()

				//time.Sleep(time.Millisecond*100)
				//fmt.Printf("Making a 1/10s sample of %v samples\n", int(sampRate/20))
				// drop out where we don't have enough sample for a complete check
				if numSamples-t < int(sampRate/10) {
					continue
				}

				for k, v := range a32[t : t+int(sampRate/10)] {
					a[k] = float64(v)
				}

				pwr, freq := spectral.Pwelch(a, float64(sampRate), &spectral.PwelchOptions{NFFT: 2048})
				myData := ""

				for k, _ := range pwr {
					myData = fmt.Sprintf("%v%v=%v,", myData, freq[k], pwr[k])
				}
				var sendMessage = webskt.WebsocketMessage{
					Event:  "NewData",
					Param1: strconv.Itoa(len(pwr)),
					Data:   myData,
				}
				_ = <-timer
				webskt.Send(sendMessage, ws)
				//fmt.Println("sending")
				//fmt.Println(time.Since(mark))
			}
		}
	}
}
