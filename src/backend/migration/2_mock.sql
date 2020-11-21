INSERT INTO patient(name, surname, email, phone_number, age)
VALUES ('John', 'Cena', 'kamil.wyszynski@wp.pl', '500600700', 20);

INSERT INTO patient(name, surname, email, phone_number, age)
VALUES ('Robert', 'Bauer', 'kamil.wyszynski@wp.pl', '500600700', 30);

INSERT INTO patient(name, surname, email, phone_number, age)
VALUES ('Margo', 'Berger', 'kamil.wyszynski@wp.pl', '500600700', 40);

INSERT INTO patient(name, surname, email, phone_number, age)
VALUES ('Johny', 'Gruber', 'kamil.wyszynski@wp.pl', '500600700', 50);

INSERT INTO patient(name, surname, email, phone_number, age)
VALUES ('David', 'Kowalski', 'kamil.wyszynski@wp.pl', '500600700', 60);

-- DOC 1
INSERT INTO specialist(id, name, surname, specialities)
VALUES ('45607504-969d-447a-b94f-e33b59beaca0','Doc#1', 'Doc#1', '{Critical Care Medicine, Diagnostic Radiology}');

INSERT INTO specialist_fee(specialist, speciality, fee_per_30_min)
values ('45607504-969d-447a-b94f-e33b59beaca0', 'Critical Care Medicine', 100);
INSERT INTO specialist_fee(specialist, speciality, fee_per_30_min)
values ('45607504-969d-447a-b94f-e33b59beaca0', 'Diagnostic Radiology', 250);

-- DOC2
INSERT INTO specialist(id, name, surname, specialities)
VALUES ('c65512c4-9809-4f8e-9d6d-3e4681fc79d6', 'Doc#1', 'Doc#1', '{Dermatology}');
INSERT INTO specialist_fee(specialist, speciality, fee_per_30_min)
values ('c65512c4-9809-4f8e-9d6d-3e4681fc79d6', 'Dermatology', 50);