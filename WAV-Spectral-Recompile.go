/*

Following information theory priciples that any signal is just a combination of sine waves.
I attempt to disect a wav, do spectral analytics and recombine a new wavform

This would prove theory, but also allow intersting compression base on minimal information rate.

I think I have the maths right, but not great, need to work on the maths.

*/

package main

import (
	"fmt"
	"math"

	//"math"
	"os"
	//"math/cmplx"
	//"github.com/mjibson/go-dsp/fft"
	//"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/spectral"
	//"github.com/mjibson/go-dsp/wav"
	"github.com/youpy/go-wav"
)

func main() {

	fileRead, err := os.Open("service-bell.wav")
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		os.Exit(0)
	}

	fileWrite, err := os.Create("service-bell-out.wav")
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		os.Exit(0)
	}

	reader := wav.NewReader(fileRead)
	readerFmt, _ := reader.Format()
	writer := wav.NewWriter(fileWrite, readerFmt.SampleRate*10, 1, readerFmt.SampleRate, readerFmt.BitsPerSample)
	sampRate := readerFmt.SampleRate

	a := make([]float64, sampRate/20) // Input amplitudes
	A := make([]float64, sampRate/20) // Output amplitudes

	fmt.Printf("WAV is %v format, of %v channels, and sample rate %v\n", readerFmt.AudioFormat, readerFmt.NumChannels, sampRate)

	for {

		var samples []wav.Sample
		var err error
		if samples, err = reader.ReadSamples(sampRate / 20); err != nil {
			break
		}

		for k, v := range samples {
			a[k] = float64(v.Values[0])
		}

		pwr, freq := spectral.Pwelch(a, float64(sampRate), &spectral.PwelchOptions{NFFT: 128})

		//fmt.Printf("this contains %v samples at %v rate, computed in %v\n\n", numSamples, sampRate, time.Since(startTime))

		// Itterate over all intervals of sample rate
		for i, _ := range A {
			// starting with an amplitude of 0
			var a float64
			b := float64(i) / float64(len(A))
			// Itterate over frequence distribution
			for k, _ := range pwr {
				// Calculate the amplutude gained from each frequency
				a = a + (pwr[k] * math.Sin(2.0*math.Pi*b*freq[k]))
			}
			A[i] = a
		}

		out := make([]wav.Sample, sampRate/20)
		// Reform back to int values
		for k, _ := range A {
			out[k].Values[0] = int(A[k])
		}

		for k, _ := range samples {
			fmt.Printf("%v,%v,%v\n", samples[k].Values, a[k], A[k])
		}
		//os.Exit(1)

		// Write data to file
		writer.WriteSamples(out)

	}
}
