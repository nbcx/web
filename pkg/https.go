package pkg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"
)

func Https(addr, dir string) {
	// var port int
	// var dir string
	// flag.IntVar(&port, "port", 8443, "HTTPS 服务端口")
	// flag.StringVar(&dir, "dir", ".", "服务目录")
	// flag.Parse()

	// 切换到服务目录
	if err := os.Chdir(dir); err != nil {
		log.Fatalf("无法进入目录 %s: %v", dir, err)
	}

	// 生成自签名证书
	cert, err := generateSelfSignedCert()
	if err != nil {
		log.Fatalf("生成自签名证书失败: %v", err)
	}

	// HTTPS 服务配置
	server := &http.Server{
		Addr:    addr, //fmt.Sprintf(":%d", port),
		Handler: http.FileServer(http.Dir(".")),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	fmt.Printf("📁 Serving %s\n", dir)
	fmt.Printf("🔐 HTTPS server running on https://localhost:%s\n", dir)

	// 启动 HTTPS 服务
	log.Fatal(server.ListenAndServeTLS("", "")) // 证书已从 tls.Config 提供
}

// generateSelfSignedCert 生成一个内存中的自签名 TLS 证书
func generateSelfSignedCert() (tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	// 设置证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			Organization: []string{"Localhost Self-Signed"},
		},
		NotBefore: time.Now().Add(-1 * time.Hour),
		NotAfter:  time.Now().AddDate(1, 0, 0), // 有效期 1 年

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// 支持 localhost 和 127.0.0.1
	template.IPAddresses = []net.IP{net.ParseIP("127.0.0.1")}
	template.DNSNames = []string{"localhost"}

	// 生成证书
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	// 编码证书和私钥为 PEM 格式
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	// 生成 tls.Certificate
	return tls.X509KeyPair(certPEM, keyPEM)
}
