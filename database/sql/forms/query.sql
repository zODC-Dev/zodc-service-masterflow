-- name: FindAllForms :many
SELECT forms.*, sqlc.embed(form_fields)
FROM forms
LEFT JOIN form_fields ON forms.id = form_fields.form_id;

-- name: CreateForm :one
INSERT INTO forms (
    file_name,
    title,
    function,
    version,
    template,
    datasheet,
    description,
    decoration
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: CreateMutiFormFields :many
INSERT INTO form_fields (
    field_id,
    icon,
    title,
    category,
    field_name,
    field_type,
    required,
    advanced_options,
    col_num,
    form_id
)
SELECT
    unnest(@field_ids::text[]) as field_id,
    unnest(@icons::text[]) as icon,
    unnest(@titles::text[]) as title,
    unnest(@categories::text[]) as category,
    unnest(@field_namess::text[]) as field_name,
    unnest(@field_types::text[]) as field_type,
    unnest(@requireds::boolean[]) as required,
    unnest(@advanced_optionss::jsonb[]) as advanced_options,
    unnest(@col_nums::int[]) as col_num,
    unnest(@form_ids::int[]) as form_id
RETURNING *;

-- name: CreateFormField :one
INSERT INTO form_fields (
    field_id,
    icon,
    title,
    category,
    field_name,
    field_type,
    required,
    advanced_options,
    col_num,
    form_id
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;