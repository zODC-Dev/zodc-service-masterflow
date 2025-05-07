-- +goose Up
INSERT INTO public.form_templates
(created_at, updated_at, deleted_at, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, 'Edit Profile Form', 'Edit Profile Form', 6, '{"Project": ["Prepify", "ZODC"], "Job title": ["Frontend Developer", "Backend Developer", "DevOps Engineer", "UI/UX Designer", "QA Engineer (Tester)", "Business Analyst (BA)"], "Seniority": ["Fresher (0-1)", "Junior (1-2)", "Mid-level (3-4)", "Senior (5+)", "", ""], "Technologies": ["Angular (14+)", "React", "Vue", "Nodejs", "Python", "Golang"]}'::jsonb, 'Edit Profile Form', 'manage_accounts|--error-40', NULL, 1, 'FORM', 'USER');
INSERT INTO public.form_templates
(created_at, updated_at, deleted_at, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, 'Hiring request form', 'Hiring request form', 1, '{"Project": ["Prepify", "ZODC"], "Job title": ["Frontend Developer", "Backend Developer", "DevOps Engineer", "UI/UX Designer", "QA Engineer (Tester)", "Business Analyst (BA)"], "Seniority": ["Fresher (0-1)", "Junior (1-2)", "Mid-level (3-4)", "Senior (5+)", "", ""], "Technologies": ["Angular (14+)", "React", "Vue", "Nodejs", "Python", "Golang"]}'::jsonb, 'Help Product Owners to formally request new team members by providing clear details about the required position, skills, timeline, and project context', 'group_add|--primary-40', NULL, 1, 'FORM', 'USER');
INSERT INTO public.form_templates
(created_at, updated_at, deleted_at, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES('2025-05-07 23:50:56.813', '2025-05-07 23:50:56.813', NULL, ' Performace Evaluation Form', ' Performace Evaluation Form', 7, NULL, ' Performace Evaluation Form', 'fact_check|--tertiary-40', NULL, 1, 'FORM', 'USER');
INSERT INTO public.form_templates
(created_at, updated_at, deleted_at, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES('2025-05-07 23:50:57.011', '2025-05-07 23:50:57.011', NULL, 'Project Information Form', 'Project Information Form', 1, NULL, 'Project Information Form', 'folder_open|--warning-40', NULL, 1, 'FORM', 'USER');
INSERT INTO public.form_templates
(created_at, updated_at, deleted_at, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES('2025-05-07 23:50:57.203', '2025-05-07 23:50:57.203', NULL, 'Sprint Retrospective Form', 'Sprint Retrospective Form', 5, NULL, 'Sprint Retrospective Form', 'insights|--secondary-40', NULL, 1, 'FORM', 'USER');
INSERT INTO public.form_templates
(created_at, updated_at, deleted_at, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES('2025-05-07 23:50:57.396', '2025-05-07 23:50:57.396', NULL, '', 'Resume Upload Form', 1, '{"Project": ["Prepify", "ZODC"], "Job title": ["Frontend Developer", "Backend Developer", "DevOps Engineer", "UI/UX Designer", "QA Engineer (Tester)", "Business Analyst (BA)"], "Seniority": ["Fresher (0-1)", "Junior (1-2)", "Mid-level (3-4)", "Senior (5+)", "", ""], "Technologies": ["Angular (14+)", "React", "Vue", "Nodejs", "Python", "Golang"]}'::jsonb, 'Resume Upload Form', 'file_upload|--info-40', NULL, 1, 'FORM', 'USER');
INSERT INTO public.form_templates
(created_at, updated_at, deleted_at, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '', 'Results of Training Session', 2, '{"Project": ["Prepify", "ZODC", "", "", "", "", "", "", "", ""], "Job title": ["Frontend Developer", "Backend Developer", "DevOps Engineer", "UI/UX Designer", "QA Engineer (Tester)", "Business Analyst (BA)", "", "", "", ""], "Seniority": ["Fresher (0-1)", "Junior (1-2)", "Mid-level (3-4)", "Senior (5+)", "", "", "", "", "", ""], "Participant": ["Trần, Mạnh Hùng", "Lê, Bảo", "Nguyễn, Hải Đăng ", "Lê, Thị Mẫn Nhi", "Dương, Ngọc Mạnh", "Đỗ, Quân ", "Văn , Phú Hòa", "Lâm, Thị Ngọc Hân", "Trần, Tấn Thành", "Dương, Hoàng Nam"], "Technologies": ["Angular (14+)", "React", "Vue", "Nodejs", "Python", "Golang", "", "", "", ""]}'::jsonb, 'Submit Issues / Results of Training Session', 'school|--warning-40', NULL, 1, 'FORM', 'USER');
INSERT INTO public.form_templates
(created_at, updated_at, deleted_at, file_name, title, category_id, data_sheet, description, decoration, template_id, current_version, tag, "type")
VALUES('2025-05-07 23:50:57.784', '2025-05-07 23:50:57.784', NULL, '', 'Leave Request Form', 8, '{"Project": ["Prepify", "ZODC", "", "", "", "", "", "", "", ""], "Job title": ["Frontend Developer", "Backend Developer", "DevOps Engineer", "UI/UX Designer", "QA Engineer (Tester)", "Business Analyst (BA)", "", "", "", ""], "Seniority": ["Fresher (0-1)", "Junior (1-2)", "Mid-level (3-4)", "Senior (5+)", "", "", "", "", "", ""], "Participant": ["Trần, Mạnh Hùng", "Lê, Bảo", "Nguyễn, Hải Đăng ", "Lê, Thị Mẫn Nhi", "Dương, Ngọc Mạnh", "Đỗ, Quân ", "Văn , Phú Hòa", "Lâm, Thị Ngọc Hân", "Trần, Tấn Thành", "Dương, Hoàng Nam"], "Technologies": ["Angular (14+)", "React", "Vue", "Nodejs", "Python", "Golang", "", "", "", ""]}'::jsonb, 'This Leave Request Form is used by employees to formally submit a request for time off from work.', 'drive_file_move|--success-40', NULL, 1, 'FORM', 'USER');




INSERT INTO public.form_template_versions
(created_at, updated_at, deleted_at, "version", form_template_id)
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, 1, 2);
INSERT INTO public.form_template_versions
(created_at, updated_at, deleted_at, "version", form_template_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, 1, 3);
INSERT INTO public.form_template_versions
(created_at, updated_at, deleted_at, "version", form_template_id)
VALUES('2025-05-07 23:50:56.813', '2025-05-07 23:50:56.813', NULL, 1, 4);
INSERT INTO public.form_template_versions
(created_at, updated_at, deleted_at, "version", form_template_id)
VALUES('2025-05-07 23:50:57.011', '2025-05-07 23:50:57.011', NULL, 1, 5);
INSERT INTO public.form_template_versions
(created_at, updated_at, deleted_at, "version", form_template_id)
VALUES('2025-05-07 23:50:57.203', '2025-05-07 23:50:57.203', NULL, 1, 6);
INSERT INTO public.form_template_versions
(created_at, updated_at, deleted_at, "version", form_template_id)
VALUES('2025-05-07 23:50:57.396', '2025-05-07 23:50:57.396', NULL, 1, 7);
INSERT INTO public.form_template_versions
(created_at, updated_at, deleted_at, "version", form_template_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, 1, 8);
INSERT INTO public.form_template_versions
(created_at, updated_at, deleted_at, "version", form_template_id)
VALUES('2025-05-07 23:50:57.784', '2025-05-07 23:50:57.784', NULL, 1, 9);




INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, '{}'::jsonb, 1, false, false, 'BASIC_FIELD', 'Title', 'title', 'eeb4335c-9cc4-46ee-a8d8-a8d087982569', 'TITLE', 'Personal Information', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Job Title", "defaultValue": null}'::jsonb, 2, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'jobTitle', 'TEXT', 'Job Title', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Professional Summary", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 3, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'professionalSummary', 'RICHTEXT', 'Professional Summary', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, '{}'::jsonb, 4, false, false, 'BASIC_FIELD', 'Title', 'title', '92c693c6-32ca-4fe6-ab4e-d0f038732e26', 'TITLE', 'Contact Information', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Phone Number", "defaultValue": null}'::jsonb, 5, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'phoneNumber', 'TEXT', 'Phone Number', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, '{"maxRows": null, "minRows": 3, "tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Location", "defaultValue": null}'::jsonb, 6, true, false, 'BASIC_FIELD', 'Text Area', 'article', 'location', 'TEXT_AREA', 'Location', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, '{}'::jsonb, 7, false, false, 'BASIC_FIELD', 'Title', 'title', 'b640d930-1ee0-4320-9d6e-8aebfc15bc33', 'TITLE', 'Skills', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, '{"options": [{"label": "Angular (14+)", "value": "Angular (14+)"}, {"label": "React", "value": "React"}, {"label": "Vue", "value": "Vue"}, {"label": "Nodejs", "value": "Nodejs"}, {"label": "Python", "value": "Python"}, {"label": "Golang", "value": "Golang"}], "tooltip": null, "helpText": null, "placeholder": "Primary Skills", "validations": [], "inputNewOption": true, "preventInputDuplicateOption": true}'::jsonb, 8, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'primarySkills', 'MULTI_SELECT', 'Primary Skills', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, '{"options": [{"label": "Angular (14+)", "value": "Angular (14+)"}, {"label": "React", "value": "React"}, {"label": "Vue", "value": "Vue"}, {"label": "Nodejs", "value": "Nodejs"}, {"label": "Python", "value": "Python"}, {"label": "Golang", "value": "Golang"}], "tooltip": null, "helpText": null, "placeholder": "Secondary Skills", "validations": [], "inputNewOption": true, "preventInputDuplicateOption": true}'::jsonb, 9, false, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'secondarySkills', 'MULTI_SELECT', 'Secondary Skills', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, '{"options": [{"label": "Fresher (0-1)", "value": "Fresher (0-1)"}, {"label": "Junior (1-2)", "value": "Junior (1-2)"}, {"label": "Mid-level (3-4)", "value": "Mid-level (3-4)"}, {"label": "Senior (5+)", "value": "Senior (5+)"}], "tooltip": null, "vertical": false, "defaultOptionIndex": null}'::jsonb, 10, true, false, 'OPTION_FIELD', 'Radio', 'radio_button_checked', 'yearsOfExperience', 'RADIO', 'Years of Experience', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, '{}'::jsonb, 11, false, false, 'BASIC_FIELD', 'Title', 'title', '3e4baa6a-191a-4c34-8bc7-4badd1873c9f', 'TITLE', 'Education & Certification', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, '{"maxRows": null, "minRows": 3, "tooltip": null, "helpText": "List one education per line in the format: Name (School, Year)", "maxChars": null, "validation": null, "placeholder": "Education", "defaultValue": null}'::jsonb, 12, false, false, 'BASIC_FIELD', 'Text Area', 'article', 'education', 'TEXT_AREA', 'Education', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.451', '2025-05-07 23:50:56.451', NULL, '{"maxRows": null, "minRows": 3, "tooltip": null, "helpText": "List one certification per line in the format: Name (Issuer, Year)", "maxChars": null, "validation": null, "placeholder": "Certification", "defaultValue": null}'::jsonb, 13, false, false, 'BASIC_FIELD', 'Text Area', 'article', 'certification', 'TEXT_AREA', 'Certification', 2);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Title', 'title', 'de0650f5-70f3-4b59-863c-a73d7caa41f9', 'TITLE', 'Requestor Information', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{"options": [{"label": "Prepify", "value": "Prepify"}, {"label": "ZODC", "value": "ZODC"}], "tooltip": null, "helpText": null, "placeholder": "Project Name", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 1, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Project Name', 'MULTI_SELECT', 'Project Name', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{"maxRows": null, "minRows": 3, "tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Project Overview", "defaultValue": null}'::jsonb, 2, true, false, 'BASIC_FIELD', 'Text Area', 'article', 'Project Overview', 'TEXT_AREA', 'Project Overview', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{}'::jsonb, 3, false, false, 'BASIC_FIELD', 'Title', 'title', 'fad5ba8e-8956-42b2-941e-a232a2695d32', 'TITLE', 'Position Details', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{"options": [{"label": "Frontend Developer", "value": "Frontend Developer"}, {"label": "Backend Developer", "value": "Backend Developer"}, {"label": "DevOps Engineer", "value": "DevOps Engineer"}, {"label": "UI/UX Designer", "value": "UI/UX Designer"}, {"label": "QA Engineer (Tester)", "value": "QA Engineer (Tester)"}, {"label": "Business Analyst (BA)", "value": "Business Analyst (BA)"}], "tooltip": null, "helpText": null, "placeholder": "Job title", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 4, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Job title', 'MULTI_SELECT', 'Job title', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{"tooltip": null, "helpText": "Enter the number of positions in order of job title and separate with commas", "maxChars": null, "validation": null, "placeholder": "Number of Positions", "defaultValue": null}'::jsonb, 5, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'Number of Positions', 'TEXT', 'Number of Positions', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{"tooltip": null, "helpText": null, "validations": [], "defaultValue": {"endDate": null, "startDate": null}, "disableWeekdays": null, "endDatePlaceholder": null, "startDatePlaceholder": null}'::jsonb, 6, true, false, 'DATE_TIME_FIELD', 'Date Range', 'date_range', 'Project Duration', 'DATE_RANGE', 'Project Duration', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Scope of work", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 7, true, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'Scope of work', 'RICHTEXT', 'Scope of work', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{}'::jsonb, 8, false, false, 'BASIC_FIELD', 'Title', 'title', 'fb16891a-7b7a-4d7a-9032-7d1a97b5fd8c', 'TITLE', 'Job Requirements', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{"options": [{"label": "Angular (14+)", "value": "Angular (14+)"}, {"label": "React", "value": "React"}, {"label": "Vue", "value": "Vue"}, {"label": "Nodejs", "value": "Nodejs"}, {"label": "Python", "value": "Python"}, {"label": "Golang", "value": "Golang"}], "tooltip": null, "helpText": null, "placeholder": "Must have Skill/Knowledge", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 9, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Must have Skill/Knowledge', 'MULTI_SELECT', 'Must have Skill/Knowledge', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{"options": [{"label": "Angular (14+)", "value": "Angular (14+)"}, {"label": "React", "value": "React"}, {"label": "Vue", "value": "Vue"}, {"label": "Nodejs", "value": "Nodejs"}, {"label": "Python", "value": "Python"}, {"label": "Golang", "value": "Golang"}], "tooltip": null, "helpText": null, "placeholder": "Nice to have Skill/Knowledge", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 10, false, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Nice to have Skill/Knowledge', 'MULTI_SELECT', 'Nice to have Skill/Knowledge', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{"options": [{"label": "Fresher (0-1)", "value": "Fresher (0-1)"}, {"label": "Junior (1-2)", "value": "Junior (1-2)"}, {"label": "Mid-level (3-4)", "value": "Mid-level (3-4)"}, {"label": "Senior (5+)", "value": "Senior (5+)"}], "tooltip": null, "vertical": false, "defaultOptionIndex": null}'::jsonb, 11, true, false, 'OPTION_FIELD', 'Radio', 'radio_button_checked', 'Years of Experience', 'RADIO', 'Years of Experience', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{}'::jsonb, 12, false, false, 'BASIC_FIELD', 'Title', 'title', '9ffa4ffc-be27-4a81-8ba7-242d356fcaf7', 'TITLE', 'Contact', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Contactor Name", "defaultValue": null}'::jsonb, 13, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'Contactor Name', 'TEXT', 'Contactor Name', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.627', '2025-05-07 23:50:56.627', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": "EMAIL", "placeholder": null, "defaultValue": null}'::jsonb, 14, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'Email', 'TEXT', 'Email', 3);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.813', '2025-05-07 23:50:56.813', NULL, '{}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Title', 'title', 'd9b11959-6775-41d5-9cd6-b024b76f7926', 'TITLE', 'Performance Rating', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.813', '2025-05-07 23:50:56.813', NULL, '{"options": [{"label": "Good", "value": "4"}, {"label": "Average", "value": "3"}, {"label": "Poor", "value": "2"}, {"label": "Very Poor", "value": "1"}], "tooltip": null, "helpText": null, "placeholder": "Code Quality"}'::jsonb, 1, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'codeQuality', 'DROPDOWN', 'Code Quality', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.813', '2025-05-07 23:50:56.813', NULL, '{"options": [{"label": "Good", "value": "4"}, {"label": "Average", "value": "3"}, {"label": "Poor", "value": "2"}, {"label": "Very Poor", "value": "1"}], "tooltip": null, "helpText": null, "placeholder": "Test Coverage"}'::jsonb, 1, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'testCoverage', 'DROPDOWN', 'Test Coverage', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.813', '2025-05-07 23:50:56.813', NULL, '{"options": [{"label": "Good", "value": "4"}, {"label": "Average", "value": "3"}, {"label": "Poor", "value": "2"}, {"label": "Very Poor", "value": "1"}], "tooltip": null, "helpText": null, "placeholder": "Documentation Quality"}'::jsonb, 2, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'documentation', 'DROPDOWN', 'Documentation Quality', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.813', '2025-05-07 23:50:56.813', NULL, '{"options": [{"label": "Good", "value": "4"}, {"label": "Average", "value": "3"}, {"label": "Poor", "value": "2"}, {"label": "Very Poor", "value": "1"}], "tooltip": null, "helpText": null, "placeholder": "Overall Performance"}'::jsonb, 2, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'overall', 'DROPDOWN', 'Overall Performance', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.813', '2025-05-07 23:50:56.813', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Key Strengths", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 3, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'strength', 'RICHTEXT', 'Key Strengths', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.813', '2025-05-07 23:50:56.813', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Areas for Improvement", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 4, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'areasForImprovement', 'RICHTEXT', 'Areas for Improvement', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:56.813', '2025-05-07 23:50:56.813', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Detailed Feedback", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 5, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'feedback', 'RICHTEXT', 'Detailed Feedback', 4);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.011', '2025-05-07 23:50:57.011', NULL, '{}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Title', 'title', '7b01d49b-b7aa-4136-9151-0876cf40e1b3', 'TITLE', 'Project Information', 5);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.011', '2025-05-07 23:50:57.011', NULL, '{"options": [{"label": "Prepify", "value": "Prepify"}, {"label": "ZODC", "value": "ZODC"}], "tooltip": null, "helpText": null, "placeholder": "Project Name", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 1, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Project Name', 'MULTI_SELECT', 'Project Name', 5);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.203', '2025-05-07 23:50:57.203', NULL, '{}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Title', 'title', '8ea31b45-5b09-413f-8dd0-88aa55324434', 'TITLE', 'Team feedback and insights', 6);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.203', '2025-05-07 23:50:57.203', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "What Went Well", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 1, true, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'whatWentWell', 'RICHTEXT', 'What Went Well', 6);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.203', '2025-05-07 23:50:57.203', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "What Could Be Improved", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 2, true, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'whatCouldBeImproved', 'RICHTEXT', 'What Could Be Improved', 6);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.203', '2025-05-07 23:50:57.203', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Action Items", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 3, true, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'actionItems', 'RICHTEXT', 'Action Items', 6);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.396', '2025-05-07 23:50:57.396', NULL, '{}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Title', 'title', '7b01d49b-b7aa-4136-9151-0876cf40e1b3', 'TITLE', 'Project Information', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.396', '2025-05-07 23:50:57.396', NULL, '{"options": [{"label": "Prepify", "value": "Prepify"}, {"label": "ZODC", "value": "ZODC"}], "tooltip": null, "helpText": null, "placeholder": "Project Name", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 1, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Project Name', 'MULTI_SELECT', 'Project Name', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.396', '2025-05-07 23:50:57.396', NULL, '{}'::jsonb, 2, false, false, 'BASIC_FIELD', 'Title', 'title', 'd89e2a18-7176-4692-9875-cc6436f3f63f', 'TITLE', 'Candidate Information', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.396', '2025-05-07 23:50:57.396', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Candidate Name", "defaultValue": null}'::jsonb, 3, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'Candidate Name', 'TEXT', 'Candidate Name', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.396', '2025-05-07 23:50:57.396', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": "EMAIL", "placeholder": "Email Address", "defaultValue": null}'::jsonb, 3, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'Email Address', 'TEXT', 'Email Address', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.396', '2025-05-07 23:50:57.396', NULL, '{"options": [{"label": "Frontend Developer", "value": "Frontend Developer"}, {"label": "Backend Developer", "value": "Backend Developer"}, {"label": "DevOps Engineer", "value": "DevOps Engineer"}, {"label": "UI/UX Designer", "value": "UI/UX Designer"}, {"label": "QA Engineer (Tester)", "value": "QA Engineer (Tester)"}, {"label": "Business Analyst (BA)", "value": "Business Analyst (BA)"}], "tooltip": null, "helpText": null, "placeholder": "Job title", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 4, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Job title', 'MULTI_SELECT', 'Job title', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.396', '2025-05-07 23:50:57.396', NULL, '{"options": [{"label": "Fresher (0-1)", "value": "Fresher (0-1)"}, {"label": "Junior (1-2)", "value": "Junior (1-2)"}, {"label": "Mid-level (3-4)", "value": "Mid-level (3-4)"}, {"label": "Senior (5+)", "value": "Senior (5+)"}], "tooltip": null, "vertical": false, "defaultOptionIndex": null}'::jsonb, 5, true, false, 'OPTION_FIELD', 'Radio', 'radio_button_checked', 'Years of Experience', 'RADIO', 'Years of Experience', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.396', '2025-05-07 23:50:57.396', NULL, '{"tooltip": null, "helpText": null, "placeholder": "Interview Date", "validations": [{"validation": "IN_THE_PAST"}], "defaultValue": null, "disableWeekdays": null}'::jsonb, 6, true, false, 'DATE_TIME_FIELD', 'Date', 'today', 'Interview Date', 'DATE', 'Interview Date', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.396', '2025-05-07 23:50:57.396', NULL, '{"tooltip": null, "helpText": null, "allowedFileType": ["pdf"]}'::jsonb, 7, true, false, 'ADVANCED_FIELD', 'Attachment', 'attach_file', 'Resume File Upload', 'ATTACHMENT', 'Resume File Upload', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.396', '2025-05-07 23:50:57.396', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Interview Notes", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 8, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'Interview Notes', 'RICHTEXT', 'Interview Notes', 7);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '{}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Title', 'title', '4cb8efba-2672-4bbe-99e1-9ecf1f682fd9', 'TITLE', '1. Basic Information', 8);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Training Session Title", "defaultValue": ""}'::jsonb, 1, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'Training Session Title', 'TEXT', 'Training Session Title', 8);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '{"tooltip": null, "helpText": null, "placeholder": "Training Date", "validations": [], "defaultValue": null, "disableWeekdays": null}'::jsonb, 2, false, false, 'DATE_TIME_FIELD', 'Date', 'today', 'Training Date', 'DATE', 'Training Date', 8);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '{"options": [{"label": "Trần, Mạnh Hùng", "value": "Trần, Mạnh Hùng"}, {"label": "Lê, Bảo", "value": "Lê, Bảo"}, {"label": "Nguyễn, Hải Đăng ", "value": "Nguyễn, Hải Đăng "}, {"label": "Lê, Thị Mẫn Nhi", "value": "Lê, Thị Mẫn Nhi"}, {"label": "Dương, Ngọc Mạnh", "value": "Dương, Ngọc Mạnh"}, {"label": "Đỗ, Quân ", "value": "Đỗ, Quân "}, {"label": "Văn , Phú Hòa", "value": "Văn , Phú Hòa"}, {"label": "Lâm, Thị Ngọc Hân", "value": "Lâm, Thị Ngọc Hân"}, {"label": "Trần, Tấn Thành", "value": "Trần, Tấn Thành"}, {"label": "Dương, Hoàng Nam", "value": "Dương, Hoàng Nam"}], "tooltip": null, "helpText": null, "placeholder": "Trainer Name", "validations": [], "inputNewOption": null, "preventInputDuplicateOption": true}'::jsonb, 3, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Trainer Name', 'MULTI_SELECT', 'Trainer Name', 8);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '{"options": [{"label": "Trần, Mạnh Hùng", "value": "Trần, Mạnh Hùng"}, {"label": "Lê, Bảo", "value": "Lê, Bảo"}, {"label": "Nguyễn, Hải Đăng ", "value": "Nguyễn, Hải Đăng "}, {"label": "Lê, Thị Mẫn Nhi", "value": "Lê, Thị Mẫn Nhi"}, {"label": "Dương, Ngọc Mạnh", "value": "Dương, Ngọc Mạnh"}, {"label": "Đỗ, Quân ", "value": "Đỗ, Quân "}, {"label": "Văn , Phú Hòa", "value": "Văn , Phú Hòa"}, {"label": "Lâm, Thị Ngọc Hân", "value": "Lâm, Thị Ngọc Hân"}, {"label": "Trần, Tấn Thành", "value": "Trần, Tấn Thành"}, {"label": "Dương, Hoàng Nam", "value": "Dương, Hoàng Nam"}], "tooltip": null, "helpText": null, "placeholder": "Participant(s) ", "validations": [], "inputNewOption": false, "preventInputDuplicateOption": true}'::jsonb, 4, true, false, 'OPTION_FIELD', 'Multi-select', 'control_point_duplicate', 'Participant(s) ', 'MULTI_SELECT', 'Participant(s) ', 8);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '{}'::jsonb, 5, false, false, 'BASIC_FIELD', 'Title', 'title', '52360b90-c9f2-4476-ab93-df705d9ae3fc', 'TITLE', '2. Session Results', 8);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Key Topics Covered", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 6, true, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'Key Topics Covered', 'RICHTEXT', 'Key Topics Covered', 8);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "placeholder": "Attendance Rate (%)", "validations": [{"validation": "GREATER_THAN", "valueValidation": "0"}], "defaultValue": null, "decimalPlaces": null}'::jsonb, 7, true, false, 'BASIC_FIELD', 'Number', 'looks_one', 'Attendance Rate (%)', 'NUMBER', 'Attendance Rate (%)', 8);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '{"options": [{"label": "High", "value": "High"}, {"label": "Medium", "value": "Medium"}, {"label": "Low", "value": "Low"}], "tooltip": null, "helpText": null, "placeholder": "Participant Engagement "}'::jsonb, 7, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'Participant Engagement ', 'DROPDOWN', 'Participant Engagement ', 8);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '{}'::jsonb, 8, false, false, 'BASIC_FIELD', 'Title', 'title', '1d3b683f-3daf-4e7a-b8c0-fe212f4f9c3c', 'TITLE', '3. Issues & Feedback', 8);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Issues Encountered ", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 9, true, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'Issues Encountered ', 'RICHTEXT', 'Issues Encountered ', 8);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '{"toolbar": [["bold", "italic"], ["underline", "strike"], ["code", "blockquote"], ["ordered_list", "bullet_list"], ["h1", "h2", "h3", "h4", "h5", "h6"], ["link", "image"], ["text_color", "background_color"], ["align_left", "align_center", "align_right", "align_justify"]], "tooltip": null, "helpText": null, "placeholder": "Suggestions for Improvement", "isShowToolbar": true, "defaultContent": ""}'::jsonb, 10, false, false, 'ADVANCED_FIELD', 'Richtext', 'format_shapes', 'Suggestions for Improvement', 'RICHTEXT', 'Suggestions for Improvement', 8);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '{}'::jsonb, 11, false, false, 'BASIC_FIELD', 'Title', 'title', '0c621767-9e5b-4b9a-a3a0-b1322eb6bcff', 'TITLE', '4. Attachments', 8);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.586', '2025-05-07 23:50:57.586', NULL, '{"tooltip": null, "helpText": null, "allowedFileType": []}'::jsonb, 12, false, false, 'ADVANCED_FIELD', 'Attachment', 'attach_file', 'Supporting Documents / Slides / Photos', 'ATTACHMENT', 'Supporting Documents / Slides / Photos', 8);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.784', '2025-05-07 23:50:57.784', NULL, '{}'::jsonb, 0, false, false, 'BASIC_FIELD', 'Title', 'title', 'd88717e9-5be8-4d24-9c0d-b619e2a5a2a3', 'TITLE', 'Personal Information', 9);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.784', '2025-05-07 23:50:57.784', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Employee Name", "defaultValue": null}'::jsonb, 1, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'Employee Name', 'TEXT', 'Employee Name', 9);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.784', '2025-05-07 23:50:57.784', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": null, "placeholder": "Employee ID", "defaultValue": null}'::jsonb, 2, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'Employee ID', 'TEXT', 'Employee ID', 9);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.784', '2025-05-07 23:50:57.784', NULL, '{"tooltip": null, "helpText": null, "maxChars": null, "validation": "EMAIL", "placeholder": "Email", "defaultValue": null}'::jsonb, 3, true, false, 'BASIC_FIELD', 'Text', 'format_color_text', 'Email', 'TEXT', 'Email', 9);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.784', '2025-05-07 23:50:57.784', NULL, '{"options": [{"label": "Prepify", "value": "Prepify"}, {"label": "ZODC", "value": "ZODC"}], "tooltip": null, "helpText": null, "placeholder": "Project Name"}'::jsonb, 4, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'Project Name', 'DROPDOWN', 'Project Name', 9);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.784', '2025-05-07 23:50:57.784', NULL, '{}'::jsonb, 5, false, false, 'BASIC_FIELD', 'Title', 'title', 'db45e5a3-ff84-48ee-88b2-f95fc8d8ca30', 'TITLE', 'Leave Details', 9);
INSERT INTO public.form_template_fields
(created_at, updated_at, deleted_at, advanced_options, col_num, required, readonly, category, title, icon, field_id, field_type, field_name, form_template_version_id)
VALUES('2025-05-07 23:50:57.784', '2025-05-07 23:50:57.784', NULL, '{"options": [{"label": "Annual Leave", "value": "Annual Leave"}, {"label": "Unpaid Leave", "value": "Unpaid Leave"}, {"label": "Sick Leave", "value": "Sick Leave"}, {"label": "Maternity Leave", "value": "Maternity Leave"}], "tooltip": null, "helpText": null, "placeholder": "Leave Type"}'::jsonb, 6, true, false, 'OPTION_FIELD', 'Dropdown', 'keyboard_arrow_down', 'Leave Type', 'DROPDOWN', 'Leave Type', 9);

-- +goose Down
