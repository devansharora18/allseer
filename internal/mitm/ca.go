package mitm

type CA struct {
	CertPath string
	KeyPath  string
}

func LoadOrCreateCA(certPath, keyPath string) (*CA, error) {
	return &CA{CertPath: certPath, KeyPath: keyPath}, nil
}
