CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- uuid provider

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
    id           uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    appointment  uuid REFERENCES appointment (id) NOT NULL,
    comment      VARCHAR                          NOT NULL,
    prescription VARCHAR                          NOT NULL,

    created_at   TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at   TIMESTAMP WITH TIME ZONE
);

CREATE TABLE payment
(
    id          uuid PRIMARY KEY         DEFAULT uuid_generate_v4(),
    appointment uuid REFERENCES appointment (id) UNIQUE NOT NULL,
    price       double precision                 NOT NULL,
    order_id    VARCHAR                          UNIQUE NOT NULL,
    status      VARCHAR                          NOT NULL,

    created_at  TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE
);