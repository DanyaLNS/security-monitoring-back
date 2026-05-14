CREATE TABLE event_sources (
    ulid VARCHAR(26) PRIMARY KEY,

    name VARCHAR(128) NOT NULL UNIQUE,
    description TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE event_types (
    ulid VARCHAR(26) PRIMARY KEY,

    code VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,

    description TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE events (
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

    occurred_at TIMESTAMP NOT NULL,

    raw_payload JSONB NOT NULL,

    normalized_payload JSONB,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE event_analysis (
    ulid VARCHAR(26) PRIMARY KEY,

    event_ulid VARCHAR(26) NOT NULL
        REFERENCES events(ulid)
        ON DELETE CASCADE,

    threat_score NUMERIC(5,2) NOT NULL,

    threat_level VARCHAR(32) NOT NULL,

    detection_method VARCHAR(64),

    analysis_summary TEXT,

    recommendations TEXT,

    analyzed_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_events_occurred_at
    ON events(occurred_at DESC);

CREATE INDEX idx_events_severity
    ON events(severity);

CREATE INDEX idx_events_status
    ON events(status);

CREATE INDEX idx_events_source_ulid
    ON events(source_ulid);

CREATE INDEX idx_events_type_ulid
    ON events(type_ulid);

CREATE INDEX idx_events_source_ip
    ON events(source_ip);

CREATE INDEX idx_events_hostname
    ON events(hostname);

CREATE INDEX idx_events_raw_payload
    ON events
    USING GIN(raw_payload);

CREATE INDEX idx_events_normalized_payload
    ON events
    USING GIN(normalized_payload);

CREATE INDEX idx_event_analysis_event_ulid
    ON event_analysis(event_ulid);

CREATE INDEX idx_event_analysis_threat_level
    ON event_analysis(threat_level);

CREATE INDEX idx_event_analysis_score
    ON event_analysis(threat_score);