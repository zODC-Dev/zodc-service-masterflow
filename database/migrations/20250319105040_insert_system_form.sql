-- +goose Up
-- Jira Form Template
INSERT INTO public.form_templates
    (id, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES
    (1, 'Jira Task Report', 'Jira Task Report', 1, NULL, 'Jira Task Report', 'task|--primary-40', NULL, 1, 'TASK', 'SYSTEM');

SELECT setval('form_templates_id_seq', (SELECT MAX(id) FROM form_templates));



INSERT INTO public.form_template_versions
    (id, "version", form_template_id)
VALUES
    (1, 1, 1);

SELECT setval('form_template_versions_id_seq', (SELECT MAX(id) FROM form_template_versions));



INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-03 19:57:02.794', '2025-05-03 19:57:02.794', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Ticket ID", "defaultValue": null}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'key', 'TEXT', 'Ticket ID', 1);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-03 19:57:02.794', '2025-05-03 19:57:02.794', NULL, '{"maxRows": null, "minRows": 3, "tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Summary", "defaultValue": null}'::jsonb, 1, false, false, 'BASIC_FIELD', 'Text Area', 'article', 'summary', 'TEXT_AREA', 'Summary', 1);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-03 19:57:02.794', '2025-05-03 19:57:02.794', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Sprint", "defaultValue": null}'::jsonb, 2, false, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'sprint::__DELIMITER__::name', 'TEXT', 'Sprint', 1);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-03 19:57:02.794', '2025-05-03 19:57:02.794', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Assignee", "defaultValue": null}'::jsonb, 2, false, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'assignee::__DELIMITER__::name', 'TEXT', 'Assignee', 1);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-03 19:57:02.794', '2025-05-03 19:57:02.794', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Status", "defaultValue": null}'::jsonb, 2, false, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'status', 'TEXT', 'Status', 1);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-03 19:57:02.794', '2025-05-03 19:57:02.794', NULL, '{"tooltip": null, "helpText": null, "placeholder": "Start Date", "validations": [], "defaultValue": null, "disableWeekdays": null}'::jsonb, 3, false, false, 'DATE_TIME_FIELD', 'Date', 'today', 'createdAt', 'DATE', 'Start Date', 1);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-03 19:57:02.794', '2025-05-03 19:57:02.794', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "placeholder": "Story Point Estimate", "validations": [], "defaultValue": null, "decimalPlaces": null}'::jsonb, 3, false, false, 'BASIC_FIELD', 'Number', 'looks_one', 'estimatePoint', 'NUMBER', 'Story Point Estimate', 1);


-- +goose Down
