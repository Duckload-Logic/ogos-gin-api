package constants

type AppointmentTypeID int

const (
	InitialInterviewAppointmentID AppointmentTypeID = iota + 1
	MentalHealthConsultationAppointmentID
	CareerGuidanceAppointmentID
	FollowUpAppointmentID
)

type AppointmentType struct {
	ID   AppointmentTypeID
	Name string
}

var AppointmentTypes = map[AppointmentTypeID]AppointmentType{
	InitialInterviewAppointmentID:         {ID: InitialInterviewAppointmentID, Name: "Initial Interview"},
	MentalHealthConsultationAppointmentID: {ID: MentalHealthConsultationAppointmentID, Name: "Mental Health Consultation"},
	CareerGuidanceAppointmentID:           {ID: CareerGuidanceAppointmentID, Name: "Career Guidance"},
	FollowUpAppointmentID:                 {ID: FollowUpAppointmentID, Name: "Follow-up"},
}
