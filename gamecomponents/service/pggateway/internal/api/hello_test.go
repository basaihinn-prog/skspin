package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestHell(t *testing.T) {
	fmt.Println(time.Now().UnixMilli())

	fmt.Println(time.UnixMilli(1710864000000))
	fmt.Println(time.UnixMilli(1710950399999))

	ul := `https://public.pgsoft-games.com/history/redirect.html?t=026d748d-097d-4dc8-9a77-5bcdc7028be1&psid=1775083500903469568&sid=1775083905272119296&gid=1489936&btt=1&type=bo&lang=zh`

	resp, _ := http.Get(ul)
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}
