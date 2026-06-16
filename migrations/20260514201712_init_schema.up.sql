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

    source_ip TEXT,
    destination_ip TEXT,

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

    severity SMALLINT NOT NULL
        CHECK (severity BETWEEN 0 AND 10),

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

INSERT INTO event_sources (
    ulid,
    name,
    description
)
VALUES
    (
        '01SRCSSHD0000000000000000',
        'sshd',
        'SSH Server Logs'
    ),
    (
        '01SRCNGINX00000000000000',
        'nginx',
        'NGINX Gateway'
    ),
    (
        '01SRCSURICATA0000000000',
        'suricata',
        'IDS Sensor'
    )
ON CONFLICT (ulid) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description;

INSERT INTO event_types (
    ulid,
    code,
    name,
    description
)
VALUES
    (
        '01TYPEAUTHFAILED00000000',
        'auth_failed',
        'Неудачная аутентификация',
        'Ошибка входа пользователя в систему'
    ),
    (
        '01TYPENETSCAN0000000000',
        'network_scan',
        'Сетевое сканирование',
        'Обнаружение признаков сканирования портов или сетевых ресурсов'
    ),
    (
        '01TYPEHTTPANOMALY000000',
        'http_anomaly',
        'HTTP-аномалия',
        'Подозрительный шаблон HTTP-запросов'
    ),
    (
        '01TYPEHTTPATTACK0000000',
        'http_attack',
        'HTTP-атака',
        'HTTP-запрос с признаками атаки на приложение'
    ),
    (
        '01TYPEHEALTHCHECK000000',
        'health_check',
        'Проверка доступности',
        'Обычный служебный запрос проверки состояния сервиса'
    )
ON CONFLICT (ulid) DO UPDATE SET
    code = EXCLUDED.code,
    name = EXCLUDED.name,
    description = EXCLUDED.description;

INSERT INTO events (
    ulid,
    source_ulid,
    type_ulid,
    severity,
    status,
    title,
    source_ip,
    destination_ip,
    hostname,
    occurred_at,
    raw_payload,
    normalized_payload,
    created_at,
    updated_at
)
VALUES
    (
        '01EVT00000000000000000001',
        '01SRCSSHD0000000000000000',
        '01TYPEAUTHFAILED00000000',
        9,
        'new',
        'SSH brute force attempt detected',
        '192.168.0.10',
        NULL,
        'server-1',
        NOW() - INTERVAL '2 minutes',
        '{
            "username": "admin",
            "attempts": 5,
            "auth_method": "password",
            "message": "SSH brute force attempt detected"
        }'::jsonb,
        NULL,
        NOW(),
        NOW()
    ),
    (
        '01EVT00000000000000000002',
        '01SRCSSHD0000000000000000',
        '01TYPEAUTHFAILED00000000',
        8,
        'new',
        'Repeated SSH authentication failure',
        '192.168.0.10',
        NULL,
        'server-1',
        NOW() - INTERVAL '3 minutes',
        '{
            "username": "admin",
            "attempts": 4,
            "auth_method": "password",
            "message": "Repeated SSH authentication failure"
        }'::jsonb,
        NULL,
        NOW(),
        NOW()
    ),
    (
        '01EVT00000000000000000003',
        '01SRCSSHD0000000000000000',
        '01TYPEAUTHFAILED00000000',
        7,
        'new',
        'Suspicious SSH login pattern',
        '192.168.0.10',
        NULL,
        'server-1',
        NOW() - INTERVAL '4 minutes',
        '{
            "username": "admin",
            "attempts": 3,
            "auth_method": "password",
            "message": "Suspicious SSH login pattern"
        }'::jsonb,
        NULL,
        NOW(),
        NOW()
    ),
    (
        '01EVT00000000000000000004',
        '01SRCSURICATA0000000000',
        '01TYPENETSCAN0000000000',
        6,
        'new',
        'Port scan detected from external IP',
        '8.8.8.8',
        NULL,
        'fw-1',
        NOW() - INTERVAL '10 minutes',
        '{
            "ports": [22, 80, 443, 5432],
            "protocol": "tcp",
            "message": "Port scan detected from external IP"
        }'::jsonb,
        NULL,
        NOW(),
        NOW()
    ),
    (
        '01EVT00000000000000000005',
        '01SRCSURICATA0000000000',
        '01TYPENETSCAN0000000000',
        5,
        'new',
        'High-frequency port probing',
        '8.8.8.8',
        NULL,
        'fw-1',
        NOW() - INTERVAL '12 minutes',
        '{
            "ports": [21, 22, 23, 25, 80, 443],
            "protocol": "tcp",
            "message": "High-frequency port probing"
        }'::jsonb,
        NULL,
        NOW(),
        NOW()
    ),
    (
        '01EVT00000000000000000006',
        '01SRCNGINX00000000000000',
        '01TYPEHTTPANOMALY000000',
        4,
        'new',
        'Suspicious HTTP request pattern',
        '10.0.0.5',
        NULL,
        'api-gateway',
        NOW() - INTERVAL '20 minutes',
        '{
            "method": "GET",
            "path": "/api/users",
            "status_code": 403,
            "message": "Suspicious HTTP request pattern"
        }'::jsonb,
        NULL,
        NOW(),
        NOW()
    ),
    (
        '01EVT00000000000000000007',
        '01SRCNGINX00000000000000',
        '01TYPEHTTPATTACK0000000',
        8,
        'new',
        'Possible injection attempt blocked',
        '10.0.0.5',
        NULL,
        'api-gateway',
        NOW() - INTERVAL '25 minutes',
        '{
            "method": "POST",
            "path": "/api/search",
            "status_code": 400,
            "attack_type": "injection",
            "message": "Possible injection attempt blocked"
        }'::jsonb,
        NULL,
        NOW(),
        NOW()
    ),
    (
        '01EVT00000000000000000008',
        '01SRCNGINX00000000000000',
        '01TYPEHEALTHCHECK000000',
        1,
        'new',
        'Normal health check request',
        '127.0.0.1',
        NULL,
        'api-gateway',
        NOW() - INTERVAL '30 minutes',
        '{
            "method": "GET",
            "path": "/health",
            "status_code": 200,
            "message": "Normal health check request"
        }'::jsonb,
        NULL,
        NOW(),
        NOW()
    )
