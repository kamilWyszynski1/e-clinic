CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- uuid provider
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TYPE genderenum AS ENUM ('XX','FEMALE', 'MALE');

CREATE TABLE patient
(
    id           uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    name         VARCHAR    NOT NULL,
    surname      VARCHAR    NOT NULL,
    email        VARCHAR    NOT NULL,
    phone_number VARCHAR    NOT NULL,
    age          int        NOT NULL,
    gender       genderenum NOT NULL      DEFAULT 'XX',

    created_at   TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at   TIMESTAMP WITH TIME ZONE
);

CREATE TYPE specialityenum AS ENUM (
    'Critical Care Medicine', 'Diagnostic Radiology', 'Dermatology', 'Cardiology' ,'Anesthesiology', 'Cardiovascular'
    );


CREATE TABLE specialist
(
    id           uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    name         VARCHAR          NOT NULL,
    surname      VARCHAR          NOT NULL,
    specialities specialityenum[] NOT NULL,

    created_at   TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at   TIMESTAMP WITH TIME ZONE
);

CREATE TABLE specialist_fee
(
    id             uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    specialist     uuid REFERENCES specialist (id) NOT NULL,
    speciality     specialityenum                  NOT NULL,
    fee_per_30_min DECIMAL                         NOT NULL
);

CREATE TYPE apoitntmentstateenum AS ENUM (
    'CREATED', -- patient created appointment, waiting for specialist to accept
    'ACCEPTED', -- specialist accepted appointment
    'TO_BE_PAID', -- specialist accepted, waiting for patient to pay
    'OK', -- patient paid for appointment
    'FINISHED', -- specialist made prescription
    'REJECTED', -- specialist rejected appointment
    'NOT_PAID' -- patient didn't pay in time
    );

CREATE TABLE appointment
(
    id             uuid PRIMARY KEY                             DEFAULT uuid_generate_v4(),
    state          apoitntmentstateenum                NOT NULL DEFAULT 'CREATED',
    patient        uuid REFERENCES patient (id)        NOT NULL,
    specialist_fee uuid REFERENCES specialist_fee (id) NOT NULL,
    scheduled_time TIMESTAMP WITH TIME ZONE            NOT NULL,
    duration       INT                                 NOT NULL, -- duration in seconds

    created_at     TIMESTAMP WITH TIME ZONE                     DEFAULT now(),
    updated_at     TIMESTAMP WITH TIME ZONE
);

CREATE TABLE appointment_form
(
    id          uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    appointment uuid REFERENCES appointment (id) NOT NULL,
    comment     VARCHAR                          NOT NULL,
    symptoms    VARCHAR[]                        NOT NULL,

    created_at  TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE
);

CREATE TABLE appointment_result
(
    id          uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    appointment uuid REFERENCES appointment (id) NOT NULL,
    comment     VARCHAR                          NOT NULL,

    created_at  TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE
);

CREATE TABLE appointment_result_prescription
(
    id                 uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    appointment_result uuid REFERENCES appointment_result (id) NOT NULL,
    drug               INT REFERENCES drug (id)                NOT NULL,
    dosing             VARCHAR                                 NOT NULL,

    created_at         TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE payment
(
    id          uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    appointment uuid REFERENCES appointment (id) UNIQUE NOT NULL,
    price       double precision                        NOT NULL,
    order_id    VARCHAR UNIQUE                          NOT NULL,
    status      VARCHAR                                 NOT NULL,

    created_at  TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE
);

CREATE TABLE drug
(
    id                  INT PRIMARY KEY,
    name                VARCHAR NOT NULL,
    type_of_preparation VARCHAR NOT NULL,
    common_name         VARCHAR NOT NULL,
    strength            VARCHAR NOT NULL, -- 4mg/5ml
    shape               VARCHAR NOT NULL
);

CREATE INDEX lowered_name_inx ON drug USING gin (name gin_trgm_ops); -- like index

CREATE TABLE substance
(
    id   uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL
);


CREATE TABLE composition
(
    id        uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    drug      INT REFERENCES drug (id)       NOT NULL,
    substance uuid REFERENCES substance (id) NOT NULL
);
