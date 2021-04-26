-- Goal schema
CREATE TABLE goal (
	goal serial NOT NULL,
	goal_name varchar(60) NOT NULL,
	unit varchar(60) NULL,
	active bool NULL DEFAULT true,
	CONSTRAINT goal_pkey PRIMARY KEY (goal)
);

-- Progress Schema
CREATE TABLE progress (
	progress serial not null,
	amount numeric(6,2) NOT NULL,
	session_date timestamp NOT NULL DEFAULT now(),
	goal int4 NOT NULL,
	CONSTRAINT progress_pkey PRIMARY KEY (progress),
    CONSTRAINT goal_tracker_goal_fkey FOREIGN KEY (goal) REFERENCES public.goal(goal)
);
