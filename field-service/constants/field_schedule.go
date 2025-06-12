package constants

type FieldScheduleStatusName string
type FieldScheduleStatus int

const (
	Available FieldScheduleStatus = 100
	Booked    FieldScheduleStatus = 200

	AvailableString FieldScheduleStatusName = "Available"
	BookedString    FieldScheduleStatusName = "Booked"
)

var mapFieldScheduleStatusIntToString = map[FieldScheduleStatus]FieldScheduleStatusName{
	Available: AvailableString,
	Booked:    BookedString,
}

var mapFileScheduleStatusStringToInt = map[FieldScheduleStatusName]FieldScheduleStatus{
	AvailableString: Available,
	BookedString:    Booked,
}

func (f FieldScheduleStatus) GetStatusString() FieldScheduleStatusName {
	return mapFieldScheduleStatusIntToString[f]
}

func (f FieldScheduleStatusName) GetStatusInt() FieldScheduleStatus {
	return mapFileScheduleStatusStringToInt[f]
}
