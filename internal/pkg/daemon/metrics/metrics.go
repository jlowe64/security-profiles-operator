/*
Copyright 2021 The Kubernetes Authors.

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

package metrics

import (
	"fmt"
	"net/http"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Metrics proxy required permissions
// +kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews;subjectaccessreviews,verbs=create
// +kubebuilder:rbac:groups=authorization.k8s.io,resources=subjectaccessreviews,verbs=create

const (
	metricNamespace               = "security_profiles_operator"
	metricNameSeccompProfile      = "seccomp_profile_total"
	metricNameSelinuxProfile      = "selinux_profile_total"
	metricNameSeccompProfileError = "seccomp_profile_error_total"
	metricNameSelinuxProfileError = "selinux_profile_error_total"
	metricLabelValueProfileUpdate = "update"
	metricLabelValueProfileDelete = "delete"
	metricLabelOperation          = "operation"
	metricsLabelReason            = "reason"

	// HandlerPath is the default path for serving metrics.
	HandlerPath = "/metrics-spod"
)

// Metrics is the main structure of this package.
type Metrics struct {
	impl                      impl
	log                       logr.Logger
	metricSeccompProfile      *prometheus.CounterVec
	metricSeccompProfileError *prometheus.CounterVec
	metricSelinuxProfile      *prometheus.CounterVec
	metricSelinuxProfileError *prometheus.CounterVec
}

// New returns a new Metrics instance.
func New() *Metrics {
	return &Metrics{
		impl: &defaultImpl{},
		log:  ctrl.Log.WithName("metrics"),
		metricSeccompProfile: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:      metricNameSeccompProfile,
				Namespace: metricNamespace,
				Help:      "Counter about seccomp profile operations.",
			},
			[]string{metricLabelOperation},
		),
		metricSeccompProfileError: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:      metricNameSeccompProfileError,
				Namespace: metricNamespace,
				Help:      "Counter about seccomp profile errors.",
			},
			[]string{metricsLabelReason},
		),
		metricSelinuxProfile: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:      metricNameSelinuxProfile,
				Namespace: metricNamespace,
				Help:      "Counter about selinux profile operations.",
			},
			[]string{metricLabelOperation},
		),
		metricSelinuxProfileError: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:      metricNameSelinuxProfileError,
				Namespace: metricNamespace,
				Help:      "Counter about selinux profile errors.",
			},
			[]string{metricsLabelReason},
		),
	}
}

// Register iterates over all available metrics and registers them.
func (m *Metrics) Register() error {
	for name, collector := range map[string]prometheus.Collector{
		metricNameSeccompProfile:      m.metricSeccompProfile,
		metricNameSeccompProfileError: m.metricSeccompProfileError,
		metricNameSelinuxProfile:      m.metricSelinuxProfile,
		metricNameSelinuxProfileError: m.metricSelinuxProfileError,
	} {
		m.log.Info(fmt.Sprintf("Registering metric: %s", name))
		if err := m.impl.Register(collector); err != nil {
			return errors.Wrapf(err, "register collector for %s metric", name)
		}
	}
	return nil
}

// Handler creates an HTTP handler for the metrics.
func (m *Metrics) Handler() http.Handler {
	handler := &http.ServeMux{}
	handler.Handle(HandlerPath, promhttp.Handler())
	return handler
}

// IncSeccompProfileUpdate increments the seccomp profile update counter.
func (m *Metrics) IncSeccompProfileUpdate() {
	m.metricSeccompProfile.
		WithLabelValues(metricLabelValueProfileUpdate).Inc()
}

// IncSeccompProfileDelete increments the seccomp profile deletion counter.
func (m *Metrics) IncSeccompProfileDelete() {
	m.metricSeccompProfile.
		WithLabelValues(metricLabelValueProfileDelete).Inc()
}

// IncSeccompProfileError increments the seccomp profile error counter for the
// provided reason.
func (m *Metrics) IncSeccompProfileError(reason event.Reason) {
	m.metricSeccompProfileError.WithLabelValues(string(reason)).Inc()
}

// IncSelinuxProfileUpdate increments the selinux profile update counter.
func (m *Metrics) IncSelinuxProfileUpdate() {
	m.metricSelinuxProfile.
		WithLabelValues(metricLabelValueProfileUpdate).Inc()
}

// IncSelinuxProfileDelete increments the selinux profile deletion counter.
func (m *Metrics) IncSelinuxProfileDelete() {
	m.metricSelinuxProfile.
		WithLabelValues(metricLabelValueProfileDelete).Inc()
}

// IncSelinuxProfileError increments the selinux profile error counter for the
// provided reason.
func (m *Metrics) IncSelinuxProfileError(reason event.Reason) {
	m.metricSelinuxProfileError.WithLabelValues(string(reason)).Inc()
}
