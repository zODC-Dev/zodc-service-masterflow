-- +goose Up
INSERT INTO public.workflows
    (id, created_at, updated_at, deleted_at, user_id, title, "type", category_id, description, decoration, project_key, current_version, is_archived)
VALUES
    (1, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 35, 'Requestor Information', 'GENERAL', 1, 'Requestor Name', 'person_add|--success-40', '', 1, false);

SELECT setval('workflows_id_seq', (SELECT MAX(id) FROM workflows));

-- Workflow Template Version
INSERT INTO public.workflow_versions
    (id, created_at, updated_at, deleted_at, "version", has_sub_workflow, workflow_id)
VALUES
    (1, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 1, false, 1);

SELECT setval('workflow_versions_id_seq', (SELECT MAX(id) FROM workflow_versions));

-- Request
INSERT INTO public.requests
    (id, created_at, updated_at, deleted_at, user_id, "key", last_update_user_id, status, title, is_template, sprint_id, parent_id, progress, started_at, completed_at, canceled_at, terminated_at, workflow_version_id)
VALUES
    (1, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 35, 1, 35, 'TO_DO', 'Requestor Information', true, NULL, NULL, 0.0, NULL, NULL, NULL, NULL, 1);

SELECT setval('requests_id_seq', (SELECT MAX(id) FROM requests));

-- Nodes
INSERT INTO public.nodes
    (id, created_at, updated_at, deleted_at, x, y, width, height, "level", "key", jira_key, title, assignee_id, sub_request_id, "type", status, is_current, jira_link_url, estimate_point, planned_start_time, planned_end_time, actual_start_time, actual_end_time, body, subject, cc_emails, to_emails, bcc_emails, is_approved, is_rejected, end_type, task_assigned_requester, task_assigned_assignee, task_assigned_participants, task_started_requester, task_started_assignee, task_started_participants, task_completed_requester, task_completed_assignee, task_completed_participants, parent_id, request_id, form_template_id, form_data_id)
