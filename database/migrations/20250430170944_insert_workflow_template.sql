-- +goose Up
INSERT INTO public.workflows
(created_at, updated_at, deleted_at, user_id, title, "type", category_id, description, decoration, project_key, current_version, is_archived)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 37, 'Recruitment Workflow 2025', 'GENERAL', 1, 'Recruit project members as required by the product owner', 'account_tree|--primary-40', '', 1, false);
INSERT INTO public.workflows
(created_at, updated_at, deleted_at, user_id, title, "type", category_id, description, decoration, project_key, current_version, is_archived)
VALUES('2025-05-01 00:37:53.123', '2025-05-01 00:37:53.123', NULL, 37, 'Retrospective Workflow', 'GENERAL', 5, 'Retrospective Workflow', 'manage_history|--primary-40', '', 1, false);



INSERT INTO public.workflow_versions
(created_at, updated_at, deleted_at, "version", has_sub_workflow, workflow_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 1, false, 1);
INSERT INTO public.workflow_versions
(created_at, updated_at, deleted_at, "version", has_sub_workflow, workflow_id)
VALUES('2025-05-01 00:37:53.123', '2025-05-01 00:37:53.123', NULL, 1, false, 2);




INSERT INTO public.requests
(created_at, updated_at, deleted_at, user_id, last_update_user_id, status, title, is_template, sprint_id, parent_id, progress, started_at, completed_at, canceled_at, terminated_at, workflow_version_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 37, 37, 'TO_DO', 'Recruitment Workflow 2025', true, NULL, NULL, 0.0, NULL, NULL, NULL, NULL, 1);
INSERT INTO public.requests
(created_at, updated_at, deleted_at, user_id, last_update_user_id, status, title, is_template, sprint_id, parent_id, progress, started_at, completed_at, canceled_at, terminated_at, workflow_version_id)
VALUES('2025-05-01 00:37:53.123', '2025-05-01 00:37:53.123', NULL, 37, 37, 'TO_DO', 'Retrospective Workflow', true, NULL, NULL, 0.0, NULL, NULL, NULL, NULL, 2);





INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('6d51a99d-adcb-4e6c-b907-b98d5b51da1c', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 194.1999192010777, 265.99997238885845, 33, 44, 1, NULL, 'Start Recruitment Request', 37, NULL, 'START', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, '', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('fba094ea-49c6-4f07-9160-413fcd0b4952', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 597.9875011444092, 49.19999694824219, 37, 44, 2, NULL, 'PO Inputs Hiring Requirements', 43, NULL, 'INPUT', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, '2025-04-30 17:30:00.000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, '', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('e4d3d147-9526-43d0-8731-6ea7c2f274df', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 1311.21875, 193.1999969482422, 60, 44, 3, NULL, 'HR Reviews Recruitment Request', 41, NULL, 'APPROVAL', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, '2025-04-30 17:30:00.000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, '', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('b0eea977-cbe0-4301-af89-1a39b2525938', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 1893.999961671374, 114.99999909173931, 66, 40, 4, NULL, 'Is Approve', 37, NULL, 'CONDITION', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, '', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('706aa8e5-b904-485f-8f26-e8663779e59b', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 2363.999966757637, -195.99997929164277, 37, 44, 5, NULL, 'PO Revises and Resubmits', 43, NULL, 'INPUT', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, '2025-04-30 17:30:00.000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, '', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('1e3ab8c8-fc8c-4b93-bd91-aa26a08c0d0d', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 2942.9999382382334, -282.00000762939294, 60, 44, 7, NULL, 'HR Re-Reviews Revised Recruitment Request', 41, NULL, 'APPROVAL', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, '2025-04-30 17:30:00.000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, '', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('62659200-c59f-4693-80b9-aba08542f20c', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 3528.9998808361197, -387.99997929164124, 66, 40, 9, NULL, 'Is Approve?', 37, NULL, 'CONDITION', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, '', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('bb9a8b66-c1f7-4216-9776-30fc96c13426', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 3945.924889791577, -517.0000254313136, 26, 44, 12, NULL, 'Terminate due to HR reject', 37, NULL, 'END', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, 'TERMINATE', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('eeb1abcf-9c16-48d9-ad11-a121c8ac67c1', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 2868.0000908261236, 345.9999542236337, 60, 44, 6, NULL, 'ODC Manager Reviews Recruitment Request', 37, NULL, 'APPROVAL', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, '2025-04-30 17:30:00.000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, '', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('da8bd1bb-dc9b-4e67-815d-5e32c05d29ac', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 3493.0000392368825, 458.9999956403477, 66, 40, 8, NULL, 'Is Approve?', 37, NULL, 'CONDITION', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, '', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('99fdeccb-a46b-4e7a-a42c-02a15c87a366', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 4050.999972207199, 755.0000784737758, 37, 44, 11, NULL, 'HR Uploads Candidate Resumes', 41, NULL, 'INPUT', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, '2025-04-30 17:30:00.000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, '', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('e89c381a-9714-429e-9d05-c88fec5fdbc8', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 3945.849866140452, 209.00002179827243, 26, 44, 10, NULL, 'Terminate due to ODC Manager reject', 37, NULL, 'END', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, 'TERMINATE', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('8d94fe82-2c18-460d-952c-92d7081f1c32', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 4730.999957130062, 475.0000119890517, 60, 44, 13, NULL, 'PO Reviews Candidate Resumes', 43, NULL, 'APPROVAL', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, '2025-04-30 17:30:00.000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, '', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('6d487d02-1587-4173-a2e5-69995b5575cc', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 5322.999896094906, 470.99999673026264, 66, 40, 14, NULL, 'Is Approve?', 37, NULL, 'CONDITION', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, '', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('fb9e3c7c-8083-4bcc-972e-3bb64b1f7f80', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 5771.762554986129, 265.99989918300435, 26, 44, 15, NULL, 'Terminate due to PO reject', 37, NULL, 'END', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, 'TERMINATE', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('d45a7cb1-4197-41eb-a714-6025db015860', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 5840.999873206722, 754.9999623979872, 80, 44, 16, NULL, 'Sends Account Setup Notification', 37, NULL, 'NOTIFICATION', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '<p>The recruitment process has been successfully completed, and new members have been onboarded to the ODC team. Please proceed with setting up their accounts to ensure a smooth integration into the system.</p><p>🔹 Required Actions:</p><p>Create user accounts and grant necessary system access.</p><p>Assign appropriate permissions based on the role.</p><p>Share login credentials and onboarding instructions with the new members.</p><p>If you have any questions or require further details, please feel free to reach out.</p><p>Best regards,</p><p>ZODC System</p><p>Schaeffler ODC Team</p>', 'New Team Member Account Setup Required', '[]', '["duonghoangnam503@gmail.com"]', '[]', true, false, false, false, '', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('2917fc7c-625c-40eb-8ecd-026014f7a35a', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 6279.893494651425, 484.99997384207904, 26, 44, 17, NULL, 'Recruitment Process Completed', 37, NULL, 'END', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, 'COMPLETE', false, false, false, false, false, false, NULL, 1, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('d2acd2e0-d214-47d2-a308-b1ec36f3c16b', '2025-05-01 00:37:53.123', '2025-05-01 00:37:53.123', NULL, 215.83747100830078, 354.99998474121094, 33, 44, 1, NULL, 'Start Event', 37, NULL, 'START', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, '', false, false, false, false, false, false, NULL, 2, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('60d14554-70b7-4cb4-8e8f-4ab62f7489d8', '2025-05-01 00:37:53.123', '2025-05-01 00:37:53.123', NULL, 449.00001335144043, 111, 37, 44, 2, NULL, 'Input Node', 37, NULL, 'INPUT', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, '2025-04-30 17:30:00.000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, '', false, false, false, false, false, false, NULL, 2, NULL, NULL);
INSERT INTO public.nodes
(id, created_at, updated_at, deleted_at, x, y, width, height, "level", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, description, attach_file, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_send_approved_form, is_send_rejected_form, is_approved, is_rejected, end_type, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES('c17bbb2a-f7f5-4adf-aa31-229d3e697614', '2025-05-01 00:37:53.123', '2025-05-01 00:37:53.123', NULL, 967.7812261581421, 354.99998474121094, 26, 44, 3, NULL, 'End Event', 37, NULL, 'END', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, false, false, 'COMPLETE', false, false, false, false, false, false, NULL, 2, NULL, NULL);




INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('524eca1f-d6a8-4715-adf4-063bdbc491d5', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, '6d51a99d-adcb-4e6c-b907-b98d5b51da1c', 'fba094ea-49c6-4f07-9160-413fcd0b4952', NULL, false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('c0709c98-dc40-4cbe-9cbd-f6d4cdcee6ce', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 'fba094ea-49c6-4f07-9160-413fcd0b4952', 'e4d3d147-9526-43d0-8731-6ea7c2f274df', NULL, false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('5c325d4b-f947-42ff-9bc3-a943740b7b16', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 'e4d3d147-9526-43d0-8731-6ea7c2f274df', 'b0eea977-cbe0-4301-af89-1a39b2525938', NULL, false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('1ba48594-ca30-454a-bc4e-e4b3d319cf3b', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 'b0eea977-cbe0-4301-af89-1a39b2525938', '706aa8e5-b904-485f-8f26-e8663779e59b', 'FALSE', false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('58f2a8e4-1474-46f0-a802-2a830fef689c', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, '706aa8e5-b904-485f-8f26-e8663779e59b', '1e3ab8c8-fc8c-4b93-bd91-aa26a08c0d0d', NULL, false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('61deceec-f8a9-400d-bf7e-4bbccf785af9', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, '1e3ab8c8-fc8c-4b93-bd91-aa26a08c0d0d', '62659200-c59f-4693-80b9-aba08542f20c', NULL, false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('1faa45dc-32cd-45b5-a494-c2479a573d63', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, '62659200-c59f-4693-80b9-aba08542f20c', 'bb9a8b66-c1f7-4216-9776-30fc96c13426', 'FALSE', false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('d69ec24c-0da4-4dcd-833c-55c294188d90', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 'b0eea977-cbe0-4301-af89-1a39b2525938', 'eeb1abcf-9c16-48d9-ad11-a121c8ac67c1', 'TRUE', false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('ea945838-c4d3-4094-bf25-f61a7c2d05a4', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, '62659200-c59f-4693-80b9-aba08542f20c', 'eeb1abcf-9c16-48d9-ad11-a121c8ac67c1', 'TRUE', false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('efb23d1d-d27e-483b-b0fa-9a228bc3fb5b', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 'eeb1abcf-9c16-48d9-ad11-a121c8ac67c1', 'da8bd1bb-dc9b-4e67-815d-5e32c05d29ac', NULL, false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('4387120b-ed49-47bb-a701-81fefaaaa7d9', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 'da8bd1bb-dc9b-4e67-815d-5e32c05d29ac', 'e89c381a-9714-429e-9d05-c88fec5fdbc8', 'FALSE', false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('f391bf70-021a-4640-b964-40db3cf91d34', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 'da8bd1bb-dc9b-4e67-815d-5e32c05d29ac', '99fdeccb-a46b-4e7a-a42c-02a15c87a366', 'TRUE', false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('c30df386-397e-4b75-856d-1f4f5b7db9a4', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, '99fdeccb-a46b-4e7a-a42c-02a15c87a366', '8d94fe82-2c18-460d-952c-92d7081f1c32', NULL, false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('f968cdc1-a3bc-48fc-8664-20e477c6b675', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, '8d94fe82-2c18-460d-952c-92d7081f1c32', '6d487d02-1587-4173-a2e5-69995b5575cc', NULL, false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('56176680-3b61-4498-8673-c78aa87108be', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, '6d487d02-1587-4173-a2e5-69995b5575cc', 'fb9e3c7c-8083-4bcc-972e-3bb64b1f7f80', 'FALSE', false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('86782c79-7f04-409f-a6cf-e8e85cc11038', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, '6d487d02-1587-4173-a2e5-69995b5575cc', 'd45a7cb1-4197-41eb-a714-6025db015860', 'TRUE', false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('f821465a-e011-4d20-90f7-be731326e1b6', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 'd45a7cb1-4197-41eb-a714-6025db015860', '2917fc7c-625c-40eb-8ecd-026014f7a35a', NULL, false, 1);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('bd85bd0c-a982-416c-b20e-c6cb45d5560d', '2025-05-01 00:37:53.123', '2025-05-01 00:37:53.123', NULL, 'd2acd2e0-d214-47d2-a308-b1ec36f3c16b', '60d14554-70b7-4cb4-8e8f-4ab62f7489d8', NULL, false, 2);
INSERT INTO public.connections
(id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('e93747f7-42f7-4fae-b877-021d4e4011f0', '2025-05-01 00:37:53.123', '2025-05-01 00:37:53.123', NULL, '60d14554-70b7-4cb4-8e8f-4ab62f7489d8', 'c17bbb2a-f7f5-4adf-aa31-229d3e697614', NULL, false, 2);




INSERT INTO public.form_data
(id, created_at, updated_at, deleted_at, form_template_version_id)
VALUES('45cff36e-dc65-47b8-b95c-dbd4adfc191e', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 3);
INSERT INTO public.form_data
(id, created_at, updated_at, deleted_at, form_template_version_id)
VALUES('304e3b2b-305b-4296-a1f9-063ca38be403', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 4);
INSERT INTO public.form_data
(id, created_at, updated_at, deleted_at, form_template_version_id)
VALUES('00029630-8715-41d0-8b56-8687c2c5039a', '2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 4);
INSERT INTO public.form_data
(id, created_at, updated_at, deleted_at, form_template_version_id)
VALUES('2824699f-2c4d-48fb-b058-1c1d01907dde', '2025-05-01 00:37:53.123', '2025-05-01 00:37:53.123', NULL, 6);




INSERT INTO public.node_condition_destinations
(created_at, updated_at, deleted_at, destination_node_id, is_true, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 'eeb1abcf-9c16-48d9-ad11-a121c8ac67c1', true, 'b0eea977-cbe0-4301-af89-1a39b2525938');
INSERT INTO public.node_condition_destinations
(created_at, updated_at, deleted_at, destination_node_id, is_true, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, '706aa8e5-b904-485f-8f26-e8663779e59b', false, 'b0eea977-cbe0-4301-af89-1a39b2525938');
INSERT INTO public.node_condition_destinations
(created_at, updated_at, deleted_at, destination_node_id, is_true, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 'eeb1abcf-9c16-48d9-ad11-a121c8ac67c1', true, '62659200-c59f-4693-80b9-aba08542f20c');
INSERT INTO public.node_condition_destinations
(created_at, updated_at, deleted_at, destination_node_id, is_true, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 'bb9a8b66-c1f7-4216-9776-30fc96c13426', false, '62659200-c59f-4693-80b9-aba08542f20c');
INSERT INTO public.node_condition_destinations
(created_at, updated_at, deleted_at, destination_node_id, is_true, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, '99fdeccb-a46b-4e7a-a42c-02a15c87a366', true, 'da8bd1bb-dc9b-4e67-815d-5e32c05d29ac');
INSERT INTO public.node_condition_destinations
(created_at, updated_at, deleted_at, destination_node_id, is_true, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 'e89c381a-9714-429e-9d05-c88fec5fdbc8', false, 'da8bd1bb-dc9b-4e67-815d-5e32c05d29ac');
INSERT INTO public.node_condition_destinations
(created_at, updated_at, deleted_at, destination_node_id, is_true, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 'd45a7cb1-4197-41eb-a714-6025db015860', true, '6d487d02-1587-4173-a2e5-69995b5575cc');
INSERT INTO public.node_condition_destinations
(created_at, updated_at, deleted_at, destination_node_id, is_true, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 'fb9e3c7c-8083-4bcc-972e-3bb64b1f7f80', false, '6d487d02-1587-4173-a2e5-69995b5575cc');




INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 1, 'INPUT', 'b1c323eb-d272-4f89-b73f-45056885125e', NULL, NULL, NULL, true, false, false, false, NULL, NULL, NULL, '45cff36e-dc65-47b8-b95c-dbd4adfc191e', 3, 3, 'fba094ea-49c6-4f07-9160-413fcd0b4952');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 2, 'HIDDEN', '42ba3bb4-82ba-4f4d-aa9d-a26976f582e2', NULL, NULL, NULL, true, false, false, false, NULL, NULL, NULL, NULL, 4, 4, 'fba094ea-49c6-4f07-9160-413fcd0b4952');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 1, 'VIEW', '66ca23b0-3e8e-4034-866c-3f2424e84888', 'b1c323eb-d272-4f89-b73f-45056885125e', 43, 1, false, false, false, false, NULL, NULL, NULL, '45cff36e-dc65-47b8-b95c-dbd4adfc191e', 3, 3, 'e4d3d147-9526-43d0-8731-6ea7c2f274df');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 4, 'HIDDEN', '7d790f5c-2d41-4103-9882-ab92d30abbc4', '1c5a07a8-1ef3-4dcf-9839-e324d674ed85', 41, 2, false, false, false, false, NULL, NULL, NULL, '304e3b2b-305b-4296-a1f9-063ca38be403', 4, 4, 'e4d3d147-9526-43d0-8731-6ea7c2f274df');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 3, 'HIDDEN', '07c35c60-71d8-40eb-a9d4-603aa510febc', '0109e91f-8468-449d-9668-b0b46e9f4fc7', 41, 3, false, false, false, false, NULL, NULL, NULL, '00029630-8715-41d0-8b56-8687c2c5039a', 4, 4, 'e4d3d147-9526-43d0-8731-6ea7c2f274df');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 1, 'EDIT', '61a5586f-0107-4dbf-8062-c4e58faba0ee', 'b1c323eb-d272-4f89-b73f-45056885125e', 43, 1, true, false, false, false, NULL, NULL, NULL, '45cff36e-dc65-47b8-b95c-dbd4adfc191e', 3, 3, '706aa8e5-b904-485f-8f26-e8663779e59b');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 2, 'HIDDEN', 'c00bef9b-6ac7-435a-be25-c8f23bb40b5b', NULL, NULL, NULL, true, false, false, false, NULL, NULL, NULL, NULL, 4, 4, '706aa8e5-b904-485f-8f26-e8663779e59b');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 1, 'VIEW', 'd9475648-e95e-46ba-a110-123fbd1a0950', 'b1c323eb-d272-4f89-b73f-45056885125e', 43, 1, false, false, false, false, NULL, NULL, NULL, '45cff36e-dc65-47b8-b95c-dbd4adfc191e', 3, 3, '1e3ab8c8-fc8c-4b93-bd91-aa26a08c0d0d');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 4, 'HIDDEN', 'b0cfae1f-4164-4894-b807-95a9f7cbbcc7', '1c5a07a8-1ef3-4dcf-9839-e324d674ed85', 41, 2, false, false, false, false, NULL, NULL, NULL, '304e3b2b-305b-4296-a1f9-063ca38be403', 4, 4, '1e3ab8c8-fc8c-4b93-bd91-aa26a08c0d0d');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 3, 'HIDDEN', '8f59c0e1-7f21-40e8-8c85-b6e959dccf8a', '0109e91f-8468-449d-9668-b0b46e9f4fc7', 41, 3, false, false, false, false, NULL, NULL, NULL, '00029630-8715-41d0-8b56-8687c2c5039a', 4, 4, '1e3ab8c8-fc8c-4b93-bd91-aa26a08c0d0d');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 1, 'VIEW', 'fc80f735-07b8-4aa4-976c-f1a2eaab6d92', 'b1c323eb-d272-4f89-b73f-45056885125e', 43, 1, false, false, false, false, NULL, NULL, NULL, '45cff36e-dc65-47b8-b95c-dbd4adfc191e', 3, 3, 'eeb1abcf-9c16-48d9-ad11-a121c8ac67c1');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 4, 'HIDDEN', '89fd86fb-3de1-4c5c-9436-2cc64eba053d', '1c5a07a8-1ef3-4dcf-9839-e324d674ed85', 41, 2, false, false, false, false, NULL, NULL, NULL, '304e3b2b-305b-4296-a1f9-063ca38be403', 4, 4, 'eeb1abcf-9c16-48d9-ad11-a121c8ac67c1');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 3, 'HIDDEN', '70eb7465-ecca-4d71-8408-f62f077c095b', '0109e91f-8468-449d-9668-b0b46e9f4fc7', 41, 3, false, false, false, false, NULL, NULL, NULL, '00029630-8715-41d0-8b56-8687c2c5039a', 4, 4, 'eeb1abcf-9c16-48d9-ad11-a121c8ac67c1');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 1, 'VIEW', '1ee0506b-32eb-4c5d-9218-1746c2698478', 'b1c323eb-d272-4f89-b73f-45056885125e', 43, 1, true, false, false, false, NULL, NULL, NULL, '45cff36e-dc65-47b8-b95c-dbd4adfc191e', 3, 3, '99fdeccb-a46b-4e7a-a42c-02a15c87a366');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 2, 'INPUT', '1c5a07a8-1ef3-4dcf-9839-e324d674ed85', NULL, NULL, NULL, true, false, false, false, NULL, NULL, NULL, '304e3b2b-305b-4296-a1f9-063ca38be403', 4, 4, '99fdeccb-a46b-4e7a-a42c-02a15c87a366');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 3, 'INPUT', '0109e91f-8468-449d-9668-b0b46e9f4fc7', NULL, NULL, NULL, false, false, false, false, NULL, NULL, NULL, '00029630-8715-41d0-8b56-8687c2c5039a', 4, 4, '99fdeccb-a46b-4e7a-a42c-02a15c87a366');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 1, 'HIDDEN', '92368c18-a5d9-4cd9-bf72-aecf0e6e77c4', 'b1c323eb-d272-4f89-b73f-45056885125e', 43, 1, false, false, false, false, NULL, NULL, NULL, '45cff36e-dc65-47b8-b95c-dbd4adfc191e', 3, 3, '8d94fe82-2c18-460d-952c-92d7081f1c32');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 4, 'VIEW', 'c142bf48-e760-45fd-b5b5-34143f4aff7a', '1c5a07a8-1ef3-4dcf-9839-e324d674ed85', 41, 2, false, false, false, false, NULL, NULL, NULL, '304e3b2b-305b-4296-a1f9-063ca38be403', 4, 4, '8d94fe82-2c18-460d-952c-92d7081f1c32');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:52.898', '2025-05-01 00:37:52.898', NULL, 3, 'VIEW', '53f56c43-0709-4a6e-8e6d-2db66a5cbef2', '0109e91f-8468-449d-9668-b0b46e9f4fc7', 41, 3, false, false, false, false, NULL, NULL, NULL, '00029630-8715-41d0-8b56-8687c2c5039a', 4, 4, '8d94fe82-2c18-460d-952c-92d7081f1c32');
INSERT INTO public.node_forms
(created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, template_version_id, node_id)
VALUES('2025-05-01 00:37:53.123', '2025-05-01 00:37:53.123', NULL, 1, 'INPUT', 'fb73a10f-191b-4b39-b3c6-a349c4819260', NULL, NULL, NULL, true, false, false, false, NULL, NULL, NULL, '2824699f-2c4d-48fb-b058-1c1d01907dde', 6, 6, '60d14554-70b7-4cb4-8e8f-4ab62f7489d8');







-- +goose Down

