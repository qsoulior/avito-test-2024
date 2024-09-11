package postgres

type Option func(*Postgres)

func MaxConns(conns int32) Option {
	return func(c *Postgres) {
		c.maxConns = conns
	}
}

func DataTypes(dataTypes []string) Option {
	return func(c *Postgres) {
		c.dataTypes = dataTypes
	}
}
