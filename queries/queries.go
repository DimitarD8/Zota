package queries

const (
	InsertIntoStore = `
		INSERT INTO store (key, value)
		VALUES ($1, $2)
		ON CONFLICT (key) DO NOTHING
	`

	SelectByKey = `
		SELECT value FROM store WHERE key = $1
	`

	SelectRandomValue = `
		SELECT value FROM store ORDER BY RANDOM() LIMIT 1
	`

	DeleteByKey = `
		DELETE FROM store WHERE key = $1
	`

	DeleteRandom = `
		DELETE FROM store WHERE key IN (
			SELECT key FROM store ORDER BY random() LIMIT 1
		)
	`

	SelectAll = `
		SELECT key, value FROM store
	`

	MutateRandom = `
		UPDATE store
		SET value = md5(random()::text)
		WHERE key IN (
			SELECT key FROM store ORDER BY random() LIMIT 1
		)
	`
)
