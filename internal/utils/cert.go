package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func GenerateCertificate() error {
    // создаём шаблон сертификата
    cert := &x509.Certificate{
        // указываем уникальный номер сертификата
        SerialNumber: big.NewInt(1658),
        // заполняем базовую информацию о владельце сертификата
        Subject: pkix.Name{
            Organization: []string{"MY ORG"},
            Country:      []string{"RU"},
        },
        // разрешаем использование сертификата для 127.0.0.1 и ::1
        IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
        // сертификат верен, начиная со времени создания
        NotBefore: time.Now(),
        // время жизни сертификата — 10 лет
        NotAfter:     time.Now().AddDate(10, 0, 0),
        SubjectKeyId: []byte{1, 2, 3, 4, 6},
        // устанавливаем использование ключа для цифровой подписи,
        // а также клиентской и серверной авторизации
        ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
        KeyUsage:    x509.KeyUsageDigitalSignature,
    }

    // создаём новый приватный RSA-ключ длиной 4096 бит
    // обратите внимание, что для генерации ключа и сертификата
    // используется rand.Reader в качестве источника случайных данных
    privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
    if err != nil {
        return err
    }

    // создаём сертификат x.509
    certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
    if err != nil {
        return err
    }

    // кодируем сертификат и ключ в формате PEM, который
    // используется для хранения и обмена криптографическими ключами
    var certPEM bytes.Buffer
    err = pem.Encode(&certPEM, &pem.Block{
        Type:  "CERTIFICATE",
        Bytes: certBytes,
    })
    if err != nil {
        return err
    }

    var privateKeyPEM bytes.Buffer
    err = pem.Encode(&privateKeyPEM, &pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    })
    if err != nil {
        return err
    }

    // Сохраняем сертификат и приватный ключ в файлы ~/cert.pem и ~/private.pem
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return err
    }

    if err = os.WriteFile(filepath.Join(homeDir, "cert.pem"), certPEM.Bytes(), 0644); err != nil {
        return err
    }

    if err = os.WriteFile(filepath.Join(homeDir, "private.pem"), privateKeyPEM.Bytes(), 0644); err != nil {
        return err
    }
    return nil
}