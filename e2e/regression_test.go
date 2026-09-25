//go:build e2e

package e2e

import (
	"testing"
)

// TestRegressionDeployAndFundFixture captures a regression fixture for deploying and funding any contract fixture via the shared harness helper.
func TestRegressionDeployAndFundFixture(t *testing.T) {
	h := newHarness(t)
	deployer := h.newAccount(t, "regression_deployer")
	address := h.DeployAndFundFixture(t, deployer, "threshold_account", scAddressVecVal(t, []string{deployer.Address()}), scU32Val(1))
	if address == "" {
		t.Fatal("expected DeployAndFundFixture to return a valid deployed contract address")
	}
	t.Logf("regression fixture deployed and funded successfully at %s", address)
}
