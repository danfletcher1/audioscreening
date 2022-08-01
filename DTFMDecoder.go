/*

This play was to extract dtfm codes from a wav file

Its a play with fft and speciral analytics

*/

package main

import (
	"fmt"
	"time"

	//"math"
	"os"
	//"math/cmplx"
	//"github.com/mjibson/go-dsp/fft"
	//"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/spectral"
	"github.com/mjibson/go-dsp/wav"
)

func main() {

	filePtr, err := os.Open("01616661123.wav")
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

	a32, err := thisWav.ReadFloats(numSamples)
	if err != nil {
		fmt.Printf("Error reading in values from WAV: %v", err)
	}

	a := make([]float64, sampRate/20)
	Y := ""
	gain := float64(100)

	startTime := time.Now()

	fmt.Printf("Sample rate %v, and samples contained %v\n", sampRate, numSamples)
	for t := 0; t < numSamples; t = t + int(sampRate/20) {

		//fmt.Printf("Making a 1/10s sample of %v samples\n", int(sampRate/20))
		// drop out where we don't have enough sample for a complete check
		if numSamples-t < int(sampRate/20) {
			continue
		}

		for k, v := range a32[t : t+int(sampRate/20)] {
			a[k] = float64(v)
		}

		pwr, freq := spectral.Pwelch(a, float64(sampRate), &spectral.PwelchOptions{NFFT: 512})

		//fmt.Printf("this contains %v samples at %v rate, computed in %v\n\n", numSamples, sampRate, time.Since(startTime))

		var A, B, C, D, E, F, G, SUM float64
		for k, _ := range pwr {

			f := freq[k]
			p := pwr[k]
			SUM = SUM + p

			if f > 697*.98 && f < 697*1.02 {
				A = A + p
			} else if f > 770*.98 && f < 770*1.02 {
				B = B + p
			} else if f > 852*.98 && f < 852*1.02 {
				C = C + p
			} else if f > 941*.98 && f < 941*1.02 {
				D = D + p
			}

			if f > 1209*.98 && f < 1209*1.02 {
				E = E + p
			} else if f > 1336*.98 && f < 1336*1.02 {
				F = F + p
			} else if f > 1477*.98 && f < 1477*1.02 {
				G = G + p
			} // else if f > 941*.99 && f < 941*1.01{
			// 	D = D + p
			// }

			if (f > 660) && (f < 1500) {
				//fmt.Printf("%v,%v\n", f, p)
			}
		}

		//fmt.Println(A, B, C, D, E, F, G, SUM, "\n")

		A = A * gain / SUM
		B = B * gain / SUM
		C = C * gain / SUM
		D = D * gain / SUM
		E = E * gain / SUM
		F = F * gain / SUM
		G = G * gain / SUM
		SUM = SUM
		//D = D*float64(len(a)*len(a))

		// fmt.Println(A, B, C, D, E, F, G, SUM, "\n")
		Z := ""
		if A > 1 && E > 1 {
			Z = Z + "1"
		}
		if B > 1 && E > 1 {
			Z = Z + "4"
		}
		if C > 1 && E > 1 {
			Z = Z + "7"
		}
		if D > 1 && E > 1 {
			Z = Z + "*"
		}
		if A > 1 && F > 1 {
			Z = Z + "2"
		}
		if B > 1 && F > 1 {
			Z = Z + "5"
		}
		if C > 1 && F > 1 {
			Z = Z + "8"
		}
		if D > 1 && F > 1 {
			Z = Z + "0"
		}
		if A > 1 && G > 1 {
			Z = Z + "3"
		}
		if B > 1 && G > 1 {
			Z = Z + "6"
		}
		if C > 1 && G > 1 {
			Z = Z + "9"
		}
		if D > 1 && G > 1 {
			Z = Z + "#"
		}

		if Z != Y {
			fmt.Printf("%v", Z)
			Y = Z
		}
	}

	fmt.Println("\n", time.Since(startTime))

	// for k, v := range fft.FFTReal(a){
	// 	fmt.Printf("%v, %v=%v", k, cmplx.Abs(v), k*int(sampRate)/numSamples)
	// }

	//	X := fft.FFTReal(a)

	// Print the magnitude and phase at each frequency.
	//	for i := 0; i < numSamples; i++ {
	//		r, θ := cmplx.Polar(X[i])
	//		θ *= 360.0 / (2 * math.Pi)
	//		if dsputils.Float64Equal(r, 0) {
	//			θ = 0 // (When the magnitude is close to 0, the angle is meaningless)
	//		}
	//		fmt.Printf("X(%d) = %.1f ∠ %.1f°\n", i, r, θ)
	//	}
}
