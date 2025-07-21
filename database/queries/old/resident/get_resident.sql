-- name: GetResidentFull :one
SELECT
  r.*,
  a.house_number,
  c.name AS city_name,
  k.name AS kebele_name,
  b.blood_type,
  b.face,
  d.type AS document_type,
  d.status AS document_status,
  d.url AS document_url,
  d.number AS document_number,
  e.status AS employment_status,
  e.occupation,
  e.employer_name,
  e.work_address,
  em.name AS emergency_contact_name,
  em.relation AS emergency_contact_relation,
  em.phone AS emergency_contact_phone
FROM resident r
LEFT JOIN address a ON r.address_id = a.id
LEFT JOIN city c ON a.city_id = c.id
LEFT JOIN kebele k ON a.kebele_id = k.id
LEFT JOIN biometric b ON r.id = b.resident_id
LEFT JOIN document d ON r.id = d.resident_id
LEFT JOIN employment e ON r.id = e.resident_id
LEFT JOIN emergency em ON r.id = em.resident_id
WHERE r.id = $1 AND r.deleted_at IS NULL;

-- name: SearchResidentsByName :many
SELECT id, first_name, second_name, last_name, search_vector
FROM resident
WHERE search_vector @@ to_tsquery('english', $1)
AND deleted_at IS NULL
ORDER BY ts_rank(search_vector, to_tsquery('english', $1)) DESC;