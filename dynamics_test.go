package dynamics

import (
	"fmt"
	"math"
	"testing"
)

// TESTS

func TestGenerateSineWave(t *testing.T) {
	frequency := 100.0
	amplitude := 5.0
	duration := 5.0
	sampleRate := 2000
	data := GenerateSineWave(frequency, amplitude, duration, sampleRate)

	// Check if the number of samples is correct
	expectedSamples := int(duration * float64(sampleRate))
	if len(data) != expectedSamples {
		t.Errorf("Expected %d samples, got %d", expectedSamples, len(data))
	}

	// Check if the first sample is correct (should be 0 for a sine wave)
	if math.Abs(data[0].Value) > 1e-10 {
		t.Errorf("First sample should be close to 0, got %f", data[0].Value)
	}

	// Check if the maximum amplitude is correct
	maxAmplitude := 0.0
	for _, sample := range data {
		if math.Abs(sample.Value) > maxAmplitude {
			maxAmplitude = math.Abs(sample.Value)
		}
	}
	if math.Abs(maxAmplitude-amplitude) > 1e-10 {
		t.Errorf("Expected max amplitude %f, got %f", amplitude, maxAmplitude)
	}

	// Check if the frequency is correct using zero-crossings
	zeroCrossings := 0
	for i := 1; i < len(data); i++ {
		if data[i-1].Value <= 0 && data[i].Value > 0 {
			zeroCrossings++
		}
	}
	measuredFrequency := float64(zeroCrossings) / duration
	if math.Abs(measuredFrequency-frequency) > 1.0 {
		t.Errorf("Measured frequency %f Hz is not close to expected frequency %f Hz", measuredFrequency, frequency)
	}
}

func TestRMS(t *testing.T) {
	// Generate sample data
	frequency := 200.0
	data := GenerateSineWave(frequency, 1, 5, 2000)

	// Run the test with tolerance
	result := RMS(data, frequency)
	expected := 0.7071
	tolerance := 0.0001

	if diff := math.Abs(result - expected); diff > tolerance {
		t.Errorf("RMS returned %f, expected %f (difference: %f)", result, expected, diff)
	}
}

func TestZeroCrossingRate(t *testing.T) {
	// Generate sample data
	data := GenerateSineWave(100, 1, 1, 1000)

	// Run the test
	result := ZeroCrossingRate(data)
	expected := 200.0
	tolerance := 1.0

	if diff := math.Abs(result - expected); diff > tolerance {
		t.Errorf("ZeroCrossingRate returned %f, expected %f (difference: %f)", result, expected, diff)
	}
}

func TestNegativeZeroCrossingRate(t *testing.T) {
	// Generate sample data
	data := GenerateSineWave(440, 1, 1, 1000)

	// Run the test
	result := NegativeZeroCrossingRate(data)
	expected := 440.0
	tolerance := 1.0

	if diff := math.Abs(result - expected); diff > tolerance {
		t.Errorf("NegativeZeroCrossingRate returned %f, expected %f (difference: %f)", result, expected, diff)
	}
}

func TestAnalyze(t *testing.T) {
	// Generate sample data
	data := GenerateSineWave(440, 1, 1, 1000)

	// Run the test
	rms, zcr := Analyze(data)
	expectedRMS := 0.7071
	toleranceRMS := 0.0001
	expectedZCR := 440.0
	toleranceZCR := 1.0

	if diff := math.Abs(rms - expectedRMS); diff > toleranceRMS {
		t.Errorf("Analyze RMS returned %f, expected %f (difference: %f)", rms, expectedRMS, diff)
	}

	if diff := math.Abs(zcr - expectedZCR); diff > toleranceZCR {
		t.Errorf("Analyze ZCR returned %f, expected %f (difference: %f)", zcr, expectedZCR, diff)
	}
}

func TestAnalyzeMultiChannel(t *testing.T) {
	// Generate sample data
	channel1 := GenerateSineWave(440, 1, 1, 2000)
	channel2 := GenerateSineWave(150, 2, 1, 2000)
	data := make([]MultiChannelSample, len(channel1))
	for i := range channel1 {
		data[i] = MultiChannelSample{
			Time:  channel1[i].Time,
			Value: []float64{channel1[i].Value, channel2[i].Value},
		}
	}

	// Run the test
	rms, zcr := AnalyzeMultiChannel(data)

	fmt.Printf("len(rms): %d, len(zcr): %d\n", len(rms), len(zcr))

	expectedRMS := []float64{0.7071, 1.4144}
	expectedZCR := []float64{440.0, 150.0}
	toleranceRMS := 0.0001
	toleranceZCR := 1.0

	for i := range rms {
		if diff := math.Abs(rms[i] - expectedRMS[i]); diff > toleranceRMS {
			t.Errorf("Analyze RMS returned %f, expected %f (difference: %f)", rms[i], expectedRMS[i], diff)
		}

		if diff := math.Abs(zcr[i] - expectedZCR[i]); diff > toleranceZCR {
			t.Errorf("Analyze ZCR returned %f, expected %f (difference: %f)", zcr[i], expectedZCR[i], diff)
		}
	}
}

// BENCHMARKS

func BenchmarkGenerateSineWave(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateSineWave(440, 1, 2, 1000)
	}
}

func BenchmarkRMS(b *testing.B) {
	// Generate sample data
	frequency := 200.0
	data := GenerateSineWave(frequency, 1, 5, 2000)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RMS(data, frequency)
	}
}

func BenchmarkZeroCrossingRate(b *testing.B) {
	// Generate sample data
	data := GenerateSineWave(440, 1, 1, 1000)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ZeroCrossingRate(data)
	}
}

func BenchmarkNegativeZeroCrossingRate(b *testing.B) {
	// Generate sample data
	data := GenerateSineWave(440, 1, 1, 1000)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NegativeZeroCrossingRate(data)
	}
}

func BenchmarkAnalyze(b *testing.B) {
	// Generate sample data
	data := GenerateSineWave(440, 1, 1, 1000)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Analyze(data)
	}
}

// a function that has a ticker every 1ms and adds a sample to the circular buffer, then every 100ms it prints the RMS and NZCR of the buffer
func BenchmarkCircularBuffer(b *testing.B) {
	sineWave := GenerateSineWave(440, 1, 1, 1000)
	// Create a new CircularBuffer with a size of 1000
	cb := NewCircularBuffer(1000)

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cb.Update(SingleChannelSample{Time: float64(i), Value: sineWave[i%len(sineWave)].Value})
		if i%100 == 0 {
			rms, zcr := cb.AnalyzeBuffer()
			fmt.Printf("Length: %d, RMS: %f, NZCR: %f\n", cb.count, rms, zcr)
		}
	}
}

func BenchmarkSlice(b *testing.B) {
	sineWave := GenerateSineWave(440, 1, 1, 1000)
	data := []SingleChannelSample{}

	// Run the benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data = append(data, SingleChannelSample{Time: float64(i), Value: sineWave[i%len(sineWave)].Value})
		if i%100 == 0 {
			if len(data) > 1000 {
				// Keep only the last 1000 samples
				data = data[len(data)-1000:]
			}
			rms, zcr := Analyze(data)
			fmt.Printf("Length: %d, RMS: %f, NZCR: %f\n", len(data), rms, zcr)
		}
	}
}
