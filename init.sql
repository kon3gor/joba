CREATE TABLE IF NOT EXISTS known_jobs (
	alert_id TEXT NOT NULL,
	id TEXT NOT NULL,

	PRIMARY KEY (alert_id, id)
);
