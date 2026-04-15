package sqlite

type Config struct {
	Path string
}

type Store struct {
	path string
}

func New(cfg Config) (*Store, error) {
	path := cfg.Path
	if path == "" {
		path = "allseer.db"
	}

	return &Store{path: path}, nil
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) Close() error {
	return nil
}
