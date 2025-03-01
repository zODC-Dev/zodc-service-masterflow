-- +goose Up
-- Jira Form Template
INSERT INTO public.form_templates 
    (id, file_name, title, category_id, template_id, data_sheet, description, decoration, tag, "type") 
VALUES
    (1, 'Jira System Form', 'Jira System Form', null, NULL, NULL, 'Jira System Form', 'settings|--primary-40', 'TASK', 'SYSTEM');

SELECT setval('form_templates_id_seq', (SELECT MAX(id) FROM form_templates));

--Jira Form Template Version
INSERT INTO public.form_template_versions 
    (id, "version", form_template_id) 
VALUES
    (1, 1, 1);

SELECT setval('form_template_versions_id_seq', (SELECT MAX(id) FROM form_template_versions));


-- Jira Form Template Fields
INSERT INTO public.form_template_fields 
    (advanced_options, col_num, required, category, title, icon, field_id, field_type, field_name, form_template_version_id) 
VALUES
    ('{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": null, "defaultValue": null}'::jsonb, 0, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'key', 'TEXT', 'Ticket ID', 1),
    ('{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": null, "defaultValue": null}'::jsonb, 0, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'summary', 'TEXT', 'Summary', 1),
    ('{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": null, "defaultValue": null}'::jsonb, 1, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'story_key', 'TEXT', 'User Story ID', 1),
    ('{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": null, "defaultValue": null}'::jsonb, 1, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'story_summary', 'TEXT', 'User Story', 1),
    ('{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": null, "defaultValue": null}'::jsonb, 2, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'sprint_id', 'TEXT', 'Sprint', 1),
    ('{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": null, "defaultValue": null}'::jsonb, 2, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'assignee_email', 'TEXT', 'Assignee', 1),
    ('{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": null, "defaultValue": null}'::jsonb, 2, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'status', 'TEXT', 'Status', 1);

-- +goose Down
