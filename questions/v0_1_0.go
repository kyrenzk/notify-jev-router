package questions

import (
	jev "github.com/havlan/jev-go"
)

const Version = "v0.1.0"

// Build returns the QuestionSet for a category (security uses a smaller subset).
func Build(category string) jev.Questions {
	switch category {
	case "security":
		return jev.Questions{
			"primary_strategy": jev.Choice(
				"Select the primary channel strategy for this security notification.",
				jev.StringCriteria(map[string]string{
					"email_only":     "Deliver via primary email only",
					"push_only":      "Deliver via push only",
					"email_and_push": "Deliver via email and push",
					"in_app_only":    "Deliver in-app only",
				}),
			),
			"push_fanout": jev.Choice(
				"How should push notifications fan out?",
				jev.StringCriteria(map[string]string{
					"all_devices":      "All registered devices",
					"last_active_only": "Last active device only",
					"none":             "Do not send push",
				}),
			),
		}
	default:
		return jev.Questions{
			"primary_strategy": jev.Choice(
				"Select the primary channel strategy for this notification.",
				jev.StringCriteria(map[string]string{
					"email_only":     "Deliver via primary email only",
					"push_only":      "Deliver via push only",
					"email_and_push": "Deliver via email and push",
					"in_app_only":    "Deliver in-app only",
				}),
			),
			"recovery_copy": jev.Noul("Should a copy be sent to the recovery email address?"),
			"push_fanout": jev.Choice(
				"How should push notifications fan out?",
				jev.StringCriteria(map[string]string{
					"all_devices":      "All registered devices",
					"last_active_only": "Last active device only",
					"none":             "Do not send push",
				}),
			),
			"email_mandatory": jev.Noul("Is primary email delivery mandatory for compliance?"),
		}
	}
}
