package prerun

import (
	"testing"
)

func TestValidateGraphPreRunE(t *testing.T) {
	testCases := []struct {
		name         string
		outputFormat string
		wantErr      bool
	}{
		{
			name:         "Invalid output format",
			outputFormat: "invalid",
			wantErr:      true,
		},
		{
			name:         "Valid output format",
			outputFormat: "dot",
			wantErr:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateGraphPreRunE(tc.outputFormat)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateGraphPreRunE() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
