package notifyjev

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// UserPrefs holds channel preference flags without PII.
type UserPrefs struct {
	EmailMarketing   bool `json:"email_marketing"`
	PushEnabled      bool `json:"push_enabled"`
	QuietHoursActive bool `json:"quiet_hours_active"`
}

// Capabilities summarizes what the host can deliver for this subject.
type Capabilities struct {
	EmailVerified    bool `json:"email_verified"`
	PushDeviceCount  int  `json:"push_device_count"`
	HasRecoveryEmail bool `json:"has_recovery_email"`
	SMSAvailable     bool `json:"sms_available"`
	InAppAvailable   bool `json:"in_app_available"`
}

// RoutingContext is the sole Jev state source; it must not contain PII.
type RoutingContext struct {
	EventID      string       `json:"event_id"`
	TemplateID   string       `json:"template_id"`
	Category     Category     `json:"category"`
	Severity     Severity     `json:"severity"`
	Locale       string       `json:"locale"`
	SubjectRef   string       `json:"subject_ref,omitempty"`
	UserPrefs    UserPrefs    `json:"user_prefs"`
	Capabilities Capabilities `json:"capabilities"`
}

// forbiddenJSONKeys must never appear in RoutingContext JSON (PII or disallowed identifiers).
var forbiddenJSONKeys = []string{
	"user_id",
	"email",
	"phone",
	"body",
	"message_body",
	"ip",
	"user_name",
	"device_tokens",
	"password",
	"name",
	"address",
}

// Validate checks RoutingContext invariants for v0.1.0.
func (rc *RoutingContext) Validate() error {
	if strings.TrimSpace(rc.EventID) == "" {
		return fmt.Errorf("%w: event_id is required", ErrInvalidContext)
	}
	if strings.TrimSpace(rc.TemplateID) == "" {
		return fmt.Errorf("%w: template_id is required", ErrInvalidContext)
	}
	if !rc.Category.valid() {
		return fmt.Errorf("%w: unknown category %q", ErrInvalidContext, rc.Category)
	}
	if !rc.Severity.valid() {
		return fmt.Errorf("%w: unknown severity %q", ErrInvalidContext, rc.Severity)
	}
	if strings.TrimSpace(rc.Locale) == "" {
		return fmt.Errorf("%w: locale is required", ErrInvalidContext)
	}
	if !rc.Capabilities.AnyDeliverable(rc.UserPrefs) {
		return fmt.Errorf("%w: no deliverable channel", ErrInvalidContext)
	}
	return nil
}

// AnyDeliverable reports whether at least one channel can receive a notification.
func (c Capabilities) AnyDeliverable(prefs UserPrefs) bool {
	if c.EmailVerified {
		return true
	}
	if prefs.PushEnabled && c.PushDeviceCount > 0 {
		return true
	}
	if c.SMSAvailable {
		return true
	}
	if c.InAppAvailable {
		return true
	}
	return false
}

// ValidateRoutingContextJSON decodes JSON with unknown-field rejection and forbidden key checks.
func ValidateRoutingContextJSON(data []byte) (RoutingContext, error) {
	if err := findForbiddenKeys(data); err != nil {
		return RoutingContext{}, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var rc RoutingContext
	if err := dec.Decode(&rc); err != nil {
		return RoutingContext{}, fmt.Errorf("%w: %v", ErrInvalidContext, err)
	}
	if err := rc.Validate(); err != nil {
		return RoutingContext{}, err
	}
	return rc, nil
}

func findForbiddenKeys(data []byte) error {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("%w: invalid json", ErrInvalidContext)
	}
	return walkForbiddenKeys(raw, "")
}

func walkForbiddenKeys(v any, path string) error {
	switch x := v.(type) {
	case map[string]any:
		for k, child := range x {
			full := k
			if path != "" {
				full = path + "." + k
			}
			for _, forbidden := range forbiddenJSONKeys {
				if strings.EqualFold(k, forbidden) {
					return fmt.Errorf("%w: forbidden field %q", ErrInvalidContext, full)
				}
			}
			if err := walkForbiddenKeys(child, full); err != nil {
				return err
			}
		}
	case []any:
		for i, child := range x {
			if err := walkForbiddenKeys(child, fmt.Sprintf("%s[%d]", path, i)); err != nil {
				return err
			}
		}
	}
	return nil
}

// JevState returns an allowlisted JSON object for Jev SystemOne state.
func (rc RoutingContext) JevState() map[string]any {
	return map[string]any{
		"event_id":     rc.EventID,
		"template_id":  rc.TemplateID,
		"category":     rc.Category,
		"severity":     rc.Severity,
		"locale":       rc.Locale,
		"subject_ref":  rc.SubjectRef,
		"user_prefs":   rc.UserPrefs,
		"capabilities": rc.Capabilities,
	}
}
