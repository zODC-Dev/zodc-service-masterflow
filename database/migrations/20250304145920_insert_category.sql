-- +goose Up
INSERT INTO public.categories
    ("name", "type", "key")
VALUES
    ('Recruitment', 'GENERAL', 'RECRUITMENT'),
    ('Employee Onboarding', 'GENERAL', 'EMPLOYEE_ONBOARDING'),
    ('Sprint', 'PROJECT', 'SPRINT'),
    ('Story', 'PROJECT', 'STORY'),
    ('Retrospective', 'GENERAL', 'RETROSPECTIVE'),
    ('Edit Profile', 'GENERAL', 'EDIT_PROFILE'),
    ('Performance Evaluate', 'GENERAL', 'PERFORMANCE_EVALUATE'),
    ('Human Resources (HR)', 'GENERAL', 'HR');

-- +goose Down

