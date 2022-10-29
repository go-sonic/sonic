package service

import (
	"context"

	"github.com/go-sonic/sonic/consts"
)

type BaseMFAService interface {
	UseMFA(mfaType consts.MFAType) bool
	GenerateMFAQRCode(ctx context.Context, content string) (string, error)
}

type TwoFactorTOTPMFAService interface {
	BaseMFAService
	GenerateTFACode(tfaKey string) (string, error)
	ValidateTFACode(tfaKey, tfaCode string) bool
	GenerateOTPKey(ctx context.Context, keyID string) (key, optURL string, err error)
}
