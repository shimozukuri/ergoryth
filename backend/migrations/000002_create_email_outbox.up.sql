CREATE TABLE ergoryth.email_outbox (
                                       id           UUID PRIMARY KEY,
                                       lead_id      UUID NOT NULL REFERENCES ergoryth.leads(id),
                                       status       VARCHAR(20) NOT NULL DEFAULT 'pending',
                                       attempts     INTEGER NOT NULL DEFAULT 0,
                                       last_error   TEXT,
                                       created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                       available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                       sent_at      TIMESTAMPTZ,

                                       CONSTRAINT email_outbox_status_valid
                                           CHECK (status IN ('pending', 'sent', 'failed')),

                                       CONSTRAINT email_outbox_attempts_valid
                                           CHECK (attempts >= 0)
);

CREATE INDEX email_outbox_pending_idx
    ON ergoryth.email_outbox (available_at)
    WHERE status = 'pending';