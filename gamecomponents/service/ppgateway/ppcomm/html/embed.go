package html

import (
	_ "embed"
)

//go:embed ppReplayTemplate.html
var PPReplayTemplateHtml []byte

//go:embed ppReplayTemplateMonkey.html
var PPReplayTemplateMonkeyHtml []byte
