package auth

import (
	"github.com/pquerna/otp/totp"
)

// TOTPGenerator TOTP生成器
type TOTPGenerator struct {
	issuer string
}

// NewTOTPGenerator 创建TOTP生成器
func NewTOTPGenerator(issuer string) *TOTPGenerator {
	return &TOTPGenerator{
		issuer: issuer,
	}
}

// Generate 生成TOTP密钥和二维码
func (t *TOTPGenerator) Generate(accountName string) (secret, qrCode string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      t.issuer,
		AccountName: accountName,
	})
	if err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

// Validate 验证TOTP代码
func (t *TOTPGenerator) Validate(secret, code string) bool {
	return totp.Validate(code, secret)
}
