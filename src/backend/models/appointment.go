package models

func (a Appointment) GetDurationIn30MinInterval() float64 {
	return float64(a.Duration / 60 / 30)
}
