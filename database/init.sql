CREATE TABLE IF NOT EXISTS hospitals (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS patients (
    id BIGSERIAL PRIMARY KEY,
    first_name_th VARCHAR,
    middle_name_th VARCHAR,
    last_name_th VARCHAR,
    first_name_en VARCHAR,
    middle_name_en VARCHAR,
    last_name_en VARCHAR,
    date_of_birth DATE,
    national_id VARCHAR UNIQUE,
    passport_id VARCHAR,
    phone_number VARCHAR,
    email VARCHAR,
    gender CHAR(1),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT chk_gender CHECK (gender IN ('M', 'F'))
);

CREATE TABLE IF NOT EXISTS staffs (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    hospital_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_staff_hospital
        FOREIGN KEY (hospital_id)
        REFERENCES hospitals(id)
        ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS patient_hospitals_mapping (
    id BIGSERIAL PRIMARY KEY,
    patient_id BIGINT NOT NULL,
    hospital_id BIGINT NOT NULL,
    patient_hn VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT unique_patient_hospital UNIQUE (patient_id, hospital_id),
    CONSTRAINT unique_hn_per_hospital UNIQUE (hospital_id, patient_hn),

    CONSTRAINT fk_phm_patient
        FOREIGN KEY (patient_id)
        REFERENCES patients(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_phm_hospital
        FOREIGN KEY (hospital_id)
        REFERENCES hospitals(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_passport_not_null
ON patients(passport_id)
WHERE passport_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_patients_national_id
ON patients(national_id);

CREATE INDEX IF NOT EXISTS idx_patients_passport_id
ON patients(passport_id);

CREATE INDEX IF NOT EXISTS idx_patients_not_deleted
ON patients(id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_staff_username
ON staffs(username);

CREATE INDEX IF NOT EXISTS idx_staffs_not_deleted
ON staffs(id)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_phm_hospital_id
ON patient_hospitals_mapping(hospital_id);

CREATE INDEX IF NOT EXISTS idx_phm_patient_id
ON patient_hospitals_mapping(patient_id);

INSERT INTO hospitals (id, name)
VALUES (1, 'Hospital A')
ON CONFLICT (id) DO NOTHING;