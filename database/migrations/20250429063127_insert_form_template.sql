-- +goose Up
-- Profile Form Template
INSERT INTO public.form_templates
    (id, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES
    (2, 'Edit Profile', 'Edit Profile', 6, '{"Project": ["Prepify", "ZODC"], "Job title": ["Frontend Developer", "Backend Developer", "DevOps Engineer", "UI/UX Designer", "QA Engineer (Tester)", "Business Analyst (BA)"], "Seniority": ["Fresher (0-1)", "Junior (1-2)", "Mid-level (3-4)", "Senior (5+)", "", ""], "Technologies": ["Angular (14+)", "React", "Vue", "Nodejs", "Python", "Golang"]}'::jsonb, 'Edit Profile', 'person_add|--primary-40', NULL, 1, 'FORM', 'USER');

SELECT setval('form_templates_id_seq', (SELECT MAX(id) FROM form_templates));




INSERT INTO public.form_templates
    (id,file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES
    (3, 'Performance Evaluation Form', 'Performance Evaluation Form', 7, NULL, 'Performance Evaluation Form', 'person_add|--primary-40', NULL, 1, 'FORM', 'USER');

SELECT setval('form_templates_id_seq', (SELECT MAX(id) FROM form_templates));



INSERT INTO public.form_template_versions
    (id, "version", form_template_id)
VALUES
    (2, 1, 2);

SELECT setval('form_template_versions_id_seq', (SELECT MAX(id) FROM form_template_versions));




INSERT INTO public.form_template_versions
    (id, "version", form_template_id)
VALUES
    (3, 1, 3);

SELECT setval('form_template_versions_id_seq', (SELECT MAX(id) FROM form_template_versions));




INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:19.035', '2025-04-29 13:40:19.035', NULL, '{}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Title', 'title', '003bef89-d35a-4106-838c-cc0da6e53684', 'TITLE', 'Personal Information', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:19.035', '2025-04-29 13:40:19.035', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Job Title", "defaultValue": null}'::jsonb, 1, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'jobTitle', 'TEXT', 'Job Title', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:19.035', '2025-04-29 13:40:19.035', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Professional Summary", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 2, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'professionalSummary', 'RICHTEXT', 'Professional Summary', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:19.035', '2025-04-29 13:40:19.035', NULL, '{}'::jsonb, 3, false, false, 'BASIC_FIELD', 'Title', 'title', 'e578aa36-adc6-46d1-9262-063a14a12087', 'TITLE', 'Contact Information', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:19.035', '2025-04-29 13:40:19.035', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "placeholder": "Phone Number", "validations": [], "defaultValue": null, "decimalPlaces": null}'::jsonb, 4, false, false, 'BASIC_FIELD', 'Number', 'looks_one', 'phoneNumber', 'NUMBER', 'Phone Number', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:19.035', '2025-04-29 13:40:19.035', NULL, '{"maxRows": null, "minRows": 3, "tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Location", "defaultValue": null}'::jsonb, 5, false, false, 'BASIC_FIELD', 'Text Area', 'article', 'location', 'TEXT_AREA', 'Location', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:19.035', '2025-04-29 13:40:19.035', NULL, '{}'::jsonb, 6, false, false, 'BASIC_FIELD', 'Title', 'title', '7b71f828-8fa0-4559-9a4c-f0e335a70a72', 'TITLE', 'Skills', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:19.035', '2025-04-29 13:40:19.035', NULL, '{"options": [{"label": "Angular (14+)", "value": "Angular (14+)"}, {"label": "React", "value": "React"}, {"label": "Vue", "value": "Vue"}, {"label": "Nodejs", "value": "Nodejs"}, {"label": "Python", "value": "Python"}, {"label": "Golang", "value": "Golang"}], "tooltip": null, "helpText": null, "placeholder": "Primary Skills", "validations": [], "inputNewOption": true, "preventInputDuplicateOption": true}'::jsonb, 7, false, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'primarySkills', 'MULTI_SELECT', 'Primary Skills', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:19.035', '2025-04-29 13:40:19.035', NULL, '{"options": [{"label": "Angular (14+)", "value": "Angular (14+)"}, {"label": "React", "value": "React"}, {"label": "Vue", "value": "Vue"}, {"label": "Nodejs", "value": "Nodejs"}, {"label": "Python", "value": "Python"}, {"label": "Golang", "value": "Golang"}], "tooltip": null, "helpText": null, "placeholder": "Secondary Skills", "validations": [], "inputNewOption": true, "preventInputDuplicateOption": true}'::jsonb, 8, false, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'secondarySkills', 'MULTI_SELECT', 'Secondary Skills', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:19.035', '2025-04-29 13:40:19.035', NULL, '{"options": [{"label": "Fresher (0-1)", "value": "Fresher (0-1)"}, {"label": "Junior (1-2)", "value": "Junior (1-2)"}, {"label": "Mid-level (3-4)", "value": "Mid-level (3-4)"}, {"label": "Senior (5+)", "value": "Senior (5+)"}], "tooltip": null, "vertical": false, "defaultOptionIndex": null}'::jsonb, 9, false, false, 'OPTION_FIELD', 'Radio', 'radio_button_checked', 'yearsOfExperience', 'RADIO', 'Years of Experience', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:19.035', '2025-04-29 13:40:19.035', NULL, '{}'::jsonb, 10, false, false, 'BASIC_FIELD', 'Title', 'title', '3aa91c42-8828-42a2-9dfc-eeab08fb4c46', 'TITLE', 'Education & Certification', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:19.035', '2025-04-29 13:40:19.035', NULL, '{"maxRows": null, "minRows": 3, "tooltip": null, "helpText": "List one education per line in the format: Name (School, Year)", "maxChars": null, "validation": null, "placeholder": "Education", "defaultValue": null}'::jsonb, 11, false, false, 'BASIC_FIELD', 'Text Area', 'article', 'education', 'TEXT_AREA', 'Education', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:19.035', '2025-04-29 13:40:19.035', NULL, '{"maxRows": null, "minRows": 3, "tooltip": null, "helpText": "List one certification per line in the format: Name (Issuer, Year)", "maxChars": null, "validation": null, "placeholder": "Certification", "defaultValue": null}'::jsonb, 12, false, false, 'BASIC_FIELD', 'Text Area', 'article', 'certification', 'TEXT_AREA', 'Certification', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:55.034', '2025-04-29 13:40:55.034', NULL, '{}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Title', 'title', 'baad30ec-ad2d-4986-98cb-d0ba90f3dbc2', 'TITLE', 'Performance Rating', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:55.034', '2025-04-29 13:40:55.034', NULL, '{"options": [{"label": "Good", "value": "4"}, {"label": "Average", "value": "3"}, {"label": "Poor", "value": "2"}, {"label": "Very Poor", "value": "1"}], "tooltip": null, "helpText": null, "placeholder": "Code Quality"}'::jsonb, 1, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'codeQuality', 'DROPDOWN', 'Code Quality', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:55.034', '2025-04-29 13:40:55.034', NULL, '{"options": [{"label": "Good", "value": "4"}, {"label": "Average", "value": "3"}, {"label": "Poor", "value": "2"}, {"label": "Very Poor", "value": "1"}], "tooltip": null, "helpText": null, "placeholder": "Test Coverage"}'::jsonb, 1, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'testCoverage', 'DROPDOWN', 'Test Coverage', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:55.034', '2025-04-29 13:40:55.034', NULL, '{"options": [{"label": "Good", "value": "4"}, {"label": "Average", "value": "3"}, {"label": "Poor", "value": "2"}, {"label": "Very Poor", "value": "1"}], "tooltip": null, "helpText": null, "placeholder": "Documentation Quality"}'::jsonb, 2, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'documentationQuality', 'DROPDOWN', 'Documentation Quality', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:55.034', '2025-04-29 13:40:55.034', NULL, '{"options": [{"label": "Good", "value": "4"}, {"label": "Average", "value": "3"}, {"label": "Poor", "value": "2"}, {"label": "Very Poor", "value": "1"}], "tooltip": null, "helpText": null, "placeholder": "Overall Performance"}'::jsonb, 2, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'overallPerformance', 'DROPDOWN', 'Overall Performance', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:55.034', '2025-04-29 13:40:55.034', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Key Strengths", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 3, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'keyStrengths', 'RICHTEXT', 'Key Strengths', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:55.034', '2025-04-29 13:40:55.034', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Areas for Improvement", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 4, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'areasForImprovement', 'RICHTEXT', 'Areas for Improvement', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-04-29 13:40:55.034', '2025-04-29 13:40:55.034', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Detailed Feedback", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 5, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'detailedFeedback', 'RICHTEXT', 'Detailed Feedback', 3);


-- +goose Down
