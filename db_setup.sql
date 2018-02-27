CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR NOT NULL,
	role VARCHAR NOT NULL,
	email VARCHAR,
	phone VARCHAR,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO public.users(name, role, email) VALUES ('Zeke', 'employee', 'zeke@cat.com');
INSERT INTO public.users(name, role, email, phone) VALUES ('Allie', 'employee', 'allie@cat.com', '123 456 789');
INSERT INTO public.users(name, role, phone) VALUES ('Moxie', 'employee', '123 456 789');
INSERT INTO public.users(name, role, email, phone) VALUES ('Barkly', 'employee', 'barkly@dog.com', '123 456 789');
INSERT INTO public.users(name, role, email, phone) VALUES ('Dexter', 'employee', 'dexter@dog.com', '123 456 789');
INSERT INTO public.users(name, role, email, phone) VALUES ('Katie', 'manager', 'katie@human.com', '123 456 789');
INSERT INTO public.users(name, role, email, phone) VALUES ('Iain', 'manager', 'iain@human.com', '123 456 789');

CREATE TABLE shifts (
	id BIGSERIAL PRIMARY KEY,
	manager_id BIGINT,
	FOREIGN KEY (manager_id) REFERENCES users(id)
	ON UPDATE CASCADE,
	employee_id BIGINT,
	FOREIGN KEY (employee_id) REFERENCES users(id)
	ON UPDATE CASCADE,
	break DOUBLE PRECISION,
	start_time TIMESTAMP NOT NULL,
	end_time TIMESTAMP NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- March 1
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (6, 3, 0.5, '2018-03-01 09:00:00', '2018-03-01 13:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (6, 1, 0.5, '2018-03-01 09:00:00', '2018-03-01 13:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (6, 2, 0.5, '2018-03-01 13:00:00', '2018-03-01 17:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (6, 4, 0.5, '2018-03-01 13:00:00', '2018-03-01 17:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (6, 5, 1, '2018-03-01 09:00:00', '2018-03-01 17:00:00');

-- March 2
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (7, null, 0.5, '2018-03-02 09:00:00', '2018-03-02 13:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (7, null, 0.5, '2018-03-02 09:00:00', '2018-03-02 13:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (7, null, 0.5, '2018-03-02 13:00:00', '2018-03-02 17:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (7, null, 0.5, '2018-03-02 13:00:00', '2018-03-02 17:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (7, null, 1, '2018-03-02 09:00:00', '2018-03-02 17:00:00');

-- March 5
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (6, null, 0.5, '2018-03-05 09:00:00', '2018-03-05 13:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (6, null, 0.5, '2018-03-05 09:00:00', '2018-03-05 13:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (6, null, 0.5, '2018-03-05 13:00:00', '2018-03-05 17:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (6, null, 0.5, '2018-03-05 13:00:00', '2018-03-05 17:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (6, null, 1, '2018-03-05 09:00:00', '2018-03-05 17:00:00');

-- March 6
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (7, null, 0.5, '2018-03-06 09:00:00', '2018-03-06 13:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (7, null, 0.5, '2018-03-06 09:00:00', '2018-03-06 13:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (7, null, 0.5, '2018-03-06 13:00:00', '2018-03-06 17:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (7, null, 0.5, '2018-03-06 13:00:00', '2018-03-06 17:00:00');
INSERT INTO public.shifts(manager_id, employee_id, break, start_time, end_time) VALUES (7, null, 1, '2018-03-06 09:00:00', '2018-03-06 17:00:00');

