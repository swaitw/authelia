package model

import (
	"github.com/authelia/authelia/v4/internal/utils"
)

// UserInfo represents the user information required by the web UI.
type UserInfo struct {
	// The users display name.
	DisplayName string `db:"-" json:"display_name"`

	// The preferred 2FA method.
	Method string `db:"second_factor_method" json:"method" valid:"required"`

	// True if a TOTP device has been registered.
	HasTOTP bool `db:"has_totp" json:"has_totp" valid:"required"`

	// True if a Webauthn device has been registered.
	HasWebauthn bool `db:"has_webauthn" json:"has_webauthn" valid:"required"`

	// True if a duo device has been configured as the preferred.
	HasDuo bool `db:"has_duo" json:"has_duo" valid:"required"`
}

// SetDefaultPreferred2FAMethod configures the default method based on what is configured as available and the users available methods.
func (i *UserInfo) SetDefaultPreferred2FAMethod(methods []string) (changed bool) {
	if len(methods) == 0 {
		// No point attempting to change the method if no methods are available.
		return false
	}

	before := i.Method

	totp, webauthn, duo := utils.IsStringInSlice(SecondFactorMethodTOTP, methods), utils.IsStringInSlice(SecondFactorMethodWebauthn, methods), utils.IsStringInSlice(SecondFactorMethodDuo, methods)

	if i.Method != "" && !utils.IsStringInSlice(i.Method, methods) {
		i.Method = ""
	}

	if i.Method == "" {
		switch {
		case i.HasTOTP && totp:
			i.Method = SecondFactorMethodTOTP
		case i.HasWebauthn && webauthn:
			i.Method = SecondFactorMethodWebauthn
		case i.HasDuo && duo:
			i.Method = SecondFactorMethodDuo
		case totp:
			i.Method = SecondFactorMethodTOTP
		case webauthn:
			i.Method = SecondFactorMethodWebauthn
		case duo:
			i.Method = SecondFactorMethodDuo
		}
	}

	return before != i.Method
}
