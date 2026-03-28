package audit

import (
	"context"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

// DispatchParams holds the parameters for the Dispatch helper.
type DispatchParams struct {
	Log          *LogParams
	Notification *NotificationParams
	Tx           datastore.DB
}

// LogParams holds the parameters for a log entry.
type LogParams struct {
	Category   string
	Action     string
	Message    string
	TargetID   structs.NullableString
	TargetType structs.NullableString
	Metadata   *LogMetadata
}

// NotificationParams holds the parameters for a notification.
type NotificationParams struct {
	ReceiverID structs.NullableString
	TargetID   structs.NullableString
	TargetType structs.NullableString
	Title      string
	Message    string
	Type       string
}

// Dispatch is a centralized helper to record logs and send notifications.
// It automatically extracts audit metadata from the context.
func Dispatch(
	ctx context.Context,
	logger Logger,
	notifier Notifier,
	params DispatchParams,
) {
	// 1. Extract Meta
	id, ip, ua, email, trace := ExtractMeta(ctx)

	// 2. Prepare and Record Log
	if logger != nil && params.Log != nil {
		entry := LogEntry{
			Category:   params.Log.Category,
			Action:     params.Log.Action,
			Message:    params.Log.Message,
			UserID:     structs.StringToNullableString(id),
			UserEmail:  structs.StringToNullableString(email),
			IPAddress:  structs.StringToNullableString(ip),
			UserAgent:  structs.StringToNullableString(ua),
			TraceID:    structs.StringToNullableString(trace),
			TargetID:   params.Log.TargetID,
			TargetType: params.Log.TargetType,
			Metadata:   params.Log.Metadata,
		}
		logger.Record(ctx, params.Tx, entry)
	}

	// 3. Prepare and Send Notification
	if notifier != nil && params.Notification != nil {
		receiverID := params.Notification.ReceiverID
		if !receiverID.Valid || receiverID.String == "" {
			receiverID = structs.StringToNullableString(id)
		}

		notif := NotificationEntry{
			ReceiverID: receiverID,
			ActorID:    structs.StringToNullableString(id),
			TargetID:   params.Notification.TargetID,
			TargetType: params.Notification.TargetType,
			Title:      params.Notification.Title,
			Message:    params.Notification.Message,
			Type:       params.Notification.Type,
		}
		_ = notifier.Send(ctx, notif)
	}
}
