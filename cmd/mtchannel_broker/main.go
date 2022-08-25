/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"context"

	"knative.dev/pkg/injection/sharedmain"

	"knative.dev/eventing/pkg/reconciler/broker"
	mttrigger "knative.dev/eventing/pkg/reconciler/broker/trigger"
)

const (
	component = "mt-broker-controller"
)

func main() {
	// sets up liveness and readiness probes.
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler)
	mux.HandleFunc("/readiness", handler)

	port := os.Getenv("PROBES_PORT")
	if port == "" {
		port = "8081"
	}

	go func() {
		// start the web server on port and accept requests
		log.Printf("Readiness and health check server listening on port %s", port)
		log.Fatal(http.ListenAndServe(":"+port, mux))
	}()

	sharedmain.Main(
		component,

		broker.NewController,

		mttrigger.NewController,
	)
	broker.Tracer.Shutdown(context.Background())
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
