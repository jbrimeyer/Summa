package summa

import (
	"net/http"
)

func StartHttp() {
	http.HandleFunc("/api/", handleApiRequest)
	http.Handle("/", http.FileServer(http.Dir(config.WebRoot())))

	if config.SSLEnable {
		http.ListenAndServeTLS(
			config.Listen,
			config.SSLCertFile(),
			config.SSLKeyFile(),
			nil,
		)
	} else {
		http.ListenAndServe(config.Listen, nil)
	}
}
