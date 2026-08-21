package staticproxy

import (
	"game/comm/ut"
	"log"
	"net/http"
)

var (
	jilid_rslotszs001_mux = http.NewServeMux()
	uat_wbslot_fd_mux     = http.NewServeMux()
	uat_history_api_mux   = http.NewServeMux()
)

func dispatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Token")
		w.Header().Set("Access-Control-Max-Age", "1800")
		// Access-Control-Expose-Headers: Content-Length,Access-Control-Allow-Origin

		// w.Header().Set("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin")
		w.Header().Set("Access-Control-Expose-Headers", "*")

		return
	}

	// uat-wbslot-fd-jlfafafa1.kafa010.com
	switch {
	// case strings.HasPrefix(r.Host, "jilid-rslotszs001.") ||
	// 	strings.HasPrefix(r.Host, "uat-history.") ||
	// 	strings.HasPrefix(r.Host, "uat-language-api.") ||
	// 	strings.HasPrefix(r.Host, "history.") ||
	// 	strings.HasPrefix(r.Host, "language-api.") ||
	// 	false:
	case ut.HasPrefix(r.Host, "jilid-rslotszs004.", "jilid-rslotszs003.", "history.", "uat-history.", "language-api.", "uat-language-api.", "cdn.jilid-rslotszs001.", "cdn.history.", "cdn.uat-history.", "cdn.language-api.", "cdn.uat-language-api.", "uat-web-cdn.", "cdn.uat-web-cdn.", "tadad-rslotszs001.", "cdn.tadad-rslotszs001."):
		jilid_rslotszs001_mux.ServeHTTP(w, r)
	case ut.HasPrefix(r.Host, "uat-history-api.", "history-api.", "cdn.uat-history-api.", "cdn.history-api."):
		uat_history_api_mux.ServeHTTP(w, r)
	case ut.HasPrefix(r.Host, "wbslot-fd-jlfafafa1.", "uat-wbslot-fd-jlfafafa1.", "cdn.wbslot-fd-jlfafafa1.", "cdn.uat-wbslot-fd-jlfafafa1.", "tadaapi", "cdn.tadaapi", "wbslot-fd-jlfafafa1-na", "cdn.wbslot-fd-jlfafafa1-na"):
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		uat_wbslot_fd_mux.ServeHTTP(w, r)
		// https://uat-wbslot-fd-jlfafafa1.kafa010.com/zeus/req?D=1&

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func runProxy(addr string) {
	//证书

	err := http.ListenAndServe(addr, http.HandlerFunc(dispatch))
	//err := http.ListenAndServeTLS(addr, http.HandlerFunc(dispatch))
	log.Fatalf(`ListenAndServe addr="%v", err="%v"`, addr, err.Error())
}

func StartProxy(addr string) {
	go runProxy(addr)
}

// https://jilid-rslotszs001.kafa010.com/csh/?ssoKey=19d8eede5f336e61bbaf2764fdd0b6a72ec61d19&lang=en-US&gameID=2&gs=moc.010afak.1afafaflj-df-tolsbw-tau&domain_gs=010afak&domain_platform=moc.010afak.1afafaflj-df-tolsbw-tau&be=moc.010afak.1afafaflj-df-tolsbw-tau&apiId=555&iu=true&legalLang=true
