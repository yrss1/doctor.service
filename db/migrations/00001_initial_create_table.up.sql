DO $$
    BEGIN
        -- EXTENSIONS --
        CREATE EXTENSION IF NOT EXISTS pgcrypto;

        -- SEQUENCES --
        CREATE SEQUENCE IF NOT EXISTS appointment_id_seq
            START WITH 99999999
            INCREMENT BY 1;

        -- TABLES --
        CREATE TABLE IF NOT EXISTS clinic (
                                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                              name VARCHAR(100) NOT NULL,
                                              address TEXT NOT NULL,
                                              phone VARCHAR(20),
                                              UNIQUE (name, address),
                                              UNIQUE (phone)
        );

        CREATE TABLE IF NOT EXISTS doctors (
                                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                               name VARCHAR(100) NOT NULL,
                                               specialization VARCHAR(100) NOT NULL,
                                               experience INT CHECK (experience >= 0) NOT NULL,
                                               price DECIMAL(10,2) NOT NULL,
                                               rating DECIMAL(3,2) DEFAULT 0.0 NOT NULL CHECK (rating BETWEEN 0 AND 5),
                                               address TEXT NOT NULL,
                                               phone VARCHAR(20) NOT NULL UNIQUE,
                                               clinic_id UUID REFERENCES clinic(id) ON DELETE SET NULL
        );

        CREATE TABLE IF NOT EXISTS schedule (
                                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                                doctor_id UUID NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
                                                slot_start TIMESTAMP NOT NULL,
                                                slot_end TIMESTAMP NOT NULL,
                                                is_available BOOLEAN DEFAULT TRUE,
                                                CHECK (slot_end > slot_start),
                                                UNIQUE (doctor_id, slot_start, slot_end)
        );

        CREATE TABLE IF NOT EXISTS appointments (
                                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                    id BIGINT PRIMARY KEY DEFAULT nextval('appointment_id_seq'),
                                                    doctor_id UUID NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
                                                    user_id UUID NOT NULL, -- ID пользователя из внешнего сервиса
                                                    schedule_id UUID NOT NULL REFERENCES schedule(id) ON DELETE CASCADE,
                                                    status VARCHAR(10) NOT NULL CHECK (status IN ('active', 'canceled', 'completed')),
                                                    meeting_url TEXT,
                                                    UNIQUE (schedule_id, user_id)
        );

        CREATE TABLE IF NOT EXISTS reviews (
                                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                               doctor_id UUID NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
                                               user_id UUID NOT NULL, -- ID пользователя из внешнего сервиса
                                               rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
                                               comment TEXT
        );

        -- DATA --
        INSERT INTO clinic (name, address, phone) VALUES
                                                      ('City Hospital', '123 Main St, Springfield', '123-456-7890'),
                                                      ('Central Clinic', '456 Oak St, Metropolis', '987-654-3210');

        INSERT INTO doctors (name, specialization, experience, price, rating, address, phone, clinic_id) VALUES
                                                                                                             ('Dr. John Doe', 'Cardiologist', 15, 150.00, 4.8, '123 Main St, Springfield', '111-222-3333', (SELECT id FROM clinic WHERE name = 'City Hospital')),
                                                                                                             ('Dr. Alice Smith', 'Dermatologist', 10, 100.00, 4.5, '456 Oak St, Metropolis', '444-555-6666', (SELECT id FROM clinic WHERE name = 'Central Clinic'));

        INSERT INTO schedule (doctor_id, slot_start, slot_end) VALUES
                                                                   ((SELECT id FROM doctors WHERE name = 'Dr. John Doe'), '2025-03-20 09:00:00', '2025-03-20 10:00:00'),
                                                                   ((SELECT id FROM doctors WHERE name = 'Dr. Alice Smith'), '2025-03-21 10:00:00', '2025-03-21 11:00:00');

        INSERT INTO appointments (doctor_id, user_id, schedule_id, status) VALUES
                                                                               ((SELECT id FROM doctors WHERE name = 'Dr. John Doe'), gen_random_uuid(), (SELECT id FROM schedule WHERE slot_start = '2025-03-20 09:00:00'), 'active'),
                                                                               ((SELECT id FROM doctors WHERE name = 'Dr. Alice Smith'), gen_random_uuid(), (SELECT id FROM schedule WHERE slot_start = '2025-03-21 10:00:00'), 'completed');

        INSERT INTO reviews (doctor_id, user_id, rating, comment) VALUES
                                                                      ((SELECT id FROM doctors WHERE name = 'Dr. John Doe'), gen_random_uuid(), 5, 'Excellent service!'),
                                                                      ((SELECT id FROM doctors WHERE name = 'Dr. Alice Smith'), gen_random_uuid(), 4, 'Very professional.');

        COMMIT;
    END $$;
