package db

type DBOpBase struct {
	DBName     string
	CollName   string
	FromModule string
}

func (b *DBOpBase) GetDBName() string {
	return b.DBName
}

func (b *DBOpBase) GetCollName() string {
	return b.CollName
}

func (b *DBOpBase) GetFromModule() string {
	return b.FromModule
}
