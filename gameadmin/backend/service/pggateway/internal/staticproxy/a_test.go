package staticproxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	assert.True(t, sharedIndexExp.MatchString("/shared/3c4695a542/index-100.json"))
	assert.True(t, sharedIndexExp.MatchString("/shared/3c4695a542/index-100.js"))
	assert.False(t, sharedIndexExp.MatchString("/shared/3c4695a542/index-100js"))
	assert.False(t, sharedIndexExp.MatchString("/shared/3c4695a542/index-100.html"))
	assert.False(t, sharedIndexExp.MatchString("/shared/3c4695a542/someother/index-100.html"))
	assert.False(t, sharedIndexExp.MatchString("/shared/3c4695a542/index-xxx.html"))
	assert.False(t, sharedIndexExp.MatchString("/pre/shared/3c4695a542/index-123.html"))
	assert.False(t, sharedIndexExp.MatchString("/shared/3c4695a542/index-123.js?"))
}
