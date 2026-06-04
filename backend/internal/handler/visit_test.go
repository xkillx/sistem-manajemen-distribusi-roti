package handler

import "testing"

func TestHaversine(t *testing.T) {
	distance := haversine(-6.2088, 106.8456, -6.9175, 107.6191)
	if distance < 100000 || distance > 200000 {
		t.Fatalf("Expected ~140km, got %d meters", int(distance))
	}

	distance = haversine(-6.2, 106.8, -6.2, 106.8)
	if distance != 0 {
		t.Fatalf("Same point should be 0, got %f", distance)
	}
}
