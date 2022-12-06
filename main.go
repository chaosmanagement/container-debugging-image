package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strings"
	"time"
)

var port string

func isEnvVariableTrue(variableName string) bool {
	variableValue, variableSpecified := os.LookupEnv(variableName)

	if !variableSpecified {
		return false
	}

	if variableValue == "" || variableValue == "1" || variableValue == "true" {
		return true
	}

	return false
}

func padRight(minLength int, str string) string {
	if len(str) < minLength {
		return str + strings.Repeat(" ", minLength-len(str))
	}

	return str
}

func getHostname() string {
	serverHostname, err := os.Hostname()
	if err != nil {
		serverHostname = "unknown"
	}

	return serverHostname
}

func getLocalAddresses(ctx context.Context) []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return []string{}
	}

	localAddresses := []string{}

	for _, addr := range addrs {
		prefix := netip.MustParsePrefix(addr.String())

		if prefix.Addr().IsLoopback() || prefix.Addr().IsLinkLocalUnicast() {
			continue
		}

		ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
		defer cancel()

		revNames, err := net.DefaultResolver.LookupAddr(ctx, prefix.Addr().StringExpanded())
		if err != nil {
			fmt.Println(err)
			continue
		}

		localAddresses = append(localAddresses, fmt.Sprintf("%s %s", padRight(39, prefix.Addr().StringExpanded()), strings.Join(revNames, ", ")))

	}

	return localAddresses
}

func printKV(w http.ResponseWriter, k, v string) {
	fmt.Fprintf(w, "%s: %s\n", padRight(28, k), v)
}

func handler(w http.ResponseWriter, r *http.Request) {
	renderedSomething := false

	if isEnvVariableTrue("DEBUG_HTTP") {
		renderedSomething = true

		printKV(w, "HTTP URL", r.URL.String())
		printKV(w, "HTTP Host", r.Host)
		printKV(w, "HTTP Listen port", port)
		printKV(w, "HTTP Referer", r.Referer())
		printKV(w, "HTTP User agent", r.UserAgent())

		fmt.Fprint(w, "\n")
	}

	if isEnvVariableTrue("DEBUG_SERVER") {
		renderedSomething = true

		printKV(w, "Server hostname", getHostname())

		localAddresses := getLocalAddresses(r.Context())
		if len(localAddresses) != 0 {
			for _, addr := range localAddresses {
				printKV(w, "Server's address", addr)
			}
		}

		fmt.Fprint(w, "\n")
	}

	if isEnvVariableTrue("DEBUG_CLIENT") {
		renderedSomething = true

		clientIp, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			clientIp = "unknown"
		} else {
			ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*100)
			defer cancel()

			revNames, err := net.DefaultResolver.LookupAddr(ctx, clientIp)
			if err == nil {
				clientIp = fmt.Sprintf("%s %s", clientIp, padRight(39, strings.Join(revNames, ", ")))
			}

		}
		printKV(w, "Client's IP", clientIp)
	}

	if !renderedSomething {
		fmt.Fprintf(w, "You need to enable at least one section via the environment variables for this software to be useful!\n")
	}
}

func main() {
	somethingEnabled := false
	if isEnvVariableTrue("DEBUG_HTTP") {
		somethingEnabled = true
		fmt.Println("Enabled HTTP debug section")
	}
	if isEnvVariableTrue("DEBUG_SERVER") {
		somethingEnabled = true
		fmt.Println("Enabled server debug section")
	}
	if isEnvVariableTrue("DEBUG_CLIENT") {
		somethingEnabled = true
		fmt.Println("Enabled client debug section")
	}

	if !somethingEnabled {
		fmt.Println("You need to enable at least one section for this fotware to be useful!")
	}

	var portSpecified bool
	port, portSpecified = os.LookupEnv("HTTP_PORT")
	if !portSpecified {
		port = "8080"
	}

	fmt.Printf("Will listen on %s port\n", port)

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
