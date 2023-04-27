package health

import (
	"context"
	"github.com/pkg/errors"
	"service/services/logger"
	"time"
)

type Patient struct {
	Name     string
	HealthFn func() error
}

type PatientHealthStatus struct {
	Healthy   bool
	Diagnosis string
}

type SystemHealthStatus struct {
	Healthy  bool
	Patients map[string]PatientHealthStatus
}

var Patients []Patient

var SystemHealth = SystemHealthStatus{
	Healthy:  false,
	Patients: make(map[string]PatientHealthStatus),
}

var ErrHealthCheckTimedOut = errors.New("health check timedout")

func WrapHealthFuncWithTimeouts(f func() error) func() error {
	return func() (err error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		ch := make(chan bool)
		go func() {
			err = f()
			ch <- true
			close(ch)
		}()
		select {
		case <-ctx.Done():
			return ErrHealthCheckTimedOut
		case <-ch:
			return err
		}
	}
}

func NewPatient(name string, fn func() error) {
	l := logger.For("health")
	patient := Patient{name, WrapHealthFuncWithTimeouts(fn)}
	Patients = append(Patients, patient)
	if _, registeredHealthStatus := SystemHealth.Patients[name]; !registeredHealthStatus {
		patientHealthStatus := PatientHealthStatus{
			Healthy:   false,
			Diagnosis: "not yet initialized",
		}
		SystemHealth.Patients[name] = patientHealthStatus
	} else {
		l.Info("previous patient with the name %s is already registered! overwriting existing record...", name)
	}
}

func IsHealthy() SystemHealthStatus {
	allHealthy := true
	for _, p := range Patients {
		patientHealth := SystemHealth.Patients[p.Name]
		healthError := p.HealthFn()
		if healthError != nil {
			allHealthy = false
			patientHealth.Healthy = false
			patientHealth.Diagnosis = healthError.Error()
		} else {
			patientHealth.Healthy = true
			patientHealth.Diagnosis = ":)"
		}
		SystemHealth.Patients[p.Name] = patientHealth
	}
	SystemHealth.Healthy = allHealthy
	return SystemHealth
}
