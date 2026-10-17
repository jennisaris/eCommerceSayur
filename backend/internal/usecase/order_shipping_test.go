package usecase_test

import (
	"ecommerce-backend/internal/usecase"
	"testing"
)

func TestCalculateShippingFee(t *testing.T) {
	tests := []struct {
		name        string
		weightGrams int
		expectedFee int64
	}{
		{
			name:        "0 grams (base fee applies)",
			weightGrams: 0,
			expectedFee: 15000,
		},
		{
			name:        "500 grams (under 1 kg base fee)",
			weightGrams: 500,
			expectedFee: 15000,
		},
		{
			name:        "1000 grams exact (1 kg boundary)",
			weightGrams: 1000,
			expectedFee: 15000,
		},
		{
			name:        "1001 grams (triggers 1 additional kg fee: 15000 + 8000)",
			weightGrams: 1001,
			expectedFee: 23000,
		},
		{
			name:        "1999 grams (under 2 kg: 15000 + 8000)",
			weightGrams: 1999,
			expectedFee: 23000,
		},
		{
			name:        "2000 grams exact (15000 + 8000)",
			weightGrams: 2000,
			expectedFee: 23000,
		},
		{
			name:        "2001 grams (triggers 2nd additional kg fee: 15000 + 16000)",
			weightGrams: 2001,
			expectedFee: 31000,
		},
		{
			name:        "5000 grams (4 additional kgs: 15000 + 32000)",
			weightGrams: 5000,
			expectedFee: 47000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := usecase.CalculateShippingFee(tt.weightGrams)
			if got != tt.expectedFee {
				t.Errorf("CalculateShippingFee(%d) = %d; expected %d", tt.weightGrams, got, tt.expectedFee)
			}
		})
	}
}
