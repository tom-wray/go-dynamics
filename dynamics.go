package dynamics

import (
	"math"
)

// Sample represents a single sample of data with a time and value.
type Sample struct {
	Time  float64
	Value float64
}

// Analyze calculates the Root Mean Square (RMS) and Negative Zero Crossing Rate (NZCR) of the given data.
//
// Parameters:
//   - data: A slice of Sample structs containing time and value data
//
// Returns:
//   - rms: The calculated Root Mean Square value
//   - zcr: The calculated Negative Zero Crossing Rate
func Analyze(data []Sample) (rms float64, zcr float64) {
	zcr = NegativeZeroCrossingRate(data)
	rms = RMS(data, zcr)
	return
}

// RMS calculates the Root Mean Square value of the given data.
//
// Parameters:
//   - data: A slice of Sample structs containing time and value data
//   - frequency: The frequency of the signal
//
// Returns:
//   - float64: The calculated Root Mean Square value
func RMS(data []Sample, frequency float64) float64 {
	if len(data) == 0 {
		return 0
	}
	if frequency == 0 {
		return 0
	}

	period := 1 / frequency

	duration := data[len(data)-1].Time - data[0].Time
	wholeCycles := math.Floor(duration / period)

	if wholeCycles < 1 {
		return calculateRMS(data)
	}

	// get last 1000 whole cycles, or x whole cycles if less than 1000
	cyclesToUse := math.Min(wholeCycles, 1000)

	// get the data from the start time to the end
	data = KeepXSecondsOfData(data, cyclesToUse*period)

	// calculate RMS
	return calculateRMS(data)
}

// calculateRMS calculates the Root Mean Square value of the given data.
//
// Parameters:
//   - data: A slice of Sample structs containing time and value data
//
// Returns:
//   - float64: The calculated Root Mean Square value
func calculateRMS(data []Sample) float64 {
	if len(data) == 0 {
		return 0
	}

	return calculateRMSAverage(data)
}

// calculateRMSAverage calculates the Root Mean Square value using the average method.
//
// Parameters:
//   - data: A slice of Sample structs containing time and value data
//
// Returns:
//   - float64: The calculated Root Mean Square value
func calculateRMSAverage(data []Sample) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value.Value * value.Value
	}
	mean := sum / float64(len(data))
	return math.Sqrt(mean)
}

// calculateRMSPeak calculates the Root Mean Square value using the peak method.
//
// Parameters:
//   - data: A slice of Sample structs containing time and value data
//
// Returns:
//   - float64: The calculated Root Mean Square value
func calculateRMSPeak(data []Sample) float64 {
	peak := 0.0
	for _, value := range data {
		absValue := math.Abs(value.Value)
		if absValue > peak {
			peak = absValue
		}
	}
	return peak / math.Sqrt(2)
}

// ZeroCrossingRate calculates the Zero Crossing Rate of the given data.
//
// Parameters:
//   - data: A slice of Sample structs containing time and value data
//
// Returns:
//   - float64: The calculated Zero Crossing Rate
func ZeroCrossingRate(data []Sample) float64 {
	if len(data) == 0 {
		return 0
	}

	crossings := 0
	for i := 1; i < len(data); i++ {
		if (data[i-1].Value >= 0 && data[i].Value < 0) || (data[i-1].Value <= 0 && data[i].Value > 0) {
			crossings++
		}
	}

	duration := data[len(data)-1].Time - data[0].Time
	return float64(crossings) / duration
}

// NegativeZeroCrossingRate calculates the Negative Zero Crossing Rate of the given data.
//
// Parameters:
//   - data: A slice of Sample structs containing time and value data
//
// Returns:
//   - float64: The calculated Negative Zero Crossing Rate
func NegativeZeroCrossingRate(data []Sample) float64 {
	if len(data) == 0 {
		return 0
	}

	crossings := 0
	for i := 1; i < len(data); i++ {
		// Only count crossings from positive to negative
		if data[i-1].Value >= 0 && data[i].Value < 0 {
			crossings++
		}
	}

	duration := data[len(data)-1].Time - data[0].Time
	return float64(crossings) / duration
}

// GenerateSineWave generates a sine wave with the specified parameters.
//
// Parameters:
//   - frequency: The frequency of the sine wave
//   - amplitude: The amplitude of the sine wave
//   - duration: The duration of the generated wave in seconds
//   - sampleRate: The number of samples per second
//
// Returns:
//   - []Sample: A slice of Sample structs representing the generated sine wave
func GenerateSineWave(frequency, amplitude, duration float64, sampleRate int) []Sample {
	samples := int(duration * float64(sampleRate))
	data := make([]Sample, samples)

	// Constants
	angularFrequency := 2 * math.Pi * frequency
	timeStep := 1.0 / float64(sampleRate)

	// Initialize first two samples
	data[0] = Sample{Time: 0, Value: 0}
	if samples > 1 {
		data[1] = Sample{Time: timeStep, Value: amplitude * math.Sin(angularFrequency*timeStep)}
	}

	// Recurrence coefficients
	c := 2 * math.Cos(angularFrequency*timeStep)

	// Generate sine wave using recurrence relation
	for i := 2; i < samples; i++ {
		t := float64(i) * timeStep
		// Recurrence relation: y[n] = c * y[n-1] - y[n-2]
		value := c*data[i-1].Value - data[i-2].Value
		data[i] = Sample{Time: t, Value: value}
	}

	return data
}

// KeepXSecondsOfData keeps the last X seconds of data from the given slice.
//
// Parameters:
//   - fastDataArray: A slice of Sample structs containing time and value data
//   - seconds: The number of seconds of data to keep
//
// Returns:
//   - []Sample: A slice of Sample structs containing the last X seconds of data
func KeepXSecondsOfData(fastDataArray []Sample, seconds float64) []Sample {
	if len(fastDataArray) == 0 {
		return fastDataArray
	}

	// cutoff time is x seconds ago
	cutoffTime := fastDataArray[len(fastDataArray)-1].Time - seconds

	// Linear search from the start, assuming cutoff is near the beginning
	for i, data := range fastDataArray {
		if data.Time >= cutoffTime {
			return fastDataArray[i:]
		}
	}

	// if the cutoff is not found, return an empty array
	return []Sample{}
}
