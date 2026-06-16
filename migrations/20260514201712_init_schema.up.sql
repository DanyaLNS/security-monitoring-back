CREATE TABLE IF NOT EXISTS event_sources (
    ulid VARCHAR(26) PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS event_types (
    ulid VARCHAR(26) PRIMARY KEY,
    code VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS events (
    ulid VARCHAR(26) PRIMARY KEY,

    source_ulid VARCHAR(26) NOT NULL
        REFERENCES event_sources(ulid),

    type_ulid VARCHAR(26) NOT NULL
        REFERENCES event_types(ulid),

    severity SMALLINT NOT NULL
        CHECK (severity BETWEEN 0 AND 10),

    status VARCHAR(32) NOT NULL DEFAULT 'new',

    title VARCHAR(255) NOT NULL,

    source_ip INET,
    destination_ip INET,
    hostname VARCHAR(255),

    occurred_at TIMESTAMPTZ NOT NULL,

    raw_payload JSONB NOT NULL,
    normalized_payload JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS incidents (
    ulid VARCHAR(26) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    severity SMALLINT NOT NULL CHECK (severity BETWEEN 0 AND 10),
    status VARCHAR(32) NOT NULL DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS incident_events (
    incident_ulid VARCHAR(26) NOT NULL
        REFERENCES incidents(ulid)
        ON DELETE CASCADE,

    event_ulid VARCHAR(26) NOT NULL
        REFERENCES events(ulid)
        ON DELETE CASCADE,

    PRIMARY KEY (incident_ulid, event_ulid)
);

CREATE TABLE IF NOT EXISTS event_analysis (
    ulid VARCHAR(26) PRIMARY KEY,

    event_ulid VARCHAR(26) NOT NULL
        REFERENCES events(ulid)
        ON DELETE CASCADE,

    threat_score NUMERIC(5,2) NOT NULL,
    threat_level VARCHAR(32) NOT NULL,
    detection_method VARCHAR(64),
    analysis_summary TEXT,
    recommendations TEXT,
    analyzed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_events_occurred_at
    ON events(occurred_at DESC);

CREATE INDEX IF NOT EXISTS idx_events_severity
    ON events(severity);

CREATE INDEX IF NOT EXISTS idx_events_status
    ON events(status);

CREATE INDEX IF NOT EXISTS idx_events_source_ulid
    ON events(source_ulid);

CREATE INDEX IF NOT EXISTS idx_events_type_ulid
    ON events(type_ulid);

CREATE INDEX IF NOT EXISTS idx_events_source_ip
    ON events(source_ip);

CREATE INDEX IF NOT EXISTS idx_events_hostname
    ON events(hostname);

CREATE INDEX IF NOT EXISTS idx_events_raw_payload
    ON events USING GIN(raw_payload);

CREATE INDEX IF NOT EXISTS idx_events_normalized_payload
    ON events USING GIN(normalized_payload);

CREATE INDEX IF NOT EXISTS idx_incident_events_incident_ulid
    ON incident_events(incident_ulid);

CREATE INDEX IF NOT EXISTS idx_incident_events_event_ulid
    ON incident_events(event_ulid);

CREATE INDEX IF NOT EXISTS idx_event_analysis_event_ulid
    ON event_analysis(event_ulid);

CREATE INDEX IF NOT EXISTS idx_event_analysis_threat_level
    ON event_analysis(threat_level);

CREATE INDEX IF NOT EXISTS idx_event_analysis_score
    ON event_analysis(threat_score);

INSERT INTO event_sources (ulid, name, description)
VALUES
    ('01ARZ3NDEKTSV4RRFFQ69G5FAV', 'backend-api', 'Основное серверное приложение'),
    ('01ARZ3NDEKTSV4RRFFQ69G5FAW', 'auth-service', 'Сервис аутентификации'),
    ('01ARZ3NDEKTSV4RRFFQ69G5FAX', 'gateway', 'API-шлюз')
ON CONFLICT (ulid) DO NOTHING;

INSERT INTO event_types (ulid, code, name, description)
VALUES
    ('01ARZ3NDEKTSV4RRFFQ69G5FB1', 'auth_failed', 'Неудачная аутентификация', 'Ошибка входа пользователя в систему'),
    ('01ARZ3NDEKTSV4RRFFQ69G5FB2', 'access_denied', 'Отказ в доступе', 'Попытка обращения к защищённому ресурсу'),
    ('01ARZ3NDEKTSV4RRFFQ69G5FB3', 'suspicious_request', 'Подозрительный запрос', 'Запрос с признаками потенциально вредоносной активности'),
    ('01ARZ3NDEKTSV4RRFFQ69G5FB4', 'rate_limit_exceeded', 'Превышение лимита запросов', 'Слишком большое количество запросов за короткий период')
ON CONFLICT (ulid) DO NOTHING;