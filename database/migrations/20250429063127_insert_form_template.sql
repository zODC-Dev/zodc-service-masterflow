-- +goose Up
INSERT INTO public.form_templates
(created_at, updated_at, deleted_at, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, 'Edit Profile Form', 'Edit Profile Form', 6, '{"Project": ["Prepify", "ZODC"], "Job title": ["Frontend Developer", "Backend Developer", "DevOps Engineer", "UI/UX Designer", "QA Engineer (Tester)", "Business Analyst (BA)"], "Seniority": ["Fresher (0-1)", "Junior (1-2)", "Mid-level (3-4)", "Senior (5+)", "", ""], "Technologies": ["Angular (14+)", "React", "Vue", "Nodejs", "Python", "Golang"]}'::jsonb, 'Edit Profile Form', 'manage_accounts|--error-40', NULL, 1, 'FORM', 'USER');
INSERT INTO public.form_templates
(created_at, updated_at, deleted_at, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, 'Hiring request form', 'Hiring request form', 1, '{"Project": ["Prepify", "ZODC"], "Job title": ["Frontend Developer", "Backend Developer", "DevOps Engineer", "UI/UX Designer", "QA Engineer (Tester)", "Business Analyst (BA)"], "Seniority": ["Fresher (0-1)", "Junior (1-2)", "Mid-level (3-4)", "Senior (5+)", "", ""], "Technologies": ["Angular (14+)", "React", "Vue", "Nodejs", "Python", "Golang"]}'::jsonb, 'Help Product Owners to formally request new team members by providing clear details about the required position, skills, timeline, and project context', 'group_add|--primary-40', NULL, 1, 'FORM', 'USER');
INSERT INTO public.form_templates
(created_at, updated_at, deleted_at, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES('2025-05-01 00:11:05.417', '2025-05-01 00:11:05.417', NULL, 'Resume Upload Form', 'Resume Upload Form', 1, NULL, 'Resume Upload Form', 'upload_file|--info-40', NULL, 1, 'FORM', 'USER');
INSERT INTO public.form_templates
(created_at, updated_at, deleted_at, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES('2025-05-01 00:11:05.420', '2025-05-01 00:11:05.420', NULL, 'Project Information Form', 'Project Information Form', 1, NULL, 'Project Information Form', 'folder_open|--warning-40', NULL, 1, 'FORM', 'USER');
INSERT INTO public.form_templates
(created_at, updated_at, deleted_at, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES('2025-05-01 00:11:05.437', '2025-05-01 00:11:05.437', NULL, 'Sprint Retrospective Form', 'Sprint Retrospective Form', 5, NULL, 'Sprint Retrospective Form', 'insights|--secondary-40', NULL, 1, 'FORM', 'USER');
INSERT INTO public.form_templates
(created_at, updated_at, deleted_at, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES('2025-05-01 00:11:05.456', '2025-05-01 00:11:05.456', NULL, ' Performace Evaluation Form', ' Performace Evaluation Form', 7, NULL, ' Performace Evaluation Form', 'fact_check|--tertiary-40', NULL, 1, 'FORM', 'USER');

INSERT INTO public.form_template_versions
(created_at, updated_at, deleted_at, "version", form_template_id)
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, 1, 2);
INSERT INTO public.form_template_versions
(created_at, updated_at, deleted_at, "version", form_template_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, 1, 3);
INSERT INTO public.form_template_versions
(created_at, updated_at, deleted_at, "version", form_template_id)
VALUES('2025-05-01 00:11:05.417', '2025-05-01 00:11:05.417', NULL, 1, 4);
INSERT INTO public.form_template_versions
(created_at, updated_at, deleted_at, "version", form_template_id)
VALUES('2025-05-01 00:11:05.420', '2025-05-01 00:11:05.420', NULL, 1, 5);
INSERT INTO public.form_template_versions
(created_at, updated_at, deleted_at, "version", form_template_id)
VALUES('2025-05-01 00:11:05.437', '2025-05-01 00:11:05.437', NULL, 1, 6);
INSERT INTO public.form_template_versions
(created_at, updated_at, deleted_at, "version", form_template_id)
VALUES('2025-05-01 00:11:05.456', '2025-05-01 00:11:05.456', NULL, 1, 7);

INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Title', 'title', 'de0650f5-70f3-4b59-863c-a73d7caa41f9', 'TITLE', 'Requestor Information', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{"options": [{"label": "Prepify", "value": "Prepify"}, {"label": "ZODC", "value": "ZODC"}], "tooltip": null, "helpText": null, "placeholder": "Project Name", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 1, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Project Name', 'MULTI_SELECT', 'Project Name', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{"maxRows": null, "minRows": 3, "tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Project Overview", "defaultValue": null}'::jsonb, 2, true, false, 'BASIC_FIELD', 'Text Area', 'article', 'Project Overview', 'TEXT_AREA', 'Project Overview', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{}'::jsonb, 3, false, false, 'BASIC_FIELD', 'Title', 'title', 'fad5ba8e-8956-42b2-941e-a232a2695d32', 'TITLE', 'Position Details', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{"options": [{"label": "Frontend Developer", "value": "Frontend Developer"}, {"label": "Backend Developer", "value": "Backend Developer"}, {"label": "DevOps Engineer", "value": "DevOps Engineer"}, {"label": "UI/UX Designer", "value": "UI/UX Designer"}, {"label": "QA Engineer (Tester)", "value": "QA Engineer (Tester)"}, {"label": "Business Analyst (BA)", "value": "Business Analyst (BA)"}], "tooltip": null, "helpText": null, "placeholder": "Job title", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 4, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Job title', 'MULTI_SELECT', 'Job title', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{"tooltip": null, "helpText": "Enter the number of positions in order of job title and separate with commas", "maxChars": null, "validation": null, "placeholder": "Number of Positions", "defaultValue": null}'::jsonb, 5, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'Number of Positions', 'TEXT', 'Number of Positions', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{"tooltip": null, "helpText": null, "validations": [], "defaultValue": {"endDate": null, "startDate": null}, "disableWeekdays": null, "endDatePlaceholder": null, "startDatePlaceholder": null}'::jsonb, 6, true, false, 'DATE_TIME_FIELD', 'Date Range', 'date_range', 'Project Duration', 'DATE_RANGE', 'Project Duration', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Scope of work", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 7, true, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'Scope of work', 'RICHTEXT', 'Scope of work', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{}'::jsonb, 8, false, false, 'BASIC_FIELD', 'Title', 'title', 'fb16891a-7b7a-4d7a-9032-7d1a97b5fd8c', 'TITLE', 'Job Requirements', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{"options": [{"label": "Angular (14+)", "value": "Angular (14+)"}, {"label": "React", "value": "React"}, {"label": "Vue", "value": "Vue"}, {"label": "Nodejs", "value": "Nodejs"}, {"label": "Python", "value": "Python"}, {"label": "Golang", "value": "Golang"}], "tooltip": null, "helpText": null, "placeholder": "Must have Skill/Knowledge", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 9, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Must have Skill/Knowledge', 'MULTI_SELECT', 'Must have Skill/Knowledge', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{"options": [{"label": "Angular (14+)", "value": "Angular (14+)"}, {"label": "React", "value": "React"}, {"label": "Vue", "value": "Vue"}, {"label": "Nodejs", "value": "Nodejs"}, {"label": "Python", "value": "Python"}, {"label": "Golang", "value": "Golang"}], "tooltip": null, "helpText": null, "placeholder": "Nice to have Skill/Knowledge", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 10, false, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Nice to have Skill/Knowledge', 'MULTI_SELECT', 'Nice to have Skill/Knowledge', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{"options": [{"label": "Fresher (0-1)", "value": "Fresher (0-1)"}, {"label": "Junior (1-2)", "value": "Junior (1-2)"}, {"label": "Mid-level (3-4)", "value": "Mid-level (3-4)"}, {"label": "Senior (5+)", "value": "Senior (5+)"}], "tooltip": null, "vertical": false, "defaultOptionIndex": null}'::jsonb, 11, true, false, 'OPTION_FIELD', 'Radio', 'radio_button_checked', 'Years of Experience', 'RADIO', 'Years of Experience', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{}'::jsonb, 12, false, false, 'BASIC_FIELD', 'Title', 'title', '9ffa4ffc-be27-4a81-8ba7-242d356fcaf7', 'TITLE', 'Contact', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Contactor Name", "defaultValue": null}'::jsonb, 13, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'Contactor Name', 'TEXT', 'Contactor Name', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.406', '2025-05-01 00:11:05.406', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": "EMAIL", "placeholder": null, "defaultValue": null}'::jsonb, 14, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'Email', 'TEXT', 'Email', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, '{}'::jsonb, 1, false, false, 'BASIC_FIELD', 'Title', 'title', 'eeb4335c-9cc4-46ee-a8d8-a8d087982569', 'TITLE', 'Personal Information', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Job Title", "defaultValue": null}'::jsonb, 2, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'jobTitle', 'TEXT', 'Job Title', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.420', '2025-05-01 00:11:05.420', NULL, '{}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Title', 'title', '7b01d49b-b7aa-4136-9151-0876cf40e1b3', 'TITLE', 'Project Information', 5);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Professional Summary", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 3, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'professionalSummary', 'RICHTEXT', 'Professional Summary', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, '{}'::jsonb, 4, false, false, 'BASIC_FIELD', 'Title', 'title', '92c693c6-32ca-4fe6-ab4e-d0f038732e26', 'TITLE', 'Contact Information', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Phone Number", "defaultValue": null}'::jsonb, 5, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'phoneNumber', 'TEXT', 'Phone Number', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, '{"maxRows": null, "minRows": 3, "tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Location", "defaultValue": null}'::jsonb, 6, true, false, 'BASIC_FIELD', 'Text Area', 'article', 'location', 'TEXT_AREA', 'Location', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, '{}'::jsonb, 7, false, false, 'BASIC_FIELD', 'Title', 'title', 'b640d930-1ee0-4320-9d6e-8aebfc15bc33', 'TITLE', 'Skills', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, '{"options": [{"label": "Angular (14+)", "value": "Angular (14+)"}, {"label": "React", "value": "React"}, {"label": "Vue", "value": "Vue"}, {"label": "Nodejs", "value": "Nodejs"}, {"label": "Python", "value": "Python"}, {"label": "Golang", "value": "Golang"}], "tooltip": null, "helpText": null, "placeholder": "Primary Skills", "validations": [], "inputNewOption": true, "preventInputDuplicateOption": true}'::jsonb, 8, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'primarySkills', 'MULTI_SELECT', 'Primary Skills', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, '{"options": [{"label": "Angular (14+)", "value": "Angular (14+)"}, {"label": "React", "value": "React"}, {"label": "Vue", "value": "Vue"}, {"label": "Nodejs", "value": "Nodejs"}, {"label": "Python", "value": "Python"}, {"label": "Golang", "value": "Golang"}], "tooltip": null, "helpText": null, "placeholder": "Secondary Skills", "validations": [], "inputNewOption": true, "preventInputDuplicateOption": true}'::jsonb, 9, false, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'secondarySkills', 'MULTI_SELECT', 'Secondary Skills', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, '{"options": [{"label": "Fresher (0-1)", "value": "Fresher (0-1)"}, {"label": "Junior (1-2)", "value": "Junior (1-2)"}, {"label": "Mid-level (3-4)", "value": "Mid-level (3-4)"}, {"label": "Senior (5+)", "value": "Senior (5+)"}], "tooltip": null, "vertical": false, "defaultOptionIndex": null}'::jsonb, 10, true, false, 'OPTION_FIELD', 'Radio', 'radio_button_checked', 'yearsOfExperience', 'RADIO', 'Years of Experience', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, '{}'::jsonb, 11, false, false, 'BASIC_FIELD', 'Title', 'title', '3e4baa6a-191a-4c34-8bc7-4badd1873c9f', 'TITLE', 'Education & Certification', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, '{"maxRows": null, "minRows": 3, "tooltip": null, "helpText": "List one education per line in the format: Name (School, Year)", "maxChars": null, "validation": null, "placeholder": "Education", "defaultValue": null}'::jsonb, 12, false, false, 'BASIC_FIELD', 'Text Area', 'article', 'education', 'TEXT_AREA', 'Education', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.402', '2025-05-01 00:11:05.402', NULL, '{"maxRows": null, "minRows": 3, "tooltip": null, "helpText": "List one certification per line in the format: Name (Issuer, Year)", "maxChars": null, "validation": null, "placeholder": "Certification", "defaultValue": null}'::jsonb, 13, false, false, 'BASIC_FIELD', 'Text Area', 'article', 'certification', 'TEXT_AREA', 'Certification', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.417', '2025-05-01 00:11:05.417', NULL, '{}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Title', 'title', '7b01d49b-b7aa-4136-9151-0876cf40e1b3', 'TITLE', 'Project Information', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.417', '2025-05-01 00:11:05.417', NULL, '{"options": [{"label": "Prepify", "value": "Prepify"}, {"label": "ZODC", "value": "ZODC"}], "tooltip": null, "helpText": null, "placeholder": "Project Name", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 1, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Project Name', 'MULTI_SELECT', 'Project Name', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.417', '2025-05-01 00:11:05.417', NULL, '{}'::jsonb, 2, false, false, 'BASIC_FIELD', 'Title', 'title', 'd89e2a18-7176-4692-9875-cc6436f3f63f', 'TITLE', 'Candidate Information', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.417', '2025-05-01 00:11:05.417', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Candidate Name", "defaultValue": null}'::jsonb, 3, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'Candidate Name', 'TEXT', 'Candidate Name', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.417', '2025-05-01 00:11:05.417', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": "EMAIL", "placeholder": "Email Address", "defaultValue": null}'::jsonb, 4, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'Email Address', 'TEXT', 'Email Address', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.417', '2025-05-01 00:11:05.417', NULL, '{"options": [{"label": "Frontend Developer", "value": "Frontend Developer"}, {"label": "Backend Developer", "value": "Backend Developer"}, {"label": "DevOps Engineer", "value": "DevOps Engineer"}, {"label": "UI/UX Designer", "value": "UI/UX Designer"}, {"label": "QA Engineer (Tester)", "value": "QA Engineer (Tester)"}, {"label": "Business Analyst (BA)", "value": "Business Analyst (BA)"}], "tooltip": null, "helpText": null, "placeholder": "Job title", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 5, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Job title', 'MULTI_SELECT', 'Job title', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.417', '2025-05-01 00:11:05.417', NULL, '{"options": [{"label": "Fresher (0-1)", "value": "Fresher (0-1)"}, {"label": "Junior (1-2)", "value": "Junior (1-2)"}, {"label": "Mid-level (3-4)", "value": "Mid-level (3-4)"}, {"label": "Senior (5+)", "value": "Senior (5+)"}], "tooltip": null, "vertical": false, "defaultOptionIndex": null}'::jsonb, 6, true, false, 'OPTION_FIELD', 'Radio', 'radio_button_checked', 'Years of Experience', 'RADIO', 'Years of Experience', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.417', '2025-05-01 00:11:05.417', NULL, '{"tooltip": null, "helpText": null, "placeholder": "Interview Date", "validations": [], "defaultValue": null, "disableWeekdays": null}'::jsonb, 7, true, false, 'DATE_TIME_FIELD', 'Date', 'today', 'Interview Date', 'DATE', 'Interview Date', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.417', '2025-05-01 00:11:05.417', NULL, '{"tooltip": null, "helpText": null, "allowedFileType": ["pdf"]}'::jsonb, 8, true, false, 'ADVANCED_FIELD', 'Attachment', 'attach_file', 'Resume File Upload', 'ATTACHMENT', 'Resume File Upload', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.417', '2025-05-01 00:11:05.417', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Interview Notes", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 9, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'Interview Notes', 'RICHTEXT', 'Interview Notes', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.420', '2025-05-01 00:11:05.420', NULL, '{"options": [{"label": "Prepify", "value": "Prepify"}, {"label": "ZODC", "value": "ZODC"}], "tooltip": null, "helpText": null, "placeholder": "Project Name", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 1, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Project Name', 'MULTI_SELECT', 'Project Name', 5);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.437', '2025-05-01 00:11:05.437', NULL, '{}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Title', 'title', '8ea31b45-5b09-413f-8dd0-88aa55324434', 'TITLE', 'Team feedback and insights', 6);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.437', '2025-05-01 00:11:05.437', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "What Went Well", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 1, true, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'whatWentWell', 'RICHTEXT', 'What Went Well', 6);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.437', '2025-05-01 00:11:05.437', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "What Could Be Improved", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 2, true, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'whatCouldBeImproved', 'RICHTEXT', 'What Could Be Improved', 6);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.437', '2025-05-01 00:11:05.437', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Action Items", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 3, true, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'actionItems', 'RICHTEXT', 'Action Items', 6);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.456', '2025-05-01 00:11:05.456', NULL, '{}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Title', 'title', 'd9b11959-6775-41d5-9cd6-b024b76f7926', 'TITLE', 'Performance Rating', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.456', '2025-05-01 00:11:05.456', NULL, '{"options": [{"label": "Good", "value": "4"}, {"label": "Average", "value": "3"}, {"label": "Poor", "value": "2"}, {"label": "Very Poor", "value": "1"}], "tooltip": null, "helpText": null, "placeholder": "Code Quality"}'::jsonb, 1, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'codeQuality', 'DROPDOWN', 'Code Quality', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.456', '2025-05-01 00:11:05.456', NULL, '{"options": [{"label": "Good", "value": "4"}, {"label": "Average", "value": "3"}, {"label": "Poor", "value": "2"}, {"label": "Very Poor", "value": "1"}], "tooltip": null, "helpText": null, "placeholder": "Test Coverage"}'::jsonb, 1, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'testCoverage', 'DROPDOWN', 'Test Coverage', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.456', '2025-05-01 00:11:05.456', NULL, '{"options": [{"label": "Good", "value": "4"}, {"label": "Average", "value": "3"}, {"label": "Poor", "value": "2"}, {"label": "Very Poor", "value": "1"}], "tooltip": null, "helpText": null, "placeholder": "Documentation Quality"}'::jsonb, 2, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'documentation', 'DROPDOWN', 'Documentation Quality', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.456', '2025-05-01 00:11:05.456', NULL, '{"options": [{"label": "Good", "value": "4"}, {"label": "Average", "value": "3"}, {"label": "Poor", "value": "2"}, {"label": "Very Poor", "value": "1"}], "tooltip": null, "helpText": null, "placeholder": "Overall Performance"}'::jsonb, 2, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'overall', 'DROPDOWN', 'Overall Performance', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.456', '2025-05-01 00:11:05.456', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Key Strengths", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 3, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'strengths', 'RICHTEXT', 'Key Strengths', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.456', '2025-05-01 00:11:05.456', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Areas for Improvement", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 4, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'areasForImprovement', 'RICHTEXT', 'Areas for Improvement', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-01 00:11:05.456', '2025-05-01 00:11:05.456', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Detailed Feedback", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 5, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'feedback', 'RICHTEXT', 'Detailed Feedback', 7);
-- +goose Down
