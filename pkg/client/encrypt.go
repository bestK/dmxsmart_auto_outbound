package client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

const (
	rsaPublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCtxsTrJZkxpuahl2CXxcKg5i83
0Zus/KdXWQH6aeZYQPzp0EJs2AqFcmvO1QRoE8l0+SlSGjNl5OI0E/VFsbfZ2PiP
/EhXpu2uncWfkmQM+yhMZG10LwMss+xwXQLkxLK8px4A/Vn+ei3QuI9XXqHjKE2m
k/owWSXoKBrnzY0BWwIDAQAB
-----END PUBLIC KEY-----`
)

// encryptPassword 使用RSA公钥加密密码
func encryptPassword(password string) (string, error) {
	// 解码PEM格式的公钥
	block, _ := pem.Decode([]byte(rsaPublicKey))
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the public key")
	}

	// 解析公钥
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %w", err)
	}

	// 类型断言为RSA公钥
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("not an RSA public key")
	}

	// 加密密码
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, []byte(password))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt password: %w", err)
	}

	// 将加密后的数据转换为base64编码
	return base64.StdEncoding.EncodeToString(encrypted), nil
} 