VALUES
    ('734b050b-5251-4ce8-8674-92c7d76ba7d7', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, -492.16406250000034, 389.99999999999994, 33, 44, 1, 1, NULL, 'Start Event', 35, NULL, 'START', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, '', false, false, false, false, false, false, false, false, false, NULL, 1, NULL, NULL),
    ('506bd7ab-2647-4029-b8e0-0f7739255337', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, -266.0000000000003, 333.9999999999998, 37, 44, 2, 2, NULL, 'PO Inputs Hiring Requirements', 34, NULL, 'INPUT', 'TO_DO', false, NULL, NULL, NULL, '2025-04-20 17:30:00.000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, '', false, false, false, false, false, false, false, false, false, NULL, 1, NULL, NULL),
    ('91133054-eb05-4b13-87e1-b96185c10023', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 242.00000424786674, 321.9999911140285, 60, 44, 3, 3, NULL, 'HR Reviews Recruitment Request', 41, NULL, 'APPROVAL', 'TO_DO', false, NULL, NULL, NULL, '2025-04-20 17:30:00.000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, '', false, false, false, false, false, false, false, false, false, NULL, 1, NULL, NULL),
    ('21eef422-4bc8-4c62-b6ac-e9c3f087451b', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 2056.9999605395424, 670.0000505293033, 37, 44, 8, 4, NULL, 'HR Uploads Candidate Resumes', 41, NULL, 'INPUT', 'TO_DO', false, NULL, NULL, NULL, '2025-04-20 17:30:00.000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, '', false, false, false, false, false, false, false, false, false, NULL, 1, NULL, NULL),
    ('e93a4eef-78f0-4320-af14-85c85fdb05c8', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 1047.999981611201, 669.999996258216, 60, 44, 5, 5, NULL, 'ODC Manager Reviews Recruitment Request', 37, NULL, 'APPROVAL', 'TO_DO', false, NULL, NULL, NULL, '2025-04-20 17:30:00.000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, '', false, false, false, false, false, false, false, false, false, NULL, 1, NULL, NULL),
    ('142ac5be-ad17-4a71-afb7-5dc05e3b61d7', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 788.9999928037751, 317.999983484634, 66, 40, 4, 6, NULL, 'Is Approve?', 35, NULL, 'CONDITION', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, '', false, false, false, false, false, false, false, false, false, NULL, 1, NULL, NULL),
    ('dc400afb-3061-4dfb-b375-ab8dd4982a8e', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 1708.0000058128714, 665.9999800611879, 66, 40, 7, 7, NULL, 'Is Approve?', 35, NULL, 'CONDITION', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, '', false, false, false, false, false, false, false, false, false, NULL, 1, NULL, NULL),
    ('c9769f88-ebb1-4838-adbb-ed8db3c94e37', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 2590, 670, 60, 44, 10, 8, NULL, 'PO Reviews Candidate Resumes', 34, NULL, 'APPROVAL', 'TO_DO', false, NULL, NULL, NULL, '2025-04-20 17:30:00.000', NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, '', false, false, false, false, false, false, false, false, false, NULL, 1, NULL, NULL),
    ('9b1aab53-4868-40e0-8b2e-195028461357', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 3140.9999736185027, 666.0000577953936, 66, 40, 11, 9, NULL, 'Is Approve?', 35, NULL, 'CONDITION', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, '', false, false, false, false, false, false, false, false, false, NULL, 1, NULL, NULL),
    ('287e1f89-282d-4bbe-90a9-f7c2156893a2', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 3490.000000000001, 563, 80, 44, 12, 10, NULL, 'Sends Account Setup Notification', 35, NULL, 'NOTIFICATION', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, '<p>The recruitment process has been successfully completed, and new members have been onboarded to the ODC team. Please proceed with setting up their accounts to ensure a smooth integration into the system.<br>🔹 Required Actions:<br>Create user accounts and grant necessary system access.<br>Assign appropriate permissions based on the role.<br>Share login credentials and onboarding instructions with the new members.<br>If you have any questions or require further details, please feel free to reach out.<br>Best regards,<br>ZODC System<br>Schaeffler ODC Team</p>', 'New Team Member Account Setup Required "', '[]', '["duonghoangnam503@gmail.com"]', '[]', false, false, '', false, false, false, false, false, false, false, false, false, NULL, 1, NULL, NULL),
    ('03a4d89c-8bea-453e-afdf-9f8f96dc10bc', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 3882.78125, 563, 26, 44, 14, 11, NULL, 'End Event', 35, NULL, 'END', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, 'COMPLETE', false, false, false, false, false, false, false, false, false, NULL, 1, NULL, NULL),
    ('554ff595-5077-47a7-bc4f-815018e65b91', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 3520.796874999999, 914.9999999999999, 26, 44, 13, 12, NULL, 'Terminate due to PO reject		', 35, NULL, 'END', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, 'TERMINATE', false, false, false, false, false, false, false, false, false, NULL, 1, NULL, NULL),
    ('58256941-a0da-4015-ac45-e9c781b78740', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 1991.8436931009237, 1021.0000883129715, 26, 44, 9, 13, NULL, 'Terminate due to ODC Manager reject		', 35, NULL, 'END', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, 'TERMINATE', false, false, false, false, false, false, false, false, false, NULL, 1, NULL, NULL),
    ('4d8839cc-844d-441f-9853-9b11d6db9d85', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 1127.8509043040303, 108.99527139527197, 26, 44, 6, 14, NULL, 'Terminate due to HR reject		', 35, NULL, 'END', 'TO_DO', false, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, false, 'TERMINATE', false, false, false, false, false, false, false, false, false, NULL, 1, NULL, NULL);

-- Connections
INSERT INTO public.connections
    (id, created_at, updated_at, deleted_at, from_node_id, to_node_id, "text", is_completed, request_id)
VALUES('1ef3d797-377d-4bb2-85ed-9503a8d88c7d', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, '734b050b-5251-4ce8-8674-92c7d76ba7d7', '506bd7ab-2647-4029-b8e0-0f7739255337', NULL, false, 1),
    ('c6fcfcd0-f733-4e1b-9f21-03cc8c54856d', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, '506bd7ab-2647-4029-b8e0-0f7739255337', '91133054-eb05-4b13-87e1-b96185c10023', NULL, false, 1),
    ('71284ed0-26b9-4fad-b69c-10a22c294ed6', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, '91133054-eb05-4b13-87e1-b96185c10023', '142ac5be-ad17-4a71-afb7-5dc05e3b61d7', NULL, false, 1),
    ('f1d54fb1-818e-4bd9-963f-ac9c44c2edee', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, '142ac5be-ad17-4a71-afb7-5dc05e3b61d7', 'e93a4eef-78f0-4320-af14-85c85fdb05c8', 'TRUE', false, 1),
    ('5801fc2b-12ea-4db0-bc11-b9f49fd9aea1', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 'e93a4eef-78f0-4320-af14-85c85fdb05c8', 'dc400afb-3061-4dfb-b375-ab8dd4982a8e', NULL, false, 1),
    ('0aa31e86-6d81-4b67-a556-927ebae396d3', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 'dc400afb-3061-4dfb-b375-ab8dd4982a8e', '21eef422-4bc8-4c62-b6ac-e9c3f087451b', 'TRUE', false, 1),
    ('f5dc4ea4-ca11-4b9f-925b-e3cbe41c2138', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, '21eef422-4bc8-4c62-b6ac-e9c3f087451b', 'c9769f88-ebb1-4838-adbb-ed8db3c94e37', NULL, false, 1),
    ('57dd2a57-0743-4943-bf04-1bda5f71d386', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 'c9769f88-ebb1-4838-adbb-ed8db3c94e37', '9b1aab53-4868-40e0-8b2e-195028461357', NULL, false, 1),
    ('3e96fe3c-fbdd-4eac-b107-f52d8b0f5f12', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, '9b1aab53-4868-40e0-8b2e-195028461357', '287e1f89-282d-4bbe-90a9-f7c2156893a2', 'TRUE', false, 1),
    ('e41f708a-c0de-4426-8103-515469900f7d', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, '287e1f89-282d-4bbe-90a9-f7c2156893a2', '03a4d89c-8bea-453e-afdf-9f8f96dc10bc', NULL, false, 1),
    ('2f220313-31ff-4dd7-837a-58ec831e021b', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, '9b1aab53-4868-40e0-8b2e-195028461357', '554ff595-5077-47a7-bc4f-815018e65b91', 'FALSE', false, 1),
    ('8319a8e6-34f2-4c18-bc95-e96b02ea1e3c', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 'dc400afb-3061-4dfb-b375-ab8dd4982a8e', '58256941-a0da-4015-ac45-e9c781b78740', 'FALSE', false, 1),
    ('5b926bbd-6777-4adb-9058-e6bbeafdb598', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, '142ac5be-ad17-4a71-afb7-5dc05e3b61d7', '4d8839cc-844d-441f-9853-9b11d6db9d85', 'FALSE', false, 1);

-- Form Data
INSERT INTO public.form_data
    (id, created_at, updated_at, deleted_at, form_template_version_id)
VALUES
    ('dbed42b2-16f3-4646-b7cc-500e26bd4850', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 2),
    ('ee02a0f6-0880-4da6-8e2c-f6139914169e', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 4),
    ('78b1fe9e-86e4-4bc0-a0aa-8467aa2a2d44', '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 4);

-- Node Form
INSERT INTO public.node_forms
    (id, created_at, updated_at, deleted_at, "level", "permission", "key", option_key, from_user_id, from_form_attached_position, is_original, is_approved, is_rejected, is_submitted, submitted_at, submitted_by_user_id, last_update_user_id, data_id, template_id, node_id)
VALUES
    (1, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 1, 'INPUT', 'a47488d7-182c-40fc-9da2-b100cc312591', NULL, NULL, NULL, true, false, false, false, NULL, NULL, NULL, 'dbed42b2-16f3-4646-b7cc-500e26bd4850', 2, '506bd7ab-2647-4029-b8e0-0f7739255337'),
    (2, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 2, 'HIDDEN', '30c7dcc6-89a9-47ec-af64-771017134a75', NULL, NULL, NULL, true, false, false, false, NULL, NULL, NULL, NULL, 4, '506bd7ab-2647-4029-b8e0-0f7739255337'),
    (3, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 1, 'VIEW', '0b4e174e-c61e-441e-97b3-d43829911406', 'a47488d7-182c-40fc-9da2-b100cc312591', 34, 1, false, false, false, false, NULL, NULL, NULL, 'dbed42b2-16f3-4646-b7cc-500e26bd4850', 2, '91133054-eb05-4b13-87e1-b96185c10023'),
    (4, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 3, 'VIEW', '75200fda-58fb-409c-aedd-98a94aebf5bb', '0cedd38d-ef74-4f47-9e28-e2c59b0707c1', 41, 2, false, false, false, false, NULL, NULL, NULL, 'ee02a0f6-0880-4da6-8e2c-f6139914169e', 4, '91133054-eb05-4b13-87e1-b96185c10023'),
    (5, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 3, 'VIEW', '864855c7-23f2-4052-85f0-4aa6e1303647', 'c765c1db-b122-4f4b-b652-b9906b9d11e8', 41, 3, false, false, false, false, NULL, NULL, NULL, '78b1fe9e-86e4-4bc0-a0aa-8467aa2a2d44', 4, '91133054-eb05-4b13-87e1-b96185c10023'),
    (6, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 1, 'VIEW', 'ba013839-2de5-451c-bd7d-9742436e79db', 'a47488d7-182c-40fc-9da2-b100cc312591', 34, 1, true, false, false, false, NULL, NULL, NULL, 'dbed42b2-16f3-4646-b7cc-500e26bd4850', 2, '21eef422-4bc8-4c62-b6ac-e9c3f087451b'),
    (7, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 2, 'INPUT', '0cedd38d-ef74-4f47-9e28-e2c59b0707c1', NULL, NULL, NULL, true, false, false, false, NULL, NULL, NULL, 'ee02a0f6-0880-4da6-8e2c-f6139914169e', 4, '21eef422-4bc8-4c62-b6ac-e9c3f087451b'),
    (8, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 3, 'INPUT', 'c765c1db-b122-4f4b-b652-b9906b9d11e8', NULL, NULL, NULL, false, false, false, false, NULL, NULL, NULL, '78b1fe9e-86e4-4bc0-a0aa-8467aa2a2d44', 4, '21eef422-4bc8-4c62-b6ac-e9c3f087451b'),
    (9, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 1, 'VIEW', '4dd6edd7-31df-4479-a674-1205914f0c33', 'a47488d7-182c-40fc-9da2-b100cc312591', 34, 1, false, false, false, false, NULL, NULL, NULL, 'dbed42b2-16f3-4646-b7cc-500e26bd4850', 2, 'e93a4eef-78f0-4320-af14-85c85fdb05c8'),
    (10, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 3, 'VIEW', '704cd149-5fda-42b0-a103-27f0cb236114', '0cedd38d-ef74-4f47-9e28-e2c59b0707c1', 41, 2, false, false, false, false, NULL, NULL, NULL, 'ee02a0f6-0880-4da6-8e2c-f6139914169e', 4, 'e93a4eef-78f0-4320-af14-85c85fdb05c8'),
    (11, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 3, 'VIEW', '09912afa-6976-4fa9-82fd-34eb2c575c9e', 'c765c1db-b122-4f4b-b652-b9906b9d11e8', 41, 3, false, false, false, false, NULL, NULL, NULL, '78b1fe9e-86e4-4bc0-a0aa-8467aa2a2d44', 4, 'e93a4eef-78f0-4320-af14-85c85fdb05c8'),
    (12, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 1, 'HIDDEN', '79ff9348-6979-4916-93f0-84067bebf5d2', 'a47488d7-182c-40fc-9da2-b100cc312591', 34, 1, false, false, false, false, NULL, NULL, NULL, 'dbed42b2-16f3-4646-b7cc-500e26bd4850', 2, 'c9769f88-ebb1-4838-adbb-ed8db3c94e37'),
    (13, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 3, 'VIEW', '92aeb1f9-3d08-4d11-bab1-0e6fc6b4a6b1', '0cedd38d-ef74-4f47-9e28-e2c59b0707c1', 41, 2, false, false, false, false, NULL, NULL, NULL, 'ee02a0f6-0880-4da6-8e2c-f6139914169e', 4, 'c9769f88-ebb1-4838-adbb-ed8db3c94e37'),
    (14, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 3, 'VIEW', '49f0fc93-a9d0-4594-b654-caeeed3c39bf', 'c765c1db-b122-4f4b-b652-b9906b9d11e8', 41, 3, false, false, false, false, NULL, NULL, NULL, '78b1fe9e-86e4-4bc0-a0aa-8467aa2a2d44', 4, 'c9769f88-ebb1-4838-adbb-ed8db3c94e37');

SELECT setval('node_forms_id_seq', (SELECT MAX(id) FROM node_forms));

-- Node Condition Destination
INSERT INTO public.node_condition_destinations
    (id, created_at, updated_at, deleted_at, destination_node_id, is_true, node_id)
VALUES
    (1, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, 'e93a4eef-78f0-4320-af14-85c85fdb05c8', true, '142ac5be-ad17-4a71-afb7-5dc05e3b61d7'),
    (2, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, '4d8839cc-844d-441f-9853-9b11d6db9d85', false, '142ac5be-ad17-4a71-afb7-5dc05e3b61d7'),
    (3, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, '21eef422-4bc8-4c62-b6ac-e9c3f087451b', true, 'dc400afb-3061-4dfb-b375-ab8dd4982a8e'),
    (4, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, '58256941-a0da-4015-ac45-e9c781b78740', false, 'dc400afb-3061-4dfb-b375-ab8dd4982a8e'),
    (5, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, '287e1f89-282d-4bbe-90a9-f7c2156893a2', true, '9b1aab53-4868-40e0-8b2e-195028461357'),
    (6, '2025-04-20 11:24:06.432', '2025-04-20 11:24:06.432', NULL, '554ff595-5077-47a7-bc4f-815018e65b91', false, '9b1aab53-4868-40e0-8b2e-195028461357');

SELECT setval('node_condition_destinations_id_seq', (SELECT MAX(id) FROM node_condition_destinations));


-- +goose Down

