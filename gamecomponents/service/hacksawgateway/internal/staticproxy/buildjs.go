package staticproxy

import (
	"bytes"
	"game/service/hacksawgateway/internal/gamedata"
)

func APIHost(content []byte) []byte {
	HOST2 := gamedata.Get().MyHost
	content = bytes.ReplaceAll(content, []byte("{{APIHOST}}"), []byte(HOST2))
	return content
}

func htmlScript(content []byte) []byte {
	content = bytes.ReplaceAll(content, []byte(`  <script>
        window.oncontextmenu = function (e) {
            e.stopPropagation();
            e.preventDefault();
            return false;
        };
    </script>`), []byte(`  <script>
        window.oncontextmenu = function (e) {
            e.stopPropagation();
            e.preventDefault();
            return false;
        };
    (function () {
    var getDeviceType = function () {
        var userAgent = navigator.userAgent || navigator.vendor || window.opera;
        // iOS detection
        if (/iPad|iPhone|iPod/.test(userAgent) && !window.MSStream) {
            return "mobile";
        }
        // Android detection
        else if (/android/i.test(userAgent)) {
            return "mobile";
        }
        // Windows detection
        else if (/windows/i.test(userAgent)) {
            return "desktop";
        }
        // Mac detection
        else if (/macintosh|mac os x/i.test(userAgent)) {
            return "desktop";
        } else {
            return "desktop";
        }
    };
    window["device_type"] = getDeviceType();
})()</script>`))
	return content
}
func JsScript(content []byte) []byte {
	content = bytes.ReplaceAll(content, []byte(`t.channel:"desktop"`), []byte(`t.channel:window["device_type"]`))
	return content
}
