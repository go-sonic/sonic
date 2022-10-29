package impl

import (
	"bytes"
	"context"
	"encoding/base64"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/yeqown/go-qrcode"

	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type baseMFAServiceImpl struct{}

func NewBaseMFAService() service.BaseMFAService {
	return &baseMFAServiceImpl{}
}

func (b baseMFAServiceImpl) GenerateMFAQRCode(ctx context.Context, content string) (string, error) {
	code, err := qrcode.New(content, qrcode.WithQRWidth(20), qrcode.WithBuiltinImageEncoder(qrcode.PNG_FORMAT))
	if err != nil {
		return "", xerr.NoType.Wrap(err).WithStatus(xerr.StatusInternalServerError).WithMsg("generate mfa qrCode error")
	}
	buf := bytes.Buffer{}
	err = code.SaveTo(&buf)
	if err != nil {
		return "", xerr.NoType.Wrap(err).WithStatus(xerr.StatusInternalServerError).WithMsg("generate mfa qrCode error")
	}
	imageBase64 := make([]byte, base64.StdEncoding.EncodedLen(buf.Len()))
	base64.StdEncoding.Encode(imageBase64, buf.Bytes())
	buf.Reset()
	buf.WriteString("data:image/png;base64,")
	buf.Write(imageBase64)
	return buf.String(), nil
}

func (b baseMFAServiceImpl) UseMFA(mfaType consts.MFAType) bool {
	return mfaType != consts.MFANone
}

type twoFactorTOTPMFAServiceImpl struct {
	service.BaseMFAService
}

func NewTwoFactorTOTPMFAService(baseMFA service.BaseMFAService) service.TwoFactorTOTPMFAService {
	return &twoFactorTOTPMFAServiceImpl{
		baseMFA,
	}
}

func (t *twoFactorTOTPMFAServiceImpl) GenerateTFACode(tfaKey string) (string, error) {
	tfaCode, err := totp.GenerateCode(tfaKey, time.Now())
	if err != nil {
		return "", xerr.BadParam.Wrapf(err, "generate code err with tfaKey=%v", tfaKey)
	}
	return tfaCode, nil
}

func (t *twoFactorTOTPMFAServiceImpl) ValidateTFACode(tfaKey, tfaCode string) bool {
	return totp.Validate(tfaCode, tfaKey)
}

func (t *twoFactorTOTPMFAServiceImpl) GenerateOTPKey(ctx context.Context, keyID string) (key, optURL string, err error) {
	otpKey, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "sonic",
		AccountName: keyID,
		SecretSize:  32,
	})
	if err != nil {
		return "", "", xerr.NoType.New("keyID=%v", keyID).WithStatus(xerr.StatusInternalServerError).WithMsg("generate totop key error")
	}
	return otpKey.Secret(), otpKey.URL(), nil
}
