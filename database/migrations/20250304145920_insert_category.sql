-- +goose Up
INSERT INTO public.categories
    ("name", "type", "key")
VALUES
    ('HR', 'GENERAL', 'HR'),
    ('Employee Onboarding', 'GENERAL', 'EMPLOYEE_ONBOARDING'),
    ('Project ZODC', 'PROJECT', 'PROJECT_ZODC'),
    ('Story', 'PROJECT', 'STORY');


-- +goose Down