ON CONFLICT (ulid) DO UPDATE SET
    source_ulid = EXCLUDED.source_ulid,
    type_ulid = EXCLUDED.type_ulid,
    severity = EXCLUDED.severity,
    status = EXCLUDED.status,
    title = EXCLUDED.title,
    source_ip = EXCLUDED.source_ip,
    destination_ip = EXCLUDED.destination_ip,
    hostname = EXCLUDED.hostname,
    occurred_at = EXCLUDED.occurred_at,
    raw_payload = EXCLUDED.raw_payload,
    normalized_payload = EXCLUDED.normalized_payload,
    updated_at = NOW();

INSERT INTO incidents (
    ulid,
    title,
    description,
    severity,
    status,
    created_at,
    updated_at
)
VALUES
    (
        '01INCSSHBRUTEFORCE000000',
        'SSH brute force cluster',
        'Серия неудачных попыток входа по SSH с одного IP-адреса на сервер server-1.',
        9,
        'open',
        NOW(),
        NOW()
    ),
    (
        '01INCNETWORKSCAN00000000',
        'Network scan activity',
        'Обнаружена серия сетевых запросов, похожих на сканирование портов внешним источником.',
        6,
        'open',
        NOW(),
        NOW()
    ),
    (
        '01INCWEBATTACK000000000',
        'Suspicious web activity',
        'Подозрительная активность на уровне HTTP-запросов к API Gateway.',
        8,
        'in_progress',
        NOW(),
        NOW()
    )
ON CONFLICT (ulid) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    severity = EXCLUDED.severity,
    status = EXCLUDED.status,
    updated_at = NOW();

INSERT INTO incident_events (
    incident_ulid,
    event_ulid
)
VALUES
    (
        '01INCSSHBRUTEFORCE000000',
        '01EVT00000000000000000001'
    ),
    (
        '01INCSSHBRUTEFORCE000000',
        '01EVT00000000000000000002'
    ),
    (
        '01INCSSHBRUTEFORCE000000',
        '01EVT00000000000000000003'
    ),

    (
        '01INCNETWORKSCAN00000000',
        '01EVT00000000000000000004'
    ),
    (
        '01INCNETWORKSCAN00000000',
        '01EVT00000000000000000005'
    ),

    (
        '01INCWEBATTACK000000000',
        '01EVT00000000000000000006'
    ),
    (
        '01INCWEBATTACK000000000',
        '01EVT00000000000000000007'
    )
ON CONFLICT (
    incident_ulid,
    event_ulid
) DO NOTHING;