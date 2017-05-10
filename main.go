package main

import (
	"fmt"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/spectral"
	"github.com/mjibson/go-dsp/wav"
	"os"
)

func main() {
	var err error
	file, err := os.Open("dodo_song.wav")
	fmt.Printf("err: %v\n", err)
	fmt.Println(fft.FFTReal([]float64{1, 2, 3}))
	wavReader, err := wav.New(file)
	fmt.Printf("err: %v\n", err)
	pwelchOpt := &spectral.PwelchOptions{
		Scale_off: true,
	}
	var silence bool = true
	clips := [][]float64{}
	currentClip := -1
	for {
		samples32, err := wavReader.ReadFloats(44100 / 100)
		if err != nil {
			break
		}
		samples64 := make([]float64, len(samples32))
		for i, s := range samples32 {
			samples64[i] = float64(s)
		}
		p, _ := spectral.Pwelch(samples64, 44100, pwelchOpt)
		var total float64
		for _, pv := range p {
			total += pv
		}
		averagePower := total / float64(len(p))
		fmt.Printf("average p: %v\n", averagePower)
		if averagePower > .5 {
			if silence {
				clips = append(clips, []float64{})
				currentClip++
			}
			clips[currentClip] = append(clips[currentClip], samples64...)
			silence = false
		} else {
			silence = true
		}
		fmt.Printf("clips: %v\n", len(clips))
	}
	pwelchOpt.Scale_off = false
	for i, c := range clips {
		p, f := spectral.Pwelch(clip, 44100, pwelchOpt)
	}
	fmt.Printf("wavReader: %#v\n", wavReader)
}
