package dynamics

import (
	"fmt"
	"math"
)

// CircularBuffer represents a circular buffer for storing SingleChannelSample data.
type CircularBuffer struct {
	data  []SingleChannelSample
	size  int
	head  int
	count int
}

// NewCircularBuffer creates a new CircularBuffer with the specified size.
func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		data:  make([]SingleChannelSample, size),
		size:  size,
		head:  0,
		count: 0,
	}
}

// Update adds a new sample to the circular buffer.
func (cb *CircularBuffer) Update(sample SingleChannelSample) {
	cb.data[cb.head] = sample
	cb.head = (cb.head + 1) % cb.size
	if cb.count < cb.size {
		cb.count++
	}
}

// GetData returns a slice of the data in the buffer, from oldest to newest.
func (cb *CircularBuffer) GetData() []SingleChannelSample {
	result := make([]SingleChannelSample, cb.count)
	for i := 0; i < cb.count; i++ {
		index := (cb.head - cb.count + i + cb.size) % cb.size
		result[i] = cb.data[index]
	}
	return result
}

// AnalyzeBuffer calculates the RMS and NZCR of the data stored in the circular buffer.
func (cb *CircularBuffer) AnalyzeBuffer() (rms float64, zcr float64) {
	if cb.count == 0 {
		return 0, 0
	}
	data := cb.GetData()
	return Analyze(data)
}

// Sample represents a single sample of data with a time and a generic value.
type Sample[T float64 | []float64] struct {
	Time  float64 `json:"time"`
	Value T       `json:"value"`
}

// Convenience type aliases for common use cases
type SingleChannelSample = Sample[float64]
type MultiChannelSample = Sample[[]float64]

// Analyze calculates the Root Mean Square (RMS) and Negative Zero Crossing Rate (NZCR) of the given data.
//
// Parameters:
//   - data: A slice of Sample structs containing time and value data
//
// Returns:
//   - rms: The calculated Root Mean Square value
//   - zcr: The calculated Negative Zero Crossing Rate
func Analyze(data []SingleChannelSample) (rms float64, zcr float64) {
	zcr = NegativeZeroCrossingRate(data)
	rms = RMS(data, zcr)
	return
}

// AnalyzeMultiChannel analyzes the given multi-channel data and returns the RMS and NZCR for each channel.
//
// Parameters:
//   - data: A slice of MultiChannelSample structs containing time and value data
//
// Returns:
//   - rms: A slice of float64 values representing the RMS for each channel
//   - zcr: A slice of float64 values representing the NZCR for each channel
func AnalyzeMultiChannel(data []MultiChannelSample) (rms []float64, zcr []float64) {
	// channel count is the length of the value array
	channelCount := len(data[0].Value)

	zcr = make([]float64, channelCount)
	rms = make([]float64, channelCount)

	fmt.Printf("channelCount: %d\n", channelCount)

	for i := range channelCount {
		singleChannelData := make([]SingleChannelSample, len(data))
		for j := range data {
			singleChannelData[j] = SingleChannelSample{Time: data[j].Time, Value: data[j].Value[i]}
		}
		zcr[i] = NegativeZeroCrossingRate(singleChannelData)
		rms[i] = RMS(singleChannelData, zcr[i])
	}

	fmt.Printf("Length of zcr: %d, length of rms: %d\n", len(zcr), len(rms))
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
func RMS(data []SingleChannelSample, frequency float64) float64 {
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
func calculateRMS(data []SingleChannelSample) float64 {
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
func calculateRMSAverage(data []SingleChannelSample) float64 {
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
func calculateRMSPeak(data []SingleChannelSample) float64 {
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
func ZeroCrossingRate(data []SingleChannelSample) float64 {
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
func NegativeZeroCrossingRate(data []SingleChannelSample) float64 {
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
func GenerateSineWave(frequency, amplitude, duration float64, sampleRate int) []SingleChannelSample {
	samples := int(duration * float64(sampleRate))
	data := make([]SingleChannelSample, samples)

	// Constants
	angularFrequency := 2 * math.Pi * frequency
	timeStep := 1.0 / float64(sampleRate)

	// Initialize first two samples
	data[0] = SingleChannelSample{Time: 0, Value: 0}
	if samples > 1 {
		data[1] = SingleChannelSample{Time: timeStep, Value: amplitude * math.Sin(angularFrequency*timeStep)}
	}

	// Recurrence coefficients
	c := 2 * math.Cos(angularFrequency*timeStep)

	// Generate sine wave using recurrence relation
	for i := 2; i < samples; i++ {
		t := float64(i) * timeStep
		// Recurrence relation: y[n] = c * y[n-1] - y[n-2]
		value := c*data[i-1].Value - data[i-2].Value
		data[i] = SingleChannelSample{Time: t, Value: value}
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
func KeepXSecondsOfData(fastDataArray []SingleChannelSample, seconds float64) []SingleChannelSample {
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
	return []SingleChannelSample{}
}
