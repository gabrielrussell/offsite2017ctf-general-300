package main

import (
	"fmt"
	"github.com/mjibson/go-dsp/spectral"
	"github.com/mjibson/go-dsp/wav"
	"github.com/mjibson/go-dsp/window"
	"github.com/wcharczuk/go-chart"
	"os"
)

func main() {
	var err error
	file, err := os.Open("dodo_song.wav")
	fmt.Printf("err: %v\n", err)
	wavReader, err := wav.New(file)
	fmt.Printf("err: %v\n", err)
	fmt.Printf("wavReader: %#v\n", wavReader)
	clipChan := make(chan []float64, 100)
	go func() {
		var silence bool = true
		pwelchOpt := &spectral.PwelchOptions{
			Scale_off: true,
		}
		clip := []float64{}
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
			//fmt.Printf("average p: %v\n", averagePower)
			if averagePower > .5 {
				clip = append(clip, samples64...)
				silence = false
			} else {
				if silence == false {
					clipChan <- clip
					clip = []float64{}
				}
				silence = true
			}
			//fmt.Printf("clips: %v\n", len(clips))
		}
		close(clipChan)
	}()
	pwelchOpt := &spectral.PwelchOptions{
		Scale_off: false,
		Window:    window.Bartlett,
		NFFT:      4096,
	}

	m := map[int]uint{
		1600: 5,
		2400: 4,
		3000: 3,
		3800: 2,
		4600: 1,
		5400: 0,
	}
	num := 0
	chars := []int{}
	allChars := make([]int, 64)
	for c := range clipChan {
		p, f := spectral.Pwelch(c, 44100, pwelchOpt)
		p = p[20:]
		f = f[20:]
		maxPower := float64(0)
		for i := range p {
			if p[i] > maxPower {
				maxPower = p[i]
			}
		}
		ch := 0
		for i := range p {
			p[i] /= maxPower
			if p[i] > 0.3 {
				i := (int(f[i]+50) / 100) * 100
				v, ok := m[i]
				if !ok {
					panic("unmapped frequency")
				}
				ch |= 1 << v
			}
		}
		allChars[ch]++
		chars = append(chars, ch)
		num++
	}
	fmt.Printf("%v", allChars)
	fmt.Printf("\n")
	for _, c := range chars {
		if c == 41 {
			fmt.Printf(" ")
		} else {
			fmt.Printf("%c", c+82)
		}
	}
}

func graph(fn string, x []float64, y []float64) {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true, //enables / displays the x-axis
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true, //enables / displays the y-axis
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: x,
				YValues: y,
			},
		},
	}

	img, _ := os.Create(fn)
	graph.Render(chart.PNG, img)
}
