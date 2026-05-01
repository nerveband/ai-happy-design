package contract

import "testing"

func TestDiscoveryParity(t *testing.T) {
	findings := DiscoveryFindings()
	if len(findings) == 0 {
		return
	}
	for _, finding := range findings {
		t.Errorf("%s %s %s", finding.Code, finding.Command, finding.Detail)
	}
}
