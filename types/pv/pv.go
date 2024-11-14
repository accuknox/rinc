package pv

import "time"

type Metrics struct {
	Timestamp time.Time
	PVs       PVs
}

type PV struct {
	PVC                string
	PVCNamespace       string
	Capacity           float64
	Used               float64
	Available          float64
	UtilizationPercent float64
}

type PVs []PV

func (pvs PVs) AppendCapacity(pvc, ns string, cap float64) PVs {
	var exists bool
	for idx, pv := range pvs {
		if pv.PVC == pvc && pv.PVCNamespace == ns {
			exists = true
			pvs[idx].Capacity = cap
		}
	}
	if !exists {
		pvs = append(pvs, PV{
			PVC:          pvc,
			PVCNamespace: ns,
			Capacity:     cap,
		})
	}
	return pvs
}

func (pvs PVs) AppendUsed(pvc, ns string, used float64) PVs {
	var exists bool
	for idx, pv := range pvs {
		if pv.PVC == pvc && pv.PVCNamespace == ns {
			exists = true
			pvs[idx].Used = used
		}
	}
	if !exists {
		pvs = append(pvs, PV{
			PVC:          pvc,
			PVCNamespace: ns,
			Used:         used,
		})
	}
	return pvs
}

func (pvs PVs) AppendAvailable(pvc, ns string, available float64) PVs {
	var exists bool
	for idx, pv := range pvs {
		if pv.PVC == pvc && pv.PVCNamespace == ns {
			exists = true
			pvs[idx].Available = available
		}
	}
	if !exists {
		pvs = append(pvs, PV{
			PVC:          pvc,
			PVCNamespace: ns,
			Available:    available,
		})
	}
	return pvs
}

func (pvs PVs) AppendUtilization(pvc, ns string, utilization float64) PVs {
	var exists bool
	for idx, pv := range pvs {
		if pv.PVC == pvc && pv.PVCNamespace == ns {
			exists = true
			pvs[idx].UtilizationPercent = utilization
		}
	}
	if !exists {
		pvs = append(pvs, PV{
			PVC:                pvc,
			PVCNamespace:       ns,
			UtilizationPercent: utilization,
		})
	}
	return pvs
}
