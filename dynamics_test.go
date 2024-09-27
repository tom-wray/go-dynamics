package dynamics

import (
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