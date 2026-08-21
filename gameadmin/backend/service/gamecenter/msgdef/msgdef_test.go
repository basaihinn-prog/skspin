package msgdef

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValid(t *testing.T) {
	ps := ModifyGoldPs{
		Change: -1,
		Reason: ReasonBet,
	}
	assert.True(t, ps.Valid())

	ps = ModifyGoldPs{
		Change: 1,
		Reason: ReasonBet,
	}
	assert.False(t, ps.Valid())

	ps = ModifyGoldPs{
		Change: 1,
		Reason: ReasonWin,
	}
	assert.True(t, ps.Valid())

	ps = ModifyGoldPs{
		Change: -1,
		Reason: ReasonWin,
	}
	assert.False(t, ps.Valid())

	ps = ModifyGoldPs{
		Change: 1,
		Reason: ReasonRefund,
	}
	assert.True(t, ps.Valid())

	ps = ModifyGoldPs{
		Change: -1,
		Reason: ReasonRefund,
	}
	assert.False(t, ps.Valid())

	ps = ModifyGoldPs{
		Change: math.MaxInt64,
		// Reason: ReasonRefund,
	}
	assert.False(t, ps.Valid())

	ps = ModifyGoldPs{
		Change: math.MinInt64,
		// Reason: ReasonRefund,
	}
	assert.False(t, ps.Valid())

	ps = ModifyGoldPs{
		Change: math.MinInt32,
		// Reason: ReasonRefund,
	}
	assert.True(t, ps.Valid())

}
