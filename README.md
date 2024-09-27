# Go Dynamics Signal Analysis Library

This package is a work in progress and is not yet ready for production use.

This Go package provides a small set of simple functions for analyzing and generating signal data, particularly focused on Root Mean Square (RMS) and Zero Crossing Rate (ZCR) calculations.

## Features

- Calculate Root Mean Square (RMS) of signal data
- Calculate Zero Crossing Rate (ZCR) and Negative Zero Crossing Rate (NZCR)
- Generate sine wave data
- Analyze signal data for RMS and NZCR
- Utility function to keep a specific duration of recent data

## Installation

To use this package in your Go project, you can install it using:

```
go get github.com/tom-wray/go-dynamics@latest
```

## Usage

Here's an example of how to use the package:

```go
package main

import (
	"fmt"

	"github.com/tom-wray/go-dynamics"
)

func main() {
	// Generate a sine wave with a frequency of 1000 Hz, amplitude of 1.0, duration of 1 second, and sample rate of 44100 Hz
	sineWave := dynamics.GenerateSineWave(1000, 1.0, 1.0, 44100)

	// Calculate the RMS and NZCR of the sine wave
	rms, zcr := dynamics.Analyze(sineWave)

	fmt.Printf("RMS: %f\n", rms)
	fmt.Printf("NZCR: %f\n", zcr)
}
```

## Contributing

Contributions are welcome! Please feel free to submit a PR.

## License

This project is licensed under the Unlicense. See the LICENSE file for details.
