package errors

import (
	errorField "field-services/constants/errors/field"
	errorFieldSchedule "field-services/constants/errors/fieldschedule"
)

func ErrMapping(err error) bool {
	allErrors := make([]error, 0)
	allErrors = append(append(GeneralErrors[:], errorField.FieldErrors[:]...), errorFieldSchedule.FieldScheduleErrors[:]...)

	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true
		}
	}
	return false
}
