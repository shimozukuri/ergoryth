CREATE SCHEMA IF NOT EXISTS ergoryth;

CREATE TABLE ergoryth.leads
(
    id         UUID PRIMARY KEY,

    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(254),
    phone      VARCHAR(11),

    company    VARCHAR(150) NOT NULL DEFAULT '',
    message    TEXT         NOT NULL DEFAULT '',

    status     VARCHAR(30)  NOT NULL DEFAULT 'new',
    created_at TIMESTAMPTZ  NOT NULL,

    CONSTRAINT leads_name_valid
        CHECK (
            char_length(name) BETWEEN 2 AND 100
                AND name ~ '^[[:alpha:] -]+$'
                AND name ~ '[[:alpha:]]'
        ),

    CONSTRAINT leads_email_valid
        CHECK (
            email IS NULL
            OR email ~ '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
        ),

    CONSTRAINT leads_phone_valid
        CHECK (
            phone IS NULL
            OR phone ~ '^7[0-9]{10}$'
        ),

    CONSTRAINT leads_message_valid
        CHECK (char_length(message) <= 2000),

    CONSTRAINT leads_contact_required
        CHECK (
            NULLIF(email, '') IS NOT NULL
            OR NULLIF(phone, '') IS NOT NULL
        ),

    CONSTRAINT leads_status_valid
        CHECK (status IN ('new'))
);
