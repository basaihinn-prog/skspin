package staticproxy

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestProxy(t *testing.T) {
	ul := "https://demogamesfree.pragmaticplay.net/gs2c/html5Game.do?extGame=1&symbol=vs20fruitswx&gname=Sweet%20Bonanza%201000&jurisdictionID=99&lobbyUrl=https%3A%2F%2Fclienthub.pragmaticplay.com%2Fslots%2Fgame-library%2F&mgckey=stylename@generic~SESSION@b2c1e5b5-750d-4690-a6c9-6cd9b188df15"

	resp, _ := http.Get(ul)
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}
