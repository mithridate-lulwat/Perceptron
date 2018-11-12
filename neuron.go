package main

import (
	"fmt"
	"math/rand"
	"time"
)

func max(a float64, values ...float64) float64 {
	m := a
	for _, val := range values {
		if m < val {
			m = val
		}
	}
	return m
}

type RGB struct {
	r, g, b float64
}

type Pixel struct {
	RGB
	colour string
}

func generateRGB() (float64, float64, float64) {
	r := float64(rand.Intn(256))
	g := float64(rand.Intn(256))
	b := float64(rand.Intn(256))
	if r == g || r == b || b == g {
		return generateRGB()
	}
	return r, g, b
}

func (p *Pixel) determineColour() {
	// The max is unique thanks to the generateRGB function
	m := max(p.r, p.g, p.b)
	switch m {
	case p.g:
		p.colour = "green"
	case p.b:
		p.colour = "blue"
	case p.r:
		p.colour = "red"
	default:
		fmt.Errorf("There is no colour !!!")
	}
}

type Perceptron struct {
	// A Perceptron is determined by its size, weights and activation function
	size         int
	weights      []float64
	bias         float64
	activation   func(float64) float64
	learningRate float64
}

func (p *Perceptron) updateWeights(input []float64, target float64) error {
	// Update the weigth according to input, learningRate and activation
	var newWeights []float64
	output, err := p.predict(input)
	if err != nil {
		return err
	}
	for i, w_i := range p.weights {
		newWeight := w_i + p.learningRate*(target-output)*input[i]
		newWeights = append(newWeights, newWeight)
	}
	p.bias = p.bias + p.learningRate*(target-output)
	p.weights = newWeights
	return nil
}

func (p Perceptron) predict(input []float64) (float64, error) {
	if p.size != len(input) {
		return 0, fmt.Errorf("The input hasn't the right dimension")
	}
	sum := 0.0
	sum += p.bias
	for i, w_i := range p.weights {
		sum += w_i * input[i]
	}
	output := p.activation(sum)
	return output, nil
}

func activation(f float64) float64 {
	// Implements the Heaviside activation function
	if f > 0 {
		return 1
	} else {
		return 0
	}
}

func main() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	var pixels []Pixel
	n := 100000
	// Generate the training set
	for i := 0; i < n; i++ {
		r, g, b := generateRGB()
		pixel := Pixel{RGB{r, g, b}, ""}
		pixel.determineColour()
		pixels = append(pixels, pixel)
	}
	// Create a perceptron for each colour
	defaultWeights := []float64{0.1, 0.1, 0.1}
	learningRate := 0.01
	neuronRed := Perceptron{3, defaultWeights, 0, activation, learningRate}
	neuronGreen := Perceptron{3, defaultWeights, 0, activation, learningRate}
	neuronBlue := Perceptron{3, defaultWeights, 0, activation, learningRate}
	// Training all the perceptrons
	for _, p := range pixels {
		inputVector := []float64{p.r, p.g, p.b}
		var targetRed, targetGreen, targetBlue float64
		// We set the target depending on the pixel's colour
		switch p.colour {
		case "green":
			targetGreen = 1.
		case "red":
			targetRed = 1.
		case "blue":
			targetBlue = 1.
		default:
			fmt.Errorf("This point hasn't got the right colour")
		}
		(&neuronRed).updateWeights(inputVector, targetRed)
		(&neuronGreen).updateWeights(inputVector, targetGreen)
		(&neuronBlue).updateWeights(inputVector, targetBlue)

	}
	fmt.Println(neuronRed.weights)
	fmt.Println(neuronBlue.weights)
	fmt.Println(neuronGreen.weights)
	// Generate the testing set
	var pixelsTest []Pixel
	k := 1000
	for i := 0; i < k; i++ {
		r, g, b := generateRGB()
		pixel := Pixel{RGB{r, g, b}, ""}
		pixel.determineColour()
		pixelsTest = append(pixelsTest, pixel)
	}
	// Test the Perceptron
	var actualRed, actualGreen, actualBlue int
	var pRed, pGreen, pBlue int
	var correct float64
	for _, p := range pixelsTest {
		inputVector := []float64{p.r, p.g, p.b}
		oRed, _ := neuronRed.predict(inputVector)
		oGreen, _ := neuronGreen.predict(inputVector)
		oBlue, _ := neuronBlue.predict(inputVector)
		m := max(oRed, oGreen, oBlue)
		var colourPredicted string
		switch m {
		case oRed:
			colourPredicted = "red"
			pRed++
		case oGreen:
			colourPredicted = "green"
			pGreen++
		case oBlue:
			colourPredicted = "blue"
			pBlue++
		}
		if p.colour == colourPredicted {
			correct++
		}
		switch p.colour {
		case "red":
			actualRed++
		case "green":
			actualGreen++
		case "blue":
			actualBlue++
		}

	}
	fmt.Println("Accuracy : ", correct/float64(k))
	fmt.Println("Red : ", actualRed, pRed)
	fmt.Println("Green : ", actualGreen, pGreen)
	fmt.Println("Blue : ", actualBlue, pBlue)
}
