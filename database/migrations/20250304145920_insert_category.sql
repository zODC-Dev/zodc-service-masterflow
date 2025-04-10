-- +goose Up
INSERT INTO public.categories
    ("name", "type", "key")
VALUES
    ('Recruitment', 'GENERAL', 'HR'),
    ('Employee Onboarding', 'GENERAL', 'EMPLOYEE_ONBOARDING'),
    ('Sprint', 'PROJECT', 'SPRINT'),
    ('Story', 'PROJECT', 'STORY');


-- +goose Down

