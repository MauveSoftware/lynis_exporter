// SPDX-FileCopyrightText: (c) Mauve Mailorder Software GmbH & Co. KG, 2020. Licensed under [Apache 2.0](LICENSE) license.
//
// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

const version string = "0.1.0"

var (
	showVersion      = flag.Bool("version", false, "Print version information.")
	listenAddress    = flag.String("web.listen-address", ":9730", "Address on which to expose metrics and web interface.")
	metricsPath      = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	configFile       = flag.String("config.path", "config.yml", "Path to config file")
	tlsEnabled       = flag.Bool("tls.enabled", false, "Enables TLS")
	tlsCertChainPath = flag.String("tls.cert-file", "", "Path to TLS cert file")
	tlsKeyPath       = flag.String("tls.key-file", "", "Path to TLS key file")
	cfg              *Config
)

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	startServer()
}

func printVersion() {
	fmt.Println("lynis_exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): Daniel Czerwonk")
	fmt.Println("Copyright: 2020, Mauve Mailorder Software GmbH & Co. KG, Licensed under Apache 2.0")
	fmt.Println("Metric exporter for Lynis audit results")
}

func startServer() {
	logrus.Infof("Starting Lynis exporter (Version: %s)", version)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Lynis Exporter (Version ` + version + `)</title></head>
			<body>
			<h1>Lynis Exporter by Mauve Mailorder Software</h1>
			<h2>Metrics</h2>
			<p><a href="/metrics">here</a></p>
			<h2>More information</h2>
			<p><a href="https://github.com/MauveSoftware/lynis_exporter">github.com/MauveSoftware/lynis_exporter</a></p>
			</body>
			</html>`))
	})

	var err error
	cfg, err = loadConfigFromFile(*configFile)
	if err != nil {
		logrus.Fatal(err)
	}

	http.HandleFunc(*metricsPath, errorHandler(handleMetricsRequest))

	logrus.Infof("Listening for %s on %s (TLS: %v)", *metricsPath, *listenAddress, *tlsEnabled)
	if *tlsEnabled {
		logrus.Fatal(http.ListenAndServeTLS(*listenAddress, *tlsCertChainPath, *tlsKeyPath, nil))
		return
	}

	logrus.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
			logrus.Errorln(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) error {
	reg := prometheus.NewRegistry()

	c := newCollector(cfg)
	err := c.collect()
	if err != nil {
		return err
	}

	reg.MustRegister(c)

	l := logrus.New()
	l.Level = logrus.ErrorLevel

	promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		ErrorLog:      l,
		ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(w, r)

	return nil
}
