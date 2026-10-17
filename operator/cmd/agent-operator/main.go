// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

// The agent operator applies its Agent CRD on boot, then reconciles every
// Agent into a Deployment and a Service (applications srd009).
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/Nokia-Bell-Labs/declarative-agents/operator/api/v1alpha1"
	"github.com/Nokia-Bell-Labs/declarative-agents/operator/internal/controller"
	"github.com/Nokia-Bell-Labs/declarative-agents/operator/internal/crd"
)

// crdEstablishWait bounds the wait for the API server to serve the Agent
// kind after the CRD is applied, so the controller's first watch finds it.
const crdEstablishWait = 60 * time.Second

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "agent-operator:", err)
		os.Exit(1)
	}
}

func run() error {
	metricsAddr := flag.String("metrics-bind-address", ":8080", "address the metrics endpoint binds to")
	probeAddr := flag.String("health-probe-bind-address", ":8081", "address the health probes bind to")
	flag.Parse()
	ctrl.SetLogger(zap.New())
	scheme := runtime.NewScheme()
	for _, add := range []func(*runtime.Scheme) error{clientgoscheme.AddToScheme, apiextensionsv1.AddToScheme, v1alpha1.AddToScheme} {
		if err := add(scheme); err != nil {
			return err
		}
	}
	cfg, err := ctrl.GetConfig()
	if err != nil {
		return fmt.Errorf("kubernetes config: %w", err)
	}
	boot, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), crdEstablishWait)
	defer cancel()
	if err := crd.Apply(ctx, boot); err != nil {
		return err // the CRD ships with the binary (srd009 R1.3)
	}
	if err := crd.WaitEstablished(ctx, boot); err != nil {
		return err
	}
	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsserver.Options{BindAddress: *metricsAddr},
		HealthProbeBindAddress: *probeAddr,
	})
	if err != nil {
		return fmt.Errorf("manager: %w", err)
	}
	if err := (&controller.AgentReconciler{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("controller: %w", err)
	}
	// The chart probes liveness on /healthz and readiness on /readyz; each
	// path serves only the checks registered for it.
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return err
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return err
	}
	return mgr.Start(ctrl.SetupSignalHandler())
}
