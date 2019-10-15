package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/harmony-one/harmony/core/numeric"
	"github.com/harmony-one/harmony/crypto/bls"
)

func CreateNewValidator() Validator {
	cr := CommissionRates{Rate: numeric.OneDec(), MaxRate: numeric.OneDec(), MaxChangeRate: numeric.ZeroDec()}
	c := Commission{cr, big.NewInt(300)}
	d := Description{Name: "SuperHero", Identity: "YouWillNotKnow", Website: "under_construction", Details: "N/A"}
	v := Validator{Address: common.Address{}, ValidatingPubKey: *bls.RandPrivateKey().GetPublicKey(),
		Stake: big.NewInt(500), UnbondingHeight: big.NewInt(20), MinSelfDelegation: big.NewInt(7),
		IsActive: false, Commission: c, Description: d}
	return v
}
