-- name: GetResidentFull :one
SELECT
  r.*,
  a.house_number,
  a.district,
  c.name AS city_name,
  rg.name AS region_name,
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
LEFT JOIN region rg ON c.region_id = rg.id
LEFT JOIN biometric b ON r.id = b.resident_id
LEFT JOIN document d ON r.id = d.resident_id
LEFT JOIN employment e ON r.id = e.resident_id
LEFT JOIN emergency em ON r.id = em.resident_id
WHERE r.id = $1 AND r.deleted_at IS NULL;